package remindat_test

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/enrico5b1b4/tbwrap"
	"github.com/enrico5b1b4/telegram-bot/reminder"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/mocks"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate/remindat"
	fakeBot "github.com/enrico5b1b4/telegram-bot/telegram/fakes"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	tb "gopkg.in/tucnak/telebot.v2"
)

func TestHandleRemindAt(t *testing.T) {
	handlerPattern, err := regexp.Compile(remindat.HandlePattern)
	require.NoError(t, err)
	chat := &tb.Chat{ID: int64(1)}

	t.Run("success", func(t *testing.T) {
		type TestCase struct {
			Text                 string
			ExpectedWordDateTime reminder.WordDateTime
		}

		testCases := map[string]TestCase{
			"with hours and minutes": {
				Text: "/remind me at 20:45 buy milk",
				ExpectedWordDateTime: reminder.WordDateTime{
					When:   reminder.Today,
					Hour:   20,
					Minute: 45,
				},
			},
			"with hours and minutes dot separator": {
				Text: "/remind me at 20.45 buy milk",
				ExpectedWordDateTime: reminder.WordDateTime{
					When:   reminder.Today,
					Hour:   20,
					Minute: 45,
				},
			},
		}

		for name := range testCases {
			t.Run(name, func(t *testing.T) {
				mockCtrl := gomock.NewController(t)
				defer mockCtrl.Finish()
				bot := fakeBot.NewBot()
				c := tbwrap.NewContext(bot, &tb.Message{Text: testCases[name].Text, Chat: chat}, nil, handlerPattern)
				mockReminderService := mocks.NewMockServicer(mockCtrl)
				mockReminderService.
					EXPECT().
					AddReminderOnWordDateTime(
						1,
						testCases[name].Text,
						testCases[name].ExpectedWordDateTime,
						"buy milk").
					Return(time.Now(), nil)

				err := remindat.HandleRemindAt(mockReminderService)(c)
				require.NoError(t, err)
				require.Len(t, bot.OutboundSendMessages, 1)
			})
		}
	})

	t.Run("failure", func(t *testing.T) {
		text := "/remind me at 20.45 buy milk"
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
					Minute: 45,
				},
				"buy milk").
			Return(time.Now(), errors.New("error"))

		err := remindat.HandleRemindAt(mockReminderService)(c)
		require.Error(t, err)
		require.Len(t, bot.OutboundSendMessages, 0)
	})
}
