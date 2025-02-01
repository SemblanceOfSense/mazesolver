package bot

import (
	"fmt"
	"log"
	getMaze "mazesolver/internal/getmaze"
	"mazesolver/internal/outputmaze"
	"mazesolver/internal/solvemaze"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

var (
    appId = "1335044960898252830"
    guildId = ""
)

func Run(BotToken string) {
    discord, err := discordgo.New(("Bot " + BotToken))
    if err != nil { fmt.Println("Bot 1"); log.Fatal(err) }

    _, err = discord.ApplicationCommandBulkOverwrite(appId, guildId, []*discordgo.ApplicationCommand {
        {
            Name: "solve-maze",
            Description: "solve one of tilley's mazes",
            Options: []*discordgo.ApplicationCommandOption {
                    {
                        Type: discordgo.ApplicationCommandOptionString,
                        Name: "image-link",
                        Description: "link to maze image",
                        Required: true,
                    },
                },
            },
        },
    )
    if err != nil { fmt.Println("Bot 2"); log.Fatal(err) }

    discord.AddHandler(func (
        s *discordgo.Session,
        i *discordgo.InteractionCreate,
    ) {
        if i.Type == discordgo.InteractionApplicationCommand {
            data := i.ApplicationCommandData()
            responseData := ""
            if i.Interaction.Member.User.ID == s.State.User.ID { return; }

            switch data.Name {
            case "solve-maze":
                maze, err := getMaze.GetMaze(i.ApplicationCommandData().Options[0].Value.(string))
                if err != nil {
                    responseData = "You must provide an image! Select the maze, click \"open in browser\", and copy the link of the image"
                } else {
                    if (len(maze) < 1) {
                        responseData = "You must provide an image! Select the maze, click \"open in browser\", and copy the link of the image"
                    } else {
                        points := solvemaze.FindPath(maze)
                        if (len(points) < 1) {
                            responseData = "You must provide an image! Select the maze, click \"open in browser\", and copy the link of the image"
                        } else {
                            outputmaze.EditMaze(points, "/tmp/maze.png", "/tmp/outputmaze.png")
            }}}}

            if responseData != "" {
                err = s.InteractionRespond(
                    i.Interaction,
                    &discordgo.InteractionResponse{
                        Type: discordgo.InteractionResponseChannelMessageWithSource,
                        Data: &discordgo.InteractionResponseData{
                            Flags: 1 << 6,
                            Content: responseData,
                        },
                    },
                )
            }
            fileName := "/tmp/outputmaze.png"
			f, err := os.Open(fileName)
			defer f.Close()
            err = s.InteractionRespond(
                i.Interaction,
                &discordgo.InteractionResponse{
                    Type: discordgo.InteractionResponseChannelMessageWithSource,
                    Data: &discordgo.InteractionResponseData{
                            Files: []*discordgo.File{
                            &discordgo.File{
                                Name:  fileName,
                                Reader: f,
                            },
                        },
                    },
                },
            )
            if err != nil {
                fmt.Println(err)
            }
        }
    })

    discord.AddHandler( func(
        s *discordgo.Session,
        m *discordgo.MessageCreate,
    ) {
        listOfGuilds := discord.State.Guilds
        for _, v := range listOfGuilds {
            fmt.Println(v.Name)
        }
    })

    err = discord.Open()
    if err != nil { log.Fatal(err) }

    stop := make (chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt)
    log.Println("Press Ctrl+C to Exit")
    <-stop

    err = discord.Close()
    if err != nil { log.Fatal(err) }
}
