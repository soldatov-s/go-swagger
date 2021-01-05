package swagger

import (
	"encoding/json"
	"errors"
	"strings"
	"sync"
)

type (
	// Base object for describing the builded API
	BaseAPI struct {
		Version string `json:"swagger,omitempty"`
		Info    IInfo  `json:"info,omitempty"`
		// Base path to API, it consist the prefix of endpoint path
		BasePath string `json:"basePath,omitempty"`
	}
	// High level object for describing the builded API
	Doc struct {
		BaseAPI
		// List of paths to endpoints
		Paths map[string]Methods `json:"paths,omitempty"`
		// List of definitions
		Definitions map[string]*Definition `json:"definitions,omitempty"`
	}
	// Information about the created swagger
	Info struct {
		Description    string   `json:"description,omitempty"`
		Title          string   `json:"title,omitempty"`
		TermsOfService string   `json:"termsOfService,omitempty"`
		Contact        IContact `json:"contact,omitempty"`
		License        ILicense `json:"license,omitempty"`
		Version        string   `json:"version,omitempty"`
	}
	Contact struct {
		Name  string `json:"name,omitempty"`
		URL   string `json:"url,omitempty"`
		Email string `json:"email,omitempty"`
	}
	License struct {
		Name string `json:"name,omitempty"`
		URL  string `json:"url,omitempty"`
	}
)

type ISwaggerAPI interface{}

type BasePather interface {
	SetBasePath(p string) Informer
}

type Informer interface {
	SetInfo(i IInfo) ISwaggerAPI
}

func NewSwagger() BasePather {
	return &BaseAPI{
		Version: "2.0",
	}
}

func (s *Doc) ReadDoc() string {
	return string(s.JSON())
}

func (s *Doc) JSON() (jsonData []byte) {
	jsonData, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return
	}
	return
}

func (s *BaseAPI) NewSwagger() BasePather {
	return &BaseAPI{
		Version: "2.0",
	}
}

func (s *BaseAPI) SetBasePath(p string) Informer {
	if s == nil {
		return nil
	}
	s.BasePath = strings.TrimRight(p, "/")
	return s
}

func (s *BaseAPI) SetInfo(i IInfo) ISwaggerAPI {
	if s == nil {
		return nil
	}
	s.Info = i
	return s
}

type IInfo interface {
	SetVersion(v string) IInfo
	SetContact(c IContact) IInfo
	SetDescription(d string) IInfo
	SetTitle(t string) IInfo
	SetLicense(l ILicense) IInfo
	SetTermOfService(tos string) IInfo
}

func NewInfo() IInfo {
	return &Info{
		Title:          "Swagger API",
		Description:    "This is an embedded Swagger-server.",
		TermsOfService: "http://swagger.io/terms/",
		Version:        "0.0.1",
		Contact:        NewContact(),
		License:        NewLicense(),
	}
}

func (i *Info) SetVersion(v string) IInfo {
	if i == nil {
		return nil
	}
	i.Version = v
	return i
}

func (i *Info) SetContact(c IContact) IInfo {
	if i == nil {
		return nil
	}
	i.Contact = c
	return i
}

func (i *Info) SetDescription(d string) IInfo {
	if i == nil {
		return nil
	}
	i.Description = d
	return i
}

func (i *Info) SetTitle(t string) IInfo {
	if i == nil {
		return nil
	}
	i.Title = t
	return i
}

func (i *Info) SetLicense(l ILicense) IInfo {
	if i == nil {
		return nil
	}
	i.License = l
	return i
}

func (i *Info) SetTermOfService(tos string) IInfo {
	if i == nil {
		return nil
	}
	i.TermsOfService = tos
	return i
}

type ILicense interface {
	SetURL(url string) ILicense
	SetName(n string) ILicense
}

func NewLicense() ILicense {
	return &License{
		Name: "Apache 2.0",
		URL:  "http://www.apache.org/licenses/LICENSE-2.0.html",
	}
}

func (l *License) SetURL(url string) ILicense {
	if l == nil {
		return nil
	}
	l.URL = url
	return l
}

func (l *License) SetName(n string) ILicense {
	if l == nil {
		return nil
	}
	l.Name = n
	return l
}

type IContact interface {
	SetEmail(e string) IContact
	SetName(n string) IContact
	SetURL(url string) IContact
}

func NewContact() IContact {
	return &Contact{
		Name:  "API Support",
		URL:   "https://test.test",
		Email: "support@test.test",
	}
}

func (c *Contact) SetName(n string) IContact {
	if c == nil {
		return nil
	}
	c.Name = n
	return c
}

func (c *Contact) SetURL(url string) IContact {
	if c == nil {
		return nil
	}
	c.URL = url
	return c
}

func (c *Contact) SetEmail(e string) IContact {
	if c == nil {
		return nil
	}
	c.Email = e
	return c
}

var (
	swaggerMu sync.RWMutex
	swag      map[string]Swagger
)

// Swagger is a interface to read swagger document.
type Swagger interface {
	ReadDoc() string
}

// Register registers swagger for given name.
func Register(name string, swagger Swagger) {
	if swagger == nil {
		panic("swagger is nil")
	}

	if swag == nil {
		swag = make(map[string]Swagger)
	}
	swaggerMu.Lock()
	defer swaggerMu.Unlock()

	if _, ok := swag[name]; ok {
		return
	}

	swag[name] = swagger
}

// ReadDoc reads swagger document.
func ReadDoc(name string) (string, error) {
	if swag != nil {
		swaggerMu.RLock()
		defer swaggerMu.RUnlock()
		return swag[name].ReadDoc(), nil
	}
	return "", errors.New("not yet registered swag")
}
