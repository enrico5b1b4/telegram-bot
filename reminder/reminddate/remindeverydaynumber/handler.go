package remindeverydaynumber

import (
	"strconv"

	"github.com/enrico5b1b4/tbwrap"
	"github.com/enrico5b1b4/telegram-bot/reminder"
	"github.com/enrico5b1b4/telegram-bot/reminder/reminddate"
)

type Message struct {
	Who     string `regexpGroup:"who"`
	Day     int    `regexpGroup:"day"`
	Hour    *int   `regexpGroup:"hour"`
	Minute  *int   `regexpGroup:"minute"`
	Message string `regexpGroup:"message"`
}

// nolint:lll
const HandlePattern = `\/remind (?P<who>me|chat) every (?P<day>\d{1,2})(?:(st|nd|rd|th))? of the month ?(at (?P<hour>\d{1,2}):(?P<minute>\d{1,2}))? (?P<message>.*)`

func HandleRemindEveryDayNumber(service reminddate.Servicer) func(c tbwrap.Context) error {
	return func(c tbwrap.Context) error {
		message := new(Message)
		if err := c.Bind(message); err != nil {
			return err
		}

		repeatDateTime := mapMessageToReminderDateTime(message)
		nextSchedule, err := service.AddRepeatableReminderOnDateTime(int(c.ChatID()), c.Text(), &repeatDateTime, c.Param("message"))
		if err != nil {
			return err
		}

		_, err = c.Send(reminddate.ReminderAddedSuccessMessage(c.Param("message"), nextSchedule))

		return err
	}
}

func mapMessageToReminderDateTime(m *Message) reminder.RepeatableDateTime {
	rdt := reminder.RepeatableDateTime{
		DayOfMonth: strconv.Itoa(m.Day),
		Month:      "*",
		Hour:       "9",
		Minute:     "0",
	}

	if m.Hour != nil {
		rdt.Hour = strconv.Itoa(*m.Hour)

		if m.Minute != nil {
			rdt.Minute = strconv.Itoa(*m.Minute)
		}
	}

	return rdt
}
