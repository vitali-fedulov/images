// Copyright 2018 Vitali Fedulov. All rights reserved. Use of this source code
// is governed by a MIT-style license that can be found in the LICENSE file.

package images

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"
)

// Gif saves image.RGBA to a file.
func Gif(img *image.RGBA, path string) {
	if destFile, err := os.Create(path); err != nil {
		log.Println("Cannot create file: ", path, err)
	} else {
		defer destFile.Close()
		gif.Encode(destFile, img, &gif.Options{
			NumColors: 256, Quantizer: nil, Drawer: nil})
	}
	return
}

// Png saves image.RGBA to a file.
func Png(img *image.RGBA, path string) {
	if destFile, err := os.Create(path); err != nil {
		log.Println("Cannot create file: ", path, err)
	} else {
		defer destFile.Close()
		png.Encode(destFile, img)
	}
	return
}

// Jpg saves image.RGBA to a file.
func Jpg(img *image.RGBA, path string, quality int) {
	if destFile, err := os.Create(path); err != nil {
		log.Println("Cannot create file: ", path, err)
	} else {
		defer destFile.Close()
		jpeg.Encode(destFile, img, &jpeg.Options{Quality: quality})
	}
	return
}
