package storage

import (
	"context"
	"math/rand"
)

//go:generate mockery --name=Storageer --output=mocks --case=underscore
type Storageer interface {
	GetRandomWOW(ctx context.Context) string
}

type Storage struct {
	wow []string
}

func New(wow []string) Storage {
	return Storage{
		wow: wow,
	}
}

func (ims Storage) GetRandomWOW(ctx context.Context) string {
	return ims.wow[rand.Intn(len(ims.wow))]
}
