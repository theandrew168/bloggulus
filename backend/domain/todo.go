package domain

// What all should live in here?
// Storage interfaces? Feed interfaces?
// Basically anything that the core domain depends on.
// Impls can always look it and impl the interfaces outlined in here

// Two BCs:
// 1. Admin - Works with models at a granular level, handles syncing feeds
//  	and touches models individually to ensure the correct state is stored.
// 2. Reader - Uses the app as regular user, a "post" is a title, link, blog,
//		blog link, and tags. Users can search posts by title/content. Read only.

// domain
//   admin (routed at /api/v1/admin/...)
//     models at the top-level - NewPost / LoadPost
//     storage - storage.Storage interface
//     feeds (fetch feed, fetch page) - feed.FeedFetcher, feed.PageFetcher
//     services (sync) - SyncService.SyncBlog(), SyncService.SyncAllBlogs()
//   reader (default, routed at /api/v1/...)
//	   models at the top-level - LoadPost
//     storage - storage.Storage interface (reads / searches blogs)
//	   NO feed management, NO services

// https://github.com/theandrew168/bloggulus-svelte/blob/main/src/lib/server/feed.ts
// https://github.com/theandrew168/bloggulus-svelte/blob/main/src/lib/server/fetch.ts

type FooStorage struct{}

type Storage interface {
	Blog() FooStorage
	Post() FooStorage
	Tag() FooStorage
	WithTransaction(operation func(store *Storage) error) error
	WithAtomic(operation func(store *Storage) error) error
	Atomically(operation func(store *Storage) error) error
}
