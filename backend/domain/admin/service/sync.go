package service

import (
	"github.com/theandrew168/bloggulus/backend/domain/admin/feed"
	"github.com/theandrew168/bloggulus/backend/domain/admin/storage"
)

type SyncService struct {
	store       storage.Storage
	feedFetcher feed.FeedFetcher
	pageFetcher feed.PageFetcher
}

func NewSyncService(store storage.Storage, feedFetcher feed.FeedFetcher, pageFetcher feed.PageFetcher) *SyncService {
	s := SyncService{
		store:       store,
		feedFetcher: feedFetcher,
		pageFetcher: pageFetcher,
	}
	return &s
}

func (s *SyncService) SyncAllBlogs() error {
	return nil
}
