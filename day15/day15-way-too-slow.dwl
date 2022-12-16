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

var sensors = computeRange(parseSensorData(input1))
var xRange = {
    max: sensors map [$.sensor.x + $.range, $.beacon.x] 
        then flatten($)
        then max($),
    min: sensors map [$.sensor.x - $.range, $.beacon.x] 
    then flatten($)
    then min($)
}
---
(xRange.min to xRange.max) map { x: $, y: 2000000 } 
    map {
        location: $,
        visible: visible($, sensors),
        beacon: sensors map $.beacon contains $
    }
    countBy ($.visible and ! $.beacon)
