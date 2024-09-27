package main

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/term"

	"image"
	_ "image/jpeg"
	_ "image/png"

	"net/http"
)

// WONT WORK FOR NEGATIVE VALUES, IF VALUE IS NEGATIVE SUBTRACT 1 BEFORE CONVERSION (OR AFTER I DON'T THINK THERE WILL BE DIFFERENCE)
func floor(val float64) int {
	return int(val)
}

// TODO, return number of rows and cols that printed image will have, so it will keep aspect ratio
func rescale(imgx uint16, imgy uint16, termx uint16, termy uint16) (float64, float64) {
	return float64(termx), float64(termy)
}

// simplest sampling, just get pixel with given coordinate
func getRgbFast(img *image.Image, x float64, y float64) (uint16, uint16, uint16) {
	rx := floor(float64((*img).Bounds().Dx()) * x)
	ry := floor(float64((*img).Bounds().Dy()) * y)
	r, g, b, _ := (*img).At((*img).Bounds().Min.X+rx, (*img).Bounds().Min.Y+ry).RGBA()
	return uint16(r / 0xff), uint16(g / 0xff), uint16(b / 0xff)
}

// printed pixel will be average value of pixels that original image had in this place
func getRgbAvg(img *image.Image, x float64, y float64, pxsizex float64, pxsizey float64) (uint16, uint16, uint16) {
	rx := floor(float64((*img).Bounds().Dx()) * x)
	ry := floor(float64((*img).Bounds().Dy()) * y)
	rxt := floor(float64((*img).Bounds().Dx()) * (x + pxsizex))
	ryt := floor(float64((*img).Bounds().Dy()) * (y + pxsizey))
	var r, g, b uint64
	for i := rx; i < rxt; i++ {
		for j := ry; j < ryt; j++ {
			cr, cg, cb, _ := (*img).At((*img).Bounds().Min.X+i, (*img).Bounds().Min.Y+j).RGBA()
			r += uint64(cr / 0xff)
			g += uint64(cg / 0xff)
			b += uint64(cb / 0xff)
		}
	}
	r /= uint64((rxt - rx) * (ryt - ry))
	g /= uint64((rxt - rx) * (ryt - ry))
	b /= uint64((rxt - rx) * (ryt - ry))
	return uint16(r), uint16(g), uint16(b)
}

// for debugging stuff
func getUv(x float64, y float64) (uint16, uint16, uint16) {
	return uint16(x * 255.0), uint16(y * 255.0), 0
}

func getRgb(img *image.Image, x float64, y float64, pxsizex float64, pxsizey float64, mode *string) (uint16, uint16, uint16) {
	if *mode == "fast" {
		return getRgbFast(img, x, y)
	} else if *mode == "average" {
		return getRgbAvg(img, x, y, pxsizex, pxsizey)
	}
	return getUv(x, y) // shouldn't happen
}

func main() {
	sampling := flag.String("sampling", "fast", "Sampling mode (fast/average)")
	dql := flag.Bool("dql", false, "Double quality")
	net := flag.Bool("net", false, "Loading images from network")
	ftermx := flag.Int("width", 0, "Width")
	ftermy := flag.Int("height", 0, "Height")
	flag.Parse()

	if !(*sampling == "average" || *sampling == "fast" || *sampling == "uv") {
		fmt.Fprintln(os.Stderr, "Invalid value specified for sampling")
		os.Exit(22)
	}

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "You need to specify file")
		os.Exit(22)
	}

	if flag.NArg() > 1 {
		fmt.Fprintln(os.Stderr, "Too many arguments")
		os.Exit(22)
	}

	col, row, err := term.GetSize(int(os.Stdout.Fd()))
	row -= 1 // there will be always new line on end for shell command prompt
	if *ftermx != 0 {
		col = *ftermx
	}
	if *ftermy != 0 {
		row = *ftermy
	}
	if err != nil && (col == 0 || row+1 == 0) {
		fmt.Fprintln(os.Stderr, "Error getting terminal size")
		os.Exit(1)
	}

	var img image.Image

	if !(*net) {
		file, err := os.Open(flag.Args()[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not open file %s\n", flag.Args()[0])
			os.Exit(1)
		}
		defer file.Close()

		img, _, err = image.Decode(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not decode image %s\n", flag.Args()[0])
			os.Exit(1)
		}
	} else {
		response, err := http.Get(flag.Args()[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to fetch image %s\n", flag.Args()[0])
			os.Exit(1)
		}
		defer response.Body.Close()

		if response.StatusCode != 200 {
			fmt.Fprintf(os.Stderr, "Server responded with code %d\n", response.StatusCode)
			os.Exit(1)
		}

		img, _, err = image.Decode(response.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not decode image %s\n", flag.Args()[0])
			os.Exit(1)
		}
	}
	if col < 1 || row < 1 {
		fmt.Fprintf(os.Stderr, "Invalid print size (%d, %d)\n", col, row)
	}
	c, r := rescale(uint16(img.Bounds().Dx()), uint16(img.Bounds().Dy()), uint16(col), uint16(row))
	fmt.Print("\033[0m")
	for i := 0.0; i < r; i++ {
		for j := 0.0; j < c; j++ {
			if *dql {
				bgr, bgg, bgb := getRgb(&img, j/c, i/r, 1.0/c, 0.5/r, sampling)
				fgr, fgg, fgb := getRgb(&img, j/c, (i+0.5)/r, 1.0/c, 0.5/r, sampling)
				fmt.Printf("\033[48;2;%d;%d;%dm\033[38;2;%d;%d;%dm▄", bgr, bgg, bgb, fgr, fgg, fgb)
			} else {
				bgr, bgg, bgb := getRgb(&img, j/c, i/r, 1.0/c, 1.0/r, sampling)
				fmt.Printf("\033[48;2;%d;%d;%dm ", bgr, bgg, bgb)
			}
		}
		fmt.Println("\033[0m")
	}
}
