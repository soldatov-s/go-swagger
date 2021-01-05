package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	gorillaSwagger "github.com/soldatov-s/go-swagger/gorilla-swagger"
	"github.com/soldatov-s/go-swagger/swagger"
)

func buildsrv1() {
	router1 := mux.NewRouter()

	router1v1 := router1.PathPrefix("/api/v1").Subrouter()
	router1v1.HandleFunc("/test", testGetHandler).Methods("GET")
	router1v1.HandleFunc("/testHeader", testGetHeader).Methods("GET")
	router1v1.HandleFunc("/testCookie", testGetCookie).Methods("GET")
	router1v1.HandleFunc("/test", testPostHandler).Methods("POST")
	router1v1.HandleFunc("/testWithoutSwagger", testWithoutSwaggerGetHandler).Methods("GET")
	router1v1.Handle("/testWithMiddleware", Middleware(http.HandlerFunc(testWithMiddlewareGetHandler))).Methods("GET")
	router1v1.HandleFunc("/testArray", testArrayGetHandler).Methods("GET")
	router1v1.HandleFunc("/testArrayOfStruct", testArrayOfStructGetHandler).Methods("GET")
	router1v1.HandleFunc("/testParamPath/{id}/{name}", testParamPathGetHandler).Methods("GET")

	router1v2 := router1.PathPrefix("/api/v2").Subrouter()
	router1v2.HandleFunc("/test", testGetHandler).Methods("GET")
	router1v2.HandleFunc("/test", testPostHandler).Methods("POST")
	router1v2.HandleFunc("/testWithoutSwagger", testWithoutSwaggerGetHandler).Methods("GET")
	router1v2.HandleFunc("/testArrayOfStruct", testArrayOfStructGetHandler).Methods("GET")

	// Build swagger
	err := gorillaSwagger.BuildSwagger(
		router1,
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

	// Build swagger
	err = gorillaSwagger.BuildSwagger(
		router1,
		"/swagger/",
		":1323",
		swagger.NewSwagger().
			SetBasePath("/api/v2").
			SetInfo(swagger.NewInfo().
				SetTitle("Swagger Example API v2").
				SetDescription("This is a sample embedded Swagger-server for Gorilla-server.").
				SetTermOfService("http://swagger.io/terms/").
				SetContact(swagger.NewContact()).
				SetLicense(swagger.NewLicense())),
		nil,
	)

	if err != nil {
		return
	}

	srv1 := &http.Server{
		Handler: router1,
		Addr:    ":1323",
	}

	go func() {
		log.Fatal(srv1.ListenAndServe())
	}()
}

func buildsrv2() {
	router2 := mux.NewRouter()

	router2v1 := router2.PathPrefix("/api/v1").Subrouter()
	router2v1.HandleFunc("/test", testGetHandler).Methods("GET")

	router2v2 := router2.PathPrefix("/api/v2").Subrouter()
	router2v2.HandleFunc("/test", testGetHandler).Methods("GET")
	router2v2.HandleFunc("/test", testPostHandler).Methods("POST")

	// Build swagger
	err := gorillaSwagger.BuildSwagger(
		router2,
		"/swagger/",
		":1324",
		swagger.NewSwagger().
			SetBasePath("/api/v1").
			SetInfo(swagger.NewInfo().
				SetTitle("Swagger Example API v2").
				SetDescription("This is a sample embedded Swagger-server for Gorilla-server.").
				SetTermOfService("http://swagger.io/terms/").
				SetContact(swagger.NewContact()).
				SetLicense(swagger.NewLicense())),
		nil,
	)

	if err != nil {
		return
	}

	// Build swagger
	err = gorillaSwagger.BuildSwagger(
		router2,
		"/swagger/",
		":1324",
		swagger.NewSwagger().
			SetBasePath("/api/v2").
			SetInfo(swagger.NewInfo().
				SetTitle("Swagger Example API v2").
				SetDescription("This is a sample embedded Swagger-server for Gorilla-server.").
				SetTermOfService("http://swagger.io/terms/").
				SetContact(swagger.NewContact()).
				SetLicense(swagger.NewLicense())),
		nil,
	)

	if err != nil {
		return
	}

	srv2 := &http.Server{
		Handler: router2,
		Addr:    ":1324",
	}

	log.Fatal(srv2.ListenAndServe())
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

func testGetHandler(w http.ResponseWriter, r *http.Request) {
	// Swagger
	if gorillaSwagger.IsBuildingSwagger(r) {
		gorillaSwagger.AddToSwagger(r).
			SetProduces("application/json").
			SetDescription("Test GetHandler").
			SetSummary("Test simply GET handler").
			AddInQueryParameter("string_id", "Some string ID", reflect.String, true).
			AddInQueryParameter("int_id", "Some int ID", reflect.Int64, false).
			AddResponse(http.StatusOK, "Test", &TestStruct{})
		return
	}
	// Main code of handler
	intID, err := strconv.Atoi(r.URL.Query().Get("int_id"))
	if err != nil {
		fmt.Println(err)
	}

	jData, err := json.Marshal(&TestStruct{
		Name:   r.URL.Query().Get("string_id"),
		Conter: intID,
	})
	if err != nil {
		fmt.Println(err)
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(jData)
	if err != nil {
		fmt.Println(err)
	}
}

func testGetHeader(w http.ResponseWriter, r *http.Request) {
	// Swagger
	if gorillaSwagger.IsBuildingSwagger(r) {
		gorillaSwagger.AddToSwagger(r).
			SetProduces("application/json").
			SetDescription("Test GetHeaderHandler").
			SetSummary("Test GET with header handler").
			AddInHeaderParameter("Authorization", "Authorization header", reflect.String, true).
			AddResponse(http.StatusOK, "Test", &TestHeader{})
		return
	}

	// Main code of handler
	token := r.Header.Get("Authorization")

	jData, err := json.Marshal(&TestHeader{
		Header: token,
	})
	if err != nil {
		fmt.Println(err)
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(jData)
	if err != nil {
		fmt.Println(err)
	}
}

func testGetCookie(w http.ResponseWriter, r *http.Request) {
	// Swagger
	if gorillaSwagger.IsBuildingSwagger(r) {
		gorillaSwagger.AddToSwagger(r).
			SetProduces("application/json").
			SetDescription("Test GetCookieHandler").
			SetSummary("Test GET with cookie handler").
			AddInCookieParameter("token-cookie", "token cookie", reflect.String, false).
			AddResponse(http.StatusOK, "Test", &TestCookie{})
		return
	}

	// Main code of handler
	cookie, err := r.Cookie("token-cookie")
	if err != nil {
		fmt.Println(err)
	}

	jData, err := json.Marshal(&TestCookie{
		Cookie: cookie.Value,
	})
	if err != nil {
		fmt.Println(err)
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(jData)
	if err != nil {
		fmt.Println(err)
	}
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
	TestStructWithAnonymousField
	TestString2    string `json:"string2"`
	TestStringEnum string `json:"testStringEnum" swagenum:"testEnum1,testEnum2,testEnum3"`
}

func testPostHandler(w http.ResponseWriter, r *http.Request) {
	// Swagger
	if gorillaSwagger.IsBuildingSwagger(r) {
		gorillaSwagger.AddToSwagger(r).
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
		return
	}

	// Main code of handler
	var (
		test         TestStruct
		jRequestData []byte
	)
	jRequestData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	r.Body.Close()
	err = json.Unmarshal(jRequestData, &test)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	jData, err := json.Marshal(&TestStructWithInterfaceField{Str: test.Name, InterfaceField: TestStruct{}})
	if err != nil {
		fmt.Println(err)
	}

	_, err = w.Write(jData)
	if err != nil {
		fmt.Println(err)
	}
}

func testWithoutSwaggerGetHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	jData, err := json.Marshal(&TestStruct{Name: "testExcludeSwaggerGetHandler"})
	if err != nil {
		fmt.Println(err)
	}

	_, err = w.Write(jData)
	if err != nil {
		fmt.Println(err)
	}
}

func testWithMiddlewareGetHandler(w http.ResponseWriter, r *http.Request) {
	// Swagger
	if gorillaSwagger.IsBuildingSwagger(r) {
		gorillaSwagger.AddToSwagger(r).
			SetProduces("application/json").
			SetDescription("Test WithMiddlewareGetHandler").
			SetSummary("Test Get handler with middleware").
			AddInQueryParameter("some_id", "Some ID", reflect.String, true).
			AddResponse(http.StatusOK, "Test", &TestStruct{})
		return
	}

	// Main code of handler
	w.WriteHeader(http.StatusOK)
	jData, err := json.Marshal(&TestStruct{Name: r.URL.Query().Get("some_id")})
	if err != nil {
		fmt.Println(err)
	}

	_, err = w.Write(jData)
	if err != nil {
		fmt.Println(err)
	}
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Swagger
		if gorillaSwagger.IsBuildingSwagger(r) {
			gorillaSwagger.AddToSwagger(r).
				AddResponse(http.StatusMethodNotAllowed, "Test middleware response", nil)
		}

		// Main code of handler
		if r.URL.Query().Get("some_id") == "fail" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type TestStructWithArray struct {
	ArrayField []string `json:"arrayField"`
}

func testArrayGetHandler(w http.ResponseWriter, r *http.Request) {
	// Swagger
	if gorillaSwagger.IsBuildingSwagger(r) {
		gorillaSwagger.AddToSwagger(r).
			SetProduces("application/json").
			SetDescription("Test ArrayGetHandler").
			SetSummary("Get array").
			AddResponse(http.StatusOK, "Test", &TestStructWithArray{})
		return
	}

	// Main code of handler
	w.WriteHeader(http.StatusOK)
	jData, err := json.Marshal(&TestStructWithArray{ArrayField: []string{"element1", "element2", "element3"}})
	if err != nil {
		fmt.Println(err)
	}

	_, err = w.Write(jData)
	if err != nil {
		fmt.Println(err)
	}
}

type TestEeee struct {
	Name   string `json:"Name"`
	Conter int64  `json:"conter"`
}

type TestStructWithArrayOfStructs struct {
	ArrayField []TestEeee `json:"arrayField"`
}

func testArrayOfStructGetHandler(w http.ResponseWriter, r *http.Request) {
	// Swagger
	if gorillaSwagger.IsBuildingSwagger(r) {
		gorillaSwagger.AddToSwagger(r).
			SetProduces("application/json").
			SetDescription("Test ArrayOfStruct").
			SetSummary("Get array of structs").
			AddResponse(http.StatusOK, "Test", &TestStructWithArrayOfStructs{})
		return
	}

	// Main code of handler
	w.WriteHeader(http.StatusOK)
	jData, err := json.Marshal(&TestStructWithArrayOfStructs{ArrayField: []TestEeee{
		{Name: "element1", Conter: 1},
		{Name: "element2", Conter: 2},
		{Name: "element3", Conter: 3},
	}})
	if err != nil {
		fmt.Println(err)
	}

	_, err = w.Write(jData)
	if err != nil {
		fmt.Println(err)
	}
}

func testParamPathGetHandler(w http.ResponseWriter, r *http.Request) {
	// Swagger
	if gorillaSwagger.IsBuildingSwagger(r) {
		gorillaSwagger.AddToSwagger(r).
			SetProduces("text/plain").
			SetDescription("Test ParamPathGetHandler").
			SetSummary("Test Param in Path GET handler").
			AddInPathParameter("id", "Some id", reflect.Int16).
			AddInPathParameter("name", "Some name", reflect.String).
			AddResponse(http.StatusOK, "Test", reflect.Int64)

		return
	}

	// Main code of handler
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(mux.Vars(r)["id"]))
	if err != nil {
		fmt.Println(err)
	}
}
