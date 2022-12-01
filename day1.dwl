%dw 2.0
output application/json

fun calories(elfReports) = elfReports map ($ as Number)
var elfCalories = (input1 splitBy '\n\n') 
    map sum(calories($ splitBy '\n'))

var part1 = max(elfCalories)
var part2 = sum((elfCalories orderBy (-1 * $))[0 to 2])
---
{
    part1: part1,
    part2: part2
}
