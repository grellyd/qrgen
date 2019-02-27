package cli

import "regexp"

const (
	breakLine          = "===================="
	programName        = "qrgen"
	usageCommand       = "go run main.go [ options ]"
	exampleCommand     = "go run main.go -i 'https://grellyd.com'"
	programDescription = "Generate QR codes contianing a string (aka a link) with a given level of error correction."
)

// A Command is an application instruction.
type Command struct {
	order       int
	long        string
	short       string
	description string
	regex       regexp.Regexp
	length      int
	value       string
}

// Present checks if a Command is present
func (c *Command) Present() bool {
	return c.value != ""
}

// ValidCommands to this program
func ValidCommands() (commands []*Command) {
	return []*Command{input(), ecLevel(), debug(), help()}
}

func ecLevel() *Command {
	return &Command{
		order:       1,
		long:        "--error-correction",
		short:       "-e",
		description: "The level of error correction to apply",
		regex:       *regexp.MustCompile("-e|--error-correction"),
		length:      2,
	}
}

func input() *Command {
	return &Command{
		order:       0,
		long:        "--link",
		short:       "-l",
		description: "The link to encode",
		regex:       *regexp.MustCompile("-l|--link"),
		length:      2,
	}
}

func debug() *Command {
	return &Command{
		order:       2,
		long:        "--debug",
		short:       "-d",
		description: "Run the program in debug mode, with extra logging",
		regex:       *regexp.MustCompile("-d|--debug"),
		length:      1,
	}
}

func help() *Command {
	return &Command{
		// Help is handled by the cli wrapper
		order:       -1,
		long:        "--help",
		short:       "-h",
		description: "Display this help dialogue and exit",
		regex:       *regexp.MustCompile("-h|--help"),
		length:      1,
	}
}
