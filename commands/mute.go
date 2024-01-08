package commands

import (
	"fmt"
	"gobot/config"
	"gobot/store"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

var Command_mute = Command{
	Name:        "mute",
	Description: "Mute a member",
	Permissions: discord.PermissionMuteMembers,
	Aliases:     []string{"silence", "shush", "moot"},
	Execute: func(message *events.MessageCreate, args []string) {
		memberId := GetArgument(args, 0)
		if memberId == "" {
			if message.Message.ReferencedMessage != nil {
				memberId = message.Message.ReferencedMessage.Author.ID.String()
			}
		}
		id := ParseMention(memberId)
		if id == 0 {
			CreateMessage(message, Message{Content: "Failed to parse member", Reply: true})
			return
		}
		if HasRole(message.Client(), *message.GuildID, id, config.Config.MuteRole) {
			CreateMessage(message, Message{Content: "Member is already silenced.", Reply: true})
			return
		}
		err := message.Client().Rest().AddMemberRole(*message.GuildID, id, snowflake.ID(config.Config.MuteRole))
		if err != nil {
			CreateMessage(message, Message{Content: "Failed to mute member", Reply: true})
			return
		}
		store.AddMuted(int64(id))
		var tag string
		member, err := message.Client().Rest().GetMember(*message.GuildID, id)
		if err != nil {
			tag = "Unknown#0000"
		} else {
			tag = member.User.Tag()
		}
		CreateMessage(message, Message{Content: fmt.Sprintf("Silenced member %s.", tag), Reply: true})
	},
}
