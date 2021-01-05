package gorillaswagger

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/soldatov-s/go-swagger/swagger"
)

var (
	log zerolog.Logger
)

type SwaggerKey string

// AddToSwagger - add endpoint description to the OpenAPI Specification
func AddToSwagger(r *http.Request) (method swagger.IMethod) {
	method = swagger.NewMethod()
	if m, ok := r.Context().Value(SwaggerKey("swagger")).(*swagger.IMethod); ok {
		if middleMethod, ok2 := (*m).(*swagger.Method); ok2 {
			*m = middleMethod
			return middleMethod
		}
		*m = method
	}
	return method
}

// IsBuildingSwagger - mark that we build swagger description for endpoint
func IsBuildingSwagger(r *http.Request) bool {
	return r.Context().Value(SwaggerKey("swagger")) != nil
}

// TODO: For middleware not return real name of handler
func handlerName(h interface{}) string {
	t := reflect.ValueOf(h).Type()
	if t.Kind() == reflect.Func {
		return runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
	}
	return t.String()
}

// initLogger initialize logger
func initLogger(logger *zerolog.Logger) {
	if logger != nil {
		log = *logger
		return
	}
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log = zerolog.New(output).With().Timestamp().Logger()
}

// BuildSwagger - build the OpenAPI Specification in JSON format
func BuildSwagger(router *mux.Router, swaggerPath, address string, sw swagger.ISwaggerAPI, logger *zerolog.Logger) (err error) {
	initLogger(logger)

	s := swagger.Doc{BaseAPI: *sw.(*swagger.BaseAPI)}

	log = log.With().Str("apiPath", address+s.BasePath).Logger()
	log.Info().Msg("Build swagger")

	s.Paths = make(map[string]swagger.Methods)
	s.Definitions = make(map[string]*swagger.Definition)

	// Walk walks the router and all its sub-routers
	err = router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err1 := route.GetPathTemplate()
		if err1 != nil {
			return err1
		}
		if !strings.HasPrefix(path, s.BasePath) {
			return nil
		}
		path = strings.TrimPrefix(path, s.BasePath)
		listMethods, _ := route.GetMethods()
		for _, pathMethod := range listMethods {
			var method swagger.IMethod
			ctx := context.WithValue(context.Background(), SwaggerKey("swagger"), &method)
			req := &http.Request{URL: &url.URL{}}
			// Call endpoint handler
			route.GetHandler().ServeHTTP(&swagger.EmptyWriter{}, req.WithContext(ctx))
			if m, ok := method.(*swagger.Method); ok {
				m.Parse(path, pathMethod, s)
				m.OperationID = handlerName(route.GetHandler())
			}
		}

		return nil
	})

	if err != nil {
		return
	}

	swagger.Register(address+s.BasePath, &s)

	router.PathPrefix(s.BasePath + swaggerPath).Handler(Handler(
		swagger.Fill("doc.json", address+s.BasePath), // The url pointing to API definition"
	))
	return nil
}

// ExcludeFromSwagger - middleware for exclude swagger description for selected endpoint.
func ExcludeFromSwagger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !IsBuildingSwagger(r) {
			next.ServeHTTP(w, r)
		}
	})
}
