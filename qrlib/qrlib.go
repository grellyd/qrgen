package qrlib

/*
#include <stdlib.h>
#include <qrcodegen.h>
*/
import "C"

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

	return Generate(l, ec)
}

/*
bool qrcodegen_encodeText(const char *text, uint8_t tempBuffer[], uint8_t qrcode[], enum qrcodegen_Ecc ecl, int minVersion, int maxVersion, enum qrcodegen_Mask mask, bool boostEcl);
*/

// Generate the QR Code
func Generate(l string, ec int) error {
	clink := C.CString(l)
	defer C.free(clink)
	b := make([]uint, C.qrcodegen_BUFFER_LEN_MAX)
	defer C.free(b)
	qr := make([]uint, C.qrcodegen_BUFFER_LEN_MAX)
	defer C.free(qr)
	ecl := C.enum_qrcodegen_Ecc{qrcodegen_Ecc_MEDIUM}
	defer C.free(ecl)
	mask := C.enum_qrcodegen_Mask{qrcodegen_Mask_AUTO}
	defer C.free(mask)
	ok := C.qrcodegen_encodeText(clink, &b[0], &qr[0], ecl, C.int(1), C.int(40), mask, boost)
	if !ok {
		return fmt.Errorf("Unable to encodeText")
	}
	return nil
}



/*
uint8_t C.uint(x)
size_t C.sizeof_x || uint16_t
enum qrcodegen_Ecc  -> enum_qrcodegen_Ecc
*/
