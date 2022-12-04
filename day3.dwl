%dw 2.0
import * from dw::core::Arrays
import * from dw::core::Strings
output application/json

fun findSharedItem(compartments) =
    (compartments[0] splitBy "")
    firstWith (itemType) -> (compartments every ($ contains itemType))

fun itemTypePriority(itemType) = do {
    var lowerAPriority = 1
    var upperAPriority = 27
    ---
    if (itemType == lower(itemType))
        charCode(itemType) - charCode("a") + lowerAPriority
    else
        charCode(itemType) - charCode("A") + upperAPriority
}

fun fixBadges(group) = findSharedItem(group map $.contents)
    then itemTypePriority($)

var rucksacks = 
    (input1 splitBy "\n")
    map {
        contents: $,
        compartments: ($ substringEvery sizeOf($) / 2)
    }
    map ($ ++ {
        sharedItemType: findSharedItem($.compartments)
    })

var groups = rucksacks
---
{
    part1: sum(rucksacks map itemTypePriority($.sharedItemType)),
    part2: sum(groups divideBy 3 map fixBadges($))
}