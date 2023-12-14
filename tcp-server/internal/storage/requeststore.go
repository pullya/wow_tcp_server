package storage

import (
	"context"
	"errors"
	"sync"

	"github.com/pullya/wow_tcp_server/tcp-server/internal/config"
)

type ShardFunc func(data []byte) uint32

//go:generate mockery --name=Requester --output=mocks --case=underscore
type Requester interface {
	Add(ctx context.Context, request string)
	Get(ctx context.Context, request string) (bool, error)
	Set(ctx context.Context, request string) error
}

type RequestStore struct {
	shards    map[uint32]map[string]bool
	shardFunc func(data string) uint32
	mu        *sync.Mutex
}

func NewRequestStore(shardFunc func(data string) uint32) *RequestStore {
	return &RequestStore{
		shards:    make(map[uint32]map[string]bool, config.Config.ShardsCnt),
		shardFunc: shardFunc,
		mu:        &sync.Mutex{},
	}
}

func (rs *RequestStore) Add(_ context.Context, request string) {
	shardKey := rs.shardFunc(request)

	rs.mu.Lock()
	defer rs.mu.Unlock()

	if _, ok := rs.shards[shardKey]; !ok {
		shard := make(map[string]bool)
		rs.shards[shardKey] = shard
	}

	rs.shards[shardKey][request] = false
}

func (rs *RequestStore) Get(_ context.Context, request string) (bool, error) {
	shardKey := rs.shardFunc(request)

	rs.mu.Lock()
	defer rs.mu.Unlock()

	if _, ok := rs.shards[shardKey]; !ok {
		return false, errors.New("shard not found")
	}

	if status, ok := rs.shards[shardKey][request]; ok {
		return status, nil
	}
	return false, errors.New("request not found")
}

func (rs *RequestStore) Set(_ context.Context, request string) error {
	shardKey := rs.shardFunc(request)

	rs.mu.Lock()
	defer rs.mu.Unlock()

	if _, ok := rs.shards[shardKey]; !ok {
		return errors.New("shard not found")
	}

	if _, ok := rs.shards[shardKey][request]; !ok {
		return errors.New("request not found")
	}
	rs.shards[shardKey][request] = true
	return nil
}
