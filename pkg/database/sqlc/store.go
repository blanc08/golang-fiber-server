package database

type Store interface {
	Querier
}

type SQLStore struct {
	*Queries
}

func NewStore() Store {
	return &SQLStore{
		Queries: New(),
	}
}
