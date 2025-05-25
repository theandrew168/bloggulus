# Bloggulus Domains, Models, and Rules

## Commands

1. Create Account ?
2. Sign In ?
3. Sign Out ?
4. Add Blog
   1. Idempotently add and then "Follow Blog"
5. Follow Blog
6. Unfollow Blog
7. Sync One Blog
   1. Internal (cron)
8. Sync All Blogs
   1. Internal (cron)
9. Delete Blog
   1. Admin only
10. Delete Post
   1. Admin only
11. Create Post
12. Update Post
   1. Title
   2. Content
   3. PublishedAt
13. Delete Expired Sessions

## Queries

1. List articles (anon vs account)
2. Search articles (anon vs account)
3. List blogs
   1. Includes "Follow" status
4. List blog details
   1. Admin only
5. List posts
   1. Admin only
6. List post details
   1. Admin only

## Aggregates

1. Account
   1. Session
2. Blog
   1. Post ?
3. Tag

## Entities

1. Blog (unique by URL)
2. Post (unique by URL)
3. Account (unique by username)

## Value Objects

1. Tag (many Posts have many)
2. Session (one Account has many)

## Ports / Adapters (IO Boundaries)

1. Persistent Storage (a PostgreSQL database)
2. Fetching RSS Feeds (via the internet)
3. Fetching Web Pages (via the internet)
4. OAuth redirects / callbacks (via GitHub / Google)

## Errors

1. Feed was unreachable
2. Feed was invalid
3. Page was unreachable
4. Blog does not exist
5. Blog already exists
6. Post does not exist
7. Post already exists

## Rules

1. A Blog must have a valid, non-empty URL
2. A Blog's feed URL must be unique
3. A new Blog can be created from just a Feed URL
4. A Blog can be deleted by a site admin
5. If present, a Blog request should include ETag / LastModified headers
6. If a response includes ETag / LastModified headers, they should be stored with the blog
7. A Blog shouldn't be synced more than once per hour (or two)
8. A Post must have a valid, non-empty URL
9. A Post's URL must be unique
10. A Post can be deleted by a site admin
11. If a Feed doesn't include a Post's content, it will be fetched directly
12. All outgoing HTTP requests should include a proper user agent
13. When syncing RSS data, new Posts should be created
14. When syncing RSS data, existing Posts may be updated
15. Post content should be stripped of HTML tags before storing
16. Posts should be searchable by their content
17. An Account can add Blogs (by feed URL)
18. An Account can follow Blogs
19. An Account can unfollow Blogs
