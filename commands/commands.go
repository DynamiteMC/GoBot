package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"gobot/config"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ayush6624/go-chatgpt"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

var ac = map[string]string{
	"bios": "Basic Input/Output System",
	"mcpe": "Minecraft Pocket Edition",
	"ram":  "Random Access Memory",
	"rom":  "Read Only Memory",
	"ssd":  "Solid State Drive",
	"hdd":  "Hard Disk Drive",
	"nyc":  "New York City",
	"aws":  "Amazon Web Services",
	"ai":   "Artificial Intelligence",
	"wifi": "Wireless Fidelity",
	"pc":   "Personal Computer",
	"gn":   "Good Night",
	"gm":   "Good Morning",
	"dp":   "Display Port",
	"hdmi": "High Definition Multimedia Interface",
	"pdf":  "Personal Document Format",
	"sata": "Serial Advanced Technology Attachment",
	"pci":  "Peripheral Component Interconnect",
	"pcie": "Peripheral Component Interconnect Express",
	"nvme": "Non-Volatile Memory Express",
	"tcp":  "Transmission Control Protocol",
	"udp":  "User Datagram Protocol",
	"www":  "World Wide Web",
	"http": "Hypertext Transfer Protocol",
	"js":   "JavaScript",
	"ts":   "TypeScript",
	"py":   "Python",
	"lol":  "Laugh Out Loud",
	"lmao": "Laughing My Ass Off",
	"btw":  "By The Way",
	"tbh":  "To Be Honest"
}

var color = 0x9C182C
var startTime = time.Now()
var True = true

func Point[T any](data T) *T {
	return &data
}

func HasAnyPrefix(str string, prefixes ...string) (bool, string) {
	for _, s := range prefixes {
		if strings.HasPrefix(str, s) {
			return true, s
		}
	}
	return false, ""
}

func IsAny(str string, strs ...string) (bool, string) {
	for _, s := range strs {
		if str == s {
			return true, s
		}
	}
	return false, ""
}

func GetArgument(args []string, index int) string {
	if len(args) <= index {
		return ""
	}
	return args[index]
}

type Command struct {
	Name        string
	Description string
	Execute     func(*events.MessageCreate, []string)
	Aliases     []string
	Permissions discord.Permissions
}

type Message struct {
	Content string
	Reply   bool
	Embeds  []discord.Embed
	Files   []*discord.File
}

func EditMessage(client bot.Client, channelID snowflake.ID, id snowflake.ID, message Message) (*discord.Message, error) {
	builder := discord.NewMessageUpdateBuilder().
		SetContent(message.Content).
		SetEmbeds(message.Embeds...).
		SetFiles(message.Files...).
		SetAllowedMentions(&discord.AllowedMentions{RepliedUser: false})
	return client.Rest().UpdateMessage(channelID, id, builder.Build())
}

func CreateMessage(e *events.MessageCreate, msg ...interface{}) (*discord.Message, error) {
	er := errors.New("invalid arguments")
	if len(msg) == 0 {
		return nil, er
	}
	switch message := msg[0].(type) {
	case Message:
		{
			builder := discord.NewMessageCreateBuilder().
				SetContent(message.Content).
				SetEmbeds(message.Embeds...).
				SetFiles(message.Files...).
				SetAllowedMentions(&discord.AllowedMentions{RepliedUser: false})
			if message.Reply {
				builder.SetMessageReferenceByID(e.MessageID)
			}
			return e.Client().Rest().CreateMessage(e.ChannelID, builder.Build())
		}
	case string:
		{
			builder := discord.NewMessageCreateBuilder().
				SetContent(message).
				SetAllowedMentions(&discord.AllowedMentions{RepliedUser: false})
			if len(msg) == 2 && msg[1] == true {
				builder.SetMessageReferenceByID(e.MessageID)
			}
			return e.Client().Rest().CreateMessage(e.ChannelID, builder.Build())
		}
	case discord.Embed:
		{
			{
				builder := discord.NewMessageCreateBuilder().
					SetEmbeds(message).
					SetAllowedMentions(&discord.AllowedMentions{RepliedUser: false})
				if len(msg) == 2 && msg[1] == true {
					builder.SetMessageReferenceByID(e.MessageID)
				}
				return e.Client().Rest().CreateMessage(e.ChannelID, builder.Build())
			}
		}
	default:
		return nil, er
	}
}

func HasRole(client bot.Client, guildId snowflake.ID, memberId snowflake.ID, id int64) bool {
	member, _ := client.Rest().GetMember(guildId, memberId)
	for _, role := range member.RoleIDs {
		if role == snowflake.ID(id) {
			return true
		}
	}
	return false
}

func ParseMention(mention string) snowflake.ID {
	if strings.HasPrefix(mention, "<@") && strings.HasSuffix(mention, ">") && !strings.HasPrefix(mention, "<@&") {
		mention = strings.TrimPrefix(strings.TrimSuffix(mention, ">"), "<@")
	}
	id, err := strconv.ParseInt(mention, 10, 64)
	if err != nil {
		return 0
	}
	return snowflake.ID(id)
}

var openaiClient *chatgpt.Client

