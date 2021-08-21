(ns bloggulus.web
  (:gen-class)
  (:require [clojure.pprint :refer [pprint]]
            [compojure.core :refer [defroutes GET]]
            [compojure.route :refer [resources]]
            [selmer.parser :as tmpl]
            [ring.adapter.jetty :refer [run-jetty]]
            [next.jdbc.connection :as connection]
            [bloggulus.db :as db])
  (:import (com.zaxxer.hikari HikariDataSource)))

(defn render-index [req]
  (tmpl/render-file "templates/index.html" {:authed true}))

(defn render-blogs [req]
  "blogs")

(defn render-request [req]
  (with-out-str (pprint req)))

(defroutes app
  (GET "/" [] render-index)
  (GET "/blogs" [] render-blogs)
  (GET "/request" [] render-request)
  (resources "/static" {:root "static"}))

(defn -main []
  (let [port (Integer/parseInt
               (or (System/getenv "PORT") "5000"))
        db-url (or (System/getenv "BLOGGULUS_DATABASE_URL")
                   (throw (Exception. "missing env var: BLOGGULUS_DATABASE_URL")))
        jdbc-url (db/db-url->jdbc-url db-url)
        db-spec {:jdbcUrl jdbc-url}]
    (with-open [^HikariDataSource conn (connection/->pool HikariDataSource db-spec)]
      (db/migrate conn)
      (printf "Listening on 127.0.0.1:%s\n" port)
      (flush)
      (run-jetty #'app {:host "127.0.0.1" :port port}))))
