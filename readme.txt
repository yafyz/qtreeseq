args:
    --inFrameStart int64
        first frame number

    --inFrameEnd int64
        last frame number

    --threadCount int
        thread count

    --invert
        inverts alpha

    --quadMinSize int
        min width/height of quadrant

    --quadTolerance int
        color diff tolerance

    --repRepeat int
        rep frame repeat count

    --inDir string
        input dir

    --outDir string
        output dir

    --repDir string
        replace frame dir

useful ffmpeg commands:
    gen pallete from sequence:
        ffmpeg -i %d.png -vf palettegen=reserve_transparent=1 palette.png

    gen transparent gif from sequence and pallet:
        ffmpeg -framerate 25 -i %d.png -i palette.png -lavfi paletteuse=alpha_threshold=128 -gifflags -offsetting out.gif