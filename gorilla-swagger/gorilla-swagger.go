package gorillaswagger

import (
	"html/template"
	"net/http"
	"regexp"

	"github.com/soldatov-s/go-swagger/swagger"
	swaggerFiles "github.com/swaggo/files"
)

// WrapHandler wraps swaggerFiles.Handler and returns http.HandlerFunc
var WrapHandler = Handler()

// Handler wraps `http.Handler` into `http.HandlerFunc`.
func Handler(confs ...func(c *swagger.Config)) http.HandlerFunc {
	handler := swaggerFiles.Handler

	config := &swagger.Config{
		URL: "doc.json",
	}

	for _, c := range confs {
		c(config)
	}

	// create a template with name
	t := template.New("swagger_index.html")
	index, _ := t.Parse(swagger.IndexTempl)

	var re = regexp.MustCompile(`(.*)(index\.html|doc\.json|favicon-16x16\.png|favicon-32x32\.png|/oauth2-redirect\.html|swagger-ui\.css|swagger-ui\.css\.map|swagger-ui\.js|swagger-ui\.js\.map|swagger-ui-bundle\.js|swagger-ui-bundle\.js\.map|swagger-ui-standalone-preset\.js|swagger-ui-standalone-preset\.js\.map)[\?|.]*`)

	return func(w http.ResponseWriter, r *http.Request) {
		var matches []string
		if matches = re.FindStringSubmatch(r.RequestURI); len(matches) != 3 {
			w.WriteHeader(http.StatusNotFound)
			_, err := w.Write([]byte("404 page not found"))
			if err != nil {
				log.Err(err).Msg("Error write bytes")
			}
			return
		}
		path := matches[2]
		prefix := matches[1]
		handler.Prefix = prefix

		switch path {
		case "index.html":
			proto := "http://"
			if r.TLS != nil {
				proto = "https://"
			}
			tmpConfig := &swagger.Config{
				URL:  proto + r.Host + prefix + config.URL,
				Name: config.Name,
			}
			err := index.Execute(w, tmpConfig)
			if err != nil {
				log.Err(err).Msg("Error build template")
			}
		case "doc.json":
			doc, err := swagger.ReadDoc(config.Name)
			if err != nil {
				log.Err(err).Msg("Error read doc")
			}
			_, err = w.Write([]byte(doc))
			if err != nil {
				log.Err(err).Msg("Error write bytes")
			}
		case "":
			http.Redirect(w, r, prefix+"index.html", http.StatusMovedPermanently)
		default:
			handler.ServeHTTP(w, r)
		}
	}
}
