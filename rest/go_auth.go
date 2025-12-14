package rest

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/michaelyusak/go-helper/dto"
	"github.com/michaelyusak/go-helper/entity"
	"github.com/michaelyusak/go-helper/helper"
)

type AuthRepo interface {
	ValidateToken(ctx context.Context, token string) (entity.JwtCustomClaims, error)
}

type GoAuthRepoOpt struct {
	BaseUrl string
}

type goAuthRepo struct {
	baseUrl string
	client  *resty.Client
}

func NewGoAuthRepo(opt GoAuthRepoOpt) *goAuthRepo {
	return &goAuthRepo{
		baseUrl: opt.BaseUrl,
		client:  resty.New(),
	}
}

func (r *goAuthRepo) ValidateToken(ctx context.Context, token string) (entity.JwtCustomClaims, error) {
	validateTokenRes := dto.Response{}

	resp, err := r.client.R().
		SetAuthToken(token).
		SetHeaders(helper.AuthHeadersFromContext(ctx)).
		SetResult(&validateTokenRes).
		SetError(&validateTokenRes).
		Post(r.baseUrl + "/v1/auth/validate-token")
	if err != nil {
		return entity.JwtCustomClaims{}, err
	}

	if resp.IsError() {
		return entity.JwtCustomClaims{}, fmt.Errorf("[goAuthRestRepo][resp.IsError] Error Response [status_code: %v][resp: %s]", resp.StatusCode(), string(resp.Body()))
	}

	customClaims, ok := validateTokenRes.Data.(entity.JwtCustomClaims)
	if !ok {
		return entity.JwtCustomClaims{}, fmt.Errorf("unexpected response")
	}

	return customClaims, nil
}
