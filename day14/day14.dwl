%dw 2.0
import * from dw::core::Objects

output application/json

fun straightLine(p1, p2) =
    if (p1.x == p2.x)
        (p1.y to p2.y) map ({x: p1.x, y: $})
    else
        (p1.x to p2.x) map ({x: $, y: p1.y})

fun fill(path) = (
    path reduce (point, filled) -> 
    filled match {
        case is Object -> straightLine(filled, point)
        else -> filled ++ straightLine(filled[-1], point)
    }
) distinctBy $

fun makePoint(pointText) = do {
    var coords = pointText splitBy ","
    ---
    {
        x: coords[0] as Number,
        y: coords[1] as Number
    }
}

fun cave(points) =
    points reduce (p, c = {}) ->
        addToCave(c, p)

fun addToCave(cave, p) =
    cave mergeWith {
            (p.x): ((cave[p.x as String] default []) << p.y) distinctBy $
                        // could be an object for y as well if this is slow
        }

fun down(p) = {x: p.x, y: p.y + 1}
fun downLeft(p) = {x: p.x - 1, y: p.y + 1}
fun downRight(p) = {x: p.x + 1, y: p.y + 1}

fun nextSandStep(sandUnit, cave) = do {
    var column = cave[sandUnit.x as String] default []
    @Lazy()
    var leftCol = cave[(sandUnit.x - 1) as String] default []
    @Lazy()
    var rightCol = cave[(sandUnit.x + 1) as String] default []
    ---
    if (column contains sandUnit.y + 1)
        if (leftCol contains sandUnit.y + 1)
            if (rightCol contains sandUnit.y + 1)
                sandUnit
            else downRight(sandUnit)
        else downLeft(sandUnit)
    else down(sandUnit)
}


@TailRec()
fun dropSand(state) = do {
    var sandUnitStart = state.sandUnit default {x: 500, y: 0}
    var sandUnitEnd = nextSandStep(sandUnitStart, state.cave)
    ---
    if (sandUnitEnd.y > maxDepth)
        state ++ {overflow: true}
    else if (sandUnitStart == sandUnitEnd) // rest
        {
            overflow: false,
            sand: state.sand << sandUnitEnd,
            cave: addToCave(state.cave, sandUnitEnd)
        }
    else // falling
        dropSand({
            sandUnit: sandUnitEnd,
            sand: state.sand,
            cave: state.cave
        })
}

@TailRec()
fun pourSand(state) = do {
    var dropOneUnit = dropSand(state)
    var sandSize = log(sizeOf(state.sand))
    ---
    if (state.overflow default false)
        state
    else
        pourSand(dropOneUnit)
}

var lines = input1 splitBy "\n"
var paths = lines map (
        ($ splitBy " -> ")
            map makePoint($)
)
var rocks = flatten(paths map fill($))
var maxDepth = max(flatten(rocks map $.y))
---
{
    part1: sizeOf(pourSand({sand: [], cave: cave(rocks)}).sand)
}

