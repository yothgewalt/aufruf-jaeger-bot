package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/Yothgewalt/aufruf-jaeger-bot/command/utility"
	"github.com/Yothgewalt/aufruf-jaeger-bot/config"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

var logger = config.NewZeroLog()

var (
	GuildID        = flag.String("g", "", "Tes")
	BotAccessToken = flag.String("t", "", "Access token for bot management")
	RemoveCommands = flag.Bool("rc", true, "Remove all command after shutdowning or not")
)

var (
	classroomService = utility.NewClassroomService(logger)

	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "greet",
			Description: "Bot will greeting you.",
		},
		{
			Name:        "classroom",
			Description: "Showing all option that can interact with google classroom.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "courses",
					Description: "Showing all command that interact with course.",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "all",
							Description: "Get all about course (course name and course id)",
							Type:        discordgo.ApplicationCommandOptionSubCommand,
						},
					},
					Type: discordgo.ApplicationCommandOptionSubCommandGroup,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"greet": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			greetingMessage := "What's up, Do you need any help?"

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: greetingMessage,
				},
			})

			logger.Info().Str("send", i.Member.User.ID).Msg(greetingMessage)
		},
		"classroom": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options
			content := ""

			switch options[0].Name {
			case "courses":
				options = options[0].Options

				switch options[0].Name {
				case "all":
					defaultPageSize := 10
					c, err := classroomService.Courses.List().PageSize(int64(defaultPageSize)).Do()
					logger.Info().Str("PageSize", strconv.Itoa(defaultPageSize)).Msg("all your course name has been sent to your discord server.")
					if err != nil {
						logger.Error().Err(err).Msg("error cannot request your course.")
						content = "Cannot request the courses cause our service conflict."
					}

					if len(c.Courses) > 0 {
						for index, course := range c.Courses {
							content += fmt.Sprintf("\n%v) %v (ID: %v)", index+1, course.Name, course.Id)
						}
					} else {
						logger.Error().Err(err).Msg("No courses found.")
						content = "No courses found."
					}
				default:
					content = "Oops, something went wrong.\n" + "Hol'up, you aren't supposed to see this message."
				}
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
				},
			})
		},
	}
)

func init() {
	flag.Parse()
}

func main() {
	var err error

	dg, err := discordgo.New("Bot " + *BotAccessToken)
	if err != nil {
		log.Err(err).Msg("error creating discord sesssion.")
		return
	}

	err = dg.Open()
	if err != nil {
		logger.Err(err).Msg("error opening connection.")
		return
	}

	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for index, value := range commands {
		cmd, err := dg.ApplicationCommandCreate(dg.State.User.ID, *GuildID, value)
		if err != nil {
			logger.Panic().Err(err).Msgf("Cannot create command: %v", value.Name)
		}
		registeredCommands[index] = cmd
	}

	log.Info().Msg("Aufrufjaeger bot is now running. Press CTRL-C to exit.")
	signalStop := make(chan os.Signal, 1)
	signal.Notify(signalStop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-signalStop

	if *RemoveCommands {
		logger.Info().Msg("Removing commands")

		for _, value := range registeredCommands {
			err := dg.ApplicationCommandDelete(dg.State.User.ID, "", value.ID)
			if err != nil {
				logger.Panic().Err(err).Msgf("Cannot delete command: %v", value.Name)
			}
		}
	}

	dg.Close()

	logger.Info().Msg("Gracefully shutting down")
}
