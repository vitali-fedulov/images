// Copyright 2018 Vitali Fedulov. All rights reserved. Use of this source code
// is governed by a MIT-style license that can be found in the LICENSE file.

package images

import (
	"image"
	"path"
	"testing"
)

func TestMasks(t *testing.T) {
	ms := masks()
	numMasks := (maskSize - 2) * (maskSize - 2)
	expectedNumMasks := len(ms)
	if len(ms) != numMasks {
		t.Errorf("Number of masks %d does not match expected value %d.",
			numMasks, expectedNumMasks)
	}
	for i := range ms {
		if len(ms[i]) > 3*3 {
			t.Errorf("Number of mask white pixels %d is more than 3*3.",
				len(ms[i]))
		}
	}
}

func TestHash(t *testing.T) {
	testDir := "testdata"
	testFile := "small.jpg"
	img, err := Open(path.Join(testDir, testFile))
	if err != nil {
		t.Error("Error opening image:", err)
		return
	}
	h, imgSize := Hash(img)
	if len(h) == 0 {
		t.Errorf("Number of h %d must not be 0.", len(h))
	}
	if imgSize != (image.Point{267, 200}) {
		t.Errorf(
			"Calculated imgSize %d is not equal to the size from image properties %d.",
			imgSize, image.Point{267, 200})
	}
	var allZeroOrLessCounter int
	for i := range h {
		if h[i] > 255 {
			t.Errorf("h[i] %f is larger than 255.", h[i])
			break
		}
		if h[i] <= 0 {
			allZeroOrLessCounter++
		}
	}
	if allZeroOrLessCounter == len(h) {
		t.Error("All h[i] are 0 or less.")
	}
}

func TestSimilar(t *testing.T) {
	testDir := "testdata"
	imgFiles := []string{
		"flipped.jpg", "large.jpg", "small.jpg", "distorted.jpg"}
	hashes := make([][]float32, len(imgFiles))
	imgSizeAll := make([]image.Point, len(imgFiles))
	for i := range imgFiles {
		img, err := Open(path.Join(testDir, imgFiles[i]))
		if err != nil {
			t.Error("Error opening image:", err)
			return
		}
		hashes[i], imgSizeAll[i] = Hash(img)
	}
	if !Similar(hashes[1], hashes[2], imgSizeAll[1], imgSizeAll[2]) {
		t.Errorf("Expected similarity between %s and %s.",
			imgFiles[1], imgFiles[2])
	}
	if Similar(hashes[1], hashes[0], imgSizeAll[1], imgSizeAll[0]) {
		t.Errorf("Expected non-similarity between %s and %s.",
			imgFiles[1], imgFiles[0])
	}
	if Similar(hashes[1], hashes[3], imgSizeAll[1], imgSizeAll[3]) {
		t.Errorf("Expected non-similarity between %s and %s.",
			imgFiles[1], imgFiles[3])
	}
}
