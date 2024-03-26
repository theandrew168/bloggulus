package storage

type Storage interface {
	Blog() BlogStorage
	Post() PostStorage
	Tag() TagStorage

	WithTransaction(operation func(store Storage) error) error
}
