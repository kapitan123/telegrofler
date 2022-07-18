package watermarker

import (
	"bytes"
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"time"

	"github.com/nfnt/resize"

	log "github.com/sirupsen/logrus"
)

type Watermarker struct {
}

func New() *Watermarker {
	return &Watermarker{}
}

func (w *Watermarker) Apply(bakground []byte, foreground []byte) ([]byte, error) {
	defer logDuration(time.Now())

	bgImg, _, err := image.Decode(bytes.NewReader(bakground))
	if err != nil {
		return nil, err
	}

	fgImg, err := png.Decode(bytes.NewReader(foreground))

	if err != nil {
		return nil, err
	}

	b := bgImg.Bounds()
	resImg := image.NewRGBA(b)
	draw.Draw(resImg, b, bgImg, image.Point{}, draw.Src)
	draw.Draw(resImg, fgImg.Bounds(), fgImg, image.Point{}, draw.Over)

	resBuf := new(bytes.Buffer)
	err = png.Encode(resBuf, resImg)

	if err != nil {
		return nil, err
	}

	return resBuf.Bytes(), nil
}

func (w *Watermarker) ApplyOnAxis(bakground []byte, foreground []byte, x int, y int, mouthLenght int) ([]byte, error) {
	defer logDuration(time.Now())

	bgImg, _, err := image.Decode(bytes.NewReader(bakground))
	if err != nil {
		return nil, err
	}

	fgImg, err := png.Decode(bytes.NewReader(foreground))

	if err != nil {
		return nil, err
	}

	bgbounds := bgImg.Bounds()
	fgBounds := fgImg.Bounds()

	resImg := image.NewRGBA(bgbounds)
	draw.Draw(resImg, bgbounds, bgImg, bgbounds.Min, draw.Src)

	//resizing the watermark
	defaultMouthLenght := 110.0000
	resizeMultiplier := float64(mouthLenght) / defaultMouthLenght

	newWidth := float64(fgBounds.Max.X) * resizeMultiplier
	newHeight := float64(fgBounds.Max.Y) * resizeMultiplier

	fgImg = resize.Resize(uint(newWidth), uint(newHeight), fgImg, resize.Lanczos3)

	//drawToRect := fgImg.Bounds()
	dp := image.Pt(x, y)

	r := image.Rectangle{dp, dp.Add(fgBounds.Size())}
	draw.Draw(resImg, r, fgImg, fgBounds.Min, draw.Over)

	resBuf := new(bytes.Buffer)
	err = png.Encode(resBuf, resImg)

	if err != nil {
		return nil, err
	}

	return resBuf.Bytes(), nil
}

func logDuration(invocation time.Time) {
	elapsed := time.Since(invocation)

	log.Printf("%s lasted %s", "watermarking", elapsed)
}
