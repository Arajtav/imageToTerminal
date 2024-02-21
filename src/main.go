package main

import (
    "golang.org/x/sys/unix" // terminal size
    "fmt"
    "flag"
    "math"
    "os"

    "image"
    _ "image/jpeg"
    _ "image/png"
)

// TODO, return number of rows and cols that printed image will have, so it will keep aspect ratio
func rescale(imgx uint16, imgy uint16, termx uint16, termy uint16) (uint16, uint16) {
    return termx, termy;
}

// simplest sampling, just get pixel with given coordinate
func getRgbFast(img *image.Image, x float64, y float64) (uint16, uint16, uint16) {
    rx := int(math.Floor(float64((*img).Bounds().Dx())*x));
    ry := int(math.Floor(float64((*img).Bounds().Dy())*y));
    r, g, b, _ := (*img).At((*img).Bounds().Min.X+rx, (*img).Bounds().Min.Y+ry).RGBA();
    return uint16(r/0xff), uint16(g/0xff), uint16(b/0xff);
}

// printed pixel will be average value of pixels that original image had in this place
func getRgbAvg(img* image.Image, x float64, y float64, pxsizex float64, pxsizey float64) (uint16, uint16, uint16) {
    rx := int(math.Floor(float64((*img).Bounds().Dx())*x));
    ry := int(math.Floor(float64((*img).Bounds().Dy())*y));
    rxt := int(math.Floor(float64((*img).Bounds().Dx())*(x+pxsizex)));
    ryt := int(math.Floor(float64((*img).Bounds().Dy())*(y+pxsizey)));
    var r, g, b uint64;
    for i := rx; i<rxt; i++ {
        for j := ry; j<ryt; j++ {
            cr, cg, cb, _ := (*img).At((*img).Bounds().Min.X+i, (*img).Bounds().Min.Y+j).RGBA();
            r += uint64(cr/0xff); g += uint64(cg/0xff); b += uint64(cb/0xff);
        }
    }
    r /= uint64((rxt-rx)*(ryt-ry));
    g /= uint64((rxt-rx)*(ryt-ry));
    b /= uint64((rxt-rx)*(ryt-ry));
    return uint16(r), uint16(g), uint16(b);
}

// for debugging stuff
func getUv(img* image.Image, x float64, y float64) (uint16, uint16, uint16) {
    return uint16(x*255.0), uint16(y*255.0), 0;
}

func getRgb(img *image.Image, x float64, y float64, pxsizex float64, pxsizey float64, mode *string) (uint16, uint16, uint16) {
    if *mode == "fast" {
        return getRgbFast(img, x, y);
    } else if *mode == "average" {
        return getRgbAvg(img, x, y, pxsizex, pxsizey);
    }
    return getUv(img, x, y); // this shouldn't happen because there is check for what sampling can be earlier
}

func main() {
    sampling := flag.String("sampling", "fast", "Sampling mode (fast/average)");
    dql := flag.Bool("dql", false, "Double quality");
    ftermx := flag.Uint64("width", 0, "Width");
    ftermy := flag.Uint64("height", 0, "Height");
    flag.Parse();

    if !(*sampling == "average" || *sampling == "fast" || *sampling == "uv") {
        fmt.Fprintln(os.Stderr, "Invalid value specified for sampling");
        os.Exit(22);
    }

    if flag.NArg() < 1 {
        fmt.Fprintln(os.Stderr, "You need to specify file");
        os.Exit(22);
    }

    if flag.NArg() > 1 {
        fmt.Fprintln(os.Stderr, "Too many arguments");
        os.Exit(22);
    }

    file, err := os.Open(flag.Args()[0]);
    if err != nil {
        fmt.Fprintf(os.Stderr, "Could not open file %s\n", flag.Args()[0]);
        os.Exit(1);
    }
    defer file.Close()

    img, _, err := image.Decode(file);
    if err != nil {
        fmt.Fprintf(os.Stderr, "Could not decode image %s\n", flag.Args()[0]);
        os.Exit(1);
    }

    ws, err := unix.IoctlGetWinsize(1, unix.TIOCGWINSZ);
    ws.Row -= 1;
    if *ftermx != 0 { ws.Col = uint16(*ftermx); }
    if *ftermy != 0 { ws.Row = uint16(*ftermy); }
    if err != nil {
        fmt.Fprintln(os.Stderr, "Error getting terminal size");
        os.Exit(1);
    }

    c, r := rescale(uint16(img.Bounds().Dx()), uint16(img.Bounds().Dy()), ws.Col, ws.Row);
    fmt.Print("\033[0m");
    for i := uint16(0); i<r; i++ {
        for j := uint16(0); j<c; j++ {
            if *dql {
                bgr, bgg, bgb := getRgb(&img, (1.0/float64(c))*float64(j), (1.0/float64(r))*float64(i), 1.0/float64(c), 0.5/float64(r), sampling);
                fgr, fgg, fgb := getRgb(&img, (1.0/float64(c))*float64(j), (1.0/float64(r))*(float64(i)+0.5), 1.0/float64(c), 0.5/float64(r), sampling);
                fmt.Printf("\033[48;2;%d;%d;%dm\033[38;2;%d;%d;%dm▄", bgr, bgg, bgb, fgr, fgg, fgb);
            } else {
                bgr, bgg, bgb := getRgb(&img, (1.0/float64(c))*float64(j), (1.0/float64(r))*float64(i), 1.0/float64(c), 1.0/float64(r), sampling);
                fmt.Printf("\033[48;2;%d;%d;%dm ", bgr, bgg, bgb);
            }
        }
        fmt.Println("\033[0m");
    }
}
