package storage

import (
	"github.com/google/uuid"
	"github.com/pullya/wow_tcp_server/tcp-server/internal/config"
	"github.com/spaolacci/murmur3"
)

func ShardKey(in string) uint32 {
	key := []byte(in)
	numShards := config.Config.ShardsCnt

	hashValue := HashShard(key)

	return hashValue % uint32(numShards)
}

func HashShard(data []byte) uint32 {
	hasher := murmur3.New32()
	hasher.Write(data)
	return hasher.Sum32()
}

func GenUID() string {
	id := uuid.New()

	return id.String()
}
