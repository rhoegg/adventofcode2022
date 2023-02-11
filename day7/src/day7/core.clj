(ns day7.core
  (:gen-class))
(require '[clojure.string :as s])
(defn parse-int [s] (Integer/parseInt s))

(defrecord directory-tree [parent path subdirs directories])
(def new-tree (->directory-tree nil [] [] {}))

(defn apply-command
  [working-tree command]
  (cond
    (= "" working-tree) new-tree
    (re-find #"^ls" command) (apply-command-ls working-tree command)
    (re-find #"^cd" command) (apply-command-cd working-tree command)
    :else working-tree))

(defn apply-command-ls
  [working-tree command]
  (let [results (drop 1 (s/split command #"\n"))
        contents (parse-ls-contents results)
        subdirs (parse-ls-subdirs results)]
    (assoc working-tree :subdirs subdirs :directories (assoc (:directories working-tree) (:path working-tree) contents))))

(defn apply-command-cd
  [working-tree command]
  (let [targetdir (second (s/split command #" "))]
    (if (= ".." targetdir)
      (assoc
       (:parent working-tree)
       :directories (:directories working-tree)) ; doesn't handle root dir
      (->directory-tree
       working-tree
       (conj (:path working-tree) targetdir)
       []
       (:directories working-tree)))))

(defn parse-ls-contents
  [results]
  (->>
   results
   (filter #(re-find #"^\d+" %))
   (map #(s/split % #" "))
   (map first)
   (map parse-int)))

(defn parse-ls-subdirs
  [results]
  (->>
   results
   (filter #(re-find #"^dir" %))
   (map #(s/split % #" "))
   (map second)))

(defn build-dir-tree
  [file-name]
  (->>
   file-name
   (slurp)
   (#(s/split % #"\n?\$ "))
   (reduce apply-command)
   ))

(defn sub-tree-dirs
  [tree dir]
  (let [dirs (:directories tree)]
    (->> dirs
       (keys)
       (filter #(= dir (take (count dir) %)))
       (map #(-> [% (get dirs %)]))
       (into {}))))

(defn total-dir-size
  ([tree dir]
   (->>
    (sub-tree-dirs tree dir)
    (vals)
    (flatten)
    (apply +))))

(defn part1
  [file-name]
  (let [t (build-dir-tree file-name)]
    (->> (:directories t)
         (keys)
         (map #(total-dir-size t %))
         (filter #(<= % 100000))
         (apply +))))

(defn part2
  [file-name]
  (let [t (build-dir-tree file-name)
        used-space (total-dir-size t [])
        required-space (- used-space 40000000)]
    (->> (:directories t)
         (keys)
         (map #(total-dir-size t %))
         (filter #(>= % required-space))
         (apply min))))
