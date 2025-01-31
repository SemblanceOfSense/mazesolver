package getMaze

import (
	"errors"
	"fmt"
	"image/color"
	"image/png"
	"os"
)

/*
Pixel types:
-1 - explored
0 - unexplored
5 - wall
9 - end
*/

type Point struct {
    X int;
    Y int;
    Value int
}

type Maze [][]Point

func GetMaze(imagepath string) (Maze, error) {
    returnMaze := *new(Maze)

    imagereader, err := os.Open(imagepath)
    if err != nil {
        return *new(Maze), err
    }

    image, err := png.Decode(imagereader)
    if err != nil {
        return *new(Maze), err
    }

    for y := image.Bounds().Min.Y; y < image.Bounds().Max.Y; y += 10 {
        newRow := make([]Point, 0)
        for x := image.Bounds().Min.X; x < image.Bounds().Max.X; x+= 10 {
            typ := 0
            switch image.At(x, y) {
            case color.RGBA{ R: 255, G: 0, B: 0, A: 255 }:
                typ = 0
            case color.RGBA{ R: 0, G: 0, B: 255, A: 255 }:
                typ = 9
            case color.RGBA{ R: 255, G: 255, B: 255, A: 255 }:
                typ = 0
            case color.RGBA{ R: 0, G: 0, B: 0, A: 255 }:
                typ = 5
            default:
                fmt.Println(image.At(x, y))
                return *new(Maze), errors.New("bad color")
            }
            newPoint := Point{
                X: x / 10,
                Y: y / 10,
                Value: typ,
            }
            newRow = append(newRow, newPoint)
        }
        returnMaze = append(returnMaze, newRow)
    }


    return returnMaze, err
}

func PrintMaze(maze Maze) {
        for _, v := range maze {
        fmt.Print("[")
        for _, vv := range v {
            fmt.Print(vv.Value)
            fmt.Print(", ")
        }
        fmt.Println("]")
    }
}
