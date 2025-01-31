package main

import (
	"fmt"
	"log"
	getMaze "mazesolver/internal/getmaze"
	"mazesolver/internal/outputmaze"
	"mazesolver/internal/solvemaze"
)

func main() {
    path := "/tmp/maze.png"
    maze, err := getMaze.GetMaze(path)
    if err != nil {
        log.Fatal(err)
    }

    p := solvemaze.FindPath(maze)
    if err != nil {
        log.Fatal(err)
    }

    newpath := "/tmp/outputmaze.png"
    _, err = outputmaze.EditMaze(p, path, newpath)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(newpath)
}
