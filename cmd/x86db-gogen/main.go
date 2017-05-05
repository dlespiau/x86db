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

func nasmOpcodeToPlan9(op string) string {
	/*
	 * SSE
	 */
	switch op {
	// The condition is encoded as an imm8 operand of CMPPS
	case "CMPEQPS", "CMPLTPS", "CMPLEPS", "CMPUNORDPS", "CMPNEQPS", "CMPNLTPS",
		"CMPNLEPS", "CMPORDPS":
		return "CMPPS"
	// The condition is encoded as an imm8 operand of CMPSS.
	case "CMPEQSS", "CMPLTSS", "CMPLESS", "CMPUNORDSS", "CMPNEQSS", "CMPNLTSS",
		"CMPNLESS", "CMPORDSS":
		return "CMPSS"
	case "CVTSI2SS":
		return "CVTSL2SS"
	case "CVTSS2SI":
		return "CVTSS2SL"
	case "CVTTSS2SI":
		return "CVTTSS2SL"
	}

	/*
	 * PCLMULQDQ
	 */
	switch op {
	case "PCLMULLQLQDQ", "PCLMULHQLQDQ", "PCLMULLQHQDQ", "PCLMULHQHQDQ":
		return "PCLMULQDQ"
	}

	/*
	 * SSE2
	 */
	switch op {
	// The condition is encoded as an imm8 operand of CMPPD
	case "CMPEQPD", "CMPLTPD", "CMPLEPD", "CMPUNORDPD", "CMPNEQPD", "CMPNLTPD",
		"CMPNLEPD", "CMPORDPD":
		return "CMPPD"
	// The condition is encoded as an imm8 operand of CMPSD.
	case "CMPEQSD", "CMPLTSD", "CMPLESD", "CMPUNORDSD", "CMPNEQSD", "CMPNLTSD",
		"CMPNLESD", "CMPORDSD":
		return "CMPSD"
	// D (double word) has been replaced by L (Long)
	case "PSUBD":
		return "PSUBL"
	case "MASKMOVDQU":
		return "MASKMOVOU"
	case "MOVD":
		return "MOVQ"
	case "MOVDQ2Q":
		return "MOVQ"
	case "MOVNTDQ":
		return "MOVNTO"
	case "MOVDQA":
		return "MOVO"
	case "MOVDQU":
		return "MOVOU"
	case "PSLLD":
		return "PSLLL"
	case "PSLLDQ":
		return "PSLLO"
	case "PSRAD":
		return "PSRAL"
	case "PSRLD":
		return "PSRLL"
	case "PSRLDQ":
		return "PSRLO"
	case "PADDD":
		return "PADDL"
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
