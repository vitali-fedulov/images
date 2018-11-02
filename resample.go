// Copyright 2018 Vitali Fedulov. All rights reserved. Use of this source code
// is governed by a MIT-style license that can be found in the LICENSE file.

package images

import (
	"image"
	"image/color"
)

// ResampleByNearest resizes an image by the nearest neighbour method to the
// output size outX, outY. It also returns the size inX, inY of the input image.
func ResampleByNearest(inImg image.Image, outImgSize image.Point) (
	outImg image.RGBA, inImgSize image.Point) {
	// Original image size.
	xMax, xMin := inImg.Bounds().Max.X, inImg.Bounds().Min.X
	yMax, yMin := inImg.Bounds().Max.Y, inImg.Bounds().Min.Y
	inImgSize.X = xMax - xMin
	inImgSize.Y = yMax - yMin

	// Destination rectangle.
	outRect := image.Rectangle{image.Point{0, 0}, outImgSize}
	// Color model of uint8 per color.
	outImg = *image.NewRGBA(outRect)
	var (
		r, g, b, a uint32
	)
	for x := 0; x < outImgSize.X; x++ {
		for y := 0; y < outImgSize.Y; y++ {
			r, g, b, a = inImg.At(
				x*inImgSize.X/outImgSize.X+xMin,
				y*inImgSize.Y/outImgSize.Y+yMin).RGBA()
			outImg.Set(x, y, color.RGBA{
				uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)})
		}
	}
	return outImg, inImgSize
}
