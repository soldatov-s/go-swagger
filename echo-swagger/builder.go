package echoswagger

import (
	"os"
	"time"

	// stdlib
	"net/http"
	"net/url"
	"strings"

	// local
	"github.com/soldatov-s/go-swagger/swagger"

	// other
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

var (
	log zerolog.Logger
)

// AddToSwagger - add endpoint description to the OpenAPI Specification
func AddToSwagger(ec echo.Context) (method swagger.IMethod) {
	method = swagger.NewMethod()
	if m, ok := ec.Get("swagger").(*swagger.IMethod); ok {
		if middleMethod, ok2 := (*m).(*swagger.Method); ok2 {
			*m = middleMethod
			return middleMethod
		}
		*m = method
	}
	return method
}

// IsBuildingSwagger - mark that we build swagger description for endpoint
func IsBuildingSwagger(ec echo.Context) bool {
	return ec.Get("swagger") != nil
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
func BuildSwagger(srv *echo.Echo, swaggerPath, address string, sw swagger.ISwaggerAPI, logger *zerolog.Logger) (err error) {
	initLogger(logger)

	s := swagger.Doc{BaseAPI: *sw.(*swagger.BaseAPI)}

	log = log.With().Str("apiPath", address+s.BasePath).Logger()
	log.Info().Msg("Build swagger")

	ctx := srv.AcquireContext()

	s.Paths = make(map[string]swagger.Methods)
	s.Definitions = make(map[string]*swagger.Definition)
	for _, r := range srv.Routes() {
		var method swagger.IMethod
		ctx.Reset(&http.Request{URL: &url.URL{}}, &swagger.EmptyWriter{})
		ctx.Set("swagger", &method)
		srv.Router().Find(r.Method, r.Path, ctx)
		if !strings.HasPrefix(r.Path, s.BasePath) {
			continue
		}

		path := strings.TrimPrefix(r.Path, s.BasePath)

		if path == "" || path == "/*" {
			continue
		}

		pathElements := strings.Split(path, "/")
		path = ""
		for _, el := range pathElements {
			if el == "" {
				continue
			}
			tmp := el
			if strings.Contains(el, ":") {
				tmp = strings.Replace(el, ":", "{", 1) + "}"
			}
			path = path + "/" + tmp
		}

		err = ctx.Handler()(ctx)
		if err != nil {
			log.Error().Msg("Failed build swagger for " + path)
			return
		}

		// This code adds real method name to swagger, it helps if we use go-work echo-provider
		if v, ok := ctx.Get("swMethodName").(string); ok {
			r.Name = v
		}

		if m, ok := method.(*swagger.Method); ok {
			m.Parse(path, r.Method, s)
			m.OperationID = r.Name
		}
	}

	swagger.Register(address+s.BasePath, &s)

	srv.GET(s.BasePath+swaggerPath, Handler(
		swagger.Fill("doc.json", address+s.BasePath), // The url pointing to API definition"
	))

	return nil
}

// ExcludeFromSwagger - middleware for exclude swagger description for selected endpoint.
func ExcludeFromSwagger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ec echo.Context) error {
		if !IsBuildingSwagger(ec) {
			if err := next(ec); err != nil {
				ec.Error(err)
				return nil
			}
		}
		ec.Reset(nil, nil)
		return nil
	}
}
