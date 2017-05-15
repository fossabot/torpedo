package multibot

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/nlopes/slack"
)

func (tb *TorpedoBot) RunSlackBot(apiKey, cmd_prefix string) {
	api := slack.New(apiKey)
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	botApi := &TorpedoBotAPI{}
	botApi.API = api
	botApi.Bot = tb
	botApi.CommandPrefix = cmd_prefix

	for msg := range rtm.IncomingEvents {
		logger.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			logger.Println("Infos:", ev.Info)
			logger.Println("Connection counter:", ev.ConnectionCount)
			// Replace #general with your Channel ID
			// rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "#general"))

		case *slack.MessageEvent:
			logger.Printf("Message: %v\n", ev)
			channel := ev.Channel
			incoming_message := ev.Text
			messageTS, _ := strconv.ParseFloat(ev.Timestamp, 64)
			jitter := int64(time.Now().Unix()) - int64(messageTS)
			if jitter < 10 {
				go tb.processChannelEvent(botApi, channel, incoming_message)
			}

		case *slack.PresenceChangeEvent:
			logger.Printf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			logger.Printf("Current latency: %v\n", ev.Value)

		case *slack.RTMError:
			logger.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			logger.Printf("Invalid credentials")
			return

		default:
			// Ignore other events..
			//logger.Printf("Unexpected: %v\n", msg.Data)
		}
	}

}
