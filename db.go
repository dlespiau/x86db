package x86db

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

// DB holds the list of known instructions
type DB struct {
	instructionsFile string
	Instructions     InstructionSlice
}

func NewDB(instructionsFile string) *DB {
	return &DB{
		instructionsFile: instructionsFile,
	}
}

var ignoreFlags = []string{
	"ignore", "OBSOLETE",
}

func ignoreInstruction(flag string) bool {
	for _, f := range ignoreFlags {
		if f == flag {
			return true
		}
	}
	return false
}

func (db *DB) readInstructions(r io.Reader) error {
	pattern := regexp.MustCompile(`^\s*(\S+)\s+(\S+)\s+(\S+|\[.*\])\s+(\S+)\s*$`)

	scanner := bufio.NewScanner(r)
next:
	for scanner.Scan() {
		line := scanner.Text()

		// strip comments
		idx := strings.IndexRune(line, ';')
		if idx >= 0 {
			line = line[0:idx]
		}

		if line == "" {
			continue
		}

		fields := pattern.FindSubmatch([]byte(line))
		// We want 4 fields
		if len(fields) != 5 {
			return fmt.Errorf("readInstructions: expected 4 fields got %d", len(fields)-1)
		}

		// The 4th field holds misc, comma separated, flags.
		var opSizeFlags OpSize
		var extension Extension
		for _, field := range strings.Split(string(fields[4]), ",") {
			if ignoreInstruction(field) {
				goto next
			}

			f, err := opSizeFromString(field)
			if err == nil {
				opSizeFlags |= f
				continue
			}

			e, err := extensionFromString(field)
			if err == nil {
				extension = e
				continue
			}
		}

		instruction := Instruction{
			Name:      string(fields[1]),
			Operands:  string(fields[2]),
			Pattern:   string(fields[3]),
			Flags:     string(fields[4]),
			Extension: extension,
			OpSize:    opSizeFlags,
		}
		db.Instructions = append(db.Instructions, instruction)
	}

	return nil
}

// Open loads instructions from disk.
func (db *DB) Open() error {
	f, err := os.Open(db.instructionsFile)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := db.readInstructions(f); err != nil {
		return err
	}

	return nil
}

// Close closes the DB precious resources.
func (db *DB) Close() {

}

// FindByExtension returns the list of instructions introduced as part of the
// specificied extension.
func (db *DB) FindByExtension(extension Extension) InstructionSlice {
	return db.Instructions.Where(func(insn Instruction) bool {
		return insn.Extension == extension
	})
}
