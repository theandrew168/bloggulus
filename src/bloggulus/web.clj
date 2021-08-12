(ns bloggulus.web
  (:gen-class)
  (:require [clojure.pprint :refer [pprint]]
            [compojure.core :refer [defroutes GET]]
            [compojure.route :refer [resources]]
            [selmer.parser :as tmpl]
            [org.httpkit.server :refer [run-server]]))

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
              (or (System/getenv "PORT") "5000"))]
    (printf "Listening on 127.0.0.1:%s\n" port)
    (flush)
    (run-server #'app {:host "127.0.0.1" :port port})))
