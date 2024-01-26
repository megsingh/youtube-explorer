package models

import "time"

type Video struct {
	Title           string    `json:"title"`
	ChannelId       string    `json:"channelId"`
	ChannelTitle    string    `json:"channelTitle"`
	Description     string    `json:"description"`
	PublishedAt     time.Time `json:"published_at"`
	ThumbnailUrl    string    `json:"thumbnail_url"`
	PaginationToken string    `json:"paginationToken"`
}

type PaginationResponse struct {
	Videos          []Video `json:"videos"`
	HasNext         bool    `json:"has_next"`
	PaginationToken string  `json:"token"`
}
