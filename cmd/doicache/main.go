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

	"github.com/adrg/xdg"
	"github.com/miku/doicache"
	log "github.com/sirupsen/logrus"
)

var (
	databaseDir   = flag.String("db", filepath.Join(xdg.CacheHome, "doicache", "default"), "leveldb directory")
	ttl           = flag.Duration("ttl", 24*time.Hour*240, "entry expiration")
	verbose       = flag.Bool("verbose", false, "be verbose")
	showVersion   = flag.Bool("version", false, "show version")
	dumpKeys      = flag.Bool("dk", false, "dump keys")
	dumpKeyValues = flag.Bool("dkv", false, "dump keys and redirects")
	bestEffort    = flag.Bool("b", false, "best effort, just log errors, do not hang")

	NotAvailable = "NOTAVAILABLE"
	version      = "undefined"
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

	if *dumpKeyValues {
		if err := cache.DumpKeyValues(os.Stdout); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	var reader io.Reader = os.Stdin

	if flag.NArg() > 0 {
		if _, err := os.Stat(flag.Arg(0)); os.IsNotExist(err) {
			reader = strings.NewReader(strings.Join(flag.Args(), "\n") + "\n")
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
	var status string

	for {
		s, err := br.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		s = strings.TrimSpace(s)
		status = "OK"

		v, err := cache.Resolve(s)
		switch {
		case err == doicache.ErrCannotResolve:
			status = "NOR"
		case err == doicache.ErrInvalidURL:
			status = "EURL"
		case err == doicache.ErrMissingURLValue:
			status = "MURL"
		case err != nil:
			switch t := err.(type) {
			case doicache.ProtocolError:
				status = fmt.Sprintf("H%d", t.StatusCode)
			default:
				if *bestEffort {
					log.Println(err)
					continue
				}
				log.Fatal(err)
			}
		}
		if v == "" {
			v = NotAvailable
		}
		fmt.Printf("%s\t%s\t%s\n", status, s, v)
	}
}
