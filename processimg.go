package main

import (
	"image"
	"image/color"

	"golang.org/x/image/draw"
)

func quadifyImage(src *image.RGBA, quadsrc *image.RGBA, minSize int, tolerance uint8, invert bool) *image.RGBA {
	new := image.NewRGBA(src.Rect)
	iqtroot := newImageQuadTree(src)

	expandQuadTree(iqtroot, minSize, tolerance, func(quad *ImageQuadTree, clean bool) {
		alpha := uint8(0)
		if clean {
			alpha = quad.Image.RGBAAt(quad.Image.Rect.Min.X, quad.Image.Rect.Min.Y).R
		} else {
			count := 0
			sum := 0
			for x := 0; x < quad.Image.Rect.Dx(); x++ {
				for y := 0; y < quad.Image.Rect.Dy(); y++ {
					count++
					sum += int(quad.Image.RGBAAt(quad.Image.Rect.Min.X+x, quad.Image.Rect.Min.Y+y).R)
				}
			}
			alpha = uint8(sum/count)
		}

		if !invert {
			alpha = 255-alpha
		}

		if alpha > 0 {
			mask := image.NewRGBA(quad.Image.Rect)
			draw.Draw(mask, mask.Rect, &image.Uniform{color.Alpha{alpha}}, image.ZP, draw.Src)
			draw.NearestNeighbor.Scale(new, quad.Image.Rect, quadsrc, quadsrc.Rect, draw.Over, &draw.Options{DstMask: mask})
		}
	})

	return new
}