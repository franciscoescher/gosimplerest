package handlers

import (
	"github.com/franciscoescher/gosimplerest/repository"
	"github.com/franciscoescher/gosimplerest/resource"
	"github.com/go-playground/validator/v10"
)

type LoggerInterface interface {
	Info(args ...interface{})
	Error(args ...interface{})
}

type GetHandlerFuncParams struct {
	Resource   *resource.Resource
	Validate   *validator.Validate
	Logger     LoggerInterface
	Repository repository.RepositoryInterface
}
