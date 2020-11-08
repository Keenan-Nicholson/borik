package bot

import (
	"fmt"
	"io"

	"github.com/rs/zerolog/log"
	"gopkg.in/gographics/imagick.v2/imagick"
)

// Magik runs content-aware scaling on an image.
func Magik(src []byte, dest io.Writer, scale float64) error {
	wand := imagick.NewMagickWand()
	wand.ReadImageBlob(src)

	width := wand.GetImageWidth()
	height := wand.GetImageHeight()

	log.Debug().
		Uint("src_width", width).
		Uint("src_height", height).
		Uint("dest_width", width/2).
		Uint("dest_height", height/2).
		Msg("Liquid rescaling image")
	err := wand.LiquidRescaleImage(uint(width/2), uint(height/2), scale, 0)
	if err != nil {
		return fmt.Errorf("error while attempting to liquid rescale: %w", err)
	}

	log.Debug().
		Uint("dest_width", width).
		Uint("dest_height", height).
		Uint("src_width", width/2).
		Uint("src_height", height/2).
		Msg("Returning image to original size")
	err = wand.ResizeImage(width, height, imagick.FILTER_LANCZOS, 1)
	if err != nil {
		return fmt.Errorf("error while attempting to resize image: %w", err)
	}

	_, err = dest.Write(wand.GetImageBlob())
	if err != nil {
		return fmt.Errorf("error writing image: %w", err)
	}

	return nil
}

// Arcweld destroys an image via a combination of operations.
func Arcweld(src []byte, dest io.Writer) error {
	wand := imagick.NewMagickWand()
	wand.ReadImageBlob(src)

	err := wand.EvaluateImageChannel(imagick.CHANNEL_RED, imagick.EVAL_OP_LEFT_SHIFT, 1)
	if err != nil {
		return fmt.Errorf("error left-shifting red channel: %w", err)
	}

	err = wand.ContrastStretchImage(0.3, 0.3)
	if err != nil {
		return fmt.Errorf("error contrast stretching image: %w", err)
	}

	err = wand.EvaluateImageChannel(imagick.CHANNEL_RED, imagick.EVAL_OP_THRESHOLD_BLACK, 0.9)
	if err != nil {
		return fmt.Errorf("error running threshold black: %w", err)
	}

	err = wand.SharpenImage(0, 0)
	if err != nil {
		return fmt.Errorf("error sharpening image: %w", err)
	}

	width := wand.GetImageWidth()
	height := wand.GetImageHeight()

	err = wand.LiquidRescaleImage(width/2, height/3, 1, 0)
	if err != nil {
		return fmt.Errorf("error liquid rescaling: %w", err)
	}

	width = wand.GetImageWidth()
	height = wand.GetImageHeight()

	err = wand.LiquidRescaleImage(width*2, height*3, 0.4, 0)
	if err != nil {
		return fmt.Errorf("error liquid rescaling: %w", err)
	}

	err = wand.ImplodeImage(0.2)
	if err != nil {
		return fmt.Errorf("error imploding image: %w", err)
	}

	err = wand.QuantizeImage(8, imagick.COLORSPACE_RGB, 0, true, false)
	if err != nil {
		return fmt.Errorf("error quantizing image: %w", err)
	}

	_, err = dest.Write(wand.GetImageBlob())
	if err != nil {
		return fmt.Errorf("error writing image: %w", err)
	}

	return nil
}