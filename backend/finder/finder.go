// This package contains read-only queries that require more data than just the
// normalized domain models (like articles or blogs+isFollowing). They are not
// grouped by types and instead exist as top-level methods of the Finder struct.
package finder

import "github.com/theandrew168/bloggulus/backend/postgres"

type Finder struct {
	conn postgres.Conn
}

func New(conn postgres.Conn) *Finder {
	f := Finder{
		conn: conn,
	}
	return &f
}
