package qrlib

import (
	"fmt"
	"strconv"
)
/*
Error Codes:
Level L -> 7% Error Correction Capability
Level M -> 15%
Level Q -> 25%
Level H -> 30%
*/


// GenerateCli is a CLI wrapper for Generate
func GenerateCli(args []string) error {
	l := ""
	ec := 1

	for _, s := range args {
		fmt.Println(s)
	}
	if len(args) > 0 {
		if args[0] != "" {
			l = args[0]
		} else {
			return fmt.Errorf("no link given")
		}
		if args[1] != "" {
			ec, err := strconv.ParseInt(args[1], 10, 0)
			if err != nil {
				return fmt.Errorf("error parsing error level: %s", err.Error())
			}
			if ec < 0 || ec > 3 {
				return fmt.Errorf("incorrect error level")
			}
		}
	} else {
		return fmt.Errorf("argument error")
	}

	return Generate(l, ec)
}

// Generate the QR Code
func Generate(l string, ec int) error {
	return nil
}
