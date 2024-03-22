package storage

type Storage interface {
	Blog() BlogStorage
	Post() PostStorage
	Tag() TagStorage

	Atomically(operation func(store Storage) error) error
}
