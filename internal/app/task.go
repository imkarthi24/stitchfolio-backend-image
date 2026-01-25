package app

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/config"
	"github.com/imkarthi24/sf-backend/internal/service/base"
	tsk "github.com/imkarthi24/sf-backend/internal/task"
	_log "github.com/loop-kar/pixie/log"
	"github.com/loop-kar/pixie/task"

	"github.com/newrelic/go-agent/v3/newrelic"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type Task struct {
	ChitDb      *gorm.DB
	BaseService base.BaseService
	NewRelic    *newrelic.Application
	Config      config.AppConfig
	Cron        *cron.Cron
}

// Sample task runner
func (a *Task) Start(ctx *context.Context, checkErr func(err error)) {
	_log.InitLogger(a.NewRelic)

	// Run notification task at 9AM IST
	// _, err := a.Cron.AddFunc("0 9 * * * *", func() {
	// 	// a.NotificationRunnerTask(ctx)
	// 	a.FeeCallBackReminderTask(ctx)
	// })

	// checkErr(err)

	// Run student reminder task at 10AM IST
	// _, err = a.Cron.AddFunc("0 10 * * * *", func() {
	// 	a.StudentReminderTask(ctx)
	// })

	// checkErr(err)

	a.Cron.Start()

	//_log.FromCtx(ctx).Info("Cron jobs started successfully")

}

func (a *Task) NotificationRunnerTask(ctx *context.Context) {

	param := tsk.NotificationTaskParam{
		BaseTaskParam: &task.BaseTaskParam{AbortProceesExecutionOnFailure: true},
	}

	notifTask := tsk.ProvideNotificationTask(&param, a.BaseService.NotificationService)

	jobRunner := task.ProvideJobRunner(notifTask, *param.BaseTaskParam)
	jobRunner.CreateAdHocJob(true)

}

func (a *Task) Shutdown(ctx *context.Context, checkErr func(err error)) {
	// Stop the cron scheduler
	if a.Cron != nil {
		a.Cron.Stop()
	}

	//_log.FromCtx(ctx).Info("Task Service Stopped...")
}
