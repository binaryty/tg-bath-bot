package storage

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/yellowpuki/tg-bath-bot/internal/er"
)

var ErrNoRecords = errors.New("no saved records")

type Storage interface {
	Save(ctx context.Context, r *Record) error
	IsExist(ctx context.Context, h string) (bool, error)
	LastVisit(ctx context.Context, user string) (time.Time, error)
}

type Record struct {
	EventToken string
	UserName   string
	CreatedAt  time.Time
}

// EventHash ...
func (r Record) EventHash() (string, error) {
	hash := sha1.New()

	if _, err := io.WriteString(hash, r.UserName); err != nil {
		return "", er.Wrap("can't calculate hash", err)
	}

	month := time.Now().Month().String()

	if _, err := io.WriteString(hash, month); err != nil {
		return "", er.Wrap("can't calculate hash", err)
	}

	year := strconv.Itoa(time.Now().Year())

	if _, err := io.WriteString(hash, year); err != nil {
		return "", er.Wrap("can't calculate hash", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}