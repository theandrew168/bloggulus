(ns bloggulus.db
  (:require [clojure.java.io :as io]
            [clojure.set :as set]
            [clojure.string :as string]
            [next.jdbc :as jdbc]
            [bloggulus.core :as core])
  (:import (java.net URI)))


;;
;; helpers
;;

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


;;
;; migrations
;;

(def ^:private stmt-create-migration-table
  "CREATE TABLE IF NOT EXISTS migration (
     migration_id SERIAL PRIMARY KEY,
     name TEXT NOT NULL UNIQUE
  )")

(def ^:private stmt-read-migrations
  "SELECT name FROM migration")

(def ^:private stmt-create-migration
  "INSERT INTO migration (name) VALUES (?)")

(defn- read-migrations []
  (-> "migrations/migrations.txt"
      (io/resource)
      (slurp)
      (string/split-lines)
      (set)))

(defn- read-applied-migrations [conn]
  (into #{} (map :name) (jdbc/plan conn [stmt-read-migrations])))

(defn- apply-migration [conn name]
  (println "applying migration:" name)
  (let [path (string/join "/" ["migrations" name])
        migration (slurp (io/resource path))]
    (jdbc/execute! conn [migration])
    (jdbc/execute! conn [stmt-create-migration name])))

(defn migrate
  "Apply any un-applied database migrations."
  [conn]
  (jdbc/execute! conn [stmt-create-migration-table])
  (let [migrations (read-migrations)
        applied (read-applied-migrations conn)
        missing (sort (set/difference migrations applied))]
    (doseq [migration missing]
      (apply-migration conn migration))
    missing))


;;
;; blog
;;

(defn- blog-rename
  [row]
  (set/rename-keys
   row {:blog/blog_id :blog-id
        :blog/feed_url :feed-url
        :blog/site_url :site-url
        :blog/title :title}))

(def ^:private stmt-blog-create
  "INSERT INTO blog (
    feed_url,
    site_url,
    title
  )
  VALUES (?,?,?)
  RETURNING blog_id")

(defn blog-create
  "Create a new blog."
  [conn {:keys [feed-url site-url title] :as blog}]
  (let [row (jdbc/execute-one!
             conn [stmt-blog-create feed-url site-url title])
        blog-id (:blog/blog_id row)]
    (assoc blog :blog-id blog-id)))

(def ^:private stmt-blog-read
  "SELECT
     blog_id,
     feed_url,
     site_url,
     title
   FROM blog
   WHERE blog_id = ?")

(defn blog-read
  "Read a blog by id."
  [conn blog-id]
  (let [row (jdbc/execute-one! conn [stmt-blog-read blog-id])]
    (-> row blog-rename core/map->Blog)))

(def ^:private stmt-blog-read-all
  "SELECT
     blog_id,
     feed_url,
     site_url,
     title
   FROM blog")

(defn blog-read-all
  "Read all blogs."
  [conn]
  (let [rows (jdbc/execute! conn [stmt-blog-read-all])]
    (map (comp blog-rename core/map->Blog) rows)))


;;
;; post
;;

(defn- post-rename
  [row]
  (set/rename-keys
   row {:post/post_id :post-id
        :post/blog_id :blog-id
        :post/url :url
        :post/title :title
        :post/preview :preview
        :post/updated :updated}))

(def ^:private stmt-post-create
  "INSERT INTO post (
    blog_id,
    url,
    title,
    preview,
    updated
  ) VALUES (?,?,?,?,?)
  ON CONFLICT DO NOTHING
  RETURNING post_id")

(defn post-create
  "Create a new post."
  [conn {:keys [url title preview updated blog] :as post}]
  (let [blog-id (:blog-id blog)
        row (jdbc/execute-one!
             conn [stmt-post-create blog-id url title preview updated])
        post-id (:post/post_id row)]
    (assoc post :post-id post-id)))

(def ^:private stmt-post-read-recent
  "SELECT
     post.post_id,
     post.url,
     post.title,
     post.preview,
     post.updated,
     blog.blog_id,
     blog.feed_url,
     blog.site_url,
     blog.title
   FROM post
   INNER JOIN blog
     ON blog.blog_id = post.blog_id
   ORDER BY post.updated DESC
   LIMIT ?")

(defn post-read-recent
  "Read the N most recent blog posts."
  [conn n]
  (let [rows (jdbc/execute! conn [stmt-post-read-recent n])]
    (map (comp post-rename core/map->Post) rows)))


;;
;; session
;;

(def ^:private stmt-session-delete-expired
  "DELETE
   FROM session
   WHERE expiry <= now()")

(defn session-delete-expired
  "Delete expired sessions."
  [conn]
  (jdbc/execute-one! conn [stmt-session-delete-expired]))



(comment
  (def db-url "postgresql://postgres:postgres@localhost:5432/postgres")
  (def jdbc-url (db-url->jdbc-url db-url))

  ;; works as a conn for development / testing
  (def db-spec {:jdbcUrl jdbc-url})

  (read-migrations)
  (read-applied-migrations db-spec)

  (migrate db-spec)

  (def blog {:feed-url "https://nullprogram.com/feed/"
             :site-url "https://nullprogram.com/"
             :title "null program"})
  (create-blog db-spec blog)

  (blog-read db-spec 1)
  (blog-read-all db-spec)
  (post-read-recent db-spec 5)

  .)
