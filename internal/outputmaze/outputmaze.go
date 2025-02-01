package outputmaze

import (
	"errors"
	"image"
	"image/color"
	"image/png"
	"math"
	getMaze "mazesolver/internal/getmaze"
	"os"
)

func EditMaze(points []getMaze.Point, oldPath, newPath string) (string, error) {
    imagereader, err := os.Open(oldPath)
    if err != nil {
        return "", err
    }

    image, err := png.Decode(imagereader)
    if err != nil {
        return "", err
    }

    for n, v := range points {
        var h float64 = (float64(1) - (((float64(n)) / (float64(len(points)))))) / float64(1.1)
        r, g, b := hueToRGB(h)
        updateColor(image, v, color.RGBA{uint8(r * 255), uint8(g * 255), uint8(b * 255), 255})
    }

    f, err := os.Create(newPath)
    if err != nil {
        return "", err
    }
    png.Encode(f, image)
    return newPath, nil
}

type Changeable interface {
    Set(x, y int, c color.Color)
}

func updateColor(image image.Image, p getMaze.Point, color color.Color) error {
    for i := p.Y * 10; i < p.Y * 10 + 10; i++ {
        for ii := p.X * 10; ii < p.X * 10 + 10; ii++ {
            if cimg, ok := image.(Changeable); ok {
                cimg.Set(ii, i, color)
            } else {
                return errors.New("Image not changeable")
            }
        }
    }
    return nil
}

func hueToRGB(h float64) (float64, float64, float64) {
    kr := math.Mod(5+h*6, 6)
    kg := math.Mod(3+h*6, 6)
    kb := math.Mod(1+h*6, 6)

    r := 1 - math.Max(min3(kr, 4-kr, 1), 0)
    g := 1 - math.Max(min3(kg, 4-kg, 1), 0)
    b := 1 - math.Max(min3(kb, 4-kb, 1), 0)

    return r, g, b
}

func min3(a, b, c float64) float64 {
    return math.Min(math.Min(a, b), c)
}
