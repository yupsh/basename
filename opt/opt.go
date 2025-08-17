package opt

// Custom types for parameters
type Suffix string

// Boolean flag types with constants
type MultipleFlag bool
const (
	Multiple   MultipleFlag = true
	NoMultiple MultipleFlag = false
)

type ZeroFlag bool
const (
	Zero   ZeroFlag = true
	NoZero ZeroFlag = false
)

// Flags represents the configuration options for the basename command
type Flags struct {
	Suffix   Suffix       // Suffix to remove
	Multiple MultipleFlag // Process multiple arguments
	Zero     ZeroFlag     // End output with NUL character instead of newline
}

// Configure methods for the opt system
func (s Suffix) Configure(flags *Flags)     { flags.Suffix = s }
func (m MultipleFlag) Configure(flags *Flags) { flags.Multiple = m }
func (z ZeroFlag) Configure(flags *Flags)     { flags.Zero = z }
