// Copyright 2018 Vitali Fedulov. All rights reserved. Use of this source code
// is governed by a MIT-style license that can be found in the LICENSE file.

package images

import (
	"image"
	"path"
	"testing"
)

func TestMasks(t *testing.T) {
	masks := Masks()
	numMasks := (maskSize - 2) * (maskSize - 2)
	expectedNumMasks := len(masks)
	if len(masks) != numMasks {
		t.Errorf("Number of masks %d does not match expected value %d.",
			numMasks, expectedNumMasks)
	}
	for i := range masks {
		if len(masks[i]) > 3*3 {
			t.Errorf("Number of mask white pixels %d is more than 3*3.",
				len(masks[i]))
		}
	}
}

func TestHash(t *testing.T) {
	testDir := "testdata"
	testFile := "small.jpg"
	masks := Masks()
	img, err := Open(path.Join(testDir, testFile))
	if err != nil {
		t.Error("Error opening image:", err)
		return
	}
	pSums, imgSize := Hash(img, masks)
	if len(pSums) == 0 {
		t.Errorf("Number of pSums %d must not be 0.", len(pSums))
	}
	if len(pSums) != len(masks) {
		t.Errorf("Number of pSums %d is not equal number of masks %d.",
			len(pSums), len(masks))
	}
	if imgSize != (image.Point{267, 200}) {
		t.Errorf(
			"Calculated imgSize %d is not equal to the size from image properties %d.",
			imgSize, image.Point{267, 200})
	}
	var allZeroOrLessCounter int
	for i := range pSums {
		if pSums[i] > 255 {
			t.Errorf("pSums[i] %f is larger than 255.", pSums[i])
			break
		}
		if pSums[i] <= 0 {
			allZeroOrLessCounter++
		}
	}
	if allZeroOrLessCounter == len(pSums) {
		t.Error("All pSums[i] are 0 or less.")
	}
}

func TestSimilar(t *testing.T) {
	testDir := "testdata"
	imgFiles := []string{
		"flipped.jpg", "large.jpg", "small.jpg", "distorted.jpg"}
	masks := Masks()
	pSumsAll := make([][]float32, len(imgFiles))
	imgSizeAll := make([]image.Point, len(imgFiles))
	for i := range imgFiles {
		img, err := Open(path.Join(testDir, imgFiles[i]))
		if err != nil {
			t.Error("Error opening image:", err)
			return
		}
		pSumsAll[i], imgSizeAll[i] = Hash(img, masks)
	}
	if !Similar(pSumsAll[1], pSumsAll[2], imgSizeAll[1], imgSizeAll[2]) {
		t.Errorf("Expected similarity between %s and %s.",
			imgFiles[1], imgFiles[2])
	}
	if Similar(pSumsAll[1], pSumsAll[0], imgSizeAll[1], imgSizeAll[0]) {
		t.Errorf("Expected non-similarity between %s and %s.",
			imgFiles[1], imgFiles[0])
	}
	if Similar(pSumsAll[1], pSumsAll[3], imgSizeAll[1], imgSizeAll[3]) {
		t.Errorf("Expected non-similarity between %s and %s.",
			imgFiles[1], imgFiles[3])
	}
}
