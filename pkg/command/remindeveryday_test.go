package command_test

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/enrico5b1b4/tbwrap"
	"github.com/enrico5b1b4/telegram-bot/pkg/command"
	"github.com/enrico5b1b4/telegram-bot/pkg/reminder"
	"github.com/enrico5b1b4/telegram-bot/pkg/reminder/mocks"
	fakeBot "github.com/enrico5b1b4/telegram-bot/pkg/telegram/fakes"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	tb "gopkg.in/tucnak/telebot.v2"
)

type TestCaseHandleRemindEveryDay struct {
	Text                       string
	ExpectedRepeatableDateTime *reminder.RepeatableDateTime
}

func TestHandleRemindEveryDay_Success(t *testing.T) {
	handlerPattern, err := regexp.Compile(command.HandlePatternRemindEveryDay)
	require.NoError(t, err)
	chat := &tb.Chat{ID: int64(1)}
	testCases := newTestHandleRemindEveryDayTestCases()

	for name := range testCases {
		t.Run(name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			bot := fakeBot.NewTBWrapBot()
			c := tbwrap.NewContext(bot, &tb.Message{Text: testCases[name].Text, Chat: chat}, nil, handlerPattern)
			mockReminderService := mocks.NewMockServicer(mockCtrl)
			mockReminderService.
				EXPECT().
				AddRepeatableReminderOnDateTime(
					1,
					testCases[name].Text,
					testCases[name].ExpectedRepeatableDateTime,
					"update weekly report").
				Return(reminder.NextScheduleChatTime{Time: time.Now(), Location: time.UTC}, nil)

			err := command.HandleRemindEveryDay(mockReminderService)(c)
			require.NoError(t, err)
			require.Len(t, bot.OutboundSendMessages, 1)
		})
	}
}

func newTestHandleRemindEveryDayTestCases() map[string]TestCaseHandleRemindEveryDay {
	return map[string]TestCaseHandleRemindEveryDay{
		"without hours and minutes": {
			Text: "/remind me every day update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				Hour:   "9",
				Minute: "0",
			},
		},
		"with hours and minutes": {
			Text: "/remind me every day at 23:34 update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				Hour:   "23",
				Minute: "34",
			},
		},
		"with hours and minutes dot separator": {
			Text: "/remind me every day at 23.34 update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				Hour:   "23",
				Minute: "34",
			},
		},
		"with only hour": {
			Text: "/remind me every day at 23 update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				Hour:   "23",
				Minute: "0",
			},
		},
		"with only hour pm": {
			Text: "/remind me every day at 8pm update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				Hour:   "20",
				Minute: "0",
			},
		},
		"with hour minute pm": {
			Text: "/remind me every day at 8:30pm update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				Hour:   "20",
				Minute: "30",
			},
		},
		"with hour minute pm dot separator": {
			Text: "/remind me every day at 8.30pm update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				Hour:   "20",
				Minute: "30",
			},
		},
		"with time of day and without hours and minutes": {
			Text: "/remind me every morning update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				Hour:   "9",
				Minute: "0",
			},
		},
		"with time of day and hours and minutes": {
			Text: "/remind me every evening at 23:34 update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				Hour:   "23",
				Minute: "34",
			},
		},
		"with time of day and hours and minutes dot separator": {
			Text: "/remind me every night at 23.34 update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				Hour:   "23",
				Minute: "34",
			},
		},
		"with time of day and only hour": {
			Text: "/remind me every evening at 23 update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				Hour:   "23",
				Minute: "0",
			},
		},
		"with time of day and only hour pm": {
			Text: "/remind me every night at 8pm update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				Hour:   "20",
				Minute: "0",
			},
		},
		"with time of day and hour minute pm": {
			Text: "/remind me every evening at 8:30pm update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				Hour:   "20",
				Minute: "30",
			},
		},
		"with time of day and hour minute pm dot separator": {
			Text: "/remind me every night at 8.30pm update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				Hour:   "20",
				Minute: "30",
			},
		},
		"with time of day only": {
			Text: "/remind me every evening update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				Hour:   "20",
				Minute: "0",
			},
		},
		"with weekday/weekend and without hours and minutes": {
			Text: "/remind me every weekday update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				DayOfWeek: "1-5",
				Hour:      "9",
				Minute:    "0",
			},
		},
		"with weekday/weekend and hours and minutes": {
			Text: "/remind me every weekend at 23:34 update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				DayOfWeek: "6,0",
				Hour:      "23",
				Minute:    "34",
			},
		},
		"with weekday/weekend and hours and minutes dot separator": {
			Text: "/remind me every weekday at 23.34 update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				DayOfWeek: "1-5",
				Hour:      "23",
				Minute:    "34",
			},
		},
		"with weekday/weekend and only hour": {
			Text: "/remind me every weekend at 23 update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				DayOfWeek: "6,0",
				Hour:      "23",
				Minute:    "0",
			},
		},
		"with weekday/weekend and only hour pm": {
			Text: "/remind me every weekday at 8pm update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				DayOfWeek: "1-5",
				Hour:      "20",
				Minute:    "0",
			},
		},
		"with weekday/weekend and hour minute pm": {
			Text: "/remind me every weekend at 8:30pm update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				DayOfWeek: "6,0",
				Hour:      "20",
				Minute:    "30",
			},
		},
		"with weekday/weekend and hour minute pm dot separator": {
			Text: "/remind me every weekday at 8.30pm update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				DayOfWeek: "1-5",
				Hour:      "20",
				Minute:    "30",
			},
		},
		"with weekday/weekend only": {
			Text: "/remind me every weekend update weekly report",
			ExpectedRepeatableDateTime: &reminder.RepeatableDateTime{
				DayOfWeek: "6,0",
				Hour:      "9",
				Minute:    "0",
			},
		},
	}
}

func TestHandleRemindEveryDay_Failure(t *testing.T) {
	handlerPattern, err := regexp.Compile(command.HandlePatternRemindEveryDay)
	require.NoError(t, err)

	chat := &tb.Chat{ID: int64(1)}
	text := "/remind me every day update weekly report"
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	bot := fakeBot.NewTBWrapBot()
	c := tbwrap.NewContext(bot, &tb.Message{Text: text, Chat: chat}, nil, handlerPattern)
	mockReminderService := mocks.NewMockServicer(mockCtrl)
	mockReminderService.
		EXPECT().
		AddRepeatableReminderOnDateTime(
			1,
			text,
			&reminder.RepeatableDateTime{
				Hour:   "9",
				Minute: "0",
			},
			"update weekly report").
		Return(reminder.NextScheduleChatTime{}, errors.New("error"))

	err = command.HandleRemindEveryDay(mockReminderService)(c)
	require.Error(t, err)
	require.Len(t, bot.OutboundSendMessages, 0)
}
