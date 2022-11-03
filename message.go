package broadcaster

import (
	"image"
)

type Message struct {
	Title string
	Body  string
	Image image.Image
}
