package api

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humamux"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/AntonioBR9998/go-common/humamw"
	"github.com/AntonioBR9998/gymapp/api/dtos"
	"github.com/AntonioBR9998/gymapp/internal/config"
	"github.com/AntonioBR9998/gymapp/internal/core"
)

const (
	API_CONTEXT    = "/api"
	API_V1         = "/v1"
	API_V1_BASE    = API_CONTEXT + API_V1
	USERS_ENDPOINT = "/users"
	UUID_REGEX     = "[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}"
)

type api struct {
	router  http.Handler
	service core.Core
}

type Server interface {
	Router() http.Handler
}

type APIResponse[T any] struct {
	Body T `contentType:"application/json"`
}

type APIResponseWithoutBody struct{}

func NewAPI(cfg config.Config, service core.Core) Server {
	a := &api{service: service}

	log.Traceln("creating new *mux.Router")
	r := mux.NewRouter()
	log.Traceln("creating a new Subrouter for path:", API_V1_BASE)
	apiV1 := r.PathPrefix(API_V1_BASE).Subrouter()

	humaConfig := huma.DefaultConfig("gymapp-backend", "1.0.0")
	humaConfig.Servers = []*huma.Server{
		{URL: cfg.API.GetURL() + API_V1_BASE},
	}

	// This configuration allows to remove "$schema" link in JSON response
	humaConfig.CreateHooks = nil

	// V1 API Definition
	gymApi := humamux.New(apiV1, humaConfig)

	gymApi.UseMiddleware()

	// Users endpoints
	huma.Post(gymApi, USERS_ENDPOINT, a.createUser)
	huma.Put(gymApi, USERS_ENDPOINT, a.modifyUser)
	huma.Get(gymApi, USERS_ENDPOINT+"/{id:"+UUID_REGEX+"}", a.getUserByID)
	huma.Get(gymApi, USERS_ENDPOINT, a.getUserList, humamw.UseMiddlewares(
		humamw.UsePagination(humamw.PaginationOptions(humamw.SetMaxLimit(3000))),
		humamw.SetHeaderUsingCallback("Total"),
		humamw.UseFilter(
			gymApi,
			map[string]humamw.FilterDefinition{
				"id": {Type: humamw.STRING},
			},
			[]string{"id", "whenChanged"},
		),
	))
	huma.Delete(gymApi, USERS_ENDPOINT+"/{id:"+UUID_REGEX+"}", a.deleteUser)

	a.router = r
	return a
}

func (a *api) Router() http.Handler {
	return a.router
}

// Users handlers
func (a *api) createUser(ctx context.Context, req *dtos.UserBaseRequest) (*APIResponse[*dtos.UserResponseBody], error) {
	// TODO

	return nil, nil
}

func (a *api) modifyUser(ctx context.Context, req *dtos.UserBaseRequest) (*APIResponse[*dtos.UserResponseBody], error) {
	// TODO

	return nil, nil
}

func (a *api) getUserByID(ctx context.Context, req *dtos.UserBaseRequest) (*APIResponse[*dtos.UserResponseBody], error) {
	// TODO

	return nil, nil
}

func (a *api) getUserList(ctx context.Context, req *struct{}) (*APIResponse[[]*dtos.UserResponseBody], error) {
	// TODO

	return nil, nil
}

func (a *api) deleteUser(ctx context.Context, request *dtos.UserBaseRequest) (*APIResponseWithoutBody, error) {
	// TODO

	return nil, nil
}
