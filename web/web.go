package web

import "embed"

//go:embed index.html openapi.json
var Content embed.FS
