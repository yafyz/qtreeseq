package main

import (
	"image"
)

type ImageQuadTree struct {
	Leaf  bool
	Image *image.RGBA
	Q     []*ImageQuadTree
}

func (self *ImageQuadTree) Split() bool {
	if !self.Leaf { return true }

	bounds := self.Image.Bounds()
	Dx, Dy := bounds.Dx(), bounds.Dy()
	Bx, By := bounds.Min.X, bounds.Min.Y
	Mx, My := bounds.Max.X, bounds.Max.Y

	if Dx < 2 || Dy < 2 { return false }

	self.Leaf = false
	self.Q = make([]*ImageQuadTree, 4)

	self.Q[0] = newImageQuadTree((self.Image.SubImage(image.Rect(Bx, By, Bx+Dx/2, By+Dy/2)).(*image.RGBA)))
	self.Q[1] = newImageQuadTree((self.Image.SubImage(image.Rect(Bx+Dx/2, By, Mx, By+Dy/2)).(*image.RGBA)))
	self.Q[2] = newImageQuadTree((self.Image.SubImage(image.Rect(Bx, By+Dy/2, Bx+Dx/2, My)).(*image.RGBA)))
	self.Q[3] = newImageQuadTree((self.Image.SubImage(image.Rect(Bx+Dx/2, By+Dy/2, Mx, My)).(*image.RGBA)))

	return true
}

func newImageQuadTree(image *image.RGBA) *ImageQuadTree {
	return &ImageQuadTree{
		Leaf:  true,
		Image: image,
		Q:     nil,
	}
}

func (self *ImageQuadTree) maxDepth(cdepth int) int {
	max := 0

	if !self.Leaf {
		cdepth += 1
		for _, q := range self.Q {
			temp := q.maxDepth(0)
			if temp > max {
				max = temp
			}
		}
	}
	return cdepth + max
}
