package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/dlespiau/x86db"
)

var nasmToPlan9 = map[string]string{
	//
	// SSE
	//

	// The condition is encoded as an imm8 operand of CMPPS
	"CMPEQPS":    "CMPPS",
	"CMPLTPS":    "CMPPS",
	"CMPLEPS":    "CMPPS",
	"CMPUNORDPS": "CMPPS",
	"CMPNEQPS":   "CMPPS",
	"CMPNLTPS":   "CMPPS",
	"CMPNLEPS":   "CMPPS",
	"CMPORDPS":   "CMPPS",

	// The condition is encoded as an imm8 operand of CMPSS.
	"CMPEQSS":    "CMPSS",
	"CMPLTSS":    "CMPSS",
	"CMPLESS":    "CMPSS",
	"CMPUNORDSS": "CMPSS",
	"CMPNEQSS":   "CMPSS",
	"CMPNLTSS":   "CMPSS",
	"CMPNLESS":   "CMPSS",
	"CMPORDSS":   "CMPSS",

	"CVTSI2SS":  "CVTSL2SS",
	"CVTSS2SI":  "CVTSS2SL",
	"CVTTSS2SI": "CVTTSS2SL",

	//
	// PCLMULQDQ
	//

	"PCLMULLQLQDQ": "PCLMULQDQ",
	"PCLMULHQLQDQ": "PCLMULQDQ",
	"PCLMULLQHQDQ": "PCLMULQDQ",
	"PCLMULHQHQDQ": "PCLMULQDQ",

	//
	// SSE2
	//

	// The condition is encoded as an imm8 operand of CMPPD
	"CMPEQPD":    "CMPPD",
	"CMPLTPD":    "CMPPD",
	"CMPLEPD":    "CMPPD",
	"CMPUNORDPD": "CMPPD",
	"CMPNEQPD":   "CMPPD",
	"CMPNLTPD":   "CMPPD",
	"CMPNLEPD":   "CMPPD",
	"CMPORDPD":   "CMPPD",

	// The condition is encoded as an imm8 operand of CMPSD.
	"CMPEQSD":    "CMPSD",
	"CMPLTSD":    "CMPSD",
	"CMPLESD":    "CMPSD",
	"CMPUNORDSD": "CMPSD",
	"CMPNEQSD":   "CMPSD",
	"CMPNLTSD":   "CMPSD",
	"CMPNLESD":   "CMPSD",
	"CMPORDSD":   "CMPSD",

	// D (double word) has been replaced by L (Long)
	// DQ (double quadword) has been replaced by O (Octoword)
	"MASKMOVDQU": "MASKMOVOU",
	"MOVD":       "MOVQ",
	"MOVDQ2Q":    "MOVQ",
	"MOVDQA":     "MOVO",
	"MOVDQU":     "MOVOU",
	"MOVNTDQ":    "MOVNTO",
	"PACKSSDW":   "PACKSSLW",
	"PADDD":      "PADDL",
	"PCMPEQD":    "PCMPEQL",
	"PCMPGTD":    "PCMPGTL",
	"PMADDWD":    "PMADDWL",
	"PMULUDQ":    "PMULULQ",
	"PSLLD":      "PSLLL",
	"PSLLDQ":     "PSLLO",
	"PSRAD":      "PSRAL",
	"PSRLD":      "PSRLL",
	"PSRLDQ":     "PSRLO",
	"PSUBD":      "PSUBL",
	"PUNPCKHDQ":  "PUNPCKHLQ",
	"PUNPCKHWD":  "PUNPCKHWL",
	"PUNPCKLDQ":  "PUNPCKLLQ",
	"PUNPCKLWD":  "PUNPCKLWL",

	// Conversions
	"CVTDQ2PD":  "CVTPL2PD",
	"CVTDQ2PS":  "CVTPL2PS",
	"CVTPD2DQ":  "CVTPD2PL",
	"CVTPS2DQ":  "CVTPS2PL",
	"CVTSD2SI":  "CVTSD2SL",
	"CVTSI2SD":  "CVTSL2SD",
	"CVTTPD2DQ": "CVTTPD2PL",
	"CVTTPS2DQ": "CVTTPS2PL",
	"CVTTSD2SI": "CVTTSD2SL",
}

