package task

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/service"
	"github.com/loop-kar/pixie/errs"
	"github.com/loop-kar/pixie/task"
)

type NotificationTaskParam struct {
	//BaseTaskParam provides inbuilt values by default
	*task.BaseTaskParam
	//....
	// this can have additional values
	// but doesn't have now
}

type NotificationTask struct {
	// kind of inheriting BaseTask methods, by this even if NotificationTask
	// doesnt implement methods of the interface , it'll fallback to BaseTask
	// implementation of the interface methods
	*task.BaseTask

	// holds the parameters need for the notification task
	*NotificationTaskParam

	notifSvc service.NotificationService
}

func ProvideNotificationTask(param *NotificationTaskParam, svc service.NotificationService) task.IBaseTask {
	context := context.Background()
	return &NotificationTask{
		BaseTask: &task.BaseTask{
			Param: param.BaseTaskParam,
			Ctx:   &context,
		},
		NotificationTaskParam: param,
		notifSvc:              svc,
	}
}

func (t *NotificationTask) FetchEntitySet() (bool, []task.TaskResponse, *errs.XError) {
	notifs, err := t.notifSvc.GetPendingNotifications(t.Ctx)
	if err != nil {
		return false, nil, err
	}

	res := transform(notifs)
	return true, res, nil
}

func (t *NotificationTask) ProcessEntitySet(notifications []task.TaskResponse) (bool, *errs.XError) {

	for _, item := range notifications {
		notif := item.(entities.Notification)
		t.notifSvc.SendNotification(t.Ctx, notif)
	}

	return false, nil
}

func transform(items []entities.Notification) []task.TaskResponse {
	response := make([]task.TaskResponse, len(items))
	for i := range items {
		response[i] = items[i]
	}
	return response
}
