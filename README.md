# go-swagger

Create Swagger (OpenAPI 2.0) descriptions via code helpers.

## Purposes
In ideal world always at first creates scheme and describe all contracts between client and backend.
In real world more simply create project and after describe swagger scheme. And if you at first create schemes, after some iterations you could optimize generated code and lost combability with schemes.
This framework works as wizard and help doesn't forgot some fields in swagger description. 

## Features
Generates swagger description (only JSON) at Runtime by using the reflect package.
Supported generation swagger description for Echo/Gorilla routers is supported. A separate package is implemented for each router.
The framework contains two main components:
- a builder that collects and connects the swagger web interface to the specified endpoint;
- descriptor for endpoint.

## Bulder
Starts by calling function BuildSwagger. The function must be called after adding endpoint to routs. The following parameters are supported:

| Parameter                 | Description                                           |
| ------------------------- | ----------------------------------------------------- |
| router/srv (Gorilla/Echo) | router/server (Gorilla/Echo) for which builds swagger |
| swaggerPath               | the swagger-endpoint name                             |
| address                   | address on which will be showed swagger-description   |
| sw                        | swagger-structure for adding endpoint description     |

For correct writes swagger-structure used a sequence of interfaces:

```Golang
swagger.NewSwagger().SetBasePath("/api/v1").SetInfo(...)
```

| Function    | Description                            | Example |
| ----------- | -------------------------------------- | ------- |
| SetBasePath | base path to api                       | /api/v1 |
| SetInfo     | information about http-service swagger | -       |

