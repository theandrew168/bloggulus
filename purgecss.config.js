module.exports = {
	extractors: [{
		extractor: (content) => {
			return content.match(/[A-Za-z0-9-_:\/]+/g) || []
		},
		extensions: ['tmpl', 'html', 'js', 'css']
	}]
}
