(ns bloggulus.core)

(defrecord Account [account-id username password email verified])
(defrecord AccountBlog [account-id blog-id])
(defrecord Blog [blog-id feed-url site-url title])
(defrecord Post [post-id url title preview updated blog])
(defrecord Session [session-id expiry account])
