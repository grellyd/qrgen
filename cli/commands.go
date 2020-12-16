package cli

import "regexp"

const (
	breakLine          = "===================="
	programName        = "qrgen"
	usageCommand       = "qrgen [ options ]"
	exampleCommand     = "qrgen -l 'https://grellyd.com'"
	programDescription = "Generate a QR code contianing a link."
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
	return []*Command{input(), ecLevel(), output(), debug(), scaleFactor(), help()}
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

func output() *Command {
	return &Command{
		order:       2,
		long:        "--output",
		short:       "-o",
		description: "The output image as a png",
		regex:       *regexp.MustCompile("-o|--output"),
		length:      2,
	}
}

func debug() *Command {
	return &Command{
		order:       3,
		long:        "--debug",
		short:       "-d",
		description: "Run the program in debug mode, with extra logging",
		regex:       *regexp.MustCompile("-d|--debug"),
		length:      1,
	}
}

func scaleFactor() *Command {
	return &Command{
		order:       4,
		long:        "--scaling-factor",
		short:       "-s",
		description: "The scaling factor to increase the size of the output image",
		regex:       *regexp.MustCompile("-s|--scaling-factor"),
		length:      2,
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
