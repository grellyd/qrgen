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
	"image/color"
	"os"
	"image"
	"image/png"
	"fmt"
	"strconv"
	"regexp"
//	"unsafe"
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
	l := ""
	ec := ECLL
	o := "out.png"

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
			ec, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("error parsing error level: %s", err.Error())
			}
			if ec < 0 || ec > 3 {
				return fmt.Errorf("incorrect error level")
			}
		}
		if args[2] != "" {
			o = args[2]
			ending, err :=  regexp.MatchString("\\.png$", o) 
			if err != nil {
				return fmt.Errorf("unable to match for ending: %s", err.Error())
			}
			if !ending {
				return fmt.Errorf("output must end in '.png'")
			}
		} 
	} else {
		return fmt.Errorf("argument error")
	}
	i, err := Generate(l, ec)
	if err!= nil {
		return fmt.Errorf("unable to generate error: %s", err.Error())
	}
	err = write(i, o)
	if err != nil {
		return fmt.Errorf("unable to write image: %s", err.Error())
	}
	return nil
}

// write the resulting image to disk
func write(i *image.Gray, o string) error {
	f, _ := os.Create(o)
	defer f.Close()
	png.Encode(f, i)
	return nil
}

// Generate the QR Code
func Generate(l string, ec ErrorCorrectionLevel) (i *image.Gray, err error) {
	clink := C.CString(l)
	//defer C.free(unsafe.Pointer(&clink))
	b := make([]C.uchar, C.qrcodegen_BUFFER_LEN_MAX)
	//defer C.free(unsafe.Pointer(&b))
	qr := make([]C.uchar, C.qrcodegen_BUFFER_LEN_MAX)
	//defer C.free(unsafe.Pointer(&qr))
	/*
	From qrcodegen.h
	bool qrcodegen_encodeText(const char *text, uint8_t tempBuffer[], uint8_t qrcode[], enum qrcodegen_Ecc ecl, int minVersion, int maxVersion, enum qrcodegen_Mask mask, bool boostEcl);
	*/
	ok := C.qrcodegen_encodeText(clink, &b[0], &qr[0], C.enum_qrcodegen_Ecc(1), C.int(1), C.int(40), C.enum_qrcodegen_Mask(-1), true)
	if !ok {
		return nil, fmt.Errorf("Unable to encodeText")
	}
	i, err = buildImage(qr)
	if err != nil {
		return nil, fmt.Errorf("unable to build an output image: %s", err.Error())
	}
	return i, nil
}

func buildImage(qr []C.uchar) (i *image.Gray, err error) {
	C.printQr(&qr[0])
	size := int(C.qrcodegen_getSize(&qr[0])) + 2 * BORDER
	g := image.NewGray(image.Rect(0, 0, size, size))
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			mod := C.qrcodegen_getModule(&qr[0], C.int(x - BORDER), C.int(y - BORDER))
			if mod { 
				g.Set(x, y, color.Black)
			} else {
				g.Set(x, y, color.White)
			}
		}
	}
	return g, nil
}
