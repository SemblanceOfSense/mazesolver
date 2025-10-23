package getMaze

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"net/http"
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
    Value int;
    Parent *Point
}

type Maze [][]Point

func GetURLMaze(imagepath string /*URL*/) (Maze, error) {
	returnMaze := *new(Maze)

	res, err := http.Get(imagepath)

    if err != nil {
        return returnMaze, err
    }

    data, err := io.ReadAll(res.Body)

    if err != nil {
        return returnMaze, err
    }

    defer res.Body.Close()

    err = os.WriteFile("/tmp/maze.png", data, 0755)
    if err != nil {
        return returnMaze, err
    }

	returnMaze, err = GetMaze("/tmp/maze.png")
	return  returnMaze, err
}

func GetMaze(imagepath string /*Filepath*/) (Maze, error) {
    returnMaze := *new(Maze)

    imagereader, err := os.Open(imagepath)
    if err != nil {
        return returnMaze, err
    }

    image, err := png.Decode(imagereader)
    if err != nil {
        return *new(Maze), err
    }

    psize := DeterminePixelSize(image)
    for y := image.Bounds().Min.Y; y < image.Bounds().Max.Y; y += psize {
        newRow := make([]Point, 0)
        for x := image.Bounds().Min.X; x < image.Bounds().Max.X; x += psize {
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
                X: x / psize,
                Y: y / psize,
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

func DeterminePixelSize(image image.Image) int {
    var x, y int
    color := image.At(x, y)
    for {
        if color != image.At(x, y) {
            break
        }
        x++
        y++
    }
    return x
}
