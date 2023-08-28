package main

import (
	"github.com/getkin/kin-openapi/openapi3"

	server "github.com/frutonanny/slug-service/internal/server/v1"
	"github.com/frutonanny/slug-service/internal/server/v1/handlers"
	createslugservice "github.com/frutonanny/slug-service/internal/services/create_slug"
	deleteslugservice "github.com/frutonanny/slug-service/internal/services/delete_slug"
	getUserSlugService "github.com/frutonanny/slug-service/internal/services/get_user_slug"
	modifyuserslugservice "github.com/frutonanny/slug-service/internal/services/modify_slug"
)

func initServer(
	addr string,
	swagger *openapi3.T,
	createSlugService *createslugservice.Service,
	deleteSlugService *deleteslugservice.Service,
	modifyUserSlugService *modifyuserslugservice.Service,
	getUserSlugService *getUserSlugService.Service,
) (*server.Server, error) {
	h := handlers.NewHandlers(
		createSlugService,
		deleteSlugService,
		modifyUserSlugService,
		getUserSlugService,
	)

	srv := server.New(
		addr,
		h,
		swagger,
	)

	return srv, nil
}
