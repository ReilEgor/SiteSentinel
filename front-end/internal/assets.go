package web

import "embed"

//go:embed templates/*.html
//go:embed templates/partials/*.html
var TemplateFS embed.FS

//go:embed static/*
var StaticFS embed.FS
