package web

import "embed"

//go:embed index.html openapi.json VERSION robots.txt sitemap.xml llms.txt favicon.svg
var Content embed.FS
