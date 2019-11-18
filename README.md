# Comparing images in Go

<img align="right" style="margin-left:18px" src="logo.gif">

Demo: [similar image search and clustering](https://similar.pictures) (deployed [from](https://github.com/vitali-fedulov/vitali-fedulov.github.io)).

The library compares images by perceptual similarity to find near duplicates and
resized images. There are no dependencies (only the Golang standard library is used). Supported image types: GIF, JPEG and PNG (golang.org/pkg/image/ as in October 2018).

`Similar` function gives a verdict whether 2 images are similar or not. The library also contains wrapper functions to open/save images and basic image resampling/resizing.

Documentation: [godoc](https://godoc.org/github.com/vitali-fedulov/images).

## UPDATE (November 2019)

The code has been simplified so that masks do not need to be calculated in a separate line (see the updated example below). As a result the Hash function only needs one argument, instead of two.

The comparison algorithm has been improved with ~15% additional matches and better quality.

## Example of comparing 2 photos
```go
package main

import (
	"fmt"
	"github.com/vitali-fedulov/images"
)

func main() {
	
	// Open photos.
	imgA, err := images.Open("photoA.jpg")
	if err != nil {
		panic(err)
	}
	imgB, err := images.Open("photoB.jpg")
	if err != nil {
		panic(err)
	}
	
	// Calculate hashes and image sizes.
	hashA, imgSizeA := images.Hash(imgA)
	hashB, imgSizeB := images.Hash(imgB)
	
	// Image comparison.
	if images.Similar(hashA, hashB, imgSizeA, imgSizeB) {
		fmt.Println("Images are similar.")
	} else {
		fmt.Println("Images are distinct.")
	}
}
```

## Algorithm for image comparison

[Detailed explanation with illustrations](https://vitali-fedulov.github.io/algorithm-for-perceptual-image-comparison.html).

Summary: In the algorithm images are resized to small squares of fixed size.
A number of masks representing several sample pixels are run against the resized
images to calculate average color values. Then the values are compared to
give the similarity verdict. Also image proportions are used to avoid matching
images of distinct shape.
