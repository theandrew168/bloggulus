(ns bloggulus.db
  (:require [clojure.java.io :as io]
            [clojure.set :as set]
            [clojure.string :as string]
            [next.jdbc :as jdbc])
  (:import (java.net URI)))


;; helpers


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

;; migrations

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

;; data abstraction layer

(def ^:private stmt-create-account
  "INSERT INTO account (
     username,
     password,
     email
   ) VALUES (
     ?,?,?
   ) RETURNING account_id")

(defn create-account
  [conn account]
  (let [username (:username account)
        password (:password account)
        email (:email account)
        stmt stmt-create-account
        row (jdbc/execute-one! conn [stmt username password email])
        account-id (:account/account_id row)]
    (assoc account :account-id account-id)))

(def ^:private stmt-read-recent-posts
  "SELECT
     post.post_id,
     post.blog_id,
     post.url,
     post.title,
     post.preview,
     post.updated,
     blog.title
   FROM post
   INNER JOIN blog
     ON blog.blog_id = post.blog_id
   ORDER BY post.updated DESC
   LIMIT ?")

(defn read-recent-posts
  "List the N most recent blog posts."
  [conn n]
  (let [stmt stmt-read-recent-posts
        rows (jdbc/execute! conn [stmt n])]
    (map ; apply naming conversions
     #(set/rename-keys
       % {:post/post_id :post-id
          :post/blog_id :blog-id
          :post/url :url
          :post/title :title
          :post/preview :preview
          :post/updated :updated
          :blog/title :blog-title}) rows)))

(def ^:private stmt-create-blog
  "INSERT INTO blog (
    feed_url,
    site_url,
    title
  ) VALUES (
    ?,?,?
  ) RETURNING blog_id")

(defn create-blog
  [conn blog]
  (let [feed-url (:feed-url blog)
        site-url (:site-url blog)
        title (:title blog)
        stmt stmt-create-blog
        result (jdbc/execute-one! conn [stmt feed-url site-url title])
        blog-id (:blog/blog_id result)]
    (assoc blog :blog-id blog-id)))

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

  (read-recent-posts db-spec 5)

  .)
