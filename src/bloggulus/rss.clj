(ns bloggulus.rss
  (:require [remus]))

(defn read-blog[feed-url]
  (let [result (remus/parse-url feed-url)
        feed (:feed result)]
    {:feed-url feed-url
     :site-url (:link feed)
     :title (:title feed)}))

(defn read-posts [feed-url]
  (let [result (remus/parse-url feed-url)
        feed (:feed result)
        entries (:entries feed)]
    (map #(select-keys % [:link :title :updated-date]) entries)))

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
