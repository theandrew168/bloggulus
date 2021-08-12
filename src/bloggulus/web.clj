(ns bloggulus.web
  (:gen-class)
  (:require [compojure.core :refer [defroutes GET]]
            [compojure.route :refer [files]]
            [hiccup.core :refer [html]]
            [hiccup.page :refer [html5]]
            [org.httpkit.server :refer [run-server]]))

(defn head [title]
  [:head
   [:title title]
   [:meta {:charset "utf-8"}]
   [:meta {:name "description" :content title}]
   [:meta {:name "viewport" :content "width=device-width, initial-scale=1.0"}]
   [:link {:rel "stylesheet" :href "/static/css/tailwind.min.css"}]
   [:script {:src "/static/js/alpine.min.js" :defer "true"}]
   [:style
    "@import url('https://fonts.googleapis.com/css?family=Karla:400,700&display=swap');"
    ".font-family-karla { font-family: karla; }"]])

(defn footer []
  [:footer.text-gray-100.bg-gray-800
   [:div.max-w-3xl.mx-auto.py-4
    [:div.flex.flex-col.items-center.justify-between.md:flex-row.space-y-1.md:space-y-0
     [:a.text-xl.font-bold.text-gray-100.hover:text-gray-400 {:href "/"} "Bloggulus"]
     [:a.text-gray-100.hover:text-gray-400 {:href "https://shallowbrooksoftware.com"} "Shallow Brook Software"]]]])

(defn render-index [req]
  (html5
   (head "Bloggulus - Index")
   (footer)))

(defn render-blogs [req]
  (html5
   (head "Bloggulus - Blogs")
   (footer)))

(defroutes app
  (GET "/" [] render-index)
  (GET "/blogs" [] render-blogs)
  (files "/static" {:root "static"}))

(def port (Integer/parseInt
           (or (System/getenv "PORT")
               "5000")))

(defn -main []
  (printf "Listening on 127.0.0.1:%s\n" port)
  (flush)
  (run-server #'app {:host "127.0.0.1" :port port}))

(comment

  (html [:span {:class "foo"} "bar"])
  (html [:span#foo "bar"])
  (html [:span.foo "bar"])
  (html5 [:span {:class "foo"} "bar"])

  (html (footer))
  (index)


  (html5 {:lang "en"} (head "Bloggulus"))

  ,)
