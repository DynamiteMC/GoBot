package main

import (
	"context"
	"fmt"

	"gobot/commands"
	"gobot/config"
	"gobot/store"
	"os"
	"os/signal"
	"syscall"

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
	client, err := disgo.New(config.Config.Token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuilds,
				gateway.IntentGuildMessages,
				gateway.IntentMessageContent,
			),
		),
		bot.WithCacheConfigOpts(
			cache.WithCaches(cache.FlagGuilds|cache.FlagRoles),
		),
		bot.WithEventListenerFunc(commands.Handle),
		bot.WithEventListenerFunc(func(event *events.GuildMemberJoin) {
			if store.IsCornered(int64(event.Member.User.ID)) {
				event.Client().Rest().AddMemberRole(event.GuildID, event.Member.User.ID, snowflake.ID(config.Config.DisgraceRole))
			}
			if store.IsMuted(int64(event.Member.User.ID)) {
				event.Client().Rest().AddMemberRole(event.GuildID, event.Member.User.ID, snowflake.ID(config.Config.MuteRole))
			}
		}),
		bot.WithEventListenerFunc(func(*events.Ready) {
			commands.RegisterCommands(commands.Command_mute,
				commands.Command_oq,
				commands.Command_unmute,
				commands.Command_kick,
				commands.Command_ban,
				commands.Command_corner,
				commands.Command_screenshot,
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
