package admin

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"log"
)

type editData struct {
    Data  any
	Error string
	Success string
}

//go:embed templates
var templateFiles embed.FS

func flagImgSrc(iso_code string) string {
	if len(iso_code) < 2 {
		return ""
	}

	r1 := rune(iso_code[0]-'A') + 0x1F1E6
	r2 := rune(iso_code[1]-'A') + 0x1F1E6

	return fmt.Sprintf("https://cdn.jsdelivr.net/npm/@svgmoji/openmoji@2.0.0/svg/%X-%X.svg", r1, r2)
}

func InitTemplates() (*template.Template, error) {
	log.Println("Initializing templates...")

    funcMap := template.FuncMap{
        "flagImgSrc": flagImgSrc,
    }

    tmpl, err := template.New("").Funcs(funcMap).ParseFS(templateFiles,
        "templates/pages/*.html",
		"templates/partials/*.html",
		"templates/components/*.html",
    )
    if err != nil {
        return nil, fmt.Errorf("parsing admin templates: %w", err)
    }

	log.Println("Templates initialized successfully.")

    return tmpl, nil
}

func render(w http.ResponseWriter, tmpl *template.Template, page string, data any) {
	if err := tmpl.ExecuteTemplate(w, page, data); err != nil {
		http.Error(w, "render error: "+err.Error(), http.StatusInternalServerError)
	}
}
