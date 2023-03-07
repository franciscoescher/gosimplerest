package handlers

import (
	"github.com/franciscoescher/gosimplerest/logger"
	"github.com/franciscoescher/gosimplerest/repository"
	"github.com/franciscoescher/gosimplerest/resource"
	"github.com/franciscoescher/gosimplerest/validator"
)

type GetHandlerFuncParams struct {
	Resource   *resource.Resource
	Validate   validator.Validator
	Logger     logger.Logger
	Repository repository.RepositoryInterface
}
