package cli

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/grellyd/qrgen/qrlib"
)

// Run the cli
func Run() (code int) {
	commands := ValidCommands()
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		for _, cmd := range commands {
			arg, value := getFullArg(args, i, cmd.length)
			if cmd.regex.MatchString(arg) {
				if cmd.short == "-h" {
					fmt.Println(usage(commands))
					return code
				}
				cmd.value = value
			}
		}
	}
	runArgs := convertCommands(commands)
	// Run the core program
	err := generateCli(runArgs)
	if err != nil {
		fmt.Println(err.Error())
		code = 1
	}
	return code
}

// GenerateCli is a CLI wrapper for Generate
func generateCli(args []string) error {
	var err error
	l := ""
	ec := qrlib.ECLL
	o := "out.png"
	s := 1

	if len(args) > 0 {
		if args[0] != "" {
			l = args[0]
		} else {
			return fmt.Errorf("no link given")
		}
		if args[1] != "" {
			i, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("error parsing error level: %s", err.Error())
			}

			ec = qrlib.ErrorCorrectionLevel(i)
			if ec < 0 || ec > 3 {
				return fmt.Errorf("incorrect error level")
			}
		}
		if args[2] != "" {
			o = args[2]
			ending, err := regexp.MatchString("\\.png$", o)
			if err != nil {
				return fmt.Errorf("unable to match for ending: %s", err.Error())
			}
			if !ending {
				return fmt.Errorf("output must end in '.png'")
			}
		}

		if args[4] != "" {
			s, err = strconv.Atoi(args[4])
			if err != nil {
				return fmt.Errorf("error parsing scaling factor:", err.Error())
			}
		}

	} else {
		return fmt.Errorf("argument error")
	}

	qr, err := qrlib.Generate(l, ec)
	if err != nil {
		return fmt.Errorf("unable to generate: %s", err.Error())
	}

	i, err := qrlib.BuildImage(qr, s)
	if err != nil {
		return fmt.Errorf("unable to build an output image: %s", err.Error())
	}

	err = qrlib.Write(i, o)
	if err != nil {
		return fmt.Errorf("unable to write image: %s", err.Error())
	}

	return nil
}

func getFullArg(args []string, start int, length int) (string, string) {
	if length == 1 || len(args) < (start+length) {
		return args[start], "true"
	}
	return strings.Join(args[start:start+length], " "), strings.Join(args[start+1:start+length], " ")
}

func usage(commands []*Command) (usageString string) {
	var builder strings.Builder
	builder.WriteString(breakLine)
	builder.WriteString("\n")
	builder.WriteString(programName)
	builder.WriteString("\n")
	builder.WriteString(breakLine)
	builder.WriteString("\n")
	builder.WriteString(programDescription)
	builder.WriteString("\n\n")
	fmt.Fprintf(&builder, "Usage: %s\n\n", usageCommand)
	fmt.Fprint(&builder, "Valid options:\n\n")
	for _, cmd := range commands {
		fmt.Fprintf(&builder, "[ %s | %s ]: %s\n", cmd.short, cmd.long, cmd.description)
	}
	return builder.String()
}

func convertCommands(commands []*Command) (args []string) {
	args = make([]string, len(commands))
	for _, cmd := range commands {
		if cmd.Present() {
			args[cmd.order] = cmd.value
		}
	}
	return args
}
