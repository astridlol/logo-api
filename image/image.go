package image

import (
	"bytes"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"image"
	_ "image/png"
)

func Generate(imageData []byte, color string, size int) error {
	emoji, _, err := image.Decode(bytes.NewReader(imageData))

	fmt.Println("Creating image with color", color)

	if err != nil {
		return err
	}

	emojiSize := float64(size) * 0.64
	resizedEmoji := resize.Resize(uint(emojiSize), uint(emojiSize), emoji, resize.Bicubic)

	context := gg.NewContext(size, size)
	context.DrawRectangle(0, 0, float64(size), float64(size))
	context.SetHexColor(color)
	context.Fill()
	context.DrawImageAnchored(resizedEmoji, size/2, size/2, 0.5, 0.5)

	return context.SavePNG("output.png")
}
