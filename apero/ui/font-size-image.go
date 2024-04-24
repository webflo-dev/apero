package ui

import (
	"github.com/gotk3/gotk3/gtk"
)

type FontSizeImage struct {
	*gtk.Image
	previousFontSize int
}

func NewFontSizeImageFromIconName(iconName string) *gtk.Image {

	image, _ := gtk.ImageNewFromIconName(iconName, gtk.ICON_SIZE_BUTTON)

	fontSizeImage := &FontSizeImage{
		Image:            image,
		previousFontSize: 0,
	}

	fontSizeImage.Connect("draw", func() {
		ctx, _ := fontSizeImage.GetStyleContext()
		value, _ := ctx.GetProperty("font-size", gtk.STATE_FLAG_NORMAL)
		fontSize := int(value.(float64))
		if fontSize != fontSizeImage.previousFontSize {
			fontSizeImage.SetPixelSize(fontSize)
			fontSizeImage.previousFontSize = fontSize
		}
	})

	return image
}
