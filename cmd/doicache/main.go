package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/miku/doicache"
)

var (
	databaseDir = flag.String("db", filepath.Join(doicache.UserHomeDir(), ".doicache/default"), "leveldb directory")
)

func main() {
	flag.Parse()
	cache := doicache.New(*databaseDir)
	cache.Verbose = true
	cache.TTL = 3 * time.Second
	b, err := cache.Get("10.1103/PhysRevLett.118.140402")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
