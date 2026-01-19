package repository

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/repository/scopes"
	"github.com/loop-kar/pixie/errs"
)

type NotificationRepository interface {
	CreateNotification(ctx *context.Context, notif entities.Notification) *errs.XError
	GetPendingNotifications(ctx *context.Context) ([]entities.Notification, *errs.XError)
	UpdateEmailNotificationStatus(ctx *context.Context, id uint, status entities.NotificationStatus) *errs.XError
	UpdateNotificationStatus(ctx *context.Context, id uint, status entities.NotificationStatus) *errs.XError
}

type notificationRepository struct {
	GormDAL
}

func ProvideNotificationRepository(customDB GormDAL) NotificationRepository {
	return &notificationRepository{GormDAL: customDB}

}

func (repo *notificationRepository) CreateNotification(ctx *context.Context, notif entities.Notification) *errs.XError {
	res := repo.WithDB(ctx).Create(&notif)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to create notification", res.Error)
	}
	return nil
}

func (repo *notificationRepository) GetPendingNotifications(ctx *context.Context) ([]entities.Notification, *errs.XError) {
	pendingNotifications := make([]entities.Notification, 0)

	res := repo.WithDB(ctx).
		Scopes(scopes.IsActive()).
		Preload("EmailNotifications").
		Preload("WhatsappNotifications").
		Where("status = ?", entities.NOTIF_PENDING).
		Find(&pendingNotifications)

	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to fetch pending notifications", res.Error)
	}

	return pendingNotifications, nil
}

func (repo *notificationRepository) UpdateEmailNotificationStatus(ctx *context.Context, id uint, status entities.NotificationStatus) *errs.XError {
	notif := entities.EmailNotification{
		Model:  &entities.Model{ID: id},
		Status: string(status),
	}
	res := repo.WithDB(ctx).Updates(notif)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to update Notification", res.Error)
	}
	return nil
}

func (repo *notificationRepository) UpdateNotificationStatus(ctx *context.Context, id uint, status entities.NotificationStatus) *errs.XError {
	notif := entities.Notification{
		Model:  &entities.Model{ID: id},
		Status: status,
	}
	res := repo.WithDB(ctx).Updates(notif)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to update Notification", res.Error)
	}
	return nil
}
