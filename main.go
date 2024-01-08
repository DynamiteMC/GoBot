package main

import (
	"context"
	"fmt"

	"gobot/commands"
	"gobot/config"
	"gobot/store"
	"gobot/web"
	"os"
	"os/signal"
	"syscall"

	"github.com/ayush6624/go-chatgpt"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
)

func main() {
	if !config.Config.Load() {
		fmt.Println("Please fill the config.json file then run the program again")
		os.Exit(1)
	}
	c, _ := chatgpt.NewClient(config.Config.OpenAIKey)
	client, err := disgo.New(config.Config.Token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuilds,
				gateway.IntentGuildMessages,
				gateway.IntentMessageContent,
			),
			gateway.WithBrowser("Discord iOS"),
			gateway.WithPresenceOpts(gateway.WithCustomActivity("I'm a human")),
		),
		bot.WithCacheConfigOpts(
			cache.WithCaches(cache.FlagGuilds|cache.FlagRoles),
		),
		bot.WithEventListenerFunc(func(m *events.MessageCreate) {
			commands.Handle(c, m)
		}),
		bot.WithEventListenerFunc(func(event *events.GuildMemberJoin) {
			if c, _ := store.GetCorner(int64(event.Member.User.ID)); c {
				event.Client().Rest().AddMemberRole(event.GuildID, event.Member.User.ID, snowflake.ID(config.Config.DisgraceRole))
			}
			if store.IsMuted(int64(event.Member.User.ID)) {
				event.Client().Rest().AddMemberRole(event.GuildID, event.Member.User.ID, snowflake.ID(config.Config.MuteRole))
			}
		}),
		bot.WithEventListenerFunc(func(e *events.Ready) {
			if config.Config.DynamiteServer != "" {
				go web.Connect(e.Client(), config.Config)
			}
			commands.RegisterCommands(
				commands.Command_mute,
				commands.Command_oq,
				commands.Command_unmute,
				commands.Command_kick,
				commands.Command_ban,
				commands.Command_corner,
				commands.Command_screenshot,
				commands.Command_clean,
				commands.Command_uncorner,
				commands.Command_mcdoc,
				commands.Command_help,
				commands.Command_info,
				commands.Command_go,
				commands.Command_mcplayer,
				commands.Command_mcping,
				commands.Command_fetchdb,
				commands.Command_fancify,
				commands.Command_warn,
				commands.Command_unwarn,
				commands.Command_warnings,
			)
			fmt.Println("Bot is online.")
		}),
	)
	if err != nil {
		panic(err)
	}
	if err = client.OpenGateway(context.TODO()); err != nil {
		panic(err)
	}

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s
}
