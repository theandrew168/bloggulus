package frontend

import "io/fs"

// This var exposes the compiled frontend files either via
// embedding (with the "embed" tag) or a regular dir (without any tags).
// This builds upon the following conditional embedding concept:
// https://github.com/golang/go/issues/44484#issuecomment-948137497
var Frontend fs.FS
