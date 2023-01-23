package templates

import (
	"bytes"
	"embed"
	"text/template"
)

//go:embed data/*
var f embed.FS

// GenerateUserData generate user data from template
func GenerateUserData(tplFiele string, data map[string]interface{}) (string, error) {
	// t, err := template.New("").Parse(tplFiele)
	// t, err := embed.Parse(template, &f)
	t, err := template.ParseFS(f, tplFiele)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
