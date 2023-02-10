(ns day5.core
  (:gen-class))

(require '[clojure.string :as s])

(defn parse-int [s] (Integer/parseInt s))

(defrecord move [count source target])

(defn load-input
  [file-name]
  (->
   file-name
   (slurp)
   (s/split #"\n\n")))

(defn parse-crates-stacks
  [last-line]
   (->>
     last-line
     (map-indexed #(-> [%1 %2]))
     (filter #(not= (second %) \space))
     (into {})
     (clojure.set/map-invert)))

(defn parse-individual-stack
  [lines stacks stack]
                                        ;(map #(nth % column) lines)
  (let [[name column] (find stacks stack)]
    [name (->>
           (drop-last lines)
           (map #(if (>= column (count %))
                   \space
                   (nth % column)))
           (filter #(not= \space %)))]))

(defn parse-crates
  [s]
  (let [lines (s/split s #"\n")
        stacks (parse-crates-stacks (last lines))]
    (->>
     (keys stacks)
     (map (partial parse-individual-stack lines stacks))
     (into {}))))

(defn parse-moves
  [s]
  (let [lines (s/split s #"\n")]
    (->>
     lines
     (map #(s/split % #" "))
     (map (fn [line]
            (->move
             (parse-int (nth line 1))
             (first (nth line 3))
             (first (nth line 5))))))))

(defn move-step-one-at-a-time
  [crates move]
  (let [{count :count, source :source, target :target} move
        cargo (take count (get crates source))]
    (->
     crates
     (assoc target (concat (reverse cargo) (get crates target)))
     (assoc source (drop count (get crates source)))
     )))

(defn move-step-all-at-once
  [crates move]
  (let [{count :count, source :source, target :target} move
        cargo (take count (get crates source))]
    (->
     crates
     (assoc target (concat cargo (get crates target)))
     (assoc source (drop count (get crates source)))
     )))

(defn move-crates
  [file-name, mover]
  (let [[crates-input moves-input] (load-input file-name)
        crates (parse-crates crates-input)
        moves (parse-moves moves-input)]
    (->>
     moves
     (reduce mover crates)
     (vals)
     (map #(first %)))))

(defn part1
  [file-name]
  (move-crates file-name move-step-one-at-a-time))

(defn part2
  [file-name]
  (move-crates file-name move-step-all-at-once))
