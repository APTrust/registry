package common

import (
	"path"
	"text/template"
)

var AlertTemplate *template.Template

func init() {
	pattern := path.Join(ProjectRoot(), "views", "alerts", "*.txt")
	AlertTemplate = template.Must(template.ParseGlob(pattern))
}
