package service

import (
	"context"
	"time"

	"github.com/ihezebin/jwt"

	"github.com/ihezebin/changes2queue/component/constant"
	"github.com/ihezebin/changes2queue/domain/entity"
	"github.com/ihezebin/changes2queue/domain/repository"
)

type exampleDomainServiceMock struct {
	exampleRepository repository.ExampleRepository
}

func (e *exampleDomainServiceMock) IsEmailAlreadyExists(ctx context.Context, example *entity.Example) (bool, error) {
	return false, nil
}

func (e *exampleDomainServiceMock) IsUsernameAlreadyExists(ctx context.Context, example *entity.Example) (bool, error) {
	return false, nil
}
func (e *exampleDomainServiceMock) GenerateToken(example *entity.Example) (string, error) {
	token := jwt.Default(jwt.WithOwner(example.Id), jwt.WithExpire(time.Hour))
	tokenStr, err := token.Signed(constant.TokenSecret)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func NewExampleServiceMock(exampleRepository repository.ExampleRepository) ExampleDomainService {
	return &exampleDomainServiceMock{
		exampleRepository: exampleRepository,
	}
}

var _ ExampleDomainService = (*exampleDomainServiceMock)(nil)
