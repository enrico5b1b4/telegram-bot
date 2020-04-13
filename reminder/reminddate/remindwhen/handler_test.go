package remindwhen_test

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/enrico5b1b4/tbwrap"
	"github.com/enrico5b1b4/telegram-bot/reminder"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/mocks"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/remindwhen"
	fakeBot "github.com/enrico5b1b4/telegram-bot/telegram/fakes"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	tb "gopkg.in/tucnak/telebot.v2"
)

func TestHandleRemindWhen(t *testing.T) {
	handlerPattern, err := regexp.Compile(remindwhen.HandlePattern)
	require.NoError(t, err)
	chat := &tb.Chat{ID: int64(1)}

	t.Run("success without hours and minutes", func(t *testing.T) {
		text := "/remind me tonight buy milk"
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		bot := fakeBot.NewBot()
		c := tbwrap.NewContext(bot, &tb.Message{Text: text, Chat: chat}, nil, handlerPattern)
		mockReminderService := mocks.NewMockServicer(mockCtrl)
		mockReminderService.
			EXPECT().
			AddReminderOnWordDateTime(
				1,
				text,
				reminder.WordDateTime{
					When:   reminder.Today,
					Hour:   20,
					Minute: 0,
				},
				"buy milk").
			Return(time.Now(), nil)

		err := remindwhen.HandleRemindWhen(mockReminderService)(c)
		require.NoError(t, err)
		require.Len(t, bot.OutboundSendMessages, 1)
	})

	t.Run("success with hours and minutes", func(t *testing.T) {
		textHoursMins := "/remind me tonight at 23:34 buy milk"
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		bot := fakeBot.NewBot()
		c := tbwrap.NewContext(bot, &tb.Message{Text: textHoursMins, Chat: chat}, nil, handlerPattern)
		mockReminderService := mocks.NewMockServicer(mockCtrl)
		mockReminderService.
			EXPECT().
			AddReminderOnWordDateTime(
				1,
				textHoursMins,
				reminder.WordDateTime{
					When:   reminder.Today,
					Hour:   23,
					Minute: 34,
				},
				"buy milk").
			Return(time.Now(), nil)

		err := remindwhen.HandleRemindWhen(mockReminderService)(c)
		require.NoError(t, err)
		require.Len(t, bot.OutboundSendMessages, 1)
	})

	t.Run("failure", func(t *testing.T) {
		text := "/remind me tonight buy milk"
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		bot := fakeBot.NewBot()
		c := tbwrap.NewContext(bot, &tb.Message{Text: text, Chat: chat}, nil, handlerPattern)
		mockReminderService := mocks.NewMockServicer(mockCtrl)
		mockReminderService.
			EXPECT().
			AddReminderOnWordDateTime(
				1,
				text,
				reminder.WordDateTime{
					When:   reminder.Today,
					Hour:   20,
					Minute: 0,
				},
				"buy milk").
			Return(time.Now(), errors.New("error"))

		err := remindwhen.HandleRemindWhen(mockReminderService)(c)
		require.Error(t, err)
		require.Len(t, bot.OutboundSendMessages, 0)
	})
}