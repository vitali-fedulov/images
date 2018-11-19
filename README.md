# Comparing images in Go

<img align="right" style="margin-left:18px" src="logo.gif">

The library compares images by perceptual similarity to find near duplicates and
resized images. Supported image types: GIF, JPEG and PNG (golang.org/pkg/image/ as in October 2018).

`Similar` function gives a verdict whether 2 images are similar or not. The library also contains wrapper functions to open/save images and basic image resampling/resizing.

Demo: [similar image search](https://www.similar.pictures) (JavaScript implementation).

Documentation: [godoc](https://godoc.org/github.com/vitali-fedulov/images).

## Example of comparing 2 photos
```go
package main

import (
	"fmt"
	"github.com/vitali-fedulov/images"
)

func main() {
	
	// Open and decode photos.
	imgA, err := images.Open("photoA.jpg")
	if err != nil {
		panic(err)
	}
	imgB, err := images.Open("photoB.jpg")
	if err != nil {
		panic(err)
	}
	
	// Define masks.
	masks := images.Masks()
	
	// Calculate hashes.
	hA, imgSizeA := images.Hash(imgA, masks)
	hB, imgSizeB := images.Hash(imgB, masks)
	
	// Image comparison.
	if images.Similar(hA, hB, imgSizeA, imgSizeB) {
		fmt.Println("Images are similar.")
	} else {
		fmt.Println("Images are distinct.")
	}
}
```

## Algorithm for image comparison

[Detailed explanation with illustrations](https://www.similar.pictures/algorithm-for-perceptual-image-comparison.html).

Summary: In the algorithm images are resized to small squares of fixed size.
A number of masks representing several sample pixels are run against the resized
images to calculate average color values. Then the values are compared to
give the similarity verdict. Also image proportions are used to avoid matching
images of distinct shape.
