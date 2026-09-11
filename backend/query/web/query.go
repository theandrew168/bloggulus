// This package contains read-only queries that require more data than just the
// normalized domain models (like articles or blogs+isFollowing).
package webquery

import "github.com/theandrew168/bloggulus/backend/postgres"

type Query struct {
	account *AccountQuery
	article *ArticleQuery
	blog    *BlogQuery
	post    *PostQuery
}

func New(conn postgres.Conn) *Query {
	qry := Query{
		account: NewAccount(conn),
		article: NewArticle(conn),
		blog:    NewBlog(conn),
		post:    NewPost(conn),
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

func (qry *Query) Post() *PostQuery {
	return qry.post
}
