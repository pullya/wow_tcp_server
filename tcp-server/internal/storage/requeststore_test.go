package storage

import (
	"context"
	"sync"
	"testing"
)

func TestRequestStore_Add(t *testing.T) {
	ctx := context.Background()
	sf := func(in string) uint32 {
		key := []byte(in)
		numShards := 8
		hashValue := HashShard(key)

		return hashValue % uint32(numShards)
	}
	rs := make(map[uint32]map[string]bool, 8)

	type fields struct {
		shards    map[uint32]map[string]bool
		shardFunc func(data string) uint32
		mutex     *sync.Mutex
	}
	type args struct {
		ctx     context.Context
		request string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "",
			fields: fields{
				shards:    rs,
				shardFunc: sf,
				mutex:     &sync.Mutex{},
			},
			args: args{ctx: ctx, request: "key"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := &RequestStore{
				shards:    tt.fields.shards,
				shardFunc: tt.fields.shardFunc,
				mu:        tt.fields.mutex,
			}
			rs.Add(tt.args.ctx, tt.args.request)
		})
	}
}
