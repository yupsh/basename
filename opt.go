package command

type Suffix string

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

type flags struct {
	Suffix   Suffix
	Multiple MultipleFlag
	Zero     ZeroFlag
}

func (s Suffix) Configure(flags *flags)       { flags.Suffix = s }
func (m MultipleFlag) Configure(flags *flags) { flags.Multiple = m }
func (z ZeroFlag) Configure(flags *flags)     { flags.Zero = z }
