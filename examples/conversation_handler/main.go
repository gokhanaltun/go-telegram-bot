package main

import (
	"context"
	"os"
	"os/signal"

	bot "github.com/gokhanaltun/go-telegram-bot"
	"github.com/gokhanaltun/go-telegram-bot/models"
)

const (
	firstNameStage = iota // Definition of the first name stage = 0
	lastNameStage         // Definition of the last name stage = 1
)

var firstName string // Global variable to store user's first name

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
	}

	b, err := bot.New("7081458788:AAFJJucpaSTMJv2mTKWBOBd80MiS-KbVvbY", opts...)
	if err != nil {
		panic(err)
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, start)

	stages := map[int]bot.HandlerFunc{
		firstNameStage: firstNameHandler,
		lastNameStage:  lastNameHandler,
	}

	convEnd := bot.ConversationEnd{
		Command:  "/cancel",
		Function: cancelConversation,
	}

	convHandler := bot.NewConversationHandler(stages, &convEnd)

	b.AddConversationHandler(convHandler)

	b.Start(ctx)
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "default handler.",
	})
}

// Handle /start command to start getting the user's name
func start(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SetActiveConversationStage(firstNameStage) //start the first name stage

	// Ask user to enter their name
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "enter your name or use the /cancel command to cancel.",
	})
}

// Handle the first name stage to get the user's first name
func firstNameHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	firstName = update.Message.Text

	b.SetActiveConversationStage(lastNameStage) //change stage to last name stage

	// Ask user to enter their last name
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "enter your last name",
	})
}

// Handle the last name stage to get the user's last name and send a hello message
func lastNameHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	lastName := update.Message.Text

	b.EndConversation() // end the conversation

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Hello, " + firstName + " " + lastName + " :)",
	})
}

// Handle /cancel command to end the conversation
func cancelConversation(ctx context.Context, b *bot.Bot, update *models.Update) {

	// Send a message to indicate the conversation has been cancelled
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "conversation cancelled",
	})
}