func Handle(aic *chatgpt.Client, message *events.MessageCreate) {
	openaiClient = aic
	if message.Message.Author.Bot {
		return
	}
	if strings.Contains(strings.ToLower(message.Message.Content), "mmm") {
		message.Client().Rest().AddReaction(message.ChannelID, message.MessageID, "✅")
	}

	if message.Message.Content == "no!" {
		ref := message.Message.ReferencedMessage
		if ref != nil {
			if ref.Author.ID == message.Client().ID() {
				ref1 := ref.ReferencedMessage
				if ref1 != nil && ref1.Author.ID == message.Message.Author.ID {
					message.Client().Rest().DeleteMessage(message.ChannelID, message.MessageID)
					message.Client().Rest().DeleteMessage(message.ChannelID, ref.ID)
				}
			}
		}
	}

	u, err := url.ParseRequestURI(message.Message.Content)

	if err == nil {
		var (
			repo      string
			file      string
			formatter string
			lineStart int
			lineEnd   int
			str       string
			code      string
		)
		if u.Host == "github.com" {
			sp := strings.Split(u.Path, "/")
			repo = strings.Join(sp[1:3], "/")
			if repo != "" {
				sp := strings.Split(u.Fragment, "-")
				if len(sp) > 0 {
					if strings.HasPrefix(sp[0], "L") {
						l := strings.TrimPrefix(sp[0], "L")
						lineStart, _ = strconv.Atoi(l)
					}
				}
				if len(sp) > 1 {
					if strings.HasPrefix(sp[1], "L") {
						l := strings.TrimPrefix(sp[1], "L")
						lineEnd, _ = strconv.Atoi(l)
					}
				}
				res, _ := http.Get(u.String())
				if res.StatusCode == 200 {
					q, _ := goquery.NewDocumentFromReader(res.Body)
					s := q.Find("react-app").Children().First()
					c, _ := s.Html()
					var data map[string]interface{}
					err := json.Unmarshal([]byte(strings.ReplaceAll(c, "&#34;", `"`)), &data)
					if err == nil {
						pl := data["payload"].(map[string]interface{})["blob"].(map[string]interface{})
						file = pl["displayName"].(string)
						fsp := strings.Split(file, ".")
						formatter = fsp[len(fsp)-1]
						var lines []interface{}
						if lineEnd != 0 {
							lines = pl["rawLines"].([]interface{})[lineStart-1 : lineEnd]
						} else {
							lines = pl["rawLines"].([]interface{})[lineStart-1 : lineStart]
						}
						for i, l := range lines {
							code += fmt.Sprint(l)
							if i != len(lines)-1 {
								code += "\n"
							}
						}
					}
				}
				if lineStart != 0 && code != "" {
					str = fmt.Sprintf("**%s %s**\nLine **%d**:\n```%s\n%s```", repo, file, lineStart, formatter, code)
					if lineEnd != 0 {
						str = fmt.Sprintf("**%s %s**\nLines **%d** - **%d**:\n```%s\n%s```", repo, file, lineStart, lineEnd, formatter, code)
					}
					CreateMessage(message, str, true)
				}
			}
		}
	}

	var str []string
	for _, c := range strings.Split(message.Message.Content, " ") {
		if a, ok := ac[strings.ToLower(c)]; ok {
			for _, s := range str {
				if s == a {
					continue
				}
			}
			str = append(str, a)
		}
	}

	if len(str) != 0 {
		CreateMessage(message, fmt.Sprintf("(%s)", strings.Join(str, ", ")), true)
	}

	args := strings.Split(message.Message.Content, " ")
	if !strings.HasPrefix(message.Message.Content, config.Config.InfoPrefix) {
		if strings.HasPrefix(message.Message.Content, config.Config.Prefix) {
			cmd := strings.ToLower(strings.Split(args[0], "\n")[0][len(config.Config.Prefix):])
			if len(strings.Split(args[0], "\n")) == 1 {
				args = args[1:]
			} else {
				args[0] = strings.Join(strings.Split(args[0], "\n")[1:], "\n")
			}
			command := commands[cmd]
			if command.Execute == nil {
				command = commands[aliases[cmd]]
				if command.Execute == nil {
					return
				}
			}
			message.Message.Member.GuildID = *message.GuildID
			if !message.Client().Caches().MemberPermissions(*message.Message.Member).Has(command.Permissions) {
				return
			}
			command.Execute(message, args)
		} else {
			return
		}
	} else {
		cmd := args[0][len(config.Config.InfoPrefix):]
		command := commands[cmd]
		if command.Execute == nil {
			command = commands[aliases[cmd]]
			if command.Execute == nil {
				return
			}
		}
		aliases := "None"
		if len(command.Aliases) > 0 {
			aliases = strings.Join(command.Aliases, ", ")
		}
		embed := discord.NewEmbedBuilder().
			SetTitle(command.Name).
			SetDescription(command.Description).
			AddFields(
				discord.EmbedField{
					Name:   "Aliases",
					Value:  aliases,
					Inline: &True,
				},
				discord.EmbedField{
					Name:   "Permissions Required",
					Value:  command.Permissions.String(),
					Inline: &True,
				},
			).SetColor(color).Build()
		CreateMessage(message, Message{Embeds: []discord.Embed{embed}, Reply: true})
	}
}

var commands = make(map[string]Command)
var aliases = make(map[string]string)

func RegisterCommands(cmds ...Command) {
	for _, command := range cmds {
		commands[command.Name] = command
		for _, alias := range command.Aliases {
			aliases[alias] = command.Name
		}
	}
}
