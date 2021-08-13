(ns bloggulus.db
  (:require [next.jdbc :as jdbc]
            [next.jdbc.connection :as connection])
  (:import (com.zaxxer.hikari HikariDataSource)))

(defn db-url-to-jdbc-url [db-url]
  (let [uri (java.net.URI. db-url)
        [username password] (clojure.string/split (.getUserInfo uri) #":")]
    (format
     "jdbc:%s://%s:%s%s?user=%s&password=%s"
     (.getScheme uri)
     (.getHost uri)
     (.getPort uri)
     (.getPath uri)
     username
     password)))

(defn ^:private create-migrations-table [conn]
  (jdbc/execute! conn ["
    CREATE TABLE IF NOT EXISTS migration (
      migration_id SERIAL PRIMARY KEY,
      name TEXT NOT NULL UNIQUE
    )"]))

(defn ^:private list-migrations []
  (-> "migrations/migrations.txt"
      (clojure.java.io/resource)
      (slurp)
      (clojure.string/split-lines)
      (set)))

(defn ^:private list-applied-migrations [conn]
  (let [query "SELECT name FROM migration"]
    (into #{} (map :name) (jdbc/plan conn [query]))))

(defn ^:private apply-migration [conn name]
  (println "applying migration:" name)
  (let [path (clojure.string/join "/" ["migrations" name])
        migration (slurp (clojure.java.io/resource path))
        insert "INSERT INTO migration (name) VALUES (?)"]
    (jdbc/execute! conn [migration])
    (jdbc/execute! conn [insert name])))

(defn migrate [conn]
  (create-migrations-table conn)
  (let [migrations (list-migrations)
        applied (list-applied-migrations conn)
        missing (sort (clojure.set/difference migrations applied))]
    (doall (map #(apply-migration conn %) missing))
    missing))
