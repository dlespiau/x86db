package main

import (
	"fmt"
	"log"
	"os"
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

	return op
}

func isAlreadyKnown(name string) bool {
	name = nasmOpcodeToPlan9(name)
	for _, opcode := range Anames {
		if name == opcode {
			return true
		}
	}
	return false
}

func isAlreadyTested(name string) bool {
	name = nasmOpcodeToPlan9(name)
	_, ok := testedMap[name]
	return ok
}

func main() {
	db := x86db.NewDB(os.Args[1])

	err := db.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// SSE instructions, not already known by the Go's cmd/internal/obj package.
	insns := db.FindByExtension(x86db.ExtensionSSE).Where(func(insn x86db.Instruction) bool {
		return !isAlreadyKnown(insn.Name)
	})

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	for _, insn := range insns {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", insn.Name, insn.Operands, insn.Pattern, insn.Flags)
	}
	w.Flush()
}
