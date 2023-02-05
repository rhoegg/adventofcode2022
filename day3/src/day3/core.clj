(ns day3.core
  (:gen-class))

(require '[clojure.string :as s])

(defn load-rucksacks
  [file-name]
  (->
   file-name
   (slurp)
   (s/split #"\n")))

(defn compartments
  [rucksack]
  (->
   rucksack
   (count)
   (/ 2)
   (split-at rucksack)))

(defn find-common-char
 ([single-array] (apply find-common-char single-array))
 ([s1 s2]
  (let [key-item (first s2)
        remaining-s1 (drop-while #(< (int %) (int key-item)) s1)]
    (cond
      (= (first remaining-s1) key-item) key-item
      (empty? remaining-s1) (throw (Exception. "unexpected end of compartment"))
      :else (find-common-char s2 remaining-s1))))
 ([s1 s2 s3]
  (let [key-item (first s1)]
    (if (and (s/includes? s2 (str key-item)) (s/includes? s3 (str key-item))) key-item
      (find-common-char (rest s1) s2 s3)))))


(defn find-packing-error
  [rucksack]
  (let [[l r] (compartments rucksack)]
    (find-common-char (sort l) (sort r))))

(defn item-priority
  [item]
  (cond
    (re-matches #"[a-z]" (str item)) (+ (- (int item) (int \a)) 1) ; Lowercase item types a through z have priorities 1 through 26.
    :else (+ (- (int item) (int \A)) 27)                     ; Uppercase item types A through Z have priorities 27 through 52.
    ))

(defn part1
  [file-name]
  (->>
   file-name
   (load-rucksacks)
   (map find-packing-error)
   (map item-priority)
   (reduce +)))

(defn part2
  [file-name]
  (->>
   file-name
   (load-rucksacks)
   (partition 3)
   (map find-common-char)
   (map item-priority)
   (reduce +)))
