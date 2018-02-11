// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

func init() {
	addTestCases(frameTests, frame)
}

var frameTests = []testCase{
	{
		Name: "frame.0",
		In: `package main

import (
	"image"

	"github.com/as/frame"
	"github.com/as/frame/font"
)

func main() {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	frame.New(img.Bounds(), font.NewGoMono(11), img, frame.Mono, frame.FrElastic)
}
`,
		Out: `package main

import (
	"image"

	"github.com/as/font"
	"github.com/as/frame"
)

func main() {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	frame.New(img, img.Bounds(), &frame.Config{Face: font.NewGoMono(11), Color: frame.Mono, Flag: frame.FrElastic})
}
`,
	},
}
