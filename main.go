package main

import (
	"fmt"
	"os"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func main() {
	botToken := os.Getenv("TOKEN")

	// Create Bot with debug on
	// Note: Please keep in mind that default logger may expose sensitive information, use in development only
	bot, err := telego.NewBot(botToken, telego.WithDefaultLogger(false, true))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Get bot user
	botUser, err := bot.GetMe()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Bot user: %+v\n", botUser)

	// Get updates channel
	updates, _ := bot.UpdatesViaLongPolling(nil)

	// Create bot handler
	bh, _ := th.NewBotHandler(bot, updates)

	// Register new handler with match on command `/start`
	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		// Send message
		_, _ = bot.SendMessage(tu.Message(
			tu.ID(update.Message.Chat.ID),
			fmt.Sprintf("Hello %s!", update.Message.From.FirstName),
		))
	}, th.CommandEqual("start"))

	// Register new handler with match on any command
	// Handlers will match only once and in order of registration,
	// so this handler will be called on any command except `/start` command
	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		// Send message
		_, _ = bot.SendMessage(tu.Message(
			tu.ID(update.Message.Chat.ID),
			"Unknown command, use /start",
		))
	}, th.AnyCommand())

	// Handle any message
	bh.HandleMessage(func(bot *telego.Bot, message telego.Message) {
		// Get chat ID from the message
		chatID := tu.ID(message.Chat.ID)

		// Copy sent messages back to the user
		_, _ = bot.CopyMessage(
			tu.CopyMessage(chatID, chatID, message.MessageID),
		)
	})

	// Stop handling updates on exit
	defer bh.Stop()
	defer bot.StopLongPolling()

	// Start handling updates
	bh.Start()
}
