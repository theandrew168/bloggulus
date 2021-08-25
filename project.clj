(defproject bloggulus "0.0.1"
  :description "A community for bloggers and readers"
  :url "https://github.com/theandrew168/bloggulus"
  :license {:name "MIT License"
            :url "https://spdx.org/licenses/MIT.html"}
  :dependencies [[org.clojure/clojure "1.10.1"]
                 [com.github.seancorfield/next.jdbc "1.2.689"]
                 [org.postgresql/postgresql "42.2.23"]
                 [com.zaxxer/HikariCP "4.0.3"]
                 [ring/ring-core "1.9.4"]
                 [ring/ring-jetty-adapter "1.9.4"]
                 [compojure/compojure "1.6.2"]
                 [selmer "1.12.44"]
                 [remus "0.2.2"]]
  :auto-clean false
  :target-path "target/%s"
  :profiles {:uberjar {:jvm-opts ["-Dclojure.compiler.direct-linking=true"]}
             :web {:main bloggulus.web
                   :aot [bloggulus.web]
                   :uberjar-name "bloggulus-web.jar"}
             :worker {:main bloggulus.worker
                      :aot [bloggulus.worker]
                      :uberjar-name "bloggulus-worker.jar"}
             :scheduler {:main bloggulus.scheduler
                         :aot [bloggulus.scheduler]
                         :uberjar-name "bloggulus-scheduler.jar"}})