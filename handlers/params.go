package handlers

import (
	"github.com/franciscoescher/gosimplerest/repository"
	"github.com/franciscoescher/gosimplerest/resource"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type GetHandlerFuncParams struct {
	Logger     *logrus.Logger
	Resource   *resource.Resource
	Validate   *validator.Validate
	Repository repository.RepositoryInterface
}
