%dw 2.0
import * from dw::core::Arrays
import * from dw::core::Strings
output application/json

fun parsePoint(text) = do {
    var tokens = text splitBy(", ")
    ---
    {
        x: substringAfter(tokens[0], "x=") as Number,
        y: substringAfter(tokens[1], "y=") as Number
    }
}

fun parseSensorData(data) = do {
    var lines = data splitBy "\n"
    ---
    lines map ($ splitBy ":") map {
        sensor: parsePoint(substringAfter($[0], "Sensor at ")),
        beacon: parsePoint(substringAfter($[1], " closest beacon is at "))
    }
}

fun manhattanDistance(p1, p2) = 
    abs(p1.x - p2.x) + abs(p1.y - p2.y)

fun computeRange(sensors) = sensors map
    ($ ++ {
        range: manhattanDistance($.sensor, $.beacon)
    })

fun visible(location, sensors) =
    sensors some (manhattanDistance($.sensor, location) <= $.range)

fun inRange(sensors, row) = sensors map 
    ($ ++ {
        rowInRange: -(abs($.sensor.y - row) - $.range)
    }) filter ($.rowInRange >= 0)

// TODO: make this a check isVisible (between check) and then reduce
fun visibleRowSegments(sensors, row) = inRange(sensors, row)
    map [$.sensor.x - $.rowInRange, $.sensor.x + $.rowInRange]
    orderBy $[0]
    reduce (segment, priors) -> do {
        var segments = if (priors[0] is Number) [priors]
                        else priors
        var last = segments[-1]
        ---
        if (segment[0] <= segments[-1][1])
            (segments take sizeOf(segments) - 1)
                << [last[0], max([last[1], segment[1]])]
        else
            segments << segment
    }
        

var sensors = computeRange(parseSensorData(input1))
var part1Segments = visibleRowSegments(sensors, 2000000)
---
{
    part1: part1Segments,
    part2: (0 to 4000000)  map {row: $, segments: visibleRowSegments(sensors, $) } filter ($.segments some ($[0] > 0))
}
