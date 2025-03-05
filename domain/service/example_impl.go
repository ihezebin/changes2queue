package service

import (
	"context"
	"time"

	"github.com/ihezebin/jwt"

	"github.com/ihezebin/changes2queue/component/constant"
	"github.com/ihezebin/changes2queue/domain/entity"
	"github.com/ihezebin/changes2queue/domain/repository"
)

type exampleDomainServiceImpl struct {
	exampleRepository repository.ExampleRepository
}

func (svc *exampleDomainServiceImpl) IsEmailAlreadyExists(ctx context.Context, example *entity.Example) (bool, error) {

	example, err := svc.exampleRepository.FindByEmail(ctx, example.Email)
	if err != nil {
		return false, err
	}

	return example != nil, nil
}

func (svc *exampleDomainServiceImpl) IsUsernameAlreadyExists(ctx context.Context, example *entity.Example) (bool, error) {
	example, err := svc.exampleRepository.FindByUsername(ctx, example.Username)
	if err != nil {
		return false, err
	}

	return example != nil, nil
}

func (svc *exampleDomainServiceImpl) GenerateToken(example *entity.Example) (string, error) {
	token := jwt.Default(jwt.WithOwner(example.Id), jwt.WithExpire(time.Hour))
	tokenStr, err := token.Signed(constant.TokenSecret)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func NewExampleServiceImpl() ExampleDomainService {
	return &exampleDomainServiceImpl{
		exampleRepository: repository.GetExampleRepository(),
	}
}

var _ ExampleDomainService = (*exampleDomainServiceImpl)(nil)
