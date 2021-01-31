package command_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/enrico5b1b4/tbwrap"
	"github.com/enrico5b1b4/telegram-bot/pkg/command"
	"github.com/enrico5b1b4/telegram-bot/pkg/command/mocks"
	fakeBot "github.com/enrico5b1b4/telegram-bot/pkg/telegram/fakes"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	tb "gopkg.in/tucnak/telebot.v2"
)

func TestHandleRemindDetail(t *testing.T) {
	handlerPattern, err := regexp.Compile(command.HandlePatternRemindDetail[0])
	require.NoError(t, err)
	text := "/reminddetail 2"
	chat := &tb.Chat{ID: int64(1)}

	t.Run("success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		bot := fakeBot.NewTBWrapBot()
		c := tbwrap.NewContext(bot, &tb.Message{Text: text, Chat: chat}, nil, handlerPattern)
		mockReminderService := mocks.NewMockRemindDetailServicer(mockCtrl)
		mockReminderService.
			EXPECT().
			GetReminder(1, 2).
			Return(&command.ReminderDetail{}, nil)

		err := command.HandleRemindDetail(mockReminderService, nil)(c)
		require.NoError(t, err)
		require.Len(t, bot.OutboundSendMessages, 1)
	})

	t.Run("failure", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		bot := fakeBot.NewTBWrapBot()
		c := tbwrap.NewContext(bot, &tb.Message{Text: text, Chat: chat}, nil, handlerPattern)
		mockReminderService := mocks.NewMockRemindDetailServicer(mockCtrl)
		mockReminderService.
			EXPECT().
			GetReminder(1, 2).
			Return(nil, errors.New("error"))

		err := command.HandleRemindDetail(mockReminderService, nil)(c)
		require.Error(t, err)
		require.Len(t, bot.OutboundSendMessages, 0)
	})
}
