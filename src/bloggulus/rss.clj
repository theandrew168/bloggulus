(ns bloggulus.rss
  (:require [remus]
            [bloggulus.core :as core]))

(defn read-blog
  [feed-url]
  (let [result (remus/parse-url feed-url)
        feed (:feed result)
        {:keys [link title]} feed]
    (core/->Blog nil feed-url link title)))

(defn- entry->post
  [entry]
  (let [{:keys [link title updated-date]} entry
        preview "Lorem ipsum dolor sit, amet consectetur adipisicing elit."]
    (core/->Post nil nil link title preview updated-date)))

(defn read-posts
  [feed-url]
  (let [result (remus/parse-url feed-url)
        feed (:feed result)
        entries (:entries feed)]
    (map entry->post entries)))

(comment
  (def feed-url "https://nullprogram.com/feed/")
  (def result (remus/parse-url feed-url))
  (def feed (:feed result))
  (:title feed)

  (def entries (:entries feed))

  (:description feed)
  (:feed-type feed)
  (count (:entries feed))

  (:link (nth entries 0))
  (:title (nth entries 0))
  (:updated-date (nth entries 0))

  (read-blog feed-url)
  (read-posts feed-url)

  (def bad-feed-url "https://nullprogram.com/feexxx/")
  (read-blog bad-feed-url)
  (read-posts bad-feed-url)

  )
