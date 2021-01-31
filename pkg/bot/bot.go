package bot

import (
	"log"

	"github.com/enrico5b1b4/telegram-bot/pkg/chatpreference"
	"github.com/enrico5b1b4/telegram-bot/pkg/command"
	"github.com/enrico5b1b4/telegram-bot/pkg/cron"
	"github.com/enrico5b1b4/telegram-bot/pkg/date"
	"github.com/enrico5b1b4/telegram-bot/pkg/reminder"
	"github.com/enrico5b1b4/telegram-bot/pkg/telegram"
	"go.etcd.io/bbolt"
)

type Bot struct {
	cronScheduler cron.Scheduler
	telegramBot   telegram.TBWrapBot
}

// nolint:funlen,lll
func New(
	allowedChats []int,
	database *bbolt.DB,
	telegramBot telegram.TBWrapBot,
) *Bot {
	cronScheduler := cron.NewScheduler()
	reminderStore := reminder.NewStore(database)
	chatPreferenceStore := chatpreference.NewStore(database)
	chatPreferenceService := chatpreference.NewService(chatPreferenceStore)
	remindCronFuncService := reminder.NewCronFuncService(telegramBot, cronScheduler, reminderStore, chatPreferenceStore)
	remindListService := command.NewRemindListService(reminderStore, cronScheduler, chatPreferenceStore)
	remindDeleteService := command.NewRemindeDeleteService(reminderStore, cronScheduler)
	reminderScheduler := reminder.NewScheduler(telegramBot, remindCronFuncService, reminderStore, cronScheduler, chatPreferenceStore)
	remindDateService := reminder.NewService(reminderScheduler, reminderStore, chatPreferenceStore, date.RealTimeNow)
	remindDetailService := command.NewRemindDetailService(reminderStore, cronScheduler, chatPreferenceStore)
	reminderLoader := reminder.NewLoaderService(telegramBot, cronScheduler, reminderStore, chatPreferenceStore, remindCronFuncService)
	setTimeZoneService := command.NewSetTimezoneService(chatPreferenceStore, reminderLoader)
	remindDetailButtons := command.NewRemindDetailButtons()
	remindListButtons := command.NewRemindListButtons()
	reminderCompleteButtons := reminder.NewButtons()

	chatPreferenceService.CreateDefaultChatPreferences(allowedChats)

	// check if DB exists and load schedules
	remindersLoaded, err := reminderLoader.LoadSchedulesFromDB()
	if err != nil {
		panic(err)
	}
	log.Printf("loaded %d reminders", remindersLoaded)

	telegramBot.Handle(command.HandlePatternRemindList, command.HandleRemindList(remindListService, remindListButtons))
	telegramBot.Handle(command.HandlePatternHelp, command.HandleRemindHelp())
	telegramBot.HandleMultiRegExp(command.HandlePatternRemindDetail, command.HandleRemindDetail(remindDetailService, command.NewRemindDetailButtons()))
	telegramBot.HandleMultiRegExp(command.HandlePatternRemindDelete, command.HandleRemindDelete(remindDeleteService))
	telegramBot.HandleRegExp(
		command.HandlePatternRemindDayMonth,
		command.HandleRemindDayMonth(remindDateService),
	)
	telegramBot.HandleRegExp(
		command.HandlePatternRemindDayOfWeek,
		command.HandleRemindDayOfWeek(remindDateService),
	)
	telegramBot.HandleRegExp(
		command.HandlePatternRemindEveryDayNumber,
		command.HandleRemindEveryDayNumber(remindDateService),
	)
	telegramBot.HandleRegExp(
		command.HandlePatternRemindEveryDayNumberMonth,
		command.HandleRemindEveryDayNumberMonth(remindDateService),
	)
	telegramBot.HandleRegExp(
		command.HandlePatternRemindIn,
		command.HandleRemindIn(remindDateService),
	)
	telegramBot.HandleRegExp(
		command.HandlePatternRemindEvery,
		command.HandleRemindEvery(remindDateService),
	)
	telegramBot.HandleRegExp(
		command.HandlePatternRemindWhen,
		command.HandleRemindWhen(remindDateService),
	)
	telegramBot.HandleRegExp(
		command.HandlePatternRemindEveryDayOfWeek,
		command.HandleRemindEveryDayOfWeek(remindDateService),
	)
	telegramBot.HandleRegExp(
		command.HandlePatternRemindEveryDay,
		command.HandleRemindEveryDay(remindDateService),
	)
	telegramBot.HandleRegExp(
		command.HandlePatternRemindAt,
		command.HandleRemindAt(remindDateService),
	)
	telegramBot.Handle(command.HandlePatternGetTimezone, command.HandleGetTimezone(chatPreferenceStore))
	telegramBot.HandleRegExp(command.HandlePatternSetTimezone, command.HandleSetTimezone(setTimeZoneService))

	// buttons
	telegramBot.HandleButton(
		remindDetailButtons[command.ReminderDetailCloseCommandBtn],
		command.HandleReminderDetailCloseBtn(),
	)
	telegramBot.HandleButton(
		remindDetailButtons[command.ReminderDetailDeleteBtn],
		command.HandleReminderDetailDeleteBtn(remindDetailService),
	)
	telegramBot.HandleButton(
		remindDetailButtons[command.ReminderDetailShowReminderCommandBtn],
		command.HandleReminderShowReminderCommandBtn(remindDetailService),
	)
	telegramBot.HandleButton(
		remindListButtons[command.ReminderListRemoveCompletedRemindersBtn],
		command.HandleReminderListRemoveCompletedRemindersBtn(remindListService),
	)
	telegramBot.HandleButton(
		remindListButtons[command.ReminderListCloseCommandBtn],
		command.HandleRemindListCloseBtn(),
	)
	telegramBot.HandleButton(
		reminderCompleteButtons[reminder.Snooze10MinuteBtn],
		reminder.HandleReminderSnoozeAmountDateTimeBtn(remindDateService, reminderStore, reminder.AmountDateTime{Minutes: 10}),
	)
	telegramBot.HandleButton(
		reminderCompleteButtons[reminder.Snooze20MinuteBtn],
		reminder.HandleReminderSnoozeAmountDateTimeBtn(remindDateService, reminderStore, reminder.AmountDateTime{Minutes: 20}),
	)
	telegramBot.HandleButton(
		reminderCompleteButtons[reminder.Snooze30MinuteBtn],
		reminder.HandleReminderSnoozeAmountDateTimeBtn(remindDateService, reminderStore, reminder.AmountDateTime{Minutes: 30}),
	)
	telegramBot.HandleButton(
		reminderCompleteButtons[reminder.Snooze1HourBtn],
		reminder.HandleReminderSnoozeAmountDateTimeBtn(remindDateService, reminderStore, reminder.AmountDateTime{Minutes: 60}),
	)
	telegramBot.HandleButton(
		reminderCompleteButtons[reminder.SnoozeThisAfternoonBtn],
		reminder.HandleReminderSnoozeWordDateTimeBtn(remindDateService, reminderStore, reminder.WordDateTime{
			When:   reminder.Today,
			Hour:   15,
			Minute: 0,
		}),
	)
	telegramBot.HandleButton(
		reminderCompleteButtons[reminder.SnoozeThisEveningBtn],
		reminder.HandleReminderSnoozeWordDateTimeBtn(remindDateService, reminderStore, reminder.WordDateTime{
			When:   reminder.Today,
			Hour:   20,
			Minute: 0,
		}),
	)
	telegramBot.HandleButton(
		reminderCompleteButtons[reminder.SnoozeTomorrowMorningBtn],
		reminder.HandleReminderSnoozeWordDateTimeBtn(remindDateService, reminderStore, reminder.WordDateTime{
			When:   reminder.Tomorrow,
			Hour:   9,
			Minute: 0,
		}),
	)
	telegramBot.HandleButton(
		reminderCompleteButtons[reminder.SnoozeTomorrowAfternoonBtn],
		reminder.HandleReminderSnoozeWordDateTimeBtn(remindDateService, reminderStore, reminder.WordDateTime{
			When:   reminder.Tomorrow,
			Hour:   15,
			Minute: 0,
		}),
	)
	telegramBot.HandleButton(
		reminderCompleteButtons[reminder.SnoozeTomorrowEveningBtn],
		reminder.HandleReminderSnoozeWordDateTimeBtn(remindDateService, reminderStore, reminder.WordDateTime{
			When:   reminder.Tomorrow,
			Hour:   20,
			Minute: 0,
		}),
	)
	telegramBot.HandleButton(
		reminderCompleteButtons[reminder.SnoozeBtn],
		reminder.HandleReminderSnoozeBtn(reminderStore),
	)
	telegramBot.HandleButton(
		reminderCompleteButtons[reminder.SnoozeCloseBtn],
		reminder.HandleReminderSnoozeCloseBtn(),
	)
	telegramBot.HandleButton(
		reminderCompleteButtons[reminder.CompleteBtn],
		reminder.HandleReminderCompleteBtn(remindCronFuncService, reminderStore),
	)

	return &Bot{
		cronScheduler: cronScheduler,
		telegramBot:   telegramBot,
	}
}

func (b *Bot) Start() {
	b.cronScheduler.Start()
	b.telegramBot.Start()
}