func nasmOpcodeToPlan9(op string) string {
	plan9, ok := nasmToPlan9[op]
	if ok {
		return plan9
	}
	return op
}

func isAlreadyKnown(insn *x86db.Instruction) bool {
	name := nasmOpcodeToPlan9(insn.Name)
	for _, opcode := range Anames {
		if name == opcode {
			return true
		}
	}
	return false
}

func isAlreadyTested(insn *x86db.Instruction) bool {
	name := nasmOpcodeToPlan9(insn.Name)
	_, ok := testedMap[name]
	return ok
}

func isMMXOperand(op string) bool {
	if op == "mmxreg" || op == "mmxrm" || op == "mmxrm64" {
		return true
	}
	return false
}

func isMMX(insn *x86db.Instruction) bool {
	for _, op := range insn.Operands {
		if isMMXOperand(op) {
			return true
		}
	}
	return false
}

func doHelp(insns x86db.InstructionSlice) {
	usage()
}

func doList(insns x86db.InstructionSlice) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	for _, insn := range insns {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", insn.Name,
			strings.Join(insn.Operands, ","), insn.Pattern, insn.Flags)
	}
	w.Flush()
}

var (
	filterFlags = flag.NewFlagSet("filter", flag.ExitOnError)
	extension   = filterFlags.String("extension", "",
		"select instructions by extension")
	notMMX = filterFlags.Bool("not-mmx", false,
		"do not select instructions taking MMX operands")
	known = filterFlags.Bool("known", false,
		"select instructions already known by the go assembler")
	notKnown = filterFlags.Bool("not-known", false,
		"select instructions not already known by the go assembler")
	tested = filterFlags.Bool("tested", false,
		"select instructions with test cases in the go assembler")
	notTested = filterFlags.Bool("not-tested", false,
		"select instructions with no test case in the go assembler")
)

type command struct {
	name    string
	help    string
	handler func(x86db.InstructionSlice)
}

var commands = []command{
	{"help", "print this help", nil},
	{"list", "list x86 instructions", doList},
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage:\n\n")
	fmt.Fprintf(os.Stderr, "  x86db-gogen command [options]\n\n")
	fmt.Fprintf(os.Stderr, "List of commands:\n\n")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	for _, cmd := range commands {
		fmt.Fprintf(w, "  %s\t%s\n", cmd.name, cmd.help)
	}
	w.Flush()
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Filtering options:\n\n")
	filterFlags.PrintDefaults()
}

func main() {
	db := x86db.NewDB()

	if err := filterFlags.Parse(os.Args[2:]); err != nil {
		log.Fatal(err)
	}

	err := db.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	insns := db.Instructions

	if *extension != "" {
		ext, err := x86db.ExtensionFromString(*extension)
		if err != nil {
			log.Fatal(err)
		}
		insns = insns.Where(func(insn x86db.Instruction) bool {
			return insn.Extension == ext
		})
	}

	if *notMMX {
		insns = insns.Where(func(insn x86db.Instruction) bool {
			return !isMMX(&insn)
		})
	}

	if *known || *notKnown {
		insns = insns.Where(func(insn x86db.Instruction) bool {
			k := isAlreadyKnown(&insn)
			if *notKnown {
				return !k
			}
			return k
		})
	}

	if *tested || *notTested {
		insns = insns.Where(func(insn x86db.Instruction) bool {
			t := isAlreadyTested(&insn)
			if *notTested {
				return !t
			}
			return t
		})
	}

	cmdName := os.Args[1]
	handled := false
	for _, cmd := range commands {
		if cmd.name != cmdName {
			continue
		}
		if cmd.handler == nil {
			continue
		}

		handled = true
		cmd.handler(insns)
	}

	if !handled {
		fmt.Fprintf(os.Stderr, "unknown command '%s'\n\n", cmdName)
		usage()
		os.Exit(1)
	}
}
