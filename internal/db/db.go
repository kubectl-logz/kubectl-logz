package db

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/pebble"
	"github.com/kubectl-logz/kubectl-logz/internal/types"
)

type DB struct {
	db *pebble.DB
}

func NewDb() (*DB, error) {
	db, err := pebble.Open("logs", &pebble.Options{})
	return &DB{db}, err
}

func (d *DB) Set(lc types.Ctx, entry types.Entry, file string, offset int64) error {
	key := fmt.Sprintf("%s,%d", lc.Hostname, entry.Time.UnixMilli())
	value := fmt.Sprintf("%s,%d", file, offset)
	//	log.Printf("key=%s value=%s\n", key, value)
	_, closer, err := d.db.Get([]byte(key))
	if err != nil {
		return d.db.Set([]byte(key), []byte(value), pebble.NoSync)
	}
	defer closer.Close()
	return nil
}

func (d *DB) Close() {
	d.db.Close()
}

func (d *DB) Get(lc types.Ctx, time time.Time) (string, int64, io.Closer, error) {
	key := fmt.Sprintf("%s,%d", lc.Hostname, time.UnixMilli())
	value, closer, err := d.db.Get([]byte(key))
	if err != nil {
		return "", 0, nil, fmt.Errorf("failed to find (%v,%v): %w", lc, time, err)
	}
	//  log.Printf("key=%s value=%s\n", key, value)
	parts := strings.Split(string(value), ",")
	file := parts[0]
	offset, _ := strconv.ParseInt(parts[1], 10, 64)
	return file, offset, closer, nil

}
