package templates

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

var funcMap = template.FuncMap{
	"formatDate": func(t time.Time) string {
		return t.Format("January 2, 2006")
	},
	"shortExcerpt": func(content string) string {
		text := strings.ReplaceAll(content, "\n", " ")
		if len(text) > 160 {
			return text[:160] + "..."
		}
		return text
	},
	"nl2br": func(s string) template.HTML {
		escaped := template.HTMLEscapeString(s)
		return template.HTML(strings.ReplaceAll(escaped, "\n", "<br>"))
	},
	"safeHTML": func(s string) template.HTML {
		return template.HTML(s)
	},
}

const templatesDir = "internal/templates"

func Render(w http.ResponseWriter, name string, data any) {
	files := []string{
		filepath.Join(templatesDir, "layout.html"),
		filepath.Join(templatesDir, name),
		// Always include shared partials
		filepath.Join(templatesDir, "partials/post_list.html"),
		filepath.Join(templatesDir, "partials/post_form.html"),
	}
	tmpl, err := template.New("").Funcs(funcMap).ParseFiles(files...)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		fmt.Printf("Template execute error: %v\n", err)
	}
}

func RenderPartial(w http.ResponseWriter, name string, data any) {
	tmpl, err := template.New("").Funcs(funcMap).ParseFiles(
		filepath.Join(templatesDir, name),
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	baseName := filepath.Base(name)
	templateName := strings.TrimSuffix(baseName, ".html")
	if err := tmpl.ExecuteTemplate(w, templateName, data); err != nil {
		fmt.Printf("Partial execute error: %v\n", err)
	}
}
