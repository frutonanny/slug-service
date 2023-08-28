// Package v1 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package v1

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"
	"time"

	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

// AddSlugsData defines model for AddSlugsData.
type AddSlugsData struct {
	// Название сегмента.
	Name Slug `json:"name"`

	// Время автоматического удаления пользователя из сегмента в формате RFC3339.
	Ttl *time.Time `json:"ttl,omitempty"`
}

// CreateSlugRequest defines model for CreateSlugRequest.
type CreateSlugRequest struct {
	// Название сегмента.
	Name Slug `json:"name"`

	// Опции сегмента.
	Options *struct {
		Percent *int `json:"percent,omitempty"`
	} `json:"options,omitempty"`
}

// CreateSlugResponse defines model for CreateSlugResponse.
type CreateSlugResponse struct {
	Data  *map[string]interface{} `json:"data,omitempty"`
	Error *Error                  `json:"error,omitempty"`
}

// DeleteSlugRequest defines model for DeleteSlugRequest.
type DeleteSlugRequest struct {
	// Название сегмента.
	Name Slug `json:"name"`
}

// DeleteSlugResponse defines model for DeleteSlugResponse.
type DeleteSlugResponse struct {
	Data  *map[string]interface{} `json:"data,omitempty"`
	Error *Error                  `json:"error,omitempty"`
}

// Error defines model for Error.
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// GetUserSlugsData defines model for GetUserSlugsData.
type GetUserSlugsData struct {
	Slugs []Slug `json:"slugs"`
}

// GetUserSlugsRequest defines model for GetUserSlugsRequest.
type GetUserSlugsRequest struct {
	// Идентификатор пользователя.
	UserID UserID `json:"userID"`
}

// GetUserSlugsResponse defines model for GetUserSlugsResponse.
type GetUserSlugsResponse struct {
	Data  *GetUserSlugsData `json:"data,omitempty"`
	Error *Error            `json:"error,omitempty"`
}

// ModifyUserSlugsRequest defines model for ModifyUserSlugsRequest.
type ModifyUserSlugsRequest struct {
	Add    *[]AddSlugsData `json:"add,omitempty"`
	Delete *[]Slug         `json:"delete,omitempty"`

	// Идентификатор пользователя.
	UserID UserID `json:"userID"`
}

// ModifyUserSlugsResponse defines model for ModifyUserSlugsResponse.
type ModifyUserSlugsResponse struct {
	Data  *map[string]interface{} `json:"data,omitempty"`
	Error *Error                  `json:"error,omitempty"`
}

// Название сегмента.
type Slug = string

// Идентификатор пользователя.
type UserID = openapi_types.UUID

// PostCreateSlugJSONBody defines parameters for PostCreateSlug.
type PostCreateSlugJSONBody = CreateSlugRequest

// PostDeleteSlugJSONBody defines parameters for PostDeleteSlug.
type PostDeleteSlugJSONBody = DeleteSlugRequest

// PostGetUserSlugsJSONBody defines parameters for PostGetUserSlugs.
type PostGetUserSlugsJSONBody = GetUserSlugsRequest

// PostModifyUserSlugsJSONBody defines parameters for PostModifyUserSlugs.
type PostModifyUserSlugsJSONBody = ModifyUserSlugsRequest

// PostCreateSlugJSONRequestBody defines body for PostCreateSlug for application/json ContentType.
type PostCreateSlugJSONRequestBody = PostCreateSlugJSONBody

// PostDeleteSlugJSONRequestBody defines body for PostDeleteSlug for application/json ContentType.
type PostDeleteSlugJSONRequestBody = PostDeleteSlugJSONBody

// PostGetUserSlugsJSONRequestBody defines body for PostGetUserSlugs for application/json ContentType.
type PostGetUserSlugsJSONRequestBody = PostGetUserSlugsJSONBody

