(ns bloggulus.db
  (:require [clojure.java.io :as io]
            [clojure.set :as set]
            [clojure.string :as string]
            [next.jdbc :as jdbc])
  (:import (java.net URI)))

(defn db-url->jdbc-url
  "Convert a libpq/python/golang database URL into a JDBC URL."
  [db-url]
  (let [uri (URI. db-url)
        [username password] (string/split (.getUserInfo uri) #":")]
    (format
     "jdbc:%s://%s:%s%s?user=%s&password=%s"
     (.getScheme uri)
     (.getHost uri)
     (.getPort uri)
     (.getPath uri)
     username
     password)))

(defn- create-migrations-table [conn]
  (jdbc/execute! conn ["
    CREATE TABLE IF NOT EXISTS migration (
      migration_id SERIAL PRIMARY KEY,
      name TEXT NOT NULL UNIQUE
    )"]))

(defn- read-migrations []
  (-> "migrations/migrations.txt"
      (io/resource)
      (slurp)
      (string/split-lines)
      (set)))

(defn- read-applied-migrations [conn]
  (let [query "SELECT name FROM migration"]
    (into #{} (map :name) (jdbc/plan conn [query]))))

(defn- apply-migration [conn name]
  (println "applying migration:" name)
  (let [path (string/join "/" ["migrations" name])
        migration (slurp (io/resource path))
        insert "INSERT INTO migration (name) VALUES (?)"]
    (jdbc/execute! conn [migration])
    (jdbc/execute! conn [insert name])))

(defn migrate
  "Apply any un-applied database migrations."
  [conn]
  (create-migrations-table conn)
  (let [migrations (read-migrations)
        applied (read-applied-migrations conn)
        missing (sort (set/difference migrations applied))]
    (doseq [migration missing]
      (apply-migration conn migration))
    missing))

(def ^:private query-read-recent-posts
  "SELECT
     post.*,
     blog.title
   FROM post
   INNER JOIN blog
     ON blog.blog_id = post.blog_id
   ORDER BY post.updated DESC
   LIMIT ?")

(defn create-blog
  [conn blog]
  (let [feed-url (:feed-url blog)
        site-url (:site-url blog)
        title (:title blog)
        query "INSERT INTO blog (feed_url, site_url, title) VALUES (?,?,?) RETURNING blog_id"
        result (jdbc/execute-one! conn [query feed-url site-url title])
        blog-id (:blog/blog_id result)]
    (assoc blog :blog-id blog-id)))

(defn read-recent-posts
  "List the N most recent blog posts."
  [conn n]
  (jdbc/execute! conn [query-read-recent-posts n]))

(comment
  (def db-url "postgresql://postgres:postgres@localhost:5432/postgres")
  (def jdbc-url (db-url->jdbc-url db-url))
  (def db-spec {:jdbcUrl jdbc-url})
  (def conn (jdbc/get-datasource db-spec))

  (list-migrations)
  (list-applied-migrations conn)
  (migrate conn)

  (def blog {:feed-url "https://nullprogram.com/feed/"
             :site-url "https://nullprogram.com/"
             :title "null program"})
  (create-blog conn blog)

  (read-recent-posts conn 5)

  .)
