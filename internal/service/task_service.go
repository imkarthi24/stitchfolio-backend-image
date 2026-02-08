package service

import (
	"context"

	"github.com/imkarthi24/sf-backend/internal/mapper"
	requestModel "github.com/imkarthi24/sf-backend/internal/model/request"
	responseModel "github.com/imkarthi24/sf-backend/internal/model/response"
	"github.com/imkarthi24/sf-backend/internal/repository"
	"github.com/loop-kar/pixie/errs"
)

type TaskService interface {
	SaveTask(*context.Context, requestModel.Task) *errs.XError
	UpdateTask(*context.Context, requestModel.Task, uint) *errs.XError
	Get(*context.Context, uint) (*responseModel.Task, *errs.XError)
	GetAll(*context.Context, string) ([]responseModel.Task, *errs.XError)
	Delete(*context.Context, uint) *errs.XError
}

type taskService struct {
	taskRepo   repository.TaskRepository
	mapper     mapper.Mapper
	respMapper mapper.ResponseMapper
}

func ProvideTaskService(repo repository.TaskRepository, mapper mapper.Mapper, respMapper mapper.ResponseMapper) TaskService {
	return taskService{
		taskRepo:   repo,
		mapper:     mapper,
		respMapper: respMapper,
	}
}

func (svc taskService) SaveTask(ctx *context.Context, task requestModel.Task) *errs.XError {
	dbTask, err := svc.mapper.Task(task)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to save task", err)
	}

	errr := svc.taskRepo.Create(ctx, dbTask)
	if errr != nil {
		return errr
	}
	return nil
}

func (svc taskService) UpdateTask(ctx *context.Context, task requestModel.Task, id uint) *errs.XError {
	dbTask, err := svc.mapper.Task(task)
	if err != nil {
		return errs.NewXError(errs.INVALID_REQUEST, "Unable to update task", err)
	}

	dbTask.ID = id
	errr := svc.taskRepo.Update(ctx, dbTask)
	if errr != nil {
		return errr
	}
	return nil
}

func (svc taskService) Get(ctx *context.Context, id uint) (*responseModel.Task, *errs.XError) {
	task, err := svc.taskRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	mappedTask, mapErr := svc.respMapper.Task(task)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map Task data", mapErr)
	}
	return mappedTask, nil
}

func (svc taskService) GetAll(ctx *context.Context, search string) ([]responseModel.Task, *errs.XError) {
	tasks, err := svc.taskRepo.GetAll(ctx, search)
	if err != nil {
		return nil, err
	}

	mappedTasks, mapErr := svc.respMapper.Tasks(tasks)
	if mapErr != nil {
		return nil, errs.NewXError(errs.MAPPING_ERROR, "Failed to map Task data", mapErr)
	}
	return mappedTasks, nil
}

func (svc taskService) Delete(ctx *context.Context, id uint) *errs.XError {
	err := svc.taskRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
