// Copyright 2018 Vitali Fedulov. All rights reserved. Use of this source code
// is governed by a MIT-style license that can be found in the LICENSE file.

package images

import (
	"image"
)

const (
	// Color similarity parameters.

	// Side dimension of a mask.
	maskSize = 24
	// Side dimension (in pixels) of a downsample square to reasonably well
	// approximate color area of a full size image.
	downsampleSize = 12

	// Cutoff value for color distance.
	colorDiff = 50
	// Cutoff coefficient for Euclidean distance (squared).
	euclCoeff = 0.2
	// Cutoff coefficient for color sign correlation.
	corrCoeff = 0.7

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
// redundant. In this particular implementation white pixels form 3x3 squares.
func masks() []map[image.Point]bool {
	ms := make([]map[image.Point]bool, 0)
	for x := 1; x < maskSize-1; x++ {
		for y := 1; y < maskSize-1; y++ {
			maskPixels := make(map[image.Point]bool)
			for dx := -1; dx < 2; dx++ {
				for dy := -1; dy < 2; dy++ {
					maskPixels[image.Point{x + dx, y + dy}] = true
				}
			}
			ms = append(ms, maskPixels)
		}
	}
	return ms
}

// Making masks.
var ms = masks()

// Number of masks.
var numMasks = len(ms)

// Hash calculates a slice of average color values of an image at the position
// of white pixels of a mask. One average value corresponds to one mask.
// The function also returns the original image width and height.
func Hash(img image.Image) (h []float32,
	imgSize image.Point) {
	// Image is resampled to the mask size. Since masks are square the images
	// also are made square for image comparison.
	resImg, imgSize := ResampleByNearest(img,
		image.Point{maskSize * downsampleSize, maskSize * downsampleSize})
	h = make([]float32, numMasks)
	var (
		x, y            int
		r, g, b, sum, s uint32
	)
	// For each mask.
	for i := 0; i < numMasks; i++ {
		sum, s = 0, 0
		// For each white pixel of a mask.
		for w := range ms[i] {
			x, y = w.X, w.Y
			// For each pixel of resImg corresponding to the white mask pixel
			// above.
			for m := 0; m < downsampleSize; m++ {
				for n := 0; n < downsampleSize; n++ {
					// Alpha channel is not used for image comparison.
					r, g, b, _ =
						resImg.At(x*downsampleSize+m, y*downsampleSize+n).RGBA()
					// A cycle over the mask numbers to calculate average value
					// for different color channels.
					switch i % 3 {
					case 0:
						sum += r
						s++
					case 1:
						sum += g
						s++
					case 2:
						sum += b
						s++
					}
				}
			}
		}
		h[i] = float32(sum) / float32(s*255)
	}
	return h, imgSize
}

// Euclidean distance threshold (squared).
var euclDist2 = float32(numMasks) * float32(colorDiff*colorDiff) * euclCoeff

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

	// Filter 2a. Euclidean distance.
	var sum float32
	for i := 0; i < numMasks; i++ {
		sum += (hA[i] - hB[i]) * (hA[i] - hB[i])
	}
	if sum > euclDist2 {
		return false
	}


	// Filter 3. Pixel brightness sign correlation test.
	sum = 0.0
	for i := 0; i < numMasks-1; i++ {
		if (hA[i] < hA[i+1]) && (hB[i] < hB[i+1]) ||
			(hA[i] == hA[i+1]) && (hB[i] == hB[i+1]) ||
			(hA[i] > hA[i+1]) && (hB[i] > hB[i+1]) {
			sum++
		}
	}
	if sum < float32(numMasks)*corrCoeff {
		return false
	}


	// Filter 2b. Euclidean distance with normalized histogram.
	sum = 0.0
	hA, hB = normalize(hA), normalize(hB)
	for i := 0; i < numMasks; i++ {
		sum += (hA[i] - hB[i]) * (hA[i] - hB[i])
	}
	if sum > euclDist2 {
		return false
	}

	return true
}

// normalize stretches histograms for the 3 channels of the hashes, so that
// minimum and maximum values of each are 0 and 255 correspondingly.
func normalize(h []float32) []float32 {
	normalized := make([]float32, numMasks)
	var rMin, gMin, bMin, rMax, gMax, bMax float32
	rMin, gMin, bMin = 256, 256, 256
	rMax, gMax, bMax = 0, 0, 0
	// Looking for extreme values.
	for n := 0; n < numMasks; n += 3 {
		if h[n] > rMax {
			rMax = h[n]
		}
		if h[n] < rMin {
			rMin = h[n]
		}
	}
	for n := 1; n < numMasks; n += 3 {
		if h[n] > gMax {
			gMax = h[n]
		}
		if h[n] < gMin {
			gMin = h[n]
		}
	}
	for n := 2; n < numMasks; n += 3 {
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
	for n := 0; n < numMasks; n += 3 {
		normalized[n] = (h[n] - rMin) * 255 / rMM
	}
	for n := 1; n < numMasks; n += 3 {
		normalized[n] = (h[n] - gMin) * 255 / gMM
	}
	for n := 2; n < numMasks; n += 3 {
		normalized[n] = (h[n] - bMin) * 255 / bMM
	}

	return normalized
}
