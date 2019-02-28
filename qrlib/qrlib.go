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
//	"unsafe"
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
			ec, err := strconv.Atoi(args[1])
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
	i, err := Generate(l, ec)
	if err!= nil {
		return fmt.Errorf("unable to generate error: %s", err.Error())
	}
	err = write(i)
	if err != nil {
		return fmt.Errorf("unable to write image: %s", err.Error())
	}
	return nil
}

// write the resulting image to disk
func write(i *image.Image) error{
	return nil
}

/*
bool qrcodegen_encodeText(const char *text, uint8_t tempBuffer[], uint8_t qrcode[], enum qrcodegen_Ecc ecl, int minVersion, int maxVersion, enum qrcodegen_Mask mask, bool boostEcl);
*/

// Generate the QR Code
func Generate(l string, ec int) (i *image.Image, err error) {
	clink := C.CString(l)
	//defer C.free(unsafe.Pointer(&clink))
	b := make([]C.uchar, C.qrcodegen_BUFFER_LEN_MAX)
	//defer C.free(unsafe.Pointer(&b))
	qr := make([]C.uchar, C.qrcodegen_BUFFER_LEN_MAX)
	//defer C.free(unsafe.Pointer(&qr))
	ok := C.qrcodegen_encodeText(clink, &b[0], &qr[0], C.enum_qrcodegen_Ecc(1), C.int(1), C.int(40), C.enum_qrcodegen_Mask(-1), true)
	if !ok {
		return nil, fmt.Errorf("Unable to encodeText")
	}
	size := C.qrcodegen_getSize(&qr[0])
	fmt.Printf("%v\n", size)
	fmt.Printf("%v\n", qr[:size*4])
	fmt.Printf("%v", qr)
	i, err = buildImage(qr)
	if err != nil {
		return nil, fmt.Errorf("unable to build an output image: %s", err.Error())
	}
	return i, nil
}

func buildImage(qr []C.uchar) (i *image.Image, err error) {
	C.printQr(&qr[0])
	size := int(C.qrcodegen_getSize(&qr[0]))
	g := image.NewGray(image.Rect(0, 0, size, size))
	for y := 0 ;y < size; y++ {
		for x := 0; x < size; x++ {
			mod := C.qrcodegen_getModule(&qr[0], C.int(x), C.int(y))
			if mod { 
				g.Set(x, y, color.Black)
			} else {
				g.Set(x, x, color.White)
			}
		}
	}
	f, _ := os.Create("out.png")
	defer f.Close()
	png.Encode(f, g)
	return nil, nil
}

/*
// Prints the given QR Code to the console.
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
