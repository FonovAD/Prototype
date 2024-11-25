package models

import "errors"

const (
	MAX_ORIGIN_LINK_LENGTH = 255
	MIN_ORIGIN_LINK_LENGTH = 0
)

type Link struct {
	UID                   int    `json:"uid"`
	OriginLink            string `json:"origin_link"`
	ShortLink             string `json:"short_link"`
	CreateTime            int    `json:"create_time"`
	ExpireTime            int    `json:"expire_time"`
	Status                string `json:"status"`
	ScheduledDeletionTime int    `json:"scheduled_deletion_time"`
}

var (
	OriginLinkLengthErr = errors.New("Origin link too long or too short (max length 255 characters, min length 0 characters)")
)

func (l *Link) Validate() error {
	if l.OriginLink > MAX_ORIGIN_LINK_LENGTH || l.OriginLink < MIN_ORIGIN_LINK_LENGTH {
		return OriginLinkLengthErr
	}
}
