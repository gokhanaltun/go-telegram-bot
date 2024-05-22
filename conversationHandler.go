package bot

import "github.com/gokhanaltun/go-telegram-bot/models"

type ConversationStage map[int]HandlerFunc

type ConversationEnd struct {
	Command  string
	Function HandlerFunc
}

type ConversationHandler struct {
	active        bool
	activeStageId int
	stages        ConversationStage
	end           *ConversationEnd
}

func NewConversationHandler(stages_ ConversationStage, end_ *ConversationEnd) *ConversationHandler {
	return &ConversationHandler{
		stages: stages_,
		end:    end_,
	}

}

func (b *Bot) AddConversationHandler(conversationHandler_ *ConversationHandler) {
	b.conversationHandler = conversationHandler_
}

func (b *Bot) SetActiveConversationStage(stageId_ int) {
	if !b.conversationHandler.active {
		b.conversationHandler.active = true
	}
	b.conversationHandler.activeStageId = stageId_
}

func (c *ConversationHandler) getStageFunction(upd *models.Update) HandlerFunc {
	if c.active {

		if upd.Message.Text == c.end.Command {
			c.active = false
			return c.end.Function
		}

		if hf, ok := c.stages[c.activeStageId]; ok {
			return hf
		}
	}
	return nil
}

func (b *Bot) EndConversation() {
	b.conversationHandler.active = false
}
