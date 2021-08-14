(ns bloggulus.db-test
  (:require [clojure.test :refer [deftest is]]
            [bloggulus.db :as db]))

(deftest test-db-url->jdbc-url
  (let [db-url "postgresql://postgres:postgres@localhost:5432/postgres"
        jdbc-url "jdbc:postgresql://localhost:5432/postgres?user=postgres&password=postgres"]
    (is (= (db/db-url->jdbc-url db-url)
           jdbc-url))))
