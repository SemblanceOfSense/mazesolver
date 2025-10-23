package main

import (
	"flag"
	"fmt"
	"log"
	"mazesolver/internal/bot"
	getMaze "mazesolver/internal/getmaze"
	"mazesolver/internal/outputmaze"
	"mazesolver/internal/solvemaze"
)

var BotToken string
var MazePath string

func init() {
    flag.StringVar(&BotToken, "bottoken", "", "discord bot token")
	flag.StringVar(&MazePath, "path", "", "file path to maze you want to solve")

    flag.Parse()
}

func main() {
	if BotToken != "" {
    	bot.Run(BotToken)
	} else if MazePath != "" {
		maze, err := getMaze.GetMaze(MazePath)
		if err != nil {
			log.Fatal(err)
		}
		path := solvemaze.FindPath(maze)
		newpath, err := outputmaze.EditMaze(path, MazePath, MazePath + ".out")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Maze saved to " + newpath)
	} else {
		fmt.Println("Please provide a bot token or maze file!")
	}
}
