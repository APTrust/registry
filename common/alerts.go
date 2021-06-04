package common

import (
	"fmt"
	"path"
	"path/filepath"
	"text/template"
)

var TextTemplates map[string]*template.Template

func init() {
	TextTemplates = make(map[string]*template.Template)
	pattern := path.Join(ProjectRoot(), "alert_templates", "*.txt")
	files, _ := filepath.Glob(pattern)
	for _, file := range files {
		name := fmt.Sprintf("alerts/%s", path.Base(file))
		TextTemplates[name] = template.Must(template.ParseFiles(file))
	}
}
