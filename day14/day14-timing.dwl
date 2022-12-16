%dw 2.0
import * from dw::core::Arrays
import * from dw::core::Objects
import * from dw::util::Timer

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
            (p.x): ((cave[p.x as String] default []) << p.y) 
                distinctBy $ 
                orderBy $
        }

fun down(p) = {x: p.x, y: p.y + 1}
fun downLeft(p) = {x: p.x - 1, y: p.y + 1}
fun downRight(p) = {x: p.x + 1, y: p.y + 1}

fun nextSandStep(sandUnit, cave) = do {
    var column = cave[sandUnit.x as String] default [] 
        dropWhile ($ < sandUnit.y + 1)
    @Lazy()
    var leftCol = cave[(sandUnit.x - 1) as String] default []
        dropWhile ($ < sandUnit.y + 1)
    @Lazy()
    var rightCol = cave[(sandUnit.x + 1) as String] default []
        dropWhile ($ < sandUnit.y + 1)
    var yPos = sandUnit.y + 1
    ---
    if (column[0] == yPos)
        if (leftCol[0] == yPos)
            if (rightCol[0] == yPos)
                sandUnit
            else downRight(sandUnit)
        else downLeft(sandUnit)
    else down(sandUnit)
}

@TailRec()
fun dropSand(state) = do {
    var floor = state.floor default false
    var sandUnitStart = state.sandUnit default {x: 500, y: 0}
    var computeSandUnitEnd = duration( () ->
        if (floor and sandUnitStart.y == maxDepth - 1)
            sandUnitStart
        else
            nextSandStep(sandUnitStart, state.cave))
    var sandUnitEnd = computeSandUnitEnd.result
    ---
    if (floor and (state.cave["500"][0] == 0))
        state ++ {done: true}
    else if (sandUnitEnd.y > maxDepth and (!(state.floor default false)))
        state ++ {done: true}
    else if (sandUnitStart == sandUnitEnd) // rest
        {
            sand: (state.sand default []) << sandUnitEnd,
            cave: addToCave(state.cave, sandUnitEnd),
            floor: state.floor,
            steps: state.steps,
            duration: log("milliseconds", state.duration)
        }
    else // falling
        dropSand({
            sandUnit: sandUnitEnd,
            sand: state.sand,
            cave: state.cave,
            floor: state.floor,
            steps: (state.steps default 0) + 1,
            duration: (state.duration default 0) + (computeSandUnitEnd.time)
        })
}

@TailRec()
fun pourSand(state) = do {
    var dropOneUnit = dropSand(state)
    var steps = log("steps", dropOneUnit.steps)
    var sandSize = log("sand", sizeOf(state.sand))
    ---
    if (state.done default false)
        state
    else
        pourSand((dropOneUnit - "steps") - "duration")
}

var lines = (input1 splitBy "\n") distinctBy $
var paths = lines map (
        ($ splitBy " -> ")
            map makePoint($)
)
var rocks = flatten(paths map fill($))
var maxDepth = max(flatten(rocks map $.y)) + 2
---
{
    // part1: sizeOf(pourSand({cave: cave(rocks)}).sand),
    part2: sizeOf(pourSand({cave: cave(rocks), floor: true}).sand),
    maxDepth: maxDepth
}

