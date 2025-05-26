// This package contains read-only queries that require more data than just the
// normalized domain models (like articles or blogs+isFollowing). They are not
// grouped by types and instead exist as top-level methods of the Finder struct.
package query

import "github.com/theandrew168/bloggulus/backend/postgres"

// TODO: Add queries for all read operations.

type Query struct {
	conn postgres.Conn
}

func New(conn postgres.Conn) *Query {
	qry := Query{
		conn: conn,
	}
	return &qry
}
