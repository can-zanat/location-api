package internal

type Service struct {
	store Store
}

type LocationDbStore interface {
}

func NewService(s Store) *Service {
	return &Service{store: s}
}
