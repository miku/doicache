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

	"github.com/sethgrid/pester"
	"github.com/syndtr/goleveldb/leveldb"
)

// ProtocolError keeps HTTP status codes.
type ProtocolError struct {
	Location   string
	Message    string
	StatusCode int
}

func (e ProtocolError) Error() string {
	return fmt.Sprintf("HTTP %d %s %s", e.StatusCode, e.Location)
}

// Response from doi.org/api/handles endpoint.
type Response struct {
	Handle       string `json:"handle"`
	ResponseCode int64  `json:"responseCode"`
	Values       []struct {
		Data      interface{} `json:"data"`
		Index     int64       `json:"index"`
		Timestamp string      `json:"timestamp"`
		Ttl       int64       `json:"ttl"`
		Type      string      `json:"type"`
	} `json:"values"`
}

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
	resp, err := pester.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, ProtocolError{StatusCode: resp.StatusCode, Location: u}
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

// Resolve returns the redirect URL for a given DOI.
func (c *Cache) Resolve(doi string) (string, error) {
	b, err := c.Get(doi)
	if err != nil {
		return "", err
	}
	var resp Response
	if err := json.Unmarshal(b, &resp); err != nil {
		return "", err
	}
	for _, v := range resp.Values {
		if v.Type != "URL" {
			continue
		}
		switch t := v.Data.(type) {
		case map[string]interface{}:
			if v, ok := t["value"]; ok {
				return fmt.Sprintf("%s", v), nil
			} else {
				return "", fmt.Errorf("value key missing")
			}
		default:
			return "", fmt.Errorf("unexpected payload for URL type: %T", v.Data)
		}
	}
	return "", fmt.Errorf("could not resolve %s", doi)
}
