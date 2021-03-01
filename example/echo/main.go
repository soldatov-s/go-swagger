package main

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/soldatov-s/go-swagger/echo-swagger"
	"github.com/soldatov-s/go-swagger/swagger"
)

func buildsrv1() {
	// Create http-server
	srv1 := echo.New()

	// Adds routes
	srv1v1 := srv1.Group("/api/v1")
	srv1v1.GET("/test", testGetHandler)
	srv1v1.GET("/testHeader", testGetHeader)
	srv1v1.GET("/testCookie", testGetCookie)
	srv1v1.POST("/test", testPostHandler)
	srv1v1.GET("/testWithoutSwagger", testWithoutSwaggerGetHandler)
	srv1v1.GET("/testWithMiddleware", testWithMiddlewareGetHandler, Middleware)
	srv1v1.GET("/testArray", testArrayGetHandler)
	srv1v1.GET("/testArrayOfStruct", testArrayOfStructGetHandler)
	srv1v1.GET("/testParamPath/:id/:name", testParamPathGetHandler)

	srv1v2 := srv1.Group("/api/v2")
	srv1v2.GET("/test", testGetHandler)
	srv1v2.POST("/test", testPostHandler)
	srv1v2.POST("/testfile", testPostFileHandler)
	srv1v2.GET("/testfile", testGetFileHandler)
	srv1v2.GET("/testWithoutSwagger", testWithoutSwaggerGetHandler)
	srv1v2.GET("/testArrayOfStruct", testArrayOfStructGetHandler)

	// Build swagger
	err := echoSwagger.BuildSwagger(
		srv1,
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

	// Build swagger
	err = echoSwagger.BuildSwagger(
		srv1,
		"/swagger/*",
		":1323",
		swagger.NewSwagger().
			SetBasePath("/api/v2").
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

	go func() {
		srv1.HideBanner = true
		srv1.Logger.Fatal(srv1.Start(":1323"))
	}()
}

func buildsrv2() {
	// Create http-server
	srv2 := echo.New()

	// Adds routes
	srv2v1 := srv2.Group("/api/v1")
	srv2v1.GET("/test", testGetHandler)

	srv2v2 := srv2.Group("/api/v2")
	srv2v2.GET("/test", testGetHandler)
	srv2v2.POST("/test", testPostHandler)

	// Build swagger
	err := echoSwagger.BuildSwagger(
		srv2,
		"/swagger/*",
		":1324",
		swagger.NewSwagger().
			SetBasePath("/api/v1").
			SetInfo(swagger.NewInfo().
				SetTitle("Swagger Example API v2").
				SetDescription("This is a sample embedded Swagger-server for Echo-server.").
				SetTermOfService("http://swagger.io/terms/").
				SetContact(swagger.NewContact()).
				SetLicense(swagger.NewLicense())),
		nil,
	)

	if err != nil {
		return
	}

	// Build swagger
	err = echoSwagger.BuildSwagger(
		srv2,
		"/swagger/*",
		":1324",
		swagger.NewSwagger().
			SetBasePath("/api/v2").
			SetInfo(swagger.NewInfo().
				SetTitle("Swagger Example API v2").
				SetDescription("This is a sample embedded Swagger-server for Echo-server.").
				SetTermOfService("http://swagger.io/terms/").
				SetContact(swagger.NewContact()).
				SetLicense(swagger.NewLicense())),
		nil,
	)

	if err != nil {
		return
	}

	srv2.HideBanner = true
	srv2.Logger.Fatal(srv2.Start(":1324"))
}

func main() {
	buildsrv1()
	buildsrv2()
}

type TestHeader struct {
	Header string `json:"header"`
}

type TestCookie struct {
	Cookie string `json:"cookie"`
}

type TestStruct struct {
	Name   string    `json:"name"`
	Conter int       `json:"counter"`
	Time   time.Time `json:"time"`
}

type TestStructWithInterfaceField struct {
	Str            string      `json:"str"`
	InterfaceField interface{} `json:"interfaceField"`
}

type TestStructWithInterfaceField2 struct {
	Str            string      `json:"str"`
	InterfaceField interface{} `json:"interfaceField"`
}

type ArrayOfTestStruct []TestStruct

type TestStructWithInterfaceField3 struct {
	Str            string      `json:"str"`
	InterfaceField interface{} `json:"interfaceField"`
}

func testGetHandler(ec echo.Context) error {
	// Swagger
	if echoSwagger.IsBuildingSwagger(ec) {
		echoSwagger.AddToSwagger(ec).
			SetProduces("application/json").
			SetDescription("Test GetHandler").
			SetSummary("Test simply GET handler").
			AddInQueryParameter("string_id", "Some string ID", reflect.String, true).
			AddInQueryParameter("int_id", "Some int ID", reflect.Int64, false).
			AddResponse(http.StatusOK, "Test", &TestStruct{})
		return nil
	}

	// Main code of handler
	intID, _ := strconv.Atoi(ec.QueryParam("int_id"))
	return ec.JSON(http.StatusOK, TestStruct{
		Name:   ec.QueryParam("string_id"),
		Conter: intID,
	})
}

func testGetHeader(ec echo.Context) error {
	// Swagger
	if echoSwagger.IsBuildingSwagger(ec) {
		echoSwagger.AddToSwagger(ec).
			SetProduces("application/json").
			SetDescription("Test GetHeaderHandler").
			SetSummary("Test GET with header handler").
			AddInHeaderParameter("Authorization", "Authorization header", reflect.String, true).
			AddResponse(http.StatusOK, "Test", &TestHeader{})
		return nil
	}

	// Main code of handler
	token := ec.Request().Header.Get("Authorization")
	return ec.JSON(http.StatusOK, TestHeader{
		Header: token,
	})
}

func testGetCookie(ec echo.Context) error {
	// Swagger
	if echoSwagger.IsBuildingSwagger(ec) {
		echoSwagger.AddToSwagger(ec).
			SetProduces("application/json").
			SetDescription("Test GetCookieHandler").
			SetSummary("Test GET with cookie handler").
			AddInCookieParameter("token-cookie", "token cookie", reflect.String, false).
			AddResponse(http.StatusOK, "Test", &TestCookie{})
		return nil
	}

	// Main code of handler
	cookie, err := ec.Cookie("token-cookie")
	if err != nil {
		fmt.Println(err)
		return err
	}

	return ec.JSON(http.StatusOK, TestCookie{
		Cookie: cookie.Value,
	})
}

type TestStructWithStructField struct {
	StructField      TestStruct             `json:"structField"`
	TestSkip         int                    `json:"-"`
	TestMapString    map[string]string      `json:"mapString"`
	TestMapInterface map[string]interface{} `json:"mapInterface"`
	TestIntAsString  int                    `json:"testIntAsString" swagtype:"string"`
}

type TestStructWithMap struct {
	TestMap map[string]string `json:"mapString"`
}

type TestStructWithSlicePtr struct {
	TestSliceStruct []*TestStruct `json:"slicePtrStruct"`
}

type TestStructWithPtrField struct {
	StructField     *TestStruct `json:"ptrStructField"`
	TestSkip        *int        `json:"-"`
	TestNoJSONTag   *int
	TestMapString   map[string]*string        `json:"mapPtrString"`
	TestMapStruct   map[string]*TestStruct    `json:"mapPtrStruct"`
	TestSliceString []*string                 `json:"slicePtrString"`
	TestSliceStruct []*TestStructWithSlicePtr `json:"slicePtrStruct"`
	UserUUID        *uuid.UUID                `json:"user_uuid"`
	ExpireAt        *time.Time                `json:"expire_at"`
}

type TestStructWithAnonymousField struct {
	TestStruct
	TestString1 string `json:"string1"`
}

type TestStructWithAnonymousField2 struct {
	TestString2    string `json:"string2"`
	TestStringEnum string `json:"testStringEnum" swagenum:"testEnum1,testEnum2,testEnum3"`
	TestStructWithAnonymousField
}

func testPostHandler(ec echo.Context) error {
	// Swagger
	if echoSwagger.IsBuildingSwagger(ec) {
		echoSwagger.AddToSwagger(ec).
			SetConsumes("application/json").
			SetProduces("application/json").
			SetDescription("Test PostHandler").
			SetSummary("Test simply POST handler").
			AddInBodyParameter("some_id", "Some ID", &TestStruct{}, false).
			AddResponse(http.StatusOK, "Test", &TestStructWithInterfaceField{InterfaceField: TestStruct{Name: "dsgsdgs"}}).
			AddResponse(http.StatusCreated, "Test", &TestStructWithInterfaceField2{InterfaceField: 0}).
			AddResponse(http.StatusAccepted, "Test", &TestStructWithInterfaceField3{InterfaceField: ArrayOfTestStruct{}}).
			AddResponse(http.StatusNonAuthoritativeInfo, "Test", &ArrayOfTestStruct{}).
			AddResponse(http.StatusNoContent, "Test", &TestStructWithMap{}).
			AddResponse(http.StatusResetContent, "Test", &TestStructWithPtrField{}).
			AddResponse(http.StatusPartialContent, "Test", make(map[string]TestStruct)).
			AddResponse(http.StatusMultiStatus, "Test", make(map[string]string)).
			AddResponse(http.StatusAlreadyReported, "Test", &TestStructWithAnonymousField2{}).
			AddResponse(http.StatusPaymentRequired, "Test", &TestStructWithStructField{}).
			AddResponse(http.StatusPaymentRequired, "Test 2", &TestStructWithStructField{})
		return nil
	}

	// Main code of handler
	var test TestStruct
	err := ec.Bind(&test)
	if err != nil {
		return err
	}

	return ec.JSON(http.StatusOK, TestStructWithInterfaceField{Str: test.Name, InterfaceField: TestStructWithInterfaceField{}})
}

func testPostFileHandler(ec echo.Context) error {
	// Swagger
	if echoSwagger.IsBuildingSwagger(ec) {
		echoSwagger.AddToSwagger(ec).
			SetConsumes("multipart/form-data").
			SetDescription("Test PostFileHandler").
			SetSummary("Test simply POST file handler").
			AddInFileParameter("uploadFile", "Upload file").
			AddResponse(http.StatusOK, "Test", &TestStruct{})

		return nil
	}

	// Main code of handler
	multipartForm, err := ec.MultipartForm()
	if err != nil {
		return err
	}

	fileName := ""
	for _, fileHeaders := range multipartForm.File {
		for _, fileHeader := range fileHeaders {
			fileName = fileHeader.Filename
		}
	}

	return ec.JSON(http.StatusOK, TestStruct{Name: fileName, Conter: 1, Time: time.Now()})
}

func testGetFileHandler(ec echo.Context) error {
	// Swagger
	if echoSwagger.IsBuildingSwagger(ec) {
		echoSwagger.AddToSwagger(ec).
			SetProduces("multipart/form-data").
			SetDescription("Test GetFileHandler").
			SetSummary("Test simply GET file handler").
			AddFileResponse(http.StatusOK, "Test get file")

		return nil
	}

	return ec.File("./testfile.txt")
}

func testWithoutSwaggerGetHandler(ec echo.Context) error {
	return ec.JSON(http.StatusOK, TestStruct{Name: "testExcludeSwaggerGetHandler"})
}

func testWithMiddlewareGetHandler(ec echo.Context) error {
	// Swagger
	if echoSwagger.IsBuildingSwagger(ec) {
		echoSwagger.AddToSwagger(ec).
			SetProduces("application/json").
			SetDescription("Test WithMiddlewareGetHandler").
			SetSummary("Test Get handler with middleware").
			AddInQueryParameter("some_id", "Some ID", reflect.String, true).
			AddResponse(http.StatusOK, "Test", &TestStruct{})
		return nil
	}

	// Main code of handler
	return ec.JSON(http.StatusOK, TestStruct{Name: ec.QueryParam("some_id")})
}

func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ec echo.Context) error {
		// Swagger
		if echoSwagger.IsBuildingSwagger(ec) {
			echoSwagger.AddToSwagger(ec).
				AddResponse(http.StatusMethodNotAllowed, "Test middleware response", nil)
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

type TestStructWithArray struct {
	ArrayField []string `json:"arrayField"`
}

func testArrayGetHandler(ec echo.Context) error {
	// Swagger
	if echoSwagger.IsBuildingSwagger(ec) {
		echoSwagger.AddToSwagger(ec).
			SetProduces("application/json").
			SetDescription("Test ArrayGetHandler").
			SetSummary("Get array").
			AddResponse(http.StatusOK, "Test", &TestStructWithArray{})
		return nil
	}

	// Main code of handler
	return ec.JSON(http.StatusOK, TestStructWithArray{ArrayField: []string{"element1", "element2", "element3"}})
}

type TestEeee struct {
	Name   string `json:"Name"`
	Conter int    `json:"conter"`
}

type TestStructWithArrayOfStructs struct {
	ArrayField []TestEeee `json:"arrayField"`
}

func testArrayOfStructGetHandler(ec echo.Context) error {
	// Swagger
	if echoSwagger.IsBuildingSwagger(ec) {
		echoSwagger.AddToSwagger(ec).
			SetProduces("application/json").
			SetDescription("Test ArrayOfStruct").
			SetSummary("Get array of structs").
			AddResponse(http.StatusOK, "Test", &TestStructWithArrayOfStructs{})
		return nil
	}

	// Main code of handler
	return ec.JSON(http.StatusOK, TestStructWithArrayOfStructs{ArrayField: []TestEeee{
		{Name: "element1", Conter: 1},
		{Name: "element2", Conter: 2},
		{Name: "element3", Conter: 3},
	}})
}

func testParamPathGetHandler(ec echo.Context) error {
	// Swagger
	if echoSwagger.IsBuildingSwagger(ec) {
		echoSwagger.AddToSwagger(ec).
			SetProduces("text/plain").
			SetDescription("Test ParamPathGetHandler").
			SetSummary("Test Param in Path GET handler").
			AddInPathParameter("id", "Some id", reflect.Int16).
			AddInPathParameter("name", "Some name", reflect.String).
			AddResponse(http.StatusOK, "Test", reflect.Int64)

		return nil
	}

	// Main code of handler
	return ec.String(http.StatusOK, ec.Param("id"))
}