SetInfo function accepts structure which describes information about service. In current version no hard sequence for writing information (yes, it's bad, you can write something twice). The interface for writing SetInfo contains next functions:

| Function         | Description             | Example                                                   |
| ---------------- | ----------------------- | --------------------------------------------------------- |
| SetTitle         | the service name        | Swagger Example API                                       |
| SetDescription   | the service description | This is a sample embedded Swagger-server for Echo-server. |
| SetTermOfService | URL to terms            | http://swagger.io/terms/                                  |
| SetContact       | contacts                | -                                                         |
| SetLicense       | licence                 | -                                                         |

SetContact function accepts structure which describes information about contacts. In current version no hard sequence for writing information (yes, it's bad, you can write something twice). The interface for writing SetContact contains next functions:

| Function | Description      | Example           |
| -------- | ---------------- | ----------------- |
| SetEmail | set email        | info@test.test    |
| SetName  | set support name | API Support       |
| SetURL   | set URL          | https://test.test |

SetLicense function accepts structure which describes information about licence. In current version no hard sequence for writing information (yes, it's bad, you can write something twice). The interface for writing SetLicense contains next functions:

| Function | Description      | Example                                         |
| -------- | ---------------- | ----------------------------------------------- |
| SetURL   | URL to licence   | http://www.apache.org/licenses/LICENSE-2.0.html |
| SetName  | The licence name | Apache 2.0                                      |

Example (Echo)
```Golang
func main() {
	// Create http-server
	srv := echo.New()

	// Adds routes
	v1 := srv.Group("/api/v1")
	v1.GET("/test", testGetHandler)
	v1.POST("/test", testPostHandler)
	v1.GET("/testWithoutSwagger", testWithoutSwaggerGetHandler)
	v1.GET("/testWithMiddleware", testWithMiddlewareGetHandler, Middleware)
	v1.GET("/testArray", testArrayGetHandler)
	v1.GET("/testArrayOfStruct", testArrayOfStructGetHandler)
	v1.GET("/testParamPath/:id/:name", testParamPathGetHandler)

	// Build swagger
	err := echoSwagger.BuildSwagger(
		srv,
		"/swagger/*",
		":1323",
		swagger.NewSwagger().
			SetBasePath("/api/v1/").
			SetInfo(swagger.NewInfo().
				SetTitle("Swagger Example API v1").
				SetDescription("This is a sample embedded Swagger-server for Echo-server.").
				SetTermOfService("http://swagger.io/terms/").
				SetContact(swagger.NewContact()).
				SetLicense(swagger.NewLicense())),
		nil,
	)

	if err != nil {
		return
	}

	srv.HideBanner = true
	srv.Logger.Fatal(srv.Start(":1323"))
}
```

Example (Gorilla)
```Golang
func main() {
	r := mux.NewRouter()
	v1 := r.PathPrefix("/api/v1").Subrouter()
	v1.HandleFunc("/test", testGetHandler).Methods("GET")
	v1.HandleFunc("/test", testPostHandler).Methods("POST")
	v1.HandleFunc("/testWithoutSwagger", testWithoutSwaggerGetHandler).Methods("GET")
	v1.Handle("/testWithMiddleware", Middleware(http.HandlerFunc(testWithMiddlewareGetHandler))).Methods("GET")
	v1.HandleFunc("/testArray", testArrayGetHandler).Methods("GET")
	v1.HandleFunc("/testArrayOfStruct", testArrayOfStructGetHandler).Methods("GET")
	v1.HandleFunc("/testParamPath/{id}/{name}", testParamPathGetHandler).Methods("GET")

	// Build swagger
	err := gorillaSwagger.BuildSwagger(
		r,
		"/swagger/",
		":1323",
		swagger.NewSwagger().
			SetBasePath("/api/v1/").
			SetInfo(swagger.NewInfo().
				SetTitle("Swagger Example API v1").
				SetDescription("This is a sample embedded Swagger-server for Gorilla-server.").
				SetTermOfService("http://swagger.io/terms/").
				SetContact(swagger.NewContact()).
				SetLicense(swagger.NewLicense())),
		nil,
	)

	if err != nil {
		return
	}

	srv := &http.Server{
		Handler: r,
		Addr:    ":1323",
	}

	log.Fatal(srv.ListenAndServe())
}
```

## Descriptor for endpoint
The descriptor is implemented as sequence interfaces, which maximally excludes incorrect endpoint descriptions.
The descriptor is called by the AddToSwagger function. A router context object is passed to it.
For Echo - echo.Context, for Gorilla - *http.Request.
Function AddToSwagger retruns interface with next functions:

| Function    | Description                 | Example          |
| ----------- | --------------------------- | ---------------- |
| SetConsumes | set consumed data           | application/json |
| SetProduces | set produced data           | application/json |
| AddResponse | add responses from endpoint | -                |

SetConsumes accepts array of strings with supported consumed data formats and returns interface with next functions:

| Function       | Description              | Example          |
| -------------- | ------------------------ | ---------------- |
| SetDescription | set endpoint description | Test PostHandler |
| SetProduces    | set produced format      | application/json |

SetProduces accepts array of strings with supported produced data formats and returns interface with next functions:

| Function       | Description              | Example                      |
| -------------- | ------------------------ | ---------------------------- |
| SetDescription | set endpoint description | Test simply POST PostHandler |

SetDescription returns interface with next functions:

| Function   | Description                        | Example      |
| ---------- | ---------------------------------- | ------------ |
| SetSummary | set short description for endpoint | Test handler |

SetSummary returns interface with next functions:

| Function             | Description                            | Example |
| -------------------- | -------------------------------------- | ------- |
| AddResponse          | add response from endpoint             | -       |
| AddInBodyParameter   | add description of in-body parameter   | -       |
| AddInPathParameter   | add description of in-path parameter   | -       |
| AddInQueryParameter  | add description of in-query parameter  | -       |
| AddInHeaderParameter | add description of in-header parameter | -       |
| AddInCookieParameter | add description of in-cookie parameter | -       |

* The function AddInBodyParameter intended to describe _in body_ parameter. According to the swagger 2.0 specification, there can be only one for a specific endpoint.
* The functions AddInPathParameter anf AddInQueryParameter intended to describe _in path_, _in query_, _in header_, _in cookie_ parameters. According to the swagger 2.0 specification, there may be several of them. WARNING! Cookie not supported in Swagger UI, adds it only for generation scheme. Therefore the returned by AddInBodyParameter/AddInPathParameter/AddInQueryParameter functions interface contains only next functions:

| Function             | Description                            | Example |
| -------------------- | -------------------------------------- | ------- |
| AddResponse          | add response from endpoint             | -       |
| AddInPathParameter   | add description of in-path parameter   | -       |
| AddInQueryParameter  | add description of in-query parameter  | -       |
| AddInHeaderParameter | add description of in-header parameter | -       |
| AddInCookieParameter | add description of in-cookie parameter | -       |

AddResponse returns interface which have only function AddResponse.

For processing descriptor only once use function IsBuildingSwagger.

### AddInBodyParameter
Accepts the name of parameter, its description, scheme and the flag required or not parameter.

### AddInPathParameter
Accepts the name of parameter, its description, type of parameter.

### AddInQueryParameter
Accepts the name of parameter, its description, type of parameter and the flag required or not parameter.

### AddInHeaderParameter
Accepts the name of parameter, its description, type of parameter and the flag required or not parameter.

### AddResponse
Accepts response code, its description, scheme.

## Examples of using descriptor for endpoint
Example (Gorilla)

```Golang
	// Swagger
	if gorillaSwagger.IsBuildingSwagger(r) {
		gorillaSwagger.AddToSwagger(r).
			SetConsumes("application/json").
			SetProduces("application/json").
			SetDescription("Test PostHandler").
			SetSummary("Test simply POST handler").
			AddInBodyParameter("some_id", "Some ID", &TestStruct{}, true).
			AddResponse(200, "Test", &TestStructWithInterfaceField{InterfaceField: TestStruct{Name: "dsgsdgs"}}).
			AddResponse(402, "Test", &TestStructWithStructField{})
		return
	}
```

Example (Echo)

```Golang
	// Swagger
	if echoSwagger.IsBuildingSwagger(ec) {
		echoSwagger.AddToSwagger(ec).
			SetConsumes("application/json").
			SetProduces("application/json").
			SetDescription("Test PostHandler").
			SetSummary("Test simply POST handler").
			AddInBodyParameter("some_id", "Some ID", &TestStruct{}, false).
			AddResponse(200, "Test", &TestStructWithInterfaceField{InterfaceField: TestStruct{Name: "dsgsdgs"}}).
			AddResponse(402, "Test", &TestStructWithStructField{})
		return nil
	}
```

## Supported schemes
* Nested structs:
```Golang
type TestStruct struct {
	Name   string `json:"name"`
	Conter int    `json:"counter"`
}

type TestStructWithStructField struct {
	StructField TestStruct `json:"structField"`
}
```

* Nested slices of structs:

```Golang
type TestEeee struct {
	Name   string `json:"Name"`
	Conter int    `json:"conter"`
}

type TestStructWithArrayOfStructs struct {
	ArrayField []TestEeee `json:"arrayField"`
}
```

* Nested slices of siple types.

* Nested interfaces:

```Golang
type TestStructWithInterfaceField struct {
	Str            string      `json:"str"`
	InterfaceField interface{} `json:"interfaceField"`
}
```

# Supporting Middleware
go-swagger works with middlewares correctly.
If Middleware is available, they will be processed correctly and will not be counted as a separate method at the endpoint.
You can specify a swagger description for the Middleware response if there is a condition under which the request is not forwarded further to the endpoint

Example (Gorilla):
```Golang
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Swagger
		if gorillaSwagger.IsBuildingSwagger(r) {
			gorillaSwagger.AddToSwagger(r).
				AddResponse(405, "Test middleware response", nil)
		}

		// Main code of handler
		if r.URL.Query().Get("some_id") == "fail" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	})
}
```

Example (Echo):
```Golang
func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ec echo.Context) error {
		// Swagger
		if echoSwagger.IsBuildingSwagger(ec) {
			echoSwagger.AddToSwagger(ec).
				AddResponse(405, "Test middleware response", nil)
		}

		// Main code of handler
		if ec.QueryParam("some_id") == "fail" {
			ec.Response().WriteHeader(http.StatusMethodNotAllowed)
			return nil
		}
		if err := next(ec); err != nil {
			ec.Error(err)
		}
		return nil
	}
}
```

* Exclude endpoint from swagger-description
Use middleware ExcludeFromSwagger for it.
If no descriptor for endpoint this endpoint will not included to swagger-description, but you can see some errors in log when swagger-description will be build.

# Multiple swagger-descriptions at one address
You can create multiple swagger-descriptions at one address. E.x. for "api/v1" for "api/2".

# Types substitution
During generated swagger-description you can replace one type to another.
For example, for save memory useful to use `int` instead `string` constants. But in contract will be `string`. If you do not make a substitution of the type, then the swagger description will contain an int, which is misleading. To avoid the "swagtype" tag, if it is present, the parser will replace the real type with capabilities.

```Golang
type TestStructWithStructField struct {
  TestIntAsString  int                    `json:"testIntAsString" swagtype:"string"`
}
```
In swagger description TestIntAsString will be as string.

# Enum
To write multiple possible values you mus use the tag "swagenum".
```Golang
type TestStructWithStructField struct {
  TestStringEnum string `json:"testStringEnum" swagenum:"testEnum1,testEnum2,testEnum3"`
}
```

# Examples
You can find examples of http-service in [/example/](/example/ "/example/"). After run open in browser http://localhost:1323/api/v1/swagger/index.html

