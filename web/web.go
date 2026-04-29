package web

import "embed"

//go:embed index.html openapi.json VERSION robots.txt sitemap.xml llms.txt
var Content embed.FS
