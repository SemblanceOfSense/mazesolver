package bot

import (
	"fmt"
	"log"
	getMaze "mazesolver/internal/getmaze"
	"mazesolver/internal/outputmaze"
	"mazesolver/internal/solvemaze"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"

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
                        Name: "message-link",
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

            var solutionlength string
            switch data.Name {
            case "solve-maze":
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
                messageUrl := i.ApplicationCommandData().Options[0].Value.(string);
                if strings.Contains(messageUrl, "media.discordapp.net") {
                    messageUrl = strings.Replace(messageUrl, "media.discordapp.net", "cdn.discordapp.com", 1)

                    re := regexp.MustCompile(`(hm=[a-f0-9]+).*`)
                    messageUrl = re.ReplaceAllString(messageUrl, `$1&`)
                }
                if err != nil {
                    fmt.Println(err)
                } else {
                    message, err := s.ChannelMessage(i.ChannelID, strings.Split(messageUrl, "/")[len(strings.Split(messageUrl, "/")) - 1]);
                    var maze getMaze.Maze
                    if strings.Contains(messageUrl, "cdn") {
                        maze, err = getMaze.GetMaze(messageUrl)
                    } else if err == nil && len(message.Attachments) > 1 {
                        maze, err = getMaze.GetMaze(message.Attachments[0].URL)
                    } else if err == nil {
                        maze, err = getMaze.GetMaze(message.Content)
                    }
                    if err != nil {
                        fmt.Println(err)
                        if strings.Contains(err.Error(), "403") {
                            responseData = "Acess Denied"
                        } else {
                            responseData = "You must provide a valid image! Provide the message link of a valid image."
                        }
                    } else {
                        if (len(maze) < 1) {
                            responseData = "You must provide a valid image! Provide the message link of a valid image."
                        } else {
                            points := solvemaze.FindPath(maze)
                            if (len(points) < 1) {
                                responseData = "You must provide a valid image! Provide the message link of a valid image."
                            } else {
                                outputmaze.EditMaze(points, "/tmp/maze.png", "/tmp/outputmaze.png")
                                solutionlength = strconv.Itoa(len(points))
                                maze = nil
                                points = nil
            }}}}}

            if responseData == "" {
                fileName := "/tmp/outputmaze.png"
                f, _ := os.Open(fileName)
                defer f.Close()
                _, err = s.FollowupMessageCreate(
                    i.Interaction,
                    false,
                    &discordgo.WebhookParams{
                                Content: "Length of solution is: " + solutionlength,
                                Files: []*discordgo.File{
                                &discordgo.File{
                                    Name:  fileName,
                                    Reader: f,
                            },
                        },
                    },
                )
                if err != nil {
                    if strings.Contains(err.Error(), "413") {
                        responseData = "Solution file is too large!"
                    } else {
                        responseData = err.Error()
                    }
                }
            }
            if responseData != "" {
                _, err = s.FollowupMessageCreate(
                    i.Interaction,
                    false,
                    &discordgo.WebhookParams{
                            Flags: 1 << 6,
                            Content: responseData,
                    },
                )
            }
        }
    })

    discord.AddHandler(func (
        s *discordgo.Session,
        m *discordgo.MessageCreate,
    ) {
        if (m.Content == "-solve") {
            if (m.MessageReference == nil) {
                s.ChannelMessageSendReply(m.ChannelID, "Please reply to a message with a maze", m.Reference())
            } else {
                reply, err := s.ChannelMessage(m.MessageReference.ChannelID, m.MessageReference.MessageID)
                if (err != nil) {
                    s.ChannelMessageSendReply(m.ChannelID, "Please reply to a message with a maze", m.Reference())
                } else {
                    messageUrl := reply.Content
                    if strings.Contains(messageUrl, "media.discordapp.net") {
                        messageUrl = strings.Replace(messageUrl, "media.discordapp.net", "cdn.discordapp.com", 1)

                        re := regexp.MustCompile(`(hm=[a-f0-9]+).*`)
                        messageUrl = re.ReplaceAllString(messageUrl, `$1&`)
                    }
                    if (len(reply.Attachments) > 1) {
                        s.ChannelMessageSendReply(m.ChannelID, "Too many images!", m.Reference())
                    } else {
                        var maze getMaze.Maze
                        if len(reply.Attachments) < 1 {
                            maze, err = getMaze.GetMaze(messageUrl)
                        } else {
                            maze, err = getMaze.GetMaze(reply.Attachments[0].URL)
                        }
                        if (err != nil) {
                            s.ChannelMessageSendReply(m.ChannelID, "Send a valid maze", m.Reference())
                        } else {
                            points := solvemaze.FindPath(maze)
                            _, err := outputmaze.EditMaze(points, "/tmp/maze.png", "/tmp/outputmaze.png")
                            if err != nil {
                                s.ChannelMessageSendReply(m.ChannelID, "server error", m.Reference())
                            }
                            solutionlength := strconv.Itoa(len(points))
                            maze = nil
                            points = nil
                            f, _ := os.Open("/tmp/outputmaze.png")
                            _, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
                                Content: "Length of solution is: " + solutionlength,
                                Files: []*discordgo.File{
                                    {
                                        Name: "/tmp/outputmaze.png",
                                        Reader: f,
                                    },
                                },
                                Reference: &discordgo.MessageReference{
                                    MessageID: m.Reference().MessageID,
                                    ChannelID: m.Reference().ChannelID,
                                    GuildID: m.Reference().GuildID,
                                },
                            })
                            if err != nil {
                                if strings.Contains(err.Error(), "413") {
                                    s.ChannelMessageSendReply(m.ChannelID, "Solution file is too large!", m.Reference())
                                } else {
                                    s.ChannelMessageSendReply(m.ChannelID, err.Error(), m.Reference())
                                }
                            }
                        }
                    }
                }
            }
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

func parseMessageLink(link string) (string, string, string, error) {
	parts := strings.Split(link, "/")
	if len(parts) < 7 {
		return "", "", "", fmt.Errorf("invalid message link format")
	}
	guildID := parts[4]
	channelID := parts[5]
	messageID := parts[6]
	return guildID, channelID, messageID, nil
}
