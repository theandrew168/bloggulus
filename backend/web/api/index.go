package api

import (
	"net/http"
)

const html = `
<!DOCTYPE html>
<html lang="en">

<head>
	<title>Bloggulus - A website for avid blog readers</title>

	<meta charset="utf-8" />
	<meta name="description" content="Bloggulus - A website for avid blog readers" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />

	<script type="module" src="https://unpkg.com/rapidoc/dist/rapidoc-min.js"></script>
</head>

<body>
	<rapi-doc
		spec-url="/openapi.yaml"
		render-style="read"
	></rapi-doc>
</body>

</html>
`

func (app *Application) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(html))
	}
}
