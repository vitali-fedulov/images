// Copyright 2018 Vitali Fedulov. All rights reserved. Use of this source code
// is governed by a MIT-style license that can be found in the LICENSE file.

// Package images allows image comparison by perceptual similarity.
// Supported image types are those default to the Go image package
// https://golang.org/pkg/image/ (which are GIF, JPEG and PNG in October 2018).
package images

import (
	"image"
	"os"
)

// Open opens and decodes an image file for a given path.
func Open(path string) (img image.Image, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err = image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, err
}
