# Bloggulus Domains, Models, and Rules

## Commands

2. Sign In (returns session)
3. Sign Out
5. Follow Blog
6. Unfollow Blog
7. Add / Sync Blog
   1. Internal (cron)
8. Sync All Blogs
   1. Internal (cron)
9. Delete Blog
   1. Admin only
10. Delete Post
   1. Admin only
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

## Models

1. Account
   1. Followed blogs
2. Session
3. Blog
4. Post
5. Tag

## Ports / Adapters (IO Boundaries)

1. Persistent Storage (a PostgreSQL database)
2. Fetching RSS Feeds (via the internet)
4. OAuth redirects / callbacks (via GitHub / Google)

## Errors

1. Feed was unreachable
2. Feed was invalid
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

## Deployment

1. JS code needs to be pulled and built, so just do that for both?
2. Configure MFD to:
   1. Build the backend Go API
   2. Build the frontend SvelteKit project
   3. Flip the symlink
   4. Cleanup any old dirs
3. Restart the bloggulus-api service
   1. Points to /usr/local/src/bloggulus/active/backend/bloggulus file
4. Restart the bloggulus-web service
   1. Points to /usr/local/src/bloggulus/active/frontend/build directory

Orrrr....

1. Build the FE to an output build/ directory
2. Embed the whole FE output into the single Go binary
3. At startup, write the embedded build dir to temp using os.CopyFS
4. Exec a NodeJS process to run the frontend (in a background goro)
5. Run the Go API like normal
6. This might be terrible. Doesn't handle systemd sockets (or restarts) well.

Orrrr....

1. Just do the first strat but without mfd?
2. Checkout the new commit into a fresh dir
3. Build the backend Go API
4. Build the frontend SvelteKit project
5. Flip the symlink
6. Cleanup any old dirs
7. Restart the bloggulus-api service
   1. Points to /usr/local/src/bloggulus/active/backend/bloggulus file
8. Restart the bloggulus-web service
   1. Points to /usr/local/src/bloggulus/active/frontend/build directory
9. Cleanup any old dirs