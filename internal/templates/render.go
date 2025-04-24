package templates

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func (t *Templates) ExecuteTemplate(w io.Writer, name string, data interface{}) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplates() *Templates {
	funcMap := template.FuncMap{}
	tmpl := template.New("").Funcs(funcMap)

	// Добавляем шаблоны из нескольких путей
	tmpl = template.Must(tmpl.ParseGlob("web/views/*.html"))
	tmpl = template.Must(tmpl.ParseGlob("web/views/components/*.html"))

	return &Templates{
		templates: tmpl,
	}
}

func GetTemplates(c echo.Context) *Templates {
	return c.Get("templates").(*Templates)
}
