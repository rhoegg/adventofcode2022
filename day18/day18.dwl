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

fun expandBounds(points, distance) =
    (0 to 2) map (dim) -> (
        [min(cubes map $[dim]) - distance, max(cubes map $[dim]) + distance]
    )

fun inBounds(bounds, point) =
    [0, 1, 2] as Array every (dim) ->
        point[dim] > bounds[dim][0] and point[dim] < bounds[dim][1] 

fun floodFillOutside(bounds, lava, point, filled = []) =
    if (! inBounds(bounds, point)) []
    else if ((lava contains point) or (filled contains point)) []
        else flatten(adjacents(point) 
            map floodFillOutside(bounds, lava, $, filled << point))

var lines = input1 splitBy "\n"
var cubes = lines map ($ splitBy "," map ($ as Number))

var cubeData = cubes map {
    this: $,
    adjacents: adjacents($),
    emptyAdjacents: adjacents($) countBy (! (cubes contains $))
}
var floodFillBounds = expandBounds(cubes, 2)

---
{
    part1: cubeData map $.emptyAdjacents then sum($),
    part2: floodFillOutside(floodFillBounds, cubes, [0,0,0]),
    x: inBounds(floodFillBounds, [-1,1,1])
}

