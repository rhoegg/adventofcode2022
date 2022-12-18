%dw 2.0
import * from dw::core::Arrays
output application/json

fun adjacents(p) = flatten(
    (0 to 2) map (dim) -> (
        [p[dim] - 1, p[dim] + 1] map (dx) -> (
            p reduce (x, point = []) -> point <<
                if (sizeOf(point) == dim) dx else x
        )
    )
)

var lines = input1 splitBy "\n"
var cubes = lines map ($ splitBy "," map ($ as Number))

var cubeData = cubes map {
    this: $,
    adjacents: adjacents($),
    emptyAdjacents: adjacents($) countBy (! (cubes contains $))
}
---
{
    part1: cubeData map $.emptyAdjacents then sum($)
}

