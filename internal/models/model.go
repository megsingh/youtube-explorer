package models

import "time"

type Video struct {
	Title           string    `json:"title"`
	ChannelId       string    `json:"channelid"`
	ChannelTitle    string    `json:"channeltitle"`
	Description     string    `json:"description"`
	PublishedAt     time.Time `json:"publishedat"`
	ThumbnailUrl    string    `json:"thumbnailurl"`
	PaginationToken string    `json:"paginationtoken"`
}

type PaginationResponse struct {
	Videos          []Video `json:"videos"`
	HasNext         bool    `json:"has_next"`
	PaginationToken string  `json:"token"`
}
