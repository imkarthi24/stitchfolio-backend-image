package service

import (
	"context"
	"fmt"
	"strings"

	config_cache "github.com/imkarthi24/sf-backend/internal/cache"
	"github.com/imkarthi24/sf-backend/internal/config"
	"github.com/imkarthi24/sf-backend/internal/mapper"
	requestModel "github.com/imkarthi24/sf-backend/internal/model/request"
	responseModel "github.com/imkarthi24/sf-backend/internal/model/response"
	"github.com/imkarthi24/sf-backend/internal/repository"
	"github.com/imkarthi24/sf-backend/pkg/errs"
)

type MasterConfigService interface {
	Save(ctx *context.Context, config requestModel.MasterConfig) *errs.XError
	Update(ctx *context.Context, config requestModel.MasterConfig, id uint) *errs.XError
	Get(ctx *context.Context, id uint) (*responseModel.MasterConfig, *errs.XError)
	Browse(ctx *context.Context, query string) ([]responseModel.MasterConfig, *errs.XError)
	GetByName(ctx *context.Context, name string) (string, *errs.XError)
	GetMultipleValues(ctx *context.Context, names []string) (map[string]string, *errs.XError)
}

type masterConfigService struct {
	masterConfigRepo repository.MasterConfigRepository
	mapper           mapper.Mapper
	config           config.AppConfig
	respMapper       mapper.ResponseMapper
}

func ProvideMasterConfigService(repo repository.MasterConfigRepository, mapper mapper.Mapper, config config.AppConfig, respMapper mapper.ResponseMapper) MasterConfigService {
	return &masterConfigService{
		masterConfigRepo: repo,
		mapper:           mapper,
		config:           config,
		respMapper:       respMapper,
	}
}

func (svc *masterConfigService) Save(ctx *context.Context, config requestModel.MasterConfig) *errs.XError {
	dbMasterConfig, err := svc.mapper.MasterConfig(config)

	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to save Master confiig", err)
	}

	return svc.masterConfigRepo.Create(ctx, dbMasterConfig)
}

func (svc *masterConfigService) Update(ctx *context.Context, config requestModel.MasterConfig, id uint) *errs.XError {
	existingConfig, err1 := svc.masterConfigRepo.Get(ctx, id)
	if err1 != nil {
		return err1
	}

	updateCache := func(value string) {
		cacheName := fmt.Sprintf("%s.%s", config.Type, config.Name)
		config_cache.SetValue(ctx, cacheName, value)
	}

	dbMasterConfig, err := svc.mapper.MasterConfig(config)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to map Master Config", err)
	}

	dbMasterConfig.ID = id
	dbMasterConfig.PreviousValue = existingConfig.CurrentValue

	value := config.CurrentValue
	if config.UseDefault {
		value = config.DefaultValue
	}

	err1 = svc.masterConfigRepo.Update(ctx, dbMasterConfig)
	if err1 != nil {
		return err1
	}

	go updateCache(value)
	return nil

}

func (svc *masterConfigService) Get(ctx *context.Context, id uint) (*responseModel.MasterConfig, *errs.XError) {

	dbMasterConfig, err := svc.masterConfigRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	respMasterConfig, errr := svc.respMapper.MasterConfig(dbMasterConfig)
	if errr != nil {
		return nil, errs.Wrap(errr, "Unable to map Master Config")
	}
	return respMasterConfig, nil
}

// GetByName takes value in A.B format and gives back the value
func (svc *masterConfigService) GetByName(ctx *context.Context, name string) (string, *errs.XError) {

	// first try to get from cache
	val, ok := config_cache.GetValue(ctx, name)
	if ok {
		return val.(string), nil
	}

	updateCache := func(value string) {
		config_cache.SetValue(ctx, name, value)
	}

	//if not availbale in cache , then fetch from db
	splitValues := strings.Split(name, ".")
	config, err := svc.masterConfigRepo.GetValue(ctx, splitValues[0], splitValues[1])
	if err != nil {
		return "", nil
	}

	value := config.CurrentValue
	if config.UseDefault {
		value = config.DefaultValue
	}

	//update the cache at last
	go updateCache(value)
	return value, nil

}

func (svc *masterConfigService) Browse(ctx *context.Context, query string) ([]responseModel.MasterConfig, *errs.XError) {

	masterConfigs, err := svc.masterConfigRepo.GetForBrowse(ctx, query)
	if err != nil {
		return nil, err
	}

	respMasterConfig, errr := svc.respMapper.MasterConfigs(masterConfigs)
	if errr != nil {
		return nil, errs.Wrap(errr, "Unable to map MasterConfigs")
	}

	return respMasterConfig, nil
}

func (svc *masterConfigService) GetMultipleValues(ctx *context.Context, names []string) (map[string]string, *errs.XError) {
	result := make(map[string]string)

	for _, name := range names {
		value, err := svc.GetByName(ctx, name)
		if err != nil {
			// If there's an error, we can either skip it or set empty string
			// Setting empty string to ensure all requested keys are in the response
			result[name] = ""
			continue
		}
		result[name] = value
	}

	return result, nil
}
