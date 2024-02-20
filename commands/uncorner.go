package commands

import (
	"fmt"
	"gobot/config"
	"gobot/store"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

var Command_uncorner = Command{
	Name:        "uncorner",
	Description: "Send a member out of the corner",
	Permissions: discord.PermissionManageRoles,
	Aliases:     []string{"unshame", "unbully", "unblind", "goodboy", "unlmfao", "didask", "unbruh", "unratio", "uncope", "skill"},
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
		if !HasRole(message.Client(), *message.GuildID, id, config.Config.DisgraceRole) {
			CreateMessage(message, Message{Content: "Member is not cornered.", Reply: true})
			return
		}
		is, corner := store.GetCorner(int64(id))
		if is {
			for _, role := range corner["roles"].([]int64) {
				message.Client().Rest().AddMemberRole(*message.GuildID, id, snowflake.ID(role))
			}
		}
		store.RemoveCornered(int64(id))
		err := message.Client().Rest().RemoveMemberRole(*message.GuildID, id, snowflake.ID(config.Config.DisgraceRole))
		if err != nil {
			CreateMessage(message, Message{Content: "Failed to uncorner member", Reply: true})
			return
		}
		var tag string
		member, err := message.Client().Rest().GetMember(*message.GuildID, id)
		if err != nil {
			tag = "Unknown#0000"
		} else {
			tag = member.User.Tag()
		}
		CreateMessage(message, Message{Content: fmt.Sprintf("Sent %s back from the corner.", tag), Reply: true})
	},
}
