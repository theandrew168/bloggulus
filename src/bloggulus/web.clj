(ns bloggulus.web
  (:gen-class)
  (:require [clojure.pprint :as pprint]
            [compojure.core :refer [routes GET]]
            [compojure.route :refer [resources]]
            [selmer.parser :as template]
            [ring.adapter.jetty :as server]
            [next.jdbc :as jdbc]
            [next.jdbc.connection :as connection]
            [bloggulus.db :as db])
  (:import (com.zaxxer.hikari HikariDataSource)))

(defn render-index [conn req]
  (let [posts (db/read-recent-posts conn 20)
        data {:authed true :posts posts}]
    (template/render-file "templates/index.html" data)))

(defn render-blogs [req]
  "blogs")

(defn render-request [req]
  (with-out-str (pprint/pprint req)))

(defn app-routes
  "Build app routes with the DB conn pool baked in."
  [conn]
  (routes
   (GET "/" [] (partial render-index conn))
   (GET "/blogs" [] render-blogs)
   (GET "/request" [] render-request)
   (resources "/static" {:root "static"})))

(defn -main []
  (let [port (Integer/parseInt (or (System/getenv "PORT") "5000"))
        db-url (or (System/getenv "BLOGGULUS_DATABASE_URL")
                   (throw (Exception. "missing env var: BLOGGULUS_DATABASE_URL")))
        jdbc-url (db/db-url->jdbc-url db-url)
        db-spec {:jdbcUrl jdbc-url}]
    (with-open [^HikariDataSource conn (connection/->pool HikariDataSource db-spec)]
      (db/migrate conn)
      (printf "Listening on 127.0.0.1:%s\n" port)
      (flush)
      (server/run-jetty (app-routes conn) {:host "127.0.0.1" :port port}))))

(comment
  (def db-url "postgresql://postgres:postgres@localhost:5432/postgres")
  (def jdbc-url (db/db-url->jdbc-url db-url))
  (def db-spec {:jdbcUrl jdbc-url})

  (def server
    (server/run-jetty
     (app-routes db-spec)
     {:host "127.0.0.1" :port 5000 :join? false}))
  (.stop server)

  (template/cache-off!)

  .)
