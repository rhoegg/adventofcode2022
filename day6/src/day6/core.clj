(ns day6.core
  (:gen-class))


(defn unique-chars-in-header
  [signal]
    (->
     signal
     (subs 0 4)
     (set)
     (count)))

(defn find-start-of-packet
  ([signal] (find-start-of-packet signal 0))

  ([signal noiseCount]
   (if (= 4 (unique-chars-in-header signal))
     (+ noiseCount 4) ; found it
     (find-start-of-packet (subs signal 1) (inc noiseCount)))))

(defn part1
  [file-name]
  (find-start-of-packet (slurp file-name)))
