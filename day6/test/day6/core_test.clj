(ns day6.core-test
  (:require [clojure.test :refer :all]
            [day6.core :refer :all]))

(deftest find-start-of-packet-test
  (testing "Finding start of packet"
    (is (= 7 (find-start-of-packet "mjqjpqmgbljsphdztnvjfqwrcgsmlb")))
    (is (= 5 (find-start-of-packet "bvwbjplbgvbhsrlpgdmjqwftvncz")))
    (is (= 6 (find-start-of-packet "nppdvjthqldpwncqszvftbrmjlhg")))
    (is (= 10 (find-start-of-packet "nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg")))
    (is (= 11 (find-start-of-packet "zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw")))))
