(ns day7.core-test
  (:require [clojure.test :refer :all]
            [day7.core :refer :all]))

(deftest build-tree-test
  (testing "bootstrap"
    (let [empty-example (apply-command "" "cd /")]
      (is (= new-tree empty-example))))
  (testing "ls"
    (let [first-example (apply-command new-tree "ls\ndir a\n14848514 b.txt\n8504156 c.dat\ndir d")]
      (testing "files"
        (let [contents (get (:directories first-example) [])]
          (is (= 2 (count contents)))
          (is (= 14848514 (first contents)))
          (is (= 8504156 (second contents)))))
      (testing "subdirectories"
        (is (= 2 (count (:subdirs first-example))))
        (is (= "a" (first (:subdirs first-example))))
        (is (= "d" (second (:subdirs first-example)))))))
  (testing "cd"
    (let [subdir-example (apply-command new-tree "cd a")]
      (testing "parent"
        (is (= new-tree (:parent subdir-example))))
      (testing "path"
        (is (= ["a"] (:path subdir-example))))
      (let [twodeep-example (apply-command subdir-example "cd e")]
        (testing "parent"
          (is (= subdir-example (:parent twodeep-example))))
        (testing "path"
          (is (= ["a" "e"] (:path twodeep-example))))
        (let [parent-example (apply-command twodeep-example "cd ..")]
          (testing "parent"
            (is (= new-tree (:parent parent-example))))
          (testing "path"
            (is (= ["a"] (:path parent-example)))))
        (let [twodeep-after-ls (apply-command twodeep-example "ls\n123 x.txt\n456 y.txt\ndir z")
              parent-after-ls-example (apply-command twodeep-after-ls "cd ..")]
          (testing "subdir contents"
            (is (= 2 (count (get (:directories parent-after-ls-example) ["a" "e"]))))))))))

(deftest total-size-test
  (let [example-tree (build-dir-tree "example.txt")]
    (testing "e"
      (is (= 584 (total-dir-size example-tree ["a" "e"])))
      (is (= 94853 (total-dir-size example-tree ["a"])))
      (is (= 24933642 (total-dir-size example-tree ["d"])))
      (is (= 48381165 (total-dir-size example-tree []))))))
