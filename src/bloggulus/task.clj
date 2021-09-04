(ns bloggulus.task
  (:require [next.jdbc :as jdbc]
            [next.jdbc.date-time]
            [bloggulus.db :as db]
            [bloggulus.rss :as rss]))

(def queue-name "bloggulus-queue")

(defn submit
  [conn task & args]
  (let [message {:task task :args args}]
    (wcar redis-spec
          (car-mq/enqueue queue-name message))))

(defn sync-blog
  "Sync all posts from a given blog into the database"
  [conn blog-id]
  (let [blog (db/blog-read conn blog-id)
        posts (rss/read-posts (:feed-url blog))]
    (doseq [post posts]
      (let [post (assoc post :blog-id blog-id)]
        (db/post-create conn post)))))

(defn prune-sessions
  "Delete all expired sessions from the database"
  [conn]
  (db/session-delete-expired conn))

(comment
  (def db-url "postgresql://postgres:postgres@localhost:5432/postgres")
  (def jdbc-url (db/db-url->jdbc-url db-url))
  (def db-spec {:jdbcUrl jdbc-url})

  (sync-blog db-spec 1)
  (prune-sessions db-spec)

  (require '[taoensso.carmine :as car :refer [wcar]])
  (require '[taoensso.carmine.message-queue :as car-mq])
  (def redis-url "redis://localhost:6379")
  (def redis-spec {:pool {} :spec {:uri redis-url}})

  (wcar redis-spec (car/ping))

  (defn worker-handler
    [{:keys [message attempt]}]
    {:status :success})

  (def my-worker
    (car-mq/worker redis-spec "bloggulus-queue"
                   {:handler worker-handler}))

  (car-mq/stop my-worker)

  (wcar redis-spec (car-mq/enqueue "bloggulus-queue" "my message!"))

  (submit redis-spec :prune-sessions)

  (def tasks
    {:sync-blog sync-blog
     :prune-sessions prune-sessions})

  ;; TODO: make each task take a map of args (always throw the db-conn in)
  (defn worker
    [db-conn redis-conn]
    (car-mq/worker
     redis-conn
     queue-name
     {:handler (fn [{{:keys [task args]} :message}]
                 (let [args (cons db-conn args)
                       f (get tasks task)]
                   (apply f args)))}))

  (def bloggulus-worker (worker db-spec redis-spec))
  (car-mq/stop bloggulus-worker)

  .)
