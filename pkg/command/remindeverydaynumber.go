package command

import (
	"strconv"

	"github.com/enrico5b1b4/tbwrap"
	"github.com/enrico5b1b4/telegram-bot/pkg/date"
	"github.com/enrico5b1b4/telegram-bot/pkg/reminder"
)

type MessageRemindEveryDayNumber struct {
	Day     int    `regexpGroup:"day"`
	Hour    *int   `regexpGroup:"hour"`
	Minute  int    `regexpGroup:"minute"`
	AMPM    string `regexpGroup:"ampm"`
	Message string `regexpGroup:"message"`
}

// nolint:lll
const HandlePatternRemindEveryDayNumber = `/remind me every (?P<day>\d{1,2})(?:(st|nd|rd|th))? of the month ?(at (?P<hour>\d{1,2})?((:|.)(?P<minute>\d{1,2}))??(?P<ampm>am|pm)?)? (?P<message>.*)`

func HandleRemindEveryDayNumber(service reminder.Servicer) func(c tbwrap.Context) error {
	return func(c tbwrap.Context) error {
		message := new(MessageRemindEveryDayNumber)
		if err := c.Bind(message); err != nil {
			return err
		}

		repeatDateTime := mapMessageRemindEveryDayNumberToReminderDateTime(message)
		nextSchedule, err := service.AddRepeatableReminderOnDateTime(int(c.ChatID()), c.Text(), &repeatDateTime, c.Param("message"))
		if err != nil {
			return err
		}

		_, err = c.Send(ReminderAddedSuccessMessage(c.Param("message"), nextSchedule))

		return err
	}
}

func mapMessageRemindEveryDayNumberToReminderDateTime(m *MessageRemindEveryDayNumber) reminder.RepeatableDateTime {
	rdt := reminder.RepeatableDateTime{
		DayOfMonth: strconv.Itoa(m.Day),
		Month:      "*",
		Hour:       "9",
		Minute:     "0",
	}

	if m.Hour != nil {
		hour, minute := date.ConvertTo24H(*m.Hour, m.Minute, m.AMPM)

		rdt.Hour = strconv.Itoa(hour)
		rdt.Minute = strconv.Itoa(minute)
	}

	return rdt
}
