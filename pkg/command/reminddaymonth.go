package command

import (
	"github.com/enrico5b1b4/tbwrap"
	"github.com/enrico5b1b4/telegram-bot/pkg/date"
	"github.com/enrico5b1b4/telegram-bot/pkg/reminder"
)

type MessageRemindDayMonth struct {
	Day     int    `regexpGroup:"day"`
	Month   string `regexpGroup:"month"`
	Hour    *int   `regexpGroup:"hour"`
	Minute  int    `regexpGroup:"minute"`
	AMPM    string `regexpGroup:"ampm"`
	Message string `regexpGroup:"message"`
}

// nolint:lll
const HandlePatternRemindDayMonth = `/remind me on the (?P<day>\d{1,2})(?:(st|nd|rd|th))? ?(of (?P<month>(J|j)anuary|(F|f)ebruary|(M|m)arch|(A|a)pril|(M|m)ay|(J|j)une|(J|j)uly|(A|a)ugust|(S|s)eptember|(O|o)ctober|(N|n)ovember|(D|d)ecember))? ?(at (?P<hour>\d{1,2})?((:|.)(?P<minute>\d{1,2}))??(?P<ampm>am|pm)?)? (?P<message>.*)`

func HandleRemindDayMonth(service reminder.Servicer) func(c tbwrap.Context) error {
	return func(c tbwrap.Context) error {
		message := new(MessageRemindDayMonth)
		if err := c.Bind(message); err != nil {
			return err
		}

		dateTime := mapMessageRemindDayMonthToReminderDateTime(message)
		nextSchedule, err := service.AddReminderOnDateTime(int(c.ChatID()), c.Text(), dateTime, c.Param("message"))
		if err != nil {
			return err
		}

		_, err = c.Send(ReminderAddedSuccessMessage(c.Param("message"), nextSchedule))
		return err
	}
}

func mapMessageRemindDayMonthToReminderDateTime(m *MessageRemindDayMonth) reminder.DateTime {
	dt := reminder.DateTime{
		DayOfMonth: m.Day,
		Month:      date.ToNumericMonth(m.Month),
		Hour:       9,
		Minute:     0,
	}

	if m.Hour != nil {
		hour, minute := date.ConvertTo24H(*m.Hour, m.Minute, m.AMPM)

		dt.Hour = hour
		dt.Minute = minute
	}

	return dt
}
