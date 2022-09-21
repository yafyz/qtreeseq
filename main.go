package main

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"runtime"
	"strconv"

	"golang.org/x/image/draw"
)

func main() {
	inFrameStart := int64(1)
	inFrameEnd := int64(0)

	threadCount := runtime.NumCPU()
	invert := false
	repRepeat := 1

	quadMinSize := 4
	quadTolerance := uint8(5)

	inDir := "in"
	outDir := "out"
	repDir := "rep"

	var err error

	for i := 1; i < len(os.Args); i++ {
		if os.Args[i][0:2] == "--" {
			switch os.Args[i][2:] {
				case "inFrameStart":
					inFrameStart, err = strconv.ParseInt(os.Args[i+1], 10, 64)
					assert(err)
					break
				case "inFrameEnd":
					inFrameEnd, err = strconv.ParseInt(os.Args[i+1], 10, 64)
					break
				case "threadCount":
					tc, err := strconv.ParseInt(os.Args[i+1], 10, 32)
					assert(err)
					threadCount = int(tc)
					break
				case "invert":
					invert = true
					continue
				case "quadMinSize":
					qms, err := strconv.ParseInt(os.Args[i+1], 10, 32)
					assert(err)
					quadMinSize = int(qms)
					break
				case "quadTolerance":
					qt, err := strconv.ParseInt(os.Args[i+1], 10, 32)
					assert(err)
					quadMinSize = int(qt)
					break
				case "repRepeat":
					rr, err := strconv.ParseInt(os.Args[i+1], 10, 32)
					assert(err)
					repRepeat = int(rr)
					break;
				case "inDir":
					inDir = os.Args[i+1]
					break
				case "outDir":
					outDir = os.Args[i+1]
					break
				case "repDir":
					repDir = os.Args[i+1]
					break
				default:
					fmt.Printf("Unkown option \"%s\"\n", os.Args[i])
					os.Exit(-1)
			}
			i++
		} else {
			fmt.Printf("Unknown argument \"%s\"\n", os.Args[i])
			os.Exit(-1)
		}
	}

	_, err = os.Stat(outDir)
	if err != nil {
		os.Mkdir(outDir, os.ModePerm)
	}

	repFramesDir, err := os.ReadDir(repDir)
	repFrames := make([]*image.RGBA, len(repFramesDir))
	assert(err)
	for i,e := range repFramesDir {
		f, err := os.Open(fmt.Sprintf("%s/%s", repDir, e.Name()))
		assert(err)
		repFrames[i] = rgbapls(f)
		fmt.Printf("rep frame %d - %s\n", i, e.Name())
	}

	workQueue := make(chan *frameInfo, 100)
	threadDoneChan := make(chan uint8, threadCount)

	go (func () {
		repFrameCounter := 0

		for i := inFrameStart; inFrameEnd > 0 && i < inFrameEnd; i++ {
			f, err := os.Open(fmt.Sprintf("%s/%d.png", inDir, i))
			if err != nil { break }
			workQueue <- &frameInfo{
				repframe: repFrames[(repFrameCounter % (len(repFrames)*repRepeat))/repRepeat],
				frame: rgbapls(f),
				outname: fmt.Sprintf("%s/%d.png", outDir, i),
			}
			repFrameCounter++
		}

		close(workQueue)
		fmt.Println("Load loop end")
	})()

	for i := 0; i < threadCount; i++ {
		go (func() {
			for fi := range workQueue {
				quadded := quadifyImage(fi.frame, fi.repframe, quadMinSize, quadTolerance, invert)
				f, err := os.Create(fi.outname)
				assert(err)
				png.Encode(f, quadded)
				fmt.Printf("%s done\n", fi.outname)
			}

			fmt.Printf("Work thread %d done\n", i)
			threadDoneChan<-0
		})()

	}

	for i := 0; i < threadCount; i++ {
		<-threadDoneChan
	}

	fmt.Println("All done")
}

func rgbapls(r io.Reader) *image.RGBA {
	a, _, err := image.Decode(r)
	assert(err)
	img := image.NewRGBA(a.Bounds())
	draw.Draw(img, img.Rect, a, image.ZP, draw.Src)
	return img
}

type frameInfo struct {
	repframe *image.RGBA
	frame *image.RGBA
	outname string
}

func assert(e error) {
	if e != nil {
		fmt.Println("PANIC!!!!")
		panic(e)
	}
}
