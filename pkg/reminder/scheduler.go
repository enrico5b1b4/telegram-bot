package reminder

//go:generate mockgen -source=$GOFILE -destination=mocks/${GOFILE} -package=mocks

import (
	"fmt"
	"time"

	"github.com/enrico5b1b4/telegram-bot/pkg/chatpreference"
	"github.com/enrico5b1b4/telegram-bot/pkg/cron"
	"github.com/enrico5b1b4/telegram-bot/pkg/telegram"
)

type Scheduler interface {
	AddReminder(r *Reminder) (int, error)
	GetNextScheduleTime(cronID int) (time.Time, error)
}

type SchedulerManager struct {
	reminderStore           Storer
	reminderCronFuncService CronFuncServicer
	scheduler               cron.Scheduler
	bot                     telegram.TBWrapBot
	chatPreferenceStore     chatpreference.Storer
}

func NewScheduler(
	bot telegram.TBWrapBot,
	reminderCronFuncService CronFuncServicer,
	reminderStore Storer,
	scheduler cron.Scheduler,
	chatPreferenceStore chatpreference.Storer,
) *SchedulerManager {
	return &SchedulerManager{
		bot:                     bot,
		reminderStore:           reminderStore,
		reminderCronFuncService: reminderCronFuncService,
		scheduler:               scheduler,
		chatPreferenceStore:     chatPreferenceStore,
	}
}

func (s *SchedulerManager) AddReminder(rem *Reminder) (int, error) {
	chatPreference, err := s.chatPreferenceStore.GetChatPreference(rem.Job.ChatID)
	if err != nil {
		return 0, err
	}

	schedule := fmt.Sprintf("CRON_TZ=%s %s", chatPreference.TimeZone, rem.Job.Schedule)
	reminderCronID, err := s.scheduler.Add(schedule, NewCronFunc(s.reminderCronFuncService, s.bot, rem))
	if err != nil {
		return 0, err
	}

	return reminderCronID, nil
}

func (s *SchedulerManager) GetNextScheduleTime(cronID int) (time.Time, error) {
	cronEntry := s.scheduler.GetEntryByID(cronID)

	return cronEntry.Next, nil
}
