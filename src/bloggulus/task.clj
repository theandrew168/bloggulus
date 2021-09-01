(ns bloggulus.task
  (:require [next.jdbc :as jdbc]
            [next.jdbc.date-time]
            [bloggulus.db :as db]
            [bloggulus.rss :as rss]))

(defn sync-blog
  "Sync all posts from a given blog into the database"
  [conn blog-id]
  (let [blog (db/blog-read conn blog-id)
        posts (rss/read-posts (:feed-url blog))]
    (doseq [post posts]
      (let [post (assoc post :blog-id blog-id)]
        (db/post-create conn post)))))

(defn prune-sessions
  "Delete all expired sessions from the database"
  [conn]
  (db/session-delete-expired conn))

(comment
  (def db-url "postgresql://postgres:postgres@localhost:5432/postgres")
  (def jdbc-url (db/db-url->jdbc-url db-url))
  (def db-spec {:jdbcUrl jdbc-url})

  (sync-blog db-spec 1)
  (prune-sessions db-spec)

  .)
