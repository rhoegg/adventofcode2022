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

(defn cheat-play-score
  [opponent result]
  (let [opponent-index (opponent-play opponent)
        cheat-play-index (condp = result
                           \Y opponent-index ; draw
                           \X (mod (- opponent-index 1) 3) ; lose
                           \Z (mod (+ opponent-index 1) 3))] ; win

       (+ cheat-play-index 1))) ; the play score is one more than the index

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

(defn score-part2
  [turn]
  (let [result (get (last turn) 0)
        opponent (get (first turn) 0)
        game-score (condp = result
         \X 0 ; lose
         \Y 3 ; draw
         \Z 6 ; win
         :else -1
         )]
     (+ (cheat-play-score opponent result) game-score)
    ))

(defn parse-input
  [file-name]
  (->>
   (s/split (slurp file-name) #"\n")
   (map #(s/split % #" ")))
  )

(defn part1
  [file-name]
  (->>
   file-name
   (parse-input)
   (map score-part1)
   (reduce +)))

(defn part2
  [file-name]
  (->>
   file-name
   (parse-input)
   (map score-part2)
   (reduce +)))
