// This package contains read-only queries that require more data than just the
// normalized domain models (like articles or blogs+isFollowing).
package query

import "github.com/theandrew168/bloggulus/backend/postgres"

// TODO: Add queries for all read operations.

type Query struct {
	account *AccountQuery
	article *ArticleQuery
	blog    *BlogQuery
}

func New(conn postgres.Conn) *Query {
	qry := Query{
		account: NewAccount(conn),
		article: NewArticle(conn),
		blog:    NewBlog(conn),
	}
	return &qry
}

func (qry *Query) Account() *AccountQuery {
	return qry.account
}

func (qry *Query) Article() *ArticleQuery {
	return qry.article
}

func (qry *Query) Blog() *BlogQuery {
	return qry.blog
}
