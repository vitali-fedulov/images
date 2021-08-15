# Comparing images in Go

Demo: [similar image search and clustering](https://similar.pictures) (deployed [from](https://github.com/vitali-fedulov/vitali-fedulov.github.io)).

Near duplicates and resized images can be found with the package. There are no dependencies: only the Golang standard library is used. Supported image types: GIF, JPEG and PNG (golang.org/pkg/image/ as in October 2018).

`Similar` function gives a verdict whether 2 images are similar or not. The library also contains wrapper functions to open/save images and basic image resampling/resizing.

Documentation: [godoc](https://godoc.org/github.com/vitali-fedulov/images).

## Example of comparing 2 photos
```go
package main

import (
	"fmt"
	"github.com/vitali-fedulov/images/v2"
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
