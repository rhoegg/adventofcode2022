(ns day6.core
  (:gen-class))


(defn unique-chars-in-header
  ([signal header-size]
    (->
     signal
     (subs 0 header-size)
     (set)
     (count))))

(defn find-unique-window
    ([size signal] (find-unique-window size signal 0))

    ([size signal noiseCount]
     (if (= size (unique-chars-in-header signal size))
       (+ noiseCount size) ; found it
       (find-unique-window size (subs signal 1) (inc noiseCount)))))

(defn find-start-of-packet [signal] (find-unique-window 4 signal))
(defn find-start-of-message [signal] (find-unique-window 14 signal))

(defn part1
  [file-name]
  (find-start-of-packet (slurp file-name)))

(defn part2
  [file-name]
  (find-start-of-message (slurp file-name)))
