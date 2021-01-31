package command

import (
	"time"

	"github.com/enrico5b1b4/telegram-bot/pkg/chatpreference"
	"github.com/enrico5b1b4/telegram-bot/pkg/reminder"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/${GOFILE} -package=mocks

type SetTimezoneServicer interface {
	SetTimeZone(chatID int, timezone string) error
}

type SetTimezoneService struct {
	reminderLoader      reminder.LoaderServicer
	chatPreferenceStore chatpreference.Storer
}

func NewSetTimezoneService(chatPreferenceStore chatpreference.Storer, reminderLoader reminder.LoaderServicer) *SetTimezoneService {
	return &SetTimezoneService{
		reminderLoader:      reminderLoader,
		chatPreferenceStore: chatPreferenceStore,
	}
}

func (s *SetTimezoneService) SetTimeZone(chatID int, timezone string) error {
	if err := validateTimeZone(timezone); err != nil {
		return err
	}

	if err := s.chatPreferenceStore.UpsertChatPreference(&chatpreference.ChatPreference{
		ChatID:   chatID,
		TimeZone: timezone,
	}); err != nil {
		return err
	}

	_, err := s.reminderLoader.ReloadSchedulesForChat(chatID)
	if err != nil {
		return err
	}

	return nil
}

// validateTimeZone validates input timezone
func validateTimeZone(tz string) error {
	_, err := time.LoadLocation(tz)
	if err != nil {
		return err
	}

	return nil
}
