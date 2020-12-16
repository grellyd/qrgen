package qrlib

/*
#include <stdio.h>
#include <stdlib.h>
#include <qrcodegen.h>

static void printQr(const uint8_t qrcode[]) {
	int size = qrcodegen_getSize(qrcode);
	int border = 4;
	for (int y = -border; y < size + border; y++) {
		for (int x = -border; x < size + border; x++) {
			fputs((qrcodegen_getModule(qrcode, x, y) ? "##" : "  "), stdout);
		}
		fputs("\n", stdout);
	}
	fputs("\n", stdout);
}
*/
import "C"

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"regexp"
	"strconv"
	"unsafe"
)

/*
ErrorCorrectionLevel Codes:
Level Low     -> 7% Error Correction Capability
Level Medium  -> 15%
Level Quality -> 25%
Level High    -> 30%
*/
type ErrorCorrectionLevel int

const (
	//ECLL Low
	ECLL ErrorCorrectionLevel = iota
	//ECLM Medium
	ECLM
	//ECLQ Quality
	ECLQ
	//ECLH High
	ECLH
)

const (
	// BORDER size around the QR code in pixels
	BORDER = 2
)

// GenerateCli is a CLI wrapper for Generate
func GenerateCli(args []string) error {
	var err error
	l := ""
	ec := ECLL
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

			ec = ErrorCorrectionLevel(i)
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

	qr, err := Generate(l, ec)
	if err != nil {
		return fmt.Errorf("unable to generate: %s", err.Error())
	}

	i, err := buildImage(qr, s)
	if err != nil {
		return fmt.Errorf("unable to build an output image: %s", err.Error())
	}

	err = write(i, o)
	if err != nil {
		return fmt.Errorf("unable to write image: %s", err.Error())
	}

	return nil
}

// write the resulting image to disk
func write(i image.Image, o string) error {
	f, _ := os.Create(o)
	defer f.Close()
	png.Encode(f, i)
	fmt.Printf("Written to %s\n", o)
	return nil
}

// Generate the QR Code
func Generate(l string, ec ErrorCorrectionLevel) ([]C.uchar, error) {
	clink := C.CString(l)
	defer C.free(unsafe.Pointer(clink))
	// b and qr do not need to be C.freed as they use make, an allocation function the go garbage collector is aware of, rather than a C.malloc based function.
	b := make([]C.uchar, C.qrcodegen_BUFFER_LEN_MAX)
	qr := make([]C.uchar, C.qrcodegen_BUFFER_LEN_MAX)
	/*
		From qrcodegen.h
		bool qrcodegen_encodeText(const char *text, uint8_t tempBuffer[], uint8_t qrcode[], enum qrcodegen_Ecc ecl, int minVersion, int maxVersion, enum qrcodegen_Mask mask, bool boostEcl);
	*/
	ok := C.qrcodegen_encodeText(clink, &b[0], &qr[0], C.enum_qrcodegen_Ecc(1), C.int(1), C.int(40), C.enum_qrcodegen_Mask(-1), true)
	if !ok {
		return nil, fmt.Errorf("Unable to encodeText")
	}
	return qr, nil
}

func buildImage(qr []C.uchar, factor int) (i *image.Gray, err error) {
	C.printQr(&qr[0])
	size := int(C.qrcodegen_getSize(&qr[0]))*factor + 2*BORDER*factor
	g := image.NewGray(image.Rect(0, 0, size, size))
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			isBlack := C.qrcodegen_getModule(&qr[0], C.int((x/factor)-BORDER), C.int((y/factor)-BORDER))
			if isBlack {
				g.Set(x, y, color.Black)
			} else {
				g.Set(x, y, color.White)
			}
		}
	}
	return g, nil
}
