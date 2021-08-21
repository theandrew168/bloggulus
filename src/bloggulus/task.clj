(ns bloggulus.task
  (:require [next.jdbc :as jdbc]
            [next.jdbc.date-time]  ; import for transparent PG timestamp conversions
            [bloggulus.db :as db]
            [bloggulus.rss :as rss]))

(defn sync-blog
  "Sync all posts from a given blong into the database"
  [conn blog-id]
  (let [blog (jdbc/execute-one!
              conn ["SELECT * FROM blog WHERE blog_id = ?" blog-id])
        feed-url (:blog/feed_url blog)
        posts (rss/read-posts feed-url)]
    (doseq [post posts]
      (let [preview "Lorem ipsum dolor sit, amet consectetur adipisicing elit."]
        (jdbc/execute! conn ["
            INSERT INTO post
              (blog_id, url, title, preview, updated)
            VALUES
              (?,?,?,?,?)
            ON CONFLICT DO NOTHING
          " blog-id (:link post) (:title post) preview (:updated-date post)])))))

(defn prune-sessions
  "Delete all expired sessions from the database"
  [conn]
  (jdbc/execute! conn ["DELETE FROM session WHERE expiry <= now()"]))

(comment
  (def db-url (System/getenv "BLOGGULUS_DATABASE_URL"))
  (def jdbc-url (db/db-url->jdbc-url db-url))
  (def db-spec {:jdbcUrl jdbc-url})
  (def ds (jdbc/get-datasource jdbc-url))

  (with-open [conn (jdbc/get-connection ds)]
    (sync-blog conn 1))

  .)
