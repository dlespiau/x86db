package x86db

import (
	"bufio"
	"bytes"
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

// NewDB creates a new DB object.
func NewDB() *DB {
	return NewDBFromFile("")
}

// NewDBFromFile creates a new DB object, loading the list of instructions from
// instructionsFile. The format of instructionsFile is the nasm one.
func NewDBFromFile(instructionsFile string) *DB {
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

		// The 3rd field is the instruction pattern
		pattern, err := patternFromString(string(fields[3][1 : len(fields[3])-1]))
		if err != nil {
			return err
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

			e, err := ExtensionFromString(field)
			if err == nil {
				extension = e
				continue
			}
		}

		instruction := Instruction{
			Name:      string(fields[1]),
			Operands:  strings.Split(string(fields[2]), ","),
			Pattern:   *pattern,
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
	var r io.Reader

	if db.instructionsFile == "" {
		// Use the insns.dat bundled with the package
		data, err := Asset("data/insns.dat")
		if err != nil {
			return err
		}

		r = bytes.NewReader(data)
	} else {
		// db file provided by the user
		f, err := os.Open(db.instructionsFile)
		if err != nil {
			return err
		}
		defer f.Close()

		r = f
	}

	if err := db.readInstructions(r); err != nil {
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
