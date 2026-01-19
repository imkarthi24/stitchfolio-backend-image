package repository

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/entities"
	"github.com/imkarthi24/sf-backend/internal/repository/scopes"
	"github.com/imkarthi24/sf-backend/pkg/db"
	"github.com/imkarthi24/sf-backend/pkg/errs"
)

type ChannelRepository interface {
	Save(*context.Context, *entities.Channel) *errs.XError
	Update(*context.Context, *entities.Channel) *errs.XError
	Get(*context.Context, uint) (*entities.Channel, *errs.XError)
	Delete(*context.Context, uint) *errs.XError
	GetAllChannels(*context.Context, string) ([]entities.Channel, *errs.XError)
	ChannelAutoComplete(*context.Context, string) ([]entities.Channel, *errs.XError)
}

type channelRepository struct {
	GormDAL
}

func ProvideChannelRepository(customDB GormDAL) ChannelRepository {
	return &channelRepository{GormDAL: customDB}
}

func (ur *channelRepository) Save(ctx *context.Context, channel *entities.Channel) *errs.XError {
	res := ur.WithDB(ctx).Create(&channel)
	if res.Error != nil {
		return errs.NewXError(errs.DATABASE, "Unable to save channel", res.Error)
	}

	//set channelId as the one created
	// This functionality has been done using gormCreateHooks

	// res = ur.txn.Txn(ctx).Model(entities.Channel{}).Where("id = ?", channel.ID).Update("channel_id", channel.ID)
	// if res.Error != nil {
	// 	return errs.NewXError(errs.DATABASE, "Unable to save channel", res.Error)
	// }

	return nil

}

func (ur *channelRepository) Update(ctx *context.Context, channel *entities.Channel) *errs.XError {
	return ur.GormDAL.Update(ctx, *channel)
}

func (ur *channelRepository) Get(ctx *context.Context, id uint) (*entities.Channel, *errs.XError) {
	channel := entities.Channel{}
	res := ur.WithDB(ctx).
		Preload("OwnerUser").
		Find(&channel, id)

	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to find channel", res.Error)
	}
	return &channel, nil
}

func (ur *channelRepository) Delete(ctx *context.Context, id uint) *errs.XError {

	channel := &entities.Channel{Model: &entities.Model{ID: id, IsActive: false}}

	err := ur.GormDAL.Delete(ctx, channel)

	return err

}

func (ur *channelRepository) GetAllChannels(ctx *context.Context, autoCompName string) ([]entities.Channel, *errs.XError) {
	channels := new([]entities.Channel)

	res := ur.WithDB(ctx).
		Preload("OwnerUser").
		Scopes(scopes.ChannelAutoComplete_Filter(autoCompName)).
		Scopes(scopes.IsActive()).
		Scopes(db.Paginate(ctx)).
		Find(&channels)

	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to fetch all channels", res.Error)
	}

	return *channels, nil
}

func (ur *channelRepository) ChannelAutoComplete(ctx *context.Context, autoCompName string) ([]entities.Channel, *errs.XError) {
	channels := new([]entities.Channel)

	res := ur.WithDB(ctx).
		Scopes(scopes.ChannelAutoComplete_Filter(autoCompName)).
		Select("id", "name").
		Find(&channels)

	if res.Error != nil {
		return nil, errs.NewXError(errs.DATABASE, "Unable to fetch  channels for auto complete", res.Error)
	}

	return *channels, nil
}