// PostModifyUserSlugsJSONRequestBody defines body for PostModifyUserSlugs for application/json ContentType.
type PostModifyUserSlugsJSONRequestBody = PostModifyUserSlugsJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /createSlug)
	PostCreateSlug(ctx echo.Context) error

	// (POST /deleteSlug)
	PostDeleteSlug(ctx echo.Context) error

	// (POST /getUserSlugs)
	PostGetUserSlugs(ctx echo.Context) error

	// (POST /modifyUserSlugs)
	PostModifyUserSlugs(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// PostCreateSlug converts echo context to params.
func (w *ServerInterfaceWrapper) PostCreateSlug(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostCreateSlug(ctx)
	return err
}

// PostDeleteSlug converts echo context to params.
func (w *ServerInterfaceWrapper) PostDeleteSlug(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostDeleteSlug(ctx)
	return err
}

// PostGetUserSlugs converts echo context to params.
func (w *ServerInterfaceWrapper) PostGetUserSlugs(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostGetUserSlugs(ctx)
	return err
}

// PostModifyUserSlugs converts echo context to params.
func (w *ServerInterfaceWrapper) PostModifyUserSlugs(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostModifyUserSlugs(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST(baseURL+"/createSlug", wrapper.PostCreateSlug)
	router.POST(baseURL+"/deleteSlug", wrapper.PostDeleteSlug)
	router.POST(baseURL+"/getUserSlugs", wrapper.PostGetUserSlugs)
	router.POST(baseURL+"/modifyUserSlugs", wrapper.PostModifyUserSlugs)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/7xX227bRhD9FWHaR8qkbo7DN9d2AwMNEkRxHhoYAUOOJQbiJbsrI0YgwHGC9qEBCvSp",
	"D72gf6AqFizELvMLwz8qdkndSOrSJjJgwJI4uzNzzsyc4WuwAy8MfPQFB/M1cLuNnqU+7jpOs9Nt8X1L",
	"WPJ7yIIQmXBRPfUtD+X/rxmegAlf6dN79PQSXR6HngZCdKSpg9xmbijcwAcT6Jf4nIZ0E/9coj4N4guK",
	"6Ib68QWN4h9pGL+hjxTRB4pK8Vu6pD5d05D+oZG0/0QRXcfv6YoiGqgzQ7qWD0Z0VYrf0JA+0I00jy+o",
	"X6JBKX5HUXye3j8sPfp2r1ar3d0CDfCV5YUdBBOqRrVWNnbK1TuPK4ZpyL/vQYOTgHmWABMcS2BZuB6C",
	"BuIslEe4YK7fgl5PA4Yvuy5DB8ynCTbHE6vg+Qu0hQRij6ElUMLyCF92kYvPwzVQWPICbP+gT/EPNKJR",
	"Dg2Z87zHEJmNvgrFs165XtcDs2IYGniun36bpOL6AlvIVMaZ7AoxyOTMw8DnmE/aSUssBxgyFrBVcBwo",
	"o56MYB87+MUAXpTRrJONZ3Qwtp+/3w4c5dVz/e/Qb4n2LEvjstTAQ86t1mrLTKrjY1rip6iU76E44siW",
	"TAguH8kPrkCPrz0rEkcWY9ZZLq7kylXhLKS+y5Ed7q+K5CixyvpODx/nnK0qgWW+cij+jwq5Hzjuydnq",
	"7C3HWZuNudmfY0UDR7XAZ5KrfTFCchBsvC1VRvm5+zv16UpKklQqGhZO36ni7D45fPzg2ZMHh3sHz+4f",
	"NJu79w6aoMkpPO7UaqOhzXZuvaDHjyYgZoL5lS5Tv6P4HY3oo1K/KD5fKKDz4TUaBu7UDaOM1bvPy/WK",
	"Uy9bdyrb5Xp9e7vRqNcNwzBmBbLbdZ0Cbexp4PongYLdFepmiV6piezUtbG0+/AQNDhFxpO4TyuJtqFv",
	"hS6YUNsytmpStSzRViTq9kRSFMNBUu2Z7H+jodopLiUJEV2pHSLZH4pYkVViybOHDpjwMOBiKlyQFB5y",
	"8U3gnCXj1xepZlph2HFtdVR/waXv8Ra1qpzy20BvvsYF66L6ISlnlXzVMDYSQNoxKoIMkn/NojWH5pa0",
	"72mgOxNJXI+R7Ea3HiNT4d0QI/n14ZYZKVgt1mBkBs0JI60ZZVmLk3QmvJXLdzErEQ2WTo48YbPytiHK",
	"ioT/lkkrXAdW0hb/tATMlERvXtfW4/GSIvpbvlJN+0v/D69QA73oLSqiQTHFGendEMsLdpxbJnrRmlHE",
	"9Z+F+L7P0TNHzVL0e0lRcGRSK8F8mi2CfTzFThB66ItSYgUadFkHTGgLEZq63glsq9MOuDB3jJ2KLpX2",
	"uPdvAAAA//+YotsCAxAAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
