package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/miku/doicache"
	log "github.com/sirupsen/logrus"
)

var (
	databaseDir = flag.String("db", filepath.Join(doicache.UserHomeDir(), ".doicache/default"), "leveldb directory")
	ttl         = flag.Duration("ttl", 24*time.Hour*240, "entry expiration")
	verbose     = flag.Bool("verbose", false, "be verbose")
	showVersion = flag.Bool("version", false, "show version")
	dumpKeys    = flag.Bool("dk", false, "dump keys")
	version     = "undefined"
)

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Printf("git-%s\n", version)
		os.Exit(0)
	}

	cache := doicache.New(*databaseDir)
	cache.Verbose = *verbose
	cache.TTL = *ttl

	if *dumpKeys {
		if err := cache.DumpKeys(os.Stdout); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	var reader io.Reader = os.Stdin

	if flag.NArg() > 0 {
		if _, err := os.Stat(flag.Arg(0)); os.IsNotExist(err) {
			for _, arg := range flag.Args() {
				v, err := cache.Resolve(arg)
				if err == doicache.ErrCannotResolve {
					log.Printf("cannot resolve %s, skipping", arg)
					continue
				}
				if err != nil {
					switch t := err.(type) {
					case doicache.ProtocolError:
						log.Printf("got HTTP %d, skipping", t.StatusCode)
						continue
					default:
						log.Fatal(err)
					}
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
		if err == doicache.ErrCannotResolve {
			log.Printf("cannot resolve %s, skipping", s)
			continue
		}
		if err != nil {
			switch t := err.(type) {
			case doicache.ProtocolError:
				log.Printf("got HTTP %d, skipping", t.StatusCode)
				continue
			default:
				log.Fatal(err)
			}
		}

		fmt.Println(v)
	}
}
