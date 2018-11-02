// Copyright 2018 Vitali Fedulov. All rights reserved. Use of this source code
// is governed by a MIT-style license that can be found in the LICENSE file.

package images

import (
	"image"
	"math"
	"math/rand"
)

const (
	// Color similarity parameters.

	// Side dimension of a mask.
	maskSize = 24
	// Number of white pixels in the mask surrounding a white seed pixel.
	numAggregationPixels = 4
	// Side dimension (in pixels) of a downsample square to reasonably well
	// approximate color area of a full size image.
	downsampleSize = 12
	// Cutoff value for color distance.
	distanceThreshold = 50
	// Cosine similarity squared.
	cosineSimilarity2 = 0.95

	// Geometric similarity parameters.

	// Image width base scale for image proportions filter.
	baseWidth = 100
	// Threshold of height pixels for images rescaled to baseWidth.
	heightThreshold = 10
)

// Masks generates masks, each of which will be used to calculate an image hash.
// Conceptually a mask is a black square image with few white pixels used for
// average color calculation. In the function output a mask is a map with keys
// corresponding to white pixel coordinates only, because black pixels are
// redundant.
func Masks() []map[image.Point]bool {
	masks := make([]map[image.Point]bool, 0)
	for x := 1; x < maskSize-1; x++ {
		for y := 1; y < maskSize-1; y++ {
			alreadyAddedPixels := make(map[image.Point]bool)
			// Aggregation seed pixel.
			alreadyAddedPixels[image.Point{x, y}] = true
			// Pixels randomly aggregating around the seed pixel.
			for len(alreadyAddedPixels) < numAggregationPixels+1 {
				// Aggregation points are placed within 3x3 pixel area.
				dx := rand.Intn(3) - 1
				dy := rand.Intn(3) - 1
				if outOfBound(x+dx, y+dy) || (dx == 0 && dy == 0) {
					continue
				}
				alreadyAddedPixels[image.Point{x + dx, y + dy}] = true
			}
			masks = append(masks, alreadyAddedPixels)
		}
	}
	return masks
}

// Function to check if a pixel coordinates are out of the mask boundaries.
func outOfBound(x, y int) bool {
	if x < 0 || y < 0 || x >= maskSize || y >= maskSize {
		return true
	}
	return false
}

// Hash calculates a slice of average color values of an image at the position
// of white pixels of a mask. One average value corresponds to one mask.
// The masks for the input are generated with the Masks function. The Hash
// function also returns the original image width and height.
func Hash(img image.Image, masks []map[image.Point]bool) (h []float32,
	imgSize image.Point) {
	// Image is resampled to the mask size. Since masks are square the images
	// also are made square for image comparison.
	resImg, imgSize := ResampleByNearest(img,
		image.Point{maskSize * downsampleSize, maskSize * downsampleSize})
	h = make([]float32, len(masks))
	var (
		x, y            int
		r, g, b, sum, s uint32
	)
	// For each mask.
	for i := 0; i < len(masks); i++ {
		sum, s = 0, 0
		// For each white pixel of a mask.
		for w := range masks[i] {
			x, y = w.X, w.Y
			// For each pixel of resImg corresponding to the white mask pixel
			// above.
			for m := 0; m < downsampleSize; m++ {
				for n := 0; n < downsampleSize; n++ {
					// Alpha channel is not used for image comparison.
					r, g, b, _ =
						resImg.At(x*downsampleSize+m, y*downsampleSize+n).RGBA()
					// A cycle over the mask numbers to calculate average value
					// for different color channels. Red, green and gray are
					// considered more visually signicant than blue.
					switch i % 3 {
					case 0:
						sum += r
						s++
					case 1:
						sum += g
						s++
					case 2:
						sum += r + g + b
						s += 3
					}
				}
			}
		}
		h[i] = float32(sum) / float32(s*255)
	}
	return h, imgSize
}

// Similar function gives a verdict for image A and B based on their hashes and
// sizes. The input parameters are generated with the Hash function.
func Similar(hA, hB []float32, imgSizeA, imgSizeB image.Point) bool {

	// Filter 1. Threshold for mismatching image proportions. Based on rescaling
	// all images to same baseWidth and cutoff at heightThreshold.
	xA, yA := imgSizeA.X, imgSizeA.Y
	xB, yB := imgSizeB.X, imgSizeB.Y
	if xA*yB*baseWidth+heightThreshold*xA*xB < xB*yA*baseWidth ||
		xA*yB*baseWidth > xB*yA*baseWidth+heightThreshold*xA*xB {
		return false
	}

	// Filter 2. Color distance threshold to exit early on obvious outliers.
	for i := 0; i < len(hA); i++ {
		if math.Abs(float64(hA[i])-float64(hB[i])) > distanceThreshold {
			return false
		}
	}

	// Filter 3. Cosine similarity threshold.
	var dotProduct, sumSqA, sumSqB float32
	for i := 0; i < len(hA); i++ {
		dotProduct += hA[i] * hB[i]
		sumSqA += hA[i] * hA[i]
		sumSqB += hB[i] * hB[i]
	}
	if dotProduct*dotProduct < cosineSimilarity2*sumSqA*sumSqB {
		return false
	}

	return true
}
