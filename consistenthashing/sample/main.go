package main

import (
	"fmt"
	"github.com/buraksezer/consistent"
	"github.com/cespare/xxhash"
)

type myMember string

func (m myMember) String() string {
	return string(m)
}

type hasher struct{}

func (h hasher) Sum64(data []byte) uint64 {
	return xxhash.Sum64(data)
}

func main() {
	cfg := consistent.Config{
		Hasher:            hasher{},
		PartitionCount:    7,
		ReplicationFactor: 20,
		Load:              1.25,
	}

	c := consistent.New(nil, cfg)

	node1 := myMember("node1.test.com")
	c.Add(node1)

	node2 := myMember("node2.test.com")
	c.Add(node2)

	key := []byte("my-key")

	owner := c.LocateKey(key)
	fmt.Println(owner.String())
}
