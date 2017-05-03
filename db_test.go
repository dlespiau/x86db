package x86db

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadInstructions(t *testing.T) {
	tests := []struct {
		input  string
		valid  bool
		golden Instruction
	}{
		{
			`SBB    reg_eax,sbytedword    [mi:  o32 83 /3 ib,s]    386,SM,ND`, true,
			Instruction{Name: "SBB", OpSize: OpSizeSM},
		},
		{
			`MOV    reg64,sdword          [mi:  o64 c7 /0 id,s]    X64,SM,OPT,ND`, true,
			Instruction{Name: "MOV", OpSize: OpSizeSM | OpSizeOPT},
		},
		{
			`ADDPS  xmmreg,xmmrm128       [rm:    np 0f 58 /r]     KATMAI,SSE`, true,
			Instruction{Name: "ADDPS", Extension: ExtensionSSE},
		},
	}

	for _, test := range tests {
		r := strings.NewReader(test.input)
		db := NewDB("")

		err := db.readInstructions(r)
		if !test.valid {
			assert.NotNil(t, err)
			continue
		}

		assert.Nil(t, err)
		assert.Equal(t, 1, len(db.Instructions))

		g := &test.golden
		parsed := &db.Instructions[0]

		assert.Equal(t, g.Name, parsed.Name)
		assert.Equal(t, g.OpSize, parsed.OpSize)
	}
}
