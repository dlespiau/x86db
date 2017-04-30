package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/dlespiau/x86db"
)

func isAlreadyKnown(name string) bool {
	for _, opcode := range Anames {
		if name == opcode {
			return true
		}
	}
	return false
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
