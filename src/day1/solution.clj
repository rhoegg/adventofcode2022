(ns day1.solution
  (:require [clojure.string :as s])
  (:require [clojure.core.reducers :as r]))
(defn parse-int [s] (Integer/parseInt s))
(defn parse-elves
  [path]
  (->
   (slurp path)
   (s/split #"\n\n")))
(defn parse-elf 
  [elfdata]
  (->>
   (s/split elfdata #"\n")
   (map parse-int)))
(defn parse-calories
  [path]
  (map parse-elf (parse-elves path))) 

(defrecord solution [part1 part2])
(defn solve
  [path]
  (def elf-calories (map #(r/fold + %) (parse-calories path)))
  (->solution
   (apply max elf-calories)
   (r/fold + (take 3 (sort > elf-calories)))))
(solve "day1/example.txt")
(solve "day1/input.txt")
