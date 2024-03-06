//go:build !embed

package frontend

import (
	"os"
)

func init() {
	Frontend = os.DirFS("./frontend/build")
}
