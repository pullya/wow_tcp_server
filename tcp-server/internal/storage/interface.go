package storage

import (
	"context"
	"math/rand"
)

type Storage interface {
	GetRandomWoW(ctx context.Context) string
}

type InMemStorage struct {
	wow []string
}

func NewInMemStorage(wow []string) InMemStorage {
	return InMemStorage{
		wow: wow,
	}
}

func (ims InMemStorage) GetRandomWoW(ctx context.Context) string {
	return ims.wow[rand.Intn(len(ims.wow))]
}
