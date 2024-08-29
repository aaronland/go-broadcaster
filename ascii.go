package broadcaster

// Adapted from: https://github.com/stdupp/goasciiart/tree/master
// There are other, more sophisticated ascii art generators, but
// since this is really just for debugging it didn't seem worth all
// the extra dependencies. See also:
// https://github.com/TheZoraiz/ascii-image-converter

import (
	"bytes"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"reflect"

	"github.com/nfnt/resize"
)

var asciiSTR = "MND8OZ$7I?+=~:,.."

func scaleImage(img image.Image, w int) (image.Image, int, int) {
	sz := img.Bounds()
	h := (sz.Max.Y * w * 10) / (sz.Max.X * 16)
	img = resize.Resize(uint(w), uint(h), img, resize.Lanczos3)
	return img, w, h
}

func convert2Ascii(img image.Image, w, h int) []byte {
	table := []byte(asciiSTR)
	buf := new(bytes.Buffer)

	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			g := color.GrayModel.Convert(img.At(j, i))
			y := reflect.ValueOf(g).FieldByName("Y").Uint()
			pos := int(y * 16 / 255)
			_ = buf.WriteByte(table[pos])
		}
		_ = buf.WriteByte('\n')
	}
	return buf.Bytes()
}
