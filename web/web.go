package web

import "embed"

//go:embed index.html openapi.json VERSION
var Content embed.FS
