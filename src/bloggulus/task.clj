(ns bloggulus.task
  (:require [next.jdbc :as jdbc]
            [bloggulus.db :as db]
            [bloggulus.rss :as rss]))

(defn sync-blog [conn blog-id]
  (let [blog (jdbc/execute-one!
              conn ["SELECT * FROM blog WHERE blog_id = ?" blog-id])
        feed-url (:blog/feed_url blog)
        posts (rss/read-posts feed-url)]
    posts))

(defn prune-sessions [conn]
  (jdbc/execute! conn ["DELETE FROM session WHERE expiry <= now()"]))

(comment
  (def db-url (System/getenv "BLOGGULUS_DATABASE_URL"))
  (def jdbc-url (db/db-url->jdbc-url db-url))
  (def db-spec {:jdbcUrl jdbc-url})
  (def ds (jdbc/get-datasource jdbc-url))

  (with-open [conn (jdbc/get-connection ds)]
    (sync-blog conn 1))

  #_(doseq [migration missing]
    (apply-migration conn migration))

  ,)
