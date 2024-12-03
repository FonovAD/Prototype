package models

import (
	"errors"
	"math/rand"
	"regexp"
	"time"
)

const (
	MAX_ORIGIN_LINK_LENGTH = 255
	MIN_ORIGIN_LINK_LENGTH = 0
)

const (
	STATUS_ACTIVE  = "active"
	STATUS_EXPIRED = "expired"
	STATUS_DELETED = "deleted"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type Link struct {
	UID                   int    `json:"uid"`
	OriginLink            string `json:"origin_link"`
	ShortLink             string `json:"short_link"`
	CreateTime            int64  `json:"create_time"`
	ExpireTime            int64  `json:"expire_time"`
	Status                string `json:"status"`
	ScheduledDeletionTime int64  `json:"scheduled_deletion_time"`
}

var (
	OriginLinkLengthErr = errors.New("Origin link too long or too short (max length 255 characters, min length 0 characters)")
)

func (l *Link) Validate() error {
	if len(l.OriginLink) > MAX_ORIGIN_LINK_LENGTH || len(l.OriginLink) < MIN_ORIGIN_LINK_LENGTH {
		return OriginLinkLengthErr
	}
	return nil
}

func IsValidURL(url string) bool {
	regex := regexp.MustCompile(`^(http|https|ftp)://[a-zA-Z0-9\-+&@#/%?=~_|!:,.;]*[a-zA-Z0-9\-+&@#/%=~_|]$`)
	return regex.MatchString(url)
}

func GenerateShortLink(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(bytes)
}
