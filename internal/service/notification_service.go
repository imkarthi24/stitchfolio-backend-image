package service

import (
	"context"
	"errors"

	"github.com/imkarthi24/sf-backend/internal/config"
	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/mapper"
	requestModel "github.com/imkarthi24/sf-backend/internal/model/request"
	"github.com/imkarthi24/sf-backend/internal/repository"
	"github.com/imkarthi24/sf-backend/pkg/errs"
	"github.com/imkarthi24/sf-backend/pkg/service/email"
	"github.com/imkarthi24/sf-backend/pkg/util"
)

type NotificationService interface {
	CreateEmailNotification(ctx *context.Context, notif requestModel.EmaiNotification) *errs.XError
	CreateEmailNotifications(ctx *context.Context, notifs []requestModel.EmaiNotification) *errs.XError
	GetPendingNotifications(ctx *context.Context) ([]entities.Notification, *errs.XError)

	SendNotification(ctx *context.Context, notif entities.Notification) *errs.XError
}

type notificationService struct {
	notifRepo  repository.NotificationRepository
	mapper     mapper.Mapper
	smtpConfig config.SMTPConfig
}

func ProvideNotificationService(repo repository.NotificationRepository, mapper mapper.Mapper, smtpConfig config.SMTPConfig) NotificationService {
	return &notificationService{
		notifRepo:  repo,
		mapper:     mapper,
		smtpConfig: smtpConfig,
	}
}

func (svc *notificationService) CreateEmailNotification(ctx *context.Context, email requestModel.EmaiNotification) *errs.XError {

	emailNotif, err := createEmailNotification(email)
	if err != nil {
		return errs.NewXError(errs.EMAILERROR, "Error creating Email Notification", err)
	}

	notification := createNotification(email.Notification)

	notification.AddEmailNotification(*emailNotif)

	return svc.notifRepo.CreateNotification(ctx, *notification)

}

func (svc *notificationService) CreateEmailNotifications(ctx *context.Context, notifs []requestModel.EmaiNotification) *errs.XError {

	emailNotifs := make([]entities.EmailNotification, 0)

	for _, notif := range notifs {
		if util.IsNilOrEmptyString(&notif.ToMailAddress) {
			continue
		}
		emailNotif, err := createEmailNotification(notif)
		if err != nil {
			return errs.NewXError(errs.EMAILERROR, "Error creating Email Notification", err)
		}
		emailNotifs = append(emailNotifs, *emailNotif)
	}

	notification := createNotification(notifs[0].Notification)
	notification.AddEmailNotification(emailNotifs...)

	err := svc.notifRepo.CreateNotification(ctx, *notification)

	return err
}

func (svc *notificationService) GetPendingNotifications(ctx *context.Context) ([]entities.Notification, *errs.XError) {
	return svc.notifRepo.GetPendingNotifications(ctx)
}

func (svc *notificationService) SendNotification(ctx *context.Context, notif entities.Notification) *errs.XError {

	notifStatus := entities.NOTIF_COMPLETED
	err := svc.sendEmailNotification(ctx, notif.EmailNotifications)
	if err != nil {
		notifStatus = entities.NOTIF_PARTIAL
	}

	// have to include whatsaap notifications as well
	return svc.notifRepo.UpdateNotificationStatus(ctx, notif.ID, notifStatus)

}

func createNotification(notif *requestModel.Notification) *entities.Notification {
	return &entities.Notification{
		Status:       entities.NOTIF_PENDING,
		SourceEntity: notif.SourceEntity,
		EntityId:     notif.EntityId,
	}
}

func createEmailNotification(notif requestModel.EmaiNotification) (*entities.EmailNotification, error) {

	if notif.EmailContent != nil {

		_, body, err := email.BuildEmailBody(*notif.EmailContent)
		if err != nil {
			return nil, err
		}

		return &entities.EmailNotification{
			Status:        string(entities.NOTIF_PENDING),
			ToMailAddress: notif.EmailContent.To[0],
			Subject:       notif.EmailContent.Subject,
			Body:          string(body),
		}, nil
	}

	return &entities.EmailNotification{
		Status:        string(entities.NOTIF_PENDING),
		ToMailAddress: notif.ToMailAddress,
		Subject:       notif.Subject,
		Body:          notif.Body,
	}, nil
}

func (svc *notificationService) sendEmailNotification(ctx *context.Context, emailNotifs []entities.EmailNotification) *errs.XError {

	var faulted bool
	for _, notif := range emailNotifs {
		recipients := []string{notif.ToMailAddress}
		mail := email.EmailContent{
			To:      recipients,
			Subject: notif.Subject,
			Message: &notif.Body,
		}

		err := email.SendEmail(&svc.smtpConfig, mail)
		notifStatus := entities.NOTIF_COMPLETED
		if err != nil {
			notifStatus = entities.NOTIF_FAULTED
			faulted = true
		}
		notif.Status = string(notifStatus)

		svc.notifRepo.UpdateEmailNotificationStatus(ctx, notif.ID, notifStatus)

	}

	if faulted {
		return errs.NewXError(errs.SMTPERROR, errs.SMTP_ERROR, errors.New("Error sending email notification"))
	}

	return nil
}
