(ns day2.core)
(require '[clojure.string :as s])

(defn opponent-play
  [c]
  (- (int c) (int \A)))

(defn me-play
  [c]
  (- (int c) (int \X)))

(defn score-game
  [opponent me]
  (cond
    (= me opponent) 3 ; draw
    (= (mod (- me 1) 3) opponent) 6 ; win
    :else 0 ; loss
    ))

(defn score-part1
  [turn]
  (let [mine (get (last turn) 0)
        opponent (get (first turn) 0)
        shape-score (condp = mine
         \X 1 ; rock
         \Y 2 ; paper
         \Z 3 ; scissors
         :else 0
         )]
    (+ (score-game (opponent-play opponent) (me-play mine)) shape-score)
    ))

(defn part1
  [file-name]
  (->>
   (s/split (slurp file-name) #"\n")
   (map #(s/split % #" "))
   (map score-part1)
   (reduce +)))
