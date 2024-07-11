package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"

	bot "github.com/gokhanaltun/go-telegram-bot"
	"github.com/gokhanaltun/go-telegram-bot/models"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
	}

	b, err := bot.New(os.Getenv("YOUR_TELEGRAM_BOT_TOKEN"), opts...)
	if nil != err {
		panic(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, startHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "", bot.MatchTypeWebAppData, webAppDataHandler)
	b.Start(ctx)
}

func startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {

	replyKeyboardMarkup := models.ReplyKeyboardMarkup{
		Keyboard: [][]models.KeyboardButton{
			{
				{Text: "Start WebApp", WebApp: &models.WebAppInfo{URL: "YOUR_WEB_APP_URL, See index.html"}},
			},
		},
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Hello, *" + bot.EscapeMarkdown(update.Message.From.FirstName) + "*",
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: replyKeyboardMarkup,
	})
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "default handler",
	})
}

func webAppDataHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	var result map[string]string

	err := json.Unmarshal([]byte(update.Message.WebAppData.Data), &result)
	if err != nil {
		fmt.Println(err)
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Data1: " + result["data1"] + " " + "Data2: " + result["data2"],
	})
}
