package service

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/mapper"
	requestModel "github.com/imkarthi24/sf-backend/internal/model/request"
	responseModel "github.com/imkarthi24/sf-backend/internal/model/response"
	"github.com/imkarthi24/sf-backend/internal/repository"
	"github.com/loop-kar/pixie/errs"
)

type ChannelService interface {
	SaveChannel(*context.Context, requestModel.Channel) *errs.XError
	UpdateChannel(*context.Context, requestModel.Channel, uint) *errs.XError
	Get(*context.Context, uint) (*responseModel.Channel, *errs.XError)
	GetAllChannels(*context.Context, string) ([]responseModel.Channel, *errs.XError)
	ChannelAutoComplete(*context.Context, string) ([]responseModel.ChannelAutoComplete, *errs.XError)
	Delete(*context.Context, uint) *errs.XError
}

type channelService struct {
	channelRepo repository.ChannelRepository
	userRepo    repository.UserRepository
	mapper      mapper.Mapper
	respMapper  mapper.ResponseMapper
}

func ProvideChannelService(repo repository.ChannelRepository, userRepo repository.UserRepository, mapper mapper.Mapper, respMapper mapper.ResponseMapper) ChannelService {
	return channelService{
		channelRepo: repo,
		userRepo:    userRepo,
		mapper:      mapper,
		respMapper:  respMapper,
	}
}

func (svc channelService) SaveChannel(ctx *context.Context, chanel requestModel.Channel) *errs.XError {

	dbChannel, err := svc.mapper.Channel(chanel)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to save channel", err)
	}

	errr := svc.channelRepo.Save(ctx, dbChannel)
	if errr != nil {
		return errr
	}

	// This functionality has been done using gormCreateHooks

	// userUpdate := entities.User{
	// 	Model: &entities.Model{ID: chanel.OwnerUserId, ChannelId: dbChannel.ID},
	// }

	// errr = svc.userRepo.Update(ctx, &userUpdate)
	// if errr != nil {
	// 	return errr
	// }

	return nil

}

func (svc channelService) UpdateChannel(ctx *context.Context, channel requestModel.Channel, id uint) *errs.XError {

	dbChannel, err := svc.mapper.Channel(channel)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to update channel", err)
	}

	dbChannel.ID = id
	errr := svc.channelRepo.Update(ctx, dbChannel)
	if errr != nil {
		return errr
	}
	return nil

}

func (svc channelService) Get(ctx *context.Context, id uint) (*responseModel.Channel, *errs.XError) {

	channel, err := svc.channelRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return svc.respMapper.Channel(channel), nil

}

func (svc channelService) Delete(ctx *context.Context, id uint) *errs.XError {
	return svc.channelRepo.Delete(ctx, id)
}

func (svc channelService) GetAllChannels(ctx *context.Context, autoCompName string) ([]responseModel.Channel, *errs.XError) {
	channels, err := svc.channelRepo.GetAllChannels(ctx, autoCompName)
	if err != nil {
		return nil, err
	}
	return svc.respMapper.Channels(channels), nil
}

func (svc channelService) ChannelAutoComplete(ctx *context.Context, autoCompName string) ([]responseModel.ChannelAutoComplete, *errs.XError) {
	channels, err := svc.channelRepo.ChannelAutoComplete(ctx, autoCompName)
	if err != nil {
		return nil, err
	}

	res := make([]responseModel.ChannelAutoComplete, 0)

	for _, ch := range channels {
		res = append(res, responseModel.ChannelAutoComplete{
			ChannelID: ch.ID,
			Name:      ch.Name,
		})
	}

	return res, nil
}
