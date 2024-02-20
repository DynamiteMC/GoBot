package commands

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Knetic/govaluate"
	"github.com/aquilax/go-perlin"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/wcharczuk/go-chart/v2"
)

var Command_graph = Command{
	Name:        "graph",
	Description: "Graph",
	Aliases:     []string{"g"},
	Execute: func(message *events.MessageCreate, args []string) {
		l := 50
		equation := "x"

		if a := strings.Join(args, " "); a != "" {
			equation = a
		}

		expression, err := govaluate.NewEvaluableExpressionWithFunctions(equation, map[string]govaluate.ExpressionFunction{
			"p1": func(arguments ...interface{}) (interface{}, error) {
				if len(arguments) != 5 {
					return 0, fmt.Errorf("must have 5 arguments: alpha, beta, n, seed, x")
				}
				p := perlin.NewPerlin(arguments[0].(float64), arguments[1].(float64), int32(arguments[2].(float64)), int64(arguments[3].(float64)))

				return p.Noise1D(arguments[4].(float64)), nil
			},
			"p2": func(arguments ...interface{}) (interface{}, error) {
				if len(arguments) != 6 {
					return 0, fmt.Errorf("must have 6 arguments: alpha, beta, n, seed, x, y")
				}
				p := perlin.NewPerlin(arguments[0].(float64), arguments[1].(float64), int32(arguments[2].(float64)), int64(arguments[3].(float64)))

				return p.Noise2D(arguments[4].(float64), arguments[5].(float64)), nil
			},
			"p3": func(arguments ...interface{}) (interface{}, error) {
				if len(arguments) != 7 {
					return 0, fmt.Errorf("must have 7 arguments: alpha, beta, n, seed, x, y, z")
				}
				p := perlin.NewPerlin(arguments[0].(float64), arguments[1].(float64), int32(arguments[2].(float64)), int64(arguments[3].(float64)))

				return p.Noise3D(arguments[4].(float64), arguments[5].(float64), arguments[6].(float64)), nil
			},
		})

		if err != nil {
			CreateMessage(message, fmt.Sprintf("Invalid expression: %s", err), true)
			return
		}

		var buffer bytes.Buffer
		var xs []float64
		var ys []float64

		for y := l; y >= -l; y-- {
			for x := -l; x <= l; x++ {
				res, err := expression.Evaluate(map[string]interface{}{"x": x})
				if err != nil {
					CreateMessage(message, fmt.Sprintf("Evaluation error (x=%d): %s", x, err), true)
					return
				}
				xs = append(xs, float64(x))
				ys = append(ys, res.(float64))
			}
		}

		graph := chart.Chart{
			Series: []chart.Series{
				chart.ContinuousSeries{
					XValues: xs,
					YValues: ys,
				},
			},
		}
		graph.Render(chart.PNG, &buffer)

		message.Client().Rest().CreateMessage(message.ChannelID, discord.NewMessageCreateBuilder().
			SetAllowedMentions(&discord.AllowedMentions{RepliedUser: false}).
			SetMessageReferenceByID(message.MessageID).
			AddFile("graph.png", "Line graph", &buffer).
			SetContentf("# y=%s", equation).
			Build())
	},
}
