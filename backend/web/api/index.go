package api

import (
	"net/http"
)

const redocHTML = `
<!DOCTYPE html>
<html lang="en">

<head>
	<title>Bloggulus - A website for avid blog readers</title>

	<meta charset="utf-8" />
	<meta name="description" content="Bloggulus - A website for avid blog readers" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />

	<link href="https://fonts.googleapis.com/css?family=Montserrat:300,400,700|Roboto:300,400,700" rel="stylesheet">
	<style>
		body {
			margin: 0;
			padding: 0;
		}
	</style>
	<script src="https://unpkg.com/@stoplight/elements/web-components.min.js"></script>
    <link rel="stylesheet" href="https://unpkg.com/@stoplight/elements/styles.min.css">
</head>

<body>
	<redoc spec-url='/openapi.yaml' />
	<script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"></script>
</body>

</html>
`

func (app *Application) handleIndexRedoc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(redocHTML))
	}
}

const rapidocHTML = `
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
	/>
</body>

</html>
`

func (app *Application) handleIndexRapidoc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(rapidocHTML))
	}
}

const stoplightHTML = `
<!DOCTYPE html>
<html lang="en">

<head>
	<title>Bloggulus - A website for avid blog readers</title>

	<meta charset="utf-8" />
	<meta name="description" content="Bloggulus - A website for avid blog readers" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />

	<script src="https://unpkg.com/@stoplight/elements/web-components.min.js"></script>
    <link rel="stylesheet" href="https://unpkg.com/@stoplight/elements/styles.min.css">
</head>

<body>
	<elements-api
		apiDescriptionUrl="/openapi.yaml"
		router="hash"
		layout="sidebar"
	/>
</body>

</html>
`

func (app *Application) handleIndexStoplight() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(stoplightHTML))
	}
}
