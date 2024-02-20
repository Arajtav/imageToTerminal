package main

import (
    "golang.org/x/sys/unix"
    "fmt"
    "os"
    "image"
    _ "image/jpeg"
    _ "image/png"
    "math"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Fprintln(os.Stderr, "You need to specify file");
        os.Exit(22);
    }

    if len(os.Args) > 2 {
        fmt.Fprintln(os.Stderr, "Too many args");
        os.Exit(22);
    }

    file, err := os.Open(os.Args[1]);
    if err != nil {
        fmt.Fprintf(os.Stderr, "Could not open file %s\n", os.Args[1]);
        os.Exit(1);
    }
    defer file.Close()

    img, _, err := image.Decode(file);
    if err != nil {
        fmt.Fprintf(os.Stderr, "Could not decode image %s\n", os.Args[1]);
        os.Exit(1);
    }

    ws, err := unix.IoctlGetWinsize(1, unix.TIOCGWINSZ);
    ws.Row -= 1;
    if err != nil {
        fmt.Fprintln(os.Stderr, "Error getting terminal size");
        os.Exit(1);
    }

    fmt.Print("\033[0m");
    for i := uint16(0); i<ws.Row; i++ {
        for j := uint16(0); j<ws.Col; j++ {
            var x int = img.Bounds().Min.X + int(math.Floor(float64(j)*(float64(img.Bounds().Dx())/float64(ws.Col))));
            var y int = img.Bounds().Min.Y + int(math.Floor(float64(i)*(float64(img.Bounds().Dy())/float64(ws.Row))));
            r, g, b, _ := img.At(x, y).RGBA();
            fmt.Printf("\033[48;2;%d;%d;%dm ", r/0xff, g/0xff, b/0xff);
        }
        fmt.Println("\033[0m");
    }
}
