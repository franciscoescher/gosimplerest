package handlers

import (
	"github.com/franciscoescher/gosimplerest/interfaces"
	"github.com/franciscoescher/gosimplerest/repository"
	"github.com/franciscoescher/gosimplerest/resource"
)

type GetHandlerFuncParams struct {
	Resource   *resource.Resource
	Validate   interfaces.Validator
	Logger     interfaces.Logger
	Repository repository.RepositoryInterface
}
