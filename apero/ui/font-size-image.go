package ui

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type FontSizeImage struct {
	*gtk.Image
	previousFontSize int
	connectHandle    glib.SignalHandle
}

func (image *FontSizeImage) ConnectDraw() {
	image.connectHandle = image.Connect("draw", func() {
		ctx, _ := image.GetStyleContext()
		value, _ := ctx.GetProperty("font-size", gtk.STATE_FLAG_NORMAL)
		fontSize := int(value.(float64))
		if fontSize != image.previousFontSize {
			image.SetPixelSize(fontSize)
			image.previousFontSize = fontSize
		}
	})

}

func NewFontSizeImageFromIconName(iconName string) *gtk.Image {

	image, _ := gtk.ImageNewFromIconName(iconName, gtk.ICON_SIZE_BUTTON)

	fontSizeImage := &FontSizeImage{
		Image:            image,
		previousFontSize: 0,
	}

	fontSizeImage.ConnectDraw()

	return image
}

func SetFontSizeImageFromIconName(image *gtk.Image, iconName string) {
	image.SetFromIconName(iconName, gtk.ICON_SIZE_BUTTON)
}
