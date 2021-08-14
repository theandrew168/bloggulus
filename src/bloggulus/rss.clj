(ns bloggulus.rss
  (:require [remus :refer [parse-url]]))

(comment
  (def feed-url "https://nullprogram.com/feed/")
  (def result (parse-url feed-url))
  (def feed (:feed result))
  (:title feed)

  (def entries (:entries feed))

  (:description feed)
  (:feed-type feed)
  (count (:entries feed))

  (:link (nth entries 0))
  (:title (nth entries 0))
  (:updated-date (nth entries 0))


  )
