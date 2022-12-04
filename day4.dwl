%dw 2.0
import * from dw::core::Arrays
output application/json

fun fullyContains(a1, a2) =
    (a2[0] >= a1[0]) and (a2[1] <= a1[1])

fun shouldReconsider(pair) =
    (pair[0] fullyContains pair[1])
    or (pair[1] fullyContains pair[0])

fun overlaps(pair) = do {
    var a1 = (pair[0][0] to pair[0][1])
    var a2 = (pair[1][0] to pair[1][1])
    ---
    a1 some (a2 contains $)
}
    

var sectionAssignmentPairs = input1 splitBy "\n" 
    map ($ splitBy ","
        map ($ splitBy "-"
            map ($ as Number)))
---
{
    part1: sectionAssignmentPairs countBy shouldReconsider($),
    part2: sectionAssignmentPairs countBy overlaps($)
}