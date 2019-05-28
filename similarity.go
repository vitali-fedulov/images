// Copyright 2018 Vitali Fedulov. All rights reserved. Use of this source code
// is governed by a MIT-style license that can be found in the LICENSE file.

package images

import (
	"image"
	"math"
)

const (
	// Color similarity parameters.

	// Side dimension of a mask.
	maskSize = 24
	// Side dimension (in pixels) of a downsample square to reasonably well
	// approximate color area of a full size image.
	downsampleSize = 12
	// Cutoff values for color distance.
	distanceThreshold  = 50
	totalDistanceCoeff = 2
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
// In this particular implementation white pixels form square 3x3 regions (thus
// de-facto a median filter is applied to the resampled image).
func Masks() []map[image.Point]bool {
	masks := make([]map[image.Point]bool, 0)
	for x := 1; x < maskSize-1; x++ {
		for y := 1; y < maskSize-1; y++ {
			maskPixels := make(map[image.Point]bool)
			for dx := -1; dx < 2; dx++ {
				for dy := -1; dy < 2; dy++ {
					maskPixels[image.Point{x + dx, y + dy}] = true
				}
			}
			masks = append(masks, maskPixels)
		}
	}
	return masks
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

	// Filter 2a. Color distance threshold.
	var diff, sumDiffs float64
	for i := 0; i < len(hA); i++ {
		diff = math.Abs(float64(hA[i]) - float64(hB[i]))
		if diff > distanceThreshold {
			return false
		}
		sumDiffs += diff
	}

	// Filter 2b. Cumulative color distance threshold with a coefficient.
	if sumDiffs*totalDistanceCoeff > float64(len(hA))*distanceThreshold {
		return false
	}

	// Filter 3a. Cosine similarity threshold.
	var dotProduct, sumSqA, sumSqB float32
	for i := 0; i < len(hA); i++ {
		dotProduct += hA[i] * hB[i]
		sumSqA += hA[i] * hA[i]
		sumSqB += hB[i] * hB[i]
	}
	if dotProduct*dotProduct < cosineSimilarity2*sumSqA*sumSqB {
		return false
	}

	// Filter 3b. Cosine similarity threshold with normalized histogram.
	hA, hB = normalize(hA), normalize(hB)
	dotProduct, sumSqA, sumSqB = 0, 0, 0
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

// normalize stretches histograms for the 3 channels of the hashes, so that
// minimum and maximum values of each are 0 and 255 correspondingly.
func normalize(h []float32) []float32 {
	normalized := make([]float32, len(h))
	var rMin, gMin, bMin, rMax, gMax, bMax float32
	rMin, gMin, bMin = 256, 256, 256
	rMax, gMax, bMax = 0, 0, 0
	// Looking for extreme values.
	for n := 0; n < len(h); n += 3 {
		if h[n] > rMax {
			rMax = h[n]
		}
		if h[n] < rMin {
			rMin = h[n]
		}
	}
	for n := 1; n < len(h); n += 3 {
		if h[n] > gMax {
			gMax = h[n]
		}
		if h[n] < gMin {
			gMin = h[n]
		}
	}
	for n := 2; n < len(h); n += 3 {
		if h[n] > bMax {
			bMax = h[n]
		}
		if h[n] < bMin {
			bMin = h[n]
		}
	}
	// Normalization.
	rMM := rMax - rMin
	gMM := gMax - gMin
	bMM := bMax - bMin
	for n := 0; n < len(h); n += 3 {
		normalized[n] = (h[n] - rMin) * 255 / rMM
	}
	for n := 1; n < len(h); n += 3 {
		normalized[n] = (h[n] - gMin) * 255 / gMM
	}
	for n := 2; n < len(h); n += 3 {
		normalized[n] = (h[n] - bMin) * 255 / bMM
	}

	return normalized
}
