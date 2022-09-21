package main

import (
	"math"
)

func expandQuadTree(quad *ImageQuadTree, minSize int, tolerance uint8, processLeaf (func (quad *ImageQuadTree, clean bool))) {
	bounds := quad.Image.Bounds()
	mainC := quad.Image.RGBAAt(bounds.Min.X, bounds.Min.Y)

	if quad.Image.Rect.Dx()/2 < minSize || quad.Image.Rect.Dy()/2 < minSize {
		processLeaf(quad, false)
		return
	}

	for x := 0; x < bounds.Dx(); x++ {
		for y := 0; y < bounds.Dy(); y++ {
			col := quad.Image.RGBAAt(bounds.Min.X+x, bounds.Min.Y+y)
			if math.Abs(float64(col.R) - float64(mainC.R)) > float64(tolerance) {
				if quad.Split() {
					for _, q := range quad.Q {
						expandQuadTree(q, minSize, tolerance, processLeaf)
					}
				} else {
					processLeaf(quad, false)
				}
				return
			}
		}
	}

	processLeaf(quad, true)
}
