// Copyright 2018 Vitali Fedulov. All rights reserved. Use of this source code
// is governed by a MIT-style license that can be found in the LICENSE file.

package images

import (
	"image"
	"path"
	"reflect"
	"testing"
)

func TestResampleByNearest(t *testing.T) {
	testDir := "testdata"
	tables := []struct {
		inFile     string
		inImgSize  image.Point
		outFile    string
		outImgSize image.Point
	}{
		{"original.png", image.Point{533, 400},
			"nearest100x100.png", image.Point{100, 100}},
		{"nearest100x100.png", image.Point{100, 100},
			"nearest533x400.png", image.Point{533, 400}},
	}

	for _, table := range tables {
		inImg, err := Open(path.Join(testDir, table.inFile))
		if err != nil {
			t.Error("Cannot decode", path.Join(testDir, table.inFile))
		}
		outImg, err := Open(path.Join(testDir, table.outFile))
		if err != nil {
			t.Error("Cannot decode", path.Join(testDir, table.outFile))
		}
		resampled, inImgSize := ResampleByNearest(inImg, table.outImgSize)
		if !reflect.DeepEqual(
			outImg.(*image.RGBA), &resampled) || table.inImgSize != inImgSize {
			t.Errorf(
				"Resample data do not match for %s and %s.",
				table.inFile, table.outFile)
		}
	}
}
