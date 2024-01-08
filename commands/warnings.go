package commands

import (
	"fmt"
	"gobot/store"

	"github.com/disgoorg/disgo/events"
)

var Command_warnings = Command{
	Name:        "warnings",
	Description: "Check warnings for a member",
	Execute: func(message *events.MessageCreate, args []string) {
		memberId := GetArgument(args, 0)
		if memberId == "" {
			if message.Message.ReferencedMessage != nil {
				memberId = message.Message.ReferencedMessage.Author.ID.String()
			} else {
				memberId = message.Message.Author.ID.String()
			}
		}
		id := ParseMention(memberId)
		if id == 0 {
			CreateMessage(message, Message{Content: "Failed to parse member", Reply: true})
			return
		}
		warnings := store.Warnings(int64(id))

		var tag string
		member, err := message.Client().Rest().GetMember(*message.GuildID, id)
		if err != nil {
			tag = "Unknown#0000"
		} else {
			tag = member.User.Tag()
		}
		CreateMessage(message, Message{Content: fmt.Sprintf("%s has %d warnings.", tag, warnings), Reply: true})
	},
}
