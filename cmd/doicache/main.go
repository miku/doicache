package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/miku/doicache"
)

var (
	databaseDir = flag.String("db", filepath.Join(doicache.UserHomeDir(), ".doicache/default"), "leveldb directory")
	ttl         = flag.Duration("ttl", 24*time.Hour*120, "entry expiration")
	verbose     = flag.Bool("verbose", false, "be verbose")
)

func main() {
	flag.Parse()

	cache := doicache.New(*databaseDir)
	cache.Verbose = *verbose
	cache.TTL = *ttl

	var reader io.Reader = os.Stdin

	if flag.NArg() > 0 {
		if _, err := os.Stat(flag.Arg(0)); os.IsNotExist(err) {
			for _, arg := range flag.Args() {
				v, err := cache.Resolve(arg)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(v)
			}
			os.Exit(0)
		} else {
			f, err := os.Open(flag.Arg(0))
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			reader = f
		}
	}

	br := bufio.NewReader(reader)

	for {
		s, err := br.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		s = strings.TrimSpace(s)
		v, err := cache.Resolve(s)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(v)
	}
}
