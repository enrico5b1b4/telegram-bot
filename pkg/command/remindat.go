package command

import (
	"github.com/enrico5b1b4/tbwrap"
	"github.com/enrico5b1b4/telegram-bot/pkg/date"
	"github.com/enrico5b1b4/telegram-bot/pkg/reminder"
)

type MessageRemindAt struct {
	Hour    int    `regexpGroup:"hour"`
	Minute  int    `regexpGroup:"minute"`
	AMPM    string `regexpGroup:"ampm"`
	Message string `regexpGroup:"message"`
}

// nolint:lll
const HandlePatternRemindAt = `/remind me at (?P<hour>\d{1,2})?((:|.)(?P<minute>\d{1,2}))??(?P<ampm>am|pm)? (?P<message>.*)`

func HandleRemindAt(service reminder.Servicer) func(c tbwrap.Context) error {
	return func(c tbwrap.Context) error {
		message := new(MessageRemindAt)
		if err := c.Bind(message); err != nil {
			return err
		}

		dateTime := mapMessageRemindAtToReminderWordDateTime(message)
		nextSchedule, err := service.AddReminderOnWordDateTime(int(c.ChatID()), c.Text(), dateTime, c.Param("message"))
		if err != nil {
			return err
		}

		_, err = c.Send(ReminderAddedSuccessMessage(c.Param("message"), nextSchedule))
		return err
	}
}

func mapMessageRemindAtToReminderWordDateTime(m *MessageRemindAt) reminder.WordDateTime {
	hour, minute := date.ConvertTo24H(m.Hour, m.Minute, m.AMPM)

	return reminder.WordDateTime{
		When:   reminder.Today,
		Hour:   hour,
		Minute: minute,
	}
}
