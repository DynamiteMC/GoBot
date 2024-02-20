package commands

import (
	"fmt"
	"gobot/store"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var Command_unwarn = Command{
	Name:        "unwarn",
	Description: "Remove one warning from a member",
	Permissions: discord.PermissionModerateMembers,
	Aliases:     []string{"unwarn", "cool", "unwarm", "cold"},
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
		store.Unwarn(int64(id))
		warnings := store.Warnings(int64(id))

		var tag string
		member, err := message.Client().Rest().GetMember(*message.GuildID, id)
		if err != nil {
			tag = "Unknown#0000"
		} else {
			tag = member.User.Tag()
		}
		CreateMessage(message, Message{Content: fmt.Sprintf("Removed one warning from member %s (%d/3)", tag, warnings), Reply: true})
	},
}
