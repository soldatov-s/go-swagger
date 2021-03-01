package swagger

import (
	"reflect"
	"strings"
)

type (
	// Description the REST-method of endpoint
	Method struct {
		Description string   `json:"description,omitempty"`
		Consumes    []string `json:"consumes,omitempty"`
		Produces    []string `json:"produces,omitempty"`
		Summary     string   `json:"summary,omitempty"`
		OperationID string   `json:"operationId,omitempty"`
		// The parameters of requests
		Parameters ArrayParameters `json:"parameters,omitempty"`
		// The endpoint responses
		Responses MapResponse `json:"responses,omitempty"`
		// Counters for same codes for multiple responses
		ResponseCode MapResponseCode `json:"-"`
	}

	// List of methods (GET, POST,...) for endpoint
	Methods map[string]IMethod
)

// First, programmer set Consumes/Produces/Response
type IMethod interface {
	// SetConsumes - sets the MIME types of accept data for the endpoint
	SetConsumes(c ...string) Producer
	// SetConsumes - sets the MIME types of return data for the endpoint
	SetProduces(p ...string) Consumer
	// AddResponse - adds a response
	AddResponse(сode int, description string, schema interface{}) Responser
	// AddFileParameter - adds a file response
	AddFileResponse(responseCode int, description string) Responser
}

type Responser interface {
	// AddParameter - adds a response
	AddResponse(сode int, description string, schema interface{}) Responser
	// AddFileParameter - adds a file response
	AddFileResponse(responseCode int, description string) Responser
}

type AdderInParameter interface {
	// AddInBodyParameter - adds a request in body parameter
	AddInPathParameter(name, description string, t reflect.Kind) AdderInParameter
	// AddInQueryParameter - adds a request in query parameter
	AddInQueryParameter(name, description string, t reflect.Kind, required bool) AdderInParameter
	// AddInHeaderParameter - adds a request in header parameter
	AddInHeaderParameter(name, description string, t reflect.Kind, required bool) AdderInParameter
	// AddInCookieParameter - adds a request in cookie parameter
	AddInCookieParameter(name, description string, t reflect.Kind, required bool) AdderInParameter
	// AddInFileParameter - adds a request in file parameter
	AddInFileParameter(name, description string) AdderInParameter
	Responser
}

type ParameterAndResponser interface {
	// AddParameter - adds a response
	AddResponse(Code int, description string, schema interface{}) Responser
	// AddFileParameter - adds a file response
	AddFileResponse(responseCode int, description string) Responser
	// AddInBodyParameter - adds a request in body parameter
	AddInBodyParameter(name, description string, t interface{}, required bool) AdderInParameter
	// AddInPathParameter - adds a request in path parameter
	AddInPathParameter(name, description string, t reflect.Kind) AdderInParameter
	// AddInQueryParameter - adds a request in query parameter
	AddInQueryParameter(name, description string, t reflect.Kind, required bool) AdderInParameter
	// AddInHeaderParameter - adds a request in header parameter
	AddInHeaderParameter(name, description string, t reflect.Kind, required bool) AdderInParameter
	// AddInCookieParameter - adds a request in cookie parameter
	AddInCookieParameter(name, description string, t reflect.Kind, required bool) AdderInParameter
	// AddInFileParameter - adds a request in file parameter
	AddInFileParameter(name, description string) AdderInParameter
}

type Summarer interface {
	// SetDescription - sets a description of endpoint
	SetSummary(d string) ParameterAndResponser
}

type Consumer interface {
	// SetDescription - sets a description of endpoint
	SetDescription(d string) Summarer
}

type Producer interface {
	// SetProduces - sets the MIME types of return data for the endpoint
	SetProduces(p ...string) Consumer
	// SetDescription - sets a description of endpoint
	SetDescription(d string) Summarer
}

// NewMethod - create a new instance of the Method
func NewMethod() *Method {
	return &Method{
		Description: "Unnamed handler",
	}
}

func (m *Method) SetConsumes(c ...string) Producer {
	if m == nil {
		return nil
	}
	m.Consumes = c
	return m
}

func (m *Method) SetProduces(p ...string) Consumer {
	if m == nil {
		return nil
	}
	m.Produces = p
	return m
}

func (m *Method) SetDescription(d string) Summarer {
	if m == nil {
		return nil
	}
	m.Description = d
	return m
}

func (m *Method) SetSummary(s string) ParameterAndResponser {
	if m == nil {
		return nil
	}
	m.Summary = s
	return m
}

func (m *Method) addIn(name, description string, t interface{}, required bool, inType InType) AdderInParameter {
	if m == nil {
		return nil
	}
	m.Parameters = append(m.Parameters, NewParameter(name, description, t, required, inType))
	return m
}

func (m *Method) AddInPathParameter(name, description string, t reflect.Kind) AdderInParameter {
	return m.addIn(name, description, t, true, InPath)
}

func (m *Method) AddInQueryParameter(name, description string, t reflect.Kind, required bool) AdderInParameter {
	return m.addIn(name, description, t, required, InQuery)
}

func (m *Method) AddInCookieParameter(name, description string, t reflect.Kind, required bool) AdderInParameter {
	return m.addIn(name, description, t, required, InCookie)
}

func (m *Method) AddInHeaderParameter(name, description string, t reflect.Kind, required bool) AdderInParameter {
	return m.addIn(name, description, t, required, InHeader)
}

func (m *Method) AddInBodyParameter(name, description string, t interface{}, required bool) AdderInParameter {
	return m.addIn(name, description, t, required, InBody)
}

func (m *Method) AddInFileParameter(name, description string) AdderInParameter {
	if m == nil {
		return nil
	}
	m.Parameters = append(m.Parameters, NewFileParameter(name, description))
	return m
}

func (m *Method) AddResponse(responseCode int, description string, schema interface{}) Responser {
	if m == nil {
		return nil
	}
	if m.Responses == nil {
		m.Responses = make(map[string]*Response)
	}
	response := NewResponse(description)

	if schema != nil {
		response.Schema = NewSchema(schema)
	}

	m.Responses[m.ResponseCode.Append(responseCode)] = response

	return m
}

func (m *Method) AddFileResponse(responseCode int, description string) Responser {
	if m == nil {
		return nil
	}
	if m.Responses == nil {
		m.Responses = make(map[string]*Response)
	}
	response := NewResponse(description)

	response.Schema = &Schema{
		TypeName: "file",
	}

	m.Responses[m.ResponseCode.Append(responseCode)] = response

	return m
}

func (m *Method) Parse(path, methodName string, sw Doc) {
	// Parse parameters
	for _, p := range m.Parameters {
		p.Parse(sw)
	}
	// Parse responses
	for _, r := range m.Responses {
		r.Parse(sw)
	}
	if sw.Paths[path] == nil {
		sw.Paths[path] = make(Methods)
	}
	sw.Paths[path][strings.ToLower(methodName)] = m
}
