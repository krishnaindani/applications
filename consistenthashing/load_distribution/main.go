package main

import (
	"fmt"
	"github.com/buraksezer/consistent"
	"github.com/cespare/xxhash"
	"math"
	"math/rand"
	"time"
)

type Member string

func (m Member) String() string {
	return string(m)
}

type hasher struct{}

func (h hasher) Sum64(data []byte) uint64 {
	return xxhash.Sum64(data)
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	var members []consistent.Member
	for i := 0; i < 8; i++ {
		member := fmt.Sprintf("member%d", i)
		members = append(members, Member(member))
	}

	cfg := consistent.Config{
		Hasher:            hasher{},
		PartitionCount:    271,
		ReplicationFactor: 40,
		Load:              1.2,
	}

	c := consistent.New(members, cfg)

	keyCount := 1000000
	load := (c.AverageLoad() * float64(keyCount)) / float64(cfg.PartitionCount)
	fmt.Println("Maximum key count for a member should be around this: ", math.Ceil(load))
	distribution := make(map[string]int)
	key := make([]byte, 4)
	for i := 0; i < keyCount; i++ {
		rand.Read(key)
		member := c.LocateKey(key)
		distribution[member.String()] += 1
	}

	for member, count := range distribution {
		fmt.Println("Member: ", member, "Key Count: ", count)
	}
}
