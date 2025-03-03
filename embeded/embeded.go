package embeded

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"

	"github.com/labstack/echo/v4"
)

//go:embed static/*
var static embed.FS

func Static() fs.FS {
	staticFS, err := fs.Sub(static, "static")
	if err != nil {
		panic(fmt.Sprintf("failed to create static file system: %v", err))
	}

	return staticFS
}

//go:embed layouts mails partials
var templateFiles embed.FS

type Template struct {
	views  map[string]*template.Template
	emails map[string]*template.Template
}

func (t *Template) Render(w io.Writer, vName string, data any, c echo.Context) error {
	view, ok := t.views[vName]
	if !ok {
		panic(fmt.Sprintf("invalid view name: %q", vName))
	}

	err := view.Execute(w, data)
	if err != nil {
		panic(fmt.Sprintf("failed to execute view %q: %v", vName, err))
	}

	return nil
}

func (t *Template) RenderEmail(w io.Writer, eName string, data any) error {
	email, ok := t.emails[eName]
	if !ok {
		panic(fmt.Sprintf("invalid email name: %q", eName))
	}

	err := email.Execute(w, data)
	if err != nil {
		panic(fmt.Sprintf("failed to execute email %q: %v", eName, err))
	}

	return nil
}

func (t *Template) NewView(name, base string, partials ...string) {
	if _, ok := t.views[name]; ok {
		panic(fmt.Sprintf("view with name %q already registered.", name))
	}

	all := make([]string, len(partials)+1)
	all[0] = fmt.Sprintf("layouts/%s", base)
	for i, p := range partials {
		all[i+1] = fmt.Sprintf("partials/%s", p)
	}

	t.views[name] = template.Must(template.New(base).Funcs(
		template.FuncMap{
			"safeHTML": safeHTML,
		},
	).ParseFS(templateFiles, all...))
}

func (t *Template) NewEmail(name, base string) {
	if _, ok := t.emails[name]; ok {
		panic(fmt.Sprintf("email with name %q already registered.", name))
	}

	email := template.Must(template.New(base).ParseFS(templateFiles, fmt.Sprintf("mails/%s", base)))
	t.emails[name] = email
}

func Templates() *Template {
	return &Template{
		views:  map[string]*template.Template{},
		emails: map[string]*template.Template{},
	}
}

func safeHTML(str string) template.HTML {
	return template.HTML(str)
}
