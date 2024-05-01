package admin

import (
	"github.com/theandrew168/bloggulus/backend/postgres"
)

type Storage struct {
	conn postgres.Conn

	blog    *BlogStorage
	post    *PostStorage
	tag     *TagStorage
	account *AccountStorage
	token   *TokenStorage
}

func New(conn postgres.Conn) *Storage {
	s := Storage{
		conn: conn,

		blog:    NewBlogStorage(conn),
		post:    NewPostStorage(conn),
		tag:     NewTagStorage(conn),
		account: NewAccountStoragee(conn),
		token:   NewTokenStorage(conn),
	}
	return &s
}

func (s *Storage) Blog() *BlogStorage {
	return s.blog
}

func (s *Storage) Post() *PostStorage {
	return s.post
}

func (s *Storage) Tag() *TagStorage {
	return s.tag
}

func (s *Storage) Account() *AccountStorage {
	return s.account
}

func (s *Storage) Token() *TokenStorage {
	return s.token
}
