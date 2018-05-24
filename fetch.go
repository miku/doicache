package doicache

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

// Entry to cache. Contains raw bytes of response and some metadata.
type Entry struct {
	Date time.Time
	Blob []byte
}

// Cache wraps the backend. XXX: Try to mitigate hot DNS servers by hardcoding
// a few of the doi.org IPs.
type Cache struct {
	Endpoint string
	TTL      time.Duration
	Verbose  bool
	name     string
	db       *leveldb.DB
}

// New returns a new cache read to be queried.
func New(filename string) *Cache {
	return NewTTL(filename, 0)
}

// NewTTL creates a new cache with a default expiration date.
func NewTTL(filename string, ttl time.Duration) *Cache {
	return &Cache{name: filename, TTL: ttl, Endpoint: "https://doi.org/api/handles"}
}

// openDatabase will open or create the database.
func (c *Cache) openDatabase() error {
	if c.db != nil {
		return nil
	}
	db, err := leveldb.OpenFile(c.name, nil)
	if err != nil {
		return err
	}
	c.db = db
	return nil
}

// Close the underlying resources.
func (c *Cache) Close() error {
	err := c.db.Close()
	if err != nil {
		return err
	}
	c.db = nil
	return nil
}

// Name returns the path to the database file.
func (c *Cache) Name() string {
	return c.name
}

func (c *Cache) Get(key string) ([]byte, error) {
	if err := c.openDatabase(); err != nil {
		return nil, err
	}
	b, err := c.db.Get([]byte(key), nil)
	if err == leveldb.ErrNotFound {
		if c.Verbose {
			log.Println("cache miss")
		}
		return c.fetch(key)
	}
	var entry Entry
	if err := json.Unmarshal(b, &entry); err != nil {
		return nil, err
	}
	if entry.Date.Add(c.TTL).Before(time.Now()) {
		if c.Verbose {
			log.Println("entry expired")
		}
		return c.fetch(key)
	}
	return entry.Blob, err
}

// Fetch object from server.
func (c *Cache) fetch(key string) ([]byte, error) {
	if err := c.openDatabase(); err != nil {
		return nil, err
	}
	u := fmt.Sprintf("%s/%s", strings.TrimRight(c.Endpoint, "/"), key)
	if c.Verbose {
		log.Println(u)
	}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Failed with HTTP %d: %s", resp.StatusCode, u)
	}
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, resp.Body); err != nil {
		return nil, err
	}
	entry := Entry{
		Date: time.Now(),
		Blob: buf.Bytes(),
	}
	b, err := json.Marshal(entry)
	if err != nil {
		return nil, err
	}
	if c.Verbose {
		log.Println(string(b))
	}
	// XXX: Annotate bytes with date.
	return buf.Bytes(), c.db.Put([]byte(key), b, nil)

}
