package api

import (
	_ "embed"
	"net/http"
)

const html = `
	<!DOCTYPE html>
	<html lang="en">

	<head>
		<title>Bloggulus</title>

		<meta charset="utf-8" />
		<meta name="description" content="Bloggulus - A collective index of my favorite blogs" />
		<meta name="viewport" content="initial-scale=1, width=device-width" />

		<script type="module" src="https://unpkg.com/rapidoc/dist/rapidoc-min.js"></script>
	</head>

	<body>
		<rapi-doc
			spec-url="/api/v1/openapi.yaml"
			render-style="read"
		></rapi-doc>
	</body>
	</html>
`

func HandleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(html))
	}
}
