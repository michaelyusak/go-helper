package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/michaelyusak/go-helper/dto"
	"github.com/michaelyusak/go-helper/entity"
	"github.com/michaelyusak/go-helper/helper"
)

type AuthRepo interface {
	ValidateToken(ctx context.Context, token string) (entity.JwtCustomClaims, int, error)
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

func (r *goAuthRepo) ValidateToken(ctx context.Context, token string) (entity.JwtCustomClaims, int, error) {
	validateTokenRes := dto.Response[entity.JwtCustomClaims]{}

	resp, err := r.client.R().
		SetAuthToken(token).
		SetHeaders(helper.AuthHeadersFromContext(ctx)).
		SetResult(&validateTokenRes).
		SetError(&validateTokenRes).
		Post(r.baseUrl + "/v1/auth/validate-token")
	if err != nil {
		return entity.JwtCustomClaims{}, http.StatusInternalServerError, err
	}

	if resp.IsError() || validateTokenRes.StatusCode >= http.StatusBadRequest {
		return entity.JwtCustomClaims{}, validateTokenRes.StatusCode, fmt.Errorf("[goAuthRestRepo][resp.IsError] Error Response [status_code: %v][resp: %s]", resp.StatusCode(), string(resp.Body()))
	}

	return validateTokenRes.Data, validateTokenRes.StatusCode, nil
}
