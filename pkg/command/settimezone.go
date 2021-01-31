package command

import (
	"fmt"

	"github.com/enrico5b1b4/tbwrap"
)

type MessageSetTimezone struct {
	TimeZone string `regexpGroup:"timezone"`
}

// nolint:lll
const HandlePatternSetTimezone = `/settimezone (?P<timezone>.*)`

func HandleSetTimezone(service SetTimezoneServicer) func(c tbwrap.Context) error {
	return func(c tbwrap.Context) error {
		message := new(MessageSetTimezone)
		if err := c.Bind(message); err != nil {
			return err
		}

		err := service.SetTimeZone(int(c.ChatID()), message.TimeZone)
		if err != nil {
			return err
		}

		_, err = c.Send(fmt.Sprintf("Timezone has been updated to: %s", message.TimeZone))
		return err
	}
}
