(ns day4.core
  (:gen-class))

(require '[clojure.string :as s])
(defn parse-int [s] (Integer/parseInt s))

(defrecord section [min max])

(defn load-section-assignments
  [file-name]
  (->
   file-name
   (slurp)
   (s/split #"\n")))

(defn parse-section-assignment
  [s]
  (let [[min max]
        (->>
         (s/split s #"\-")
         (map parse-int))]
    (->section min max)))

(defn parse-section-assignment-pair
  [line]
  (->>
   (s/split line #",")
   (map #(parse-section-assignment %))))

(defn parse-section-assignments
  [file-name]
  (->>
   file-name
   (load-section-assignments)
   (map #(parse-section-assignment-pair %))))

(defn fully-contains
  [subject, object]
  (and
   (<= (:min subject) (:min object))
   (>= (:max subject) (:max object))))

(defn one-fully-contains
  [pair]
  (or
   (fully-contains (first pair) (last pair))
   (fully-contains (last pair) (first pair))))

(defn pair-overlaps
  [pair]
  (not (or
   (< (:max (first pair)) (:min (second pair)))
   (< (:max (second pair)) (:min (first pair))))))

(defn part1
  [file-name]
  (->>
   file-name
   (parse-section-assignments)
   (filter one-fully-contains)
   (count)))

(defn part2
  [file-name]
  (->>
   file-name
   (parse-section-assignments)
   (filter pair-overlaps)
   (count)))
