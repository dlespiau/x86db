package x86db

import "fmt"

type OpSize uint32

const (
	OpSizeSM OpSize = 1 << iota
	OpSizeSM2
	OpSizeSB
	OpSizeSW
	OpSizeSD
	OpSizeSQ
	OpSizeSO
	OpSizeSY
	OpSizeSZ
	OpSizeSIZE
	OpSizeSX
	OpSizeAR0
	OpSizeAR1
	OpSizeAR2
	OpSizeAR3
	OpSizeAR4
	OpSizeOPT
)

type opSizeInfo struct {
	flag OpSize
	name string
	help string
}

var opSizeTab = []opSizeInfo{
	{OpSizeSM, "SM", "Size match"},
	{OpSizeSM2, "SM2", "Size match first two operands"},
	{OpSizeSB, "SB", "Unsized operands can't be non-byte"},
	{OpSizeSW, "SW", "Unsized operands can't be non-word"},
	{OpSizeSD, "SD", "Unsized operands can't be non-dword"},
	{OpSizeSQ, "SQ", "Unsized operands can't be non-qword"},
	{OpSizeSO, "SO", "Unsized operands can't be non-oword"},
	{OpSizeSY, "SY", "Unsized operands can't be non-yword"},
	{OpSizeSZ, "SZ", "Unsized operands can't be non-zword"},
	{OpSizeSIZE, "SIZE", "Unsized operands must match the bitsize"},
	{OpSizeSX, "SX", "Unsized operands not allowed"},
	{OpSizeAR0, "AR0", "SB, SW, SD applies to argument 0"},
	{OpSizeAR1, "AR1", "SB, SW, SD applies to argument 1"},
	{OpSizeAR2, "AR2", "SB, SW, SD applies to argument 2"},
	{OpSizeAR3, "AR3", "SB, SW, SD applies to argument 3"},
	{OpSizeAR4, "AR4", "SB, SW, SD applies to argument 4"},
	{OpSizeOPT, "OPT", "Optimizing assembly only"},
}

func opSizeFromString(name string) (OpSize, error) {
	for _, info := range opSizeTab {
		if info.name == name {
			return info.flag, nil
		}
	}

	return OpSize(0), fmt.Errorf("no OpSize with name '%s'", name)
}

type Extension int

const (
	ExtensionBase Extension = iota
	ExtensionFPU
	ExtensionMMX
	Extension3DNOW
	ExtensionSSE
	ExtensionSSE2
	ExtensionSSE3
	ExtensionVMX
	ExtensionSSSE3
	ExtensionSSE4A
	ExtensionSSE41
	ExtensionSSE42
	ExtensionSSE5
	ExtensionAVX
	ExtensionAVX2
	ExtensionFMA
	ExtensionBMI1
	ExtensionBMI2
	ExtensionTBM
	ExtensionRTM
	ExtensionINVPCID
	ExtensionAVX512
	ExtensionAVX512CD
	ExtensionAVX512ER
	ExtensionAVX512PF
	ExtensionMPX
	ExtensionSHA
	ExtensionPREFETCHWT1
	ExtensionAVX512VL
	ExtensionAVX512DQ
	ExtensionAVX512BW
	ExtensionAVX512IFMA
	ExtensionAVX512VBMI
)

// ExtensionInfo stores metadata about an extension.
type ExtensionInfo struct {
	Extension Extension
	Name      string
	Help      string
}

// ExtensionList is the list of known extensions.
var ExtensionList = []ExtensionInfo{
	{ExtensionFPU, "FPU", "FPU"},
	{ExtensionMMX, "MMX", "MMX"},
	{Extension3DNOW, "3DNOW", "3DNow!"},
	{ExtensionSSE, "SSE", "SSE"},
	{ExtensionSSE2, "SSE2", "SSE2"},
	{ExtensionSSE3, "SSE3", "SSE3 (PNI)"},
	{ExtensionVMX, "VMX", "VMX"},
	{ExtensionSSSE3, "SSSE3", "SSSE3"},
	{ExtensionSSE4A, "SSE4A", "AMD SSE4a"},
	{ExtensionSSE41, "SSE41", "SSE4.1"},
	{ExtensionSSE42, "SSE42", "SSE4.2"},
	{ExtensionSSE5, "SSE5", "SSE5"},
	{ExtensionAVX, "AVX", "AVX (128b)"},
	{ExtensionAVX2, "AVX2", "AVX2 (256b)"},
	{ExtensionFMA, "FMA", ""},
	{ExtensionBMI1, "BMI1", ""},
	{ExtensionBMI2, "BMI2", ""},
	{ExtensionTBM, "TBM", ""},
	{ExtensionRTM, "RTM", ""},
	{ExtensionINVPCID, "INVPCID", ""},
	{ExtensionAVX512, "AVX512", "AVX-512F (512b)"},
	{ExtensionAVX512CD, "AVX512CD", "AVX-512 Conflict Detection"},
	{ExtensionAVX512ER, "AVX512ER", "AVX-512 Exponential and Reciprocal"},
	{ExtensionAVX512PF, "AVX512PF", "AVX-512 Prefetch"},
	{ExtensionMPX, "MPX", "MPX"},
	{ExtensionSHA, "SHA", "SHA"},
	{ExtensionPREFETCHWT1, "PREFETCHWT1", "PREFETCHWT1"},
	{ExtensionAVX512VL, "AVX512VL", "AVX-512 Vector Length Orthogonality"},
	{ExtensionAVX512DQ, "AVX512DQ", "AVX-512 Dword and Qword"},
	{ExtensionAVX512BW, "AVX512BW", "AVX-512 Byte and Word"},
	{ExtensionAVX512IFMA, "AVX512IFMA", "AVX-512 IFMA instructions"},
	{ExtensionAVX512VBMI, "AVX512VBMI", "AVX-512 VBMI instructions"},
}

func ExtensionFromString(name string) (Extension, error) {
	for _, info := range ExtensionList {
		if info.Name == name {
			return info.Extension, nil
		}
	}

	return ExtensionBase, fmt.Errorf("no Extension with name '%s'", name)
}

// +gen slice:"Where"
type Instruction struct {
	Name      string
	Operands  []string
	Pattern   string
	Flags     string
	Extension Extension
	OpSize    OpSize
}

// String implements the stringer interface for Instruction
func (i *Instruction) String() string {
	return i.Name
}
