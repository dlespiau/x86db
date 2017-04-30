// Generated by: gen
// TypeWriter: slice
// Directive: +gen on Instruction

package x86db

// InstructionSlice is a slice of type Instruction. Use it where you would use []Instruction.
type InstructionSlice []Instruction

// Where returns a new InstructionSlice whose elements return true for func. See: http://clipperhouse.github.io/gen/#Where
func (rcv InstructionSlice) Where(fn func(Instruction) bool) (result InstructionSlice) {
	for _, v := range rcv {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}