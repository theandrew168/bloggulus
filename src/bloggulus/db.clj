(ns bloggulus.db
  (:require [clojure.java.io :refer [resource]]
            [clojure.set :refer [difference]]
            [clojure.string :refer [join split split-lines]]
            [next.jdbc :as jdbc]
            [next.jdbc.connection :as connection])
  (:import (com.zaxxer.hikari HikariDataSource)
           (java.net URI)))

(defn db-url->jdbc-url [db-url]
  (let [uri (URI. db-url)
        [username password] (split (.getUserInfo uri) #":")]
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

(defn- list-migrations []
  (-> "migrations/migrations.txt"
      (resource)
      (slurp)
      (split-lines)
      (set)))

(defn- list-applied-migrations [conn]
  (let [query "SELECT name FROM migration"]
    (into #{} (map :name) (jdbc/plan conn [query]))))

(defn- apply-migration [conn name]
  (println "applying migration:" name)
  (let [path (join "/" ["migrations" name])
        migration (slurp (resource path))
        insert "INSERT INTO migration (name) VALUES (?)"]
    (jdbc/execute! conn [migration])
    (jdbc/execute! conn [insert name])))

(defn migrate [conn]
  (create-migrations-table conn)
  (let [migrations (list-migrations)
        applied (list-applied-migrations conn)
        missing (sort (difference migrations applied))]
    (doall (map #(apply-migration conn %) missing))
    missing))
