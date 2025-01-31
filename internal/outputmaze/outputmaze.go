package outputmaze

import (
	"errors"
	"image"
	"image/color"
	"image/png"
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

    for _, v := range points {
        updateColor(image, v, color.RGBA{0, 255, 0, 255})
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
