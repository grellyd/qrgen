### qrgen
## Written by @grellyd

qrgen generates a qr code image from a given string. 

Using CGo bindings, the inner qrlib library uses [Nayuki's](https://github.com/nayuki) excellent [C implementation](https://github.com/nayuki/QR-Code-generator/tree/master/c) of the QR Code Generation Algorithm as specified by [ISO/IEC 18004:2015](https://www.iso.org/standard/62021.html).

A user must provide the given link. If an [error correction level](https://www.qrcode.com/en/about/error_correction.html) is not given, level M (15% error) will be chosen.

See the help dialogue for usage information.

## Notes

Due to gonew's space limitation, qrgen can only handle space-less strings.
