package commands

import (
	"fmt"
	"gobot/store"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var Command_warn = Command{
	Name:        "warn",
	Description: "Warn a member",
	Permissions: discord.PermissionModerateMembers,
	Aliases:     []string{"warm", "hot"},
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
		store.Warn(int64(id))
		warnings := store.Warnings(int64(id))

		var tag string
		member, err := message.Client().Rest().GetMember(*message.GuildID, id)
		if err != nil {
			tag = "Unknown#0000"
		} else {
			tag = member.User.Tag()
		}
		var ex string
		if warnings > 3 {
			store.SetWarnings(int64(id), 0)
			ex = ". They just got banned!!!"
			message.Client().Rest().AddBan(*message.GuildID, id, 0)
		}
		if warnings == 3 {
			ex = ". Next warning they get banned!!"
		}
		CreateMessage(message, Message{Content: fmt.Sprintf("Warned member %s (%d/3)%s", tag, warnings, ex), Reply: true})
	},
}
