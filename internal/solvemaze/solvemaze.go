package solvemaze

import (
	"mazesolver/internal/getmaze"
	"mazesolver/util/queue"
)

func FindPath(maze getMaze.Maze) []getMaze.Point {
    returnSlice := make([]getMaze.Point, 0)
    //debugMaze := GetDebugMaze(maze)
    //findPathRecursive(maze[1][1], maze, &returnSlice, debugMaze)
    findPathNonRecursive(maze, maze[1][1], &returnSlice)
    return returnSlice
}

func findPathRecursive(next getMaze.Point, maze getMaze.Maze, slice *[]getMaze.Point, debugMaze getMaze.Maze) bool {
    maze[next.Y][next.X].Value = -1
    debugMaze[next.Y][next.X].Value = 1
    for _, v := range getAdjacent(next, maze) {
        switch v.Value {
        case -1:
            continue
        case 5:
            continue
        case 0:
            if (findPathRecursive(v, maze, slice, debugMaze)) {
                *slice = append(*slice, v)
                return true
            }
        case 9:
            *slice = append(*slice, v)
            return true
        }
    }
    return false
}

func getAdjacent(point getMaze.Point, maze getMaze.Maze) []getMaze.Point {
    returnSlice := make([]getMaze.Point, 0)
    returnSlice = append(returnSlice, maze[point.Y + 1][point.X])
    returnSlice = append(returnSlice, maze[point.Y - 1][point.X])
    returnSlice = append(returnSlice, maze[point.Y][point.X + 1])
    returnSlice = append(returnSlice, maze[point.Y][point.X - 1])
    return returnSlice
}

func GetDebugMaze(maze getMaze.Maze) getMaze.Maze {
    a := make(getMaze.Maze, len(maze))
    for i := range a {
        a[i] = make([]getMaze.Point, len(maze[0]))
    }
    return a
}

func findPathNonRecursive(maze getMaze.Maze, root getMaze.Point, slice *[]getMaze.Point) {
    next := buildMaze(maze, root)
    for next != nil {
        *slice = append(*slice, *next)
        next = next.Parent
    }
}

func buildMaze(maze getMaze.Maze, root getMaze.Point) *getMaze.Point {
    Q := make(queue.Queue[getMaze.Point], 0)
    root.Value = -1
    Q.Enqueue(root)
    for len(Q) > 0 {
        v := Q.Dequeue()
        if v.Value == 9 {
            return (&v)
        }
        for _, w := range getAdjacent(v, maze) {
            if (w.Value == 0 || w.Value == 9) {
                maze[w.Y][w.X].Value = -1
                w.Parent = &v
                Q.Enqueue(w)
            }
        }
    }
    return nil
}

func contains(maze getMaze.Maze, val int) bool {
    for _, v := range maze {
        for _, vv := range v {
            if vv.Value == val {
                return true
            }
        }
    }
    return false
}
