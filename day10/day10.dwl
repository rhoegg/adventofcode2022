%dw 2.0
import * from dw::core::Arrays

output application/csv

fun cycles(command) =
    command[0] match {
        case "addx" -> 2
        else -> 1
    }

fun execute(commands, x) = do {
    var command = commands[0]
    var rest = commands drop 1
    var nextX = command[0] match {
        case "addx" -> x + (command[1])
        else -> x
    }
    ---
    ((1 to cycles(command)) map x) 
       ++ if(isEmpty(rest)) [] else execute(rest, nextX)  
}

fun keyFrame(cycle) = ((cycle - 20) mod 40) == 0

fun part1(commands) = 
    (commands map (x, i) -> {cycle: i + 1, x: x, signalStrength: (i + 1) * x})
    filter keyFrame($.cycle)
    then sum($.signalStrength)

fun draw(commands) = 
    (((commands map (x, i) -> (x - (i mod 40)) match {
        case -1 ->  "#"
        case 0 -> "#"
        case 1 -> "#"
        else -> "."
    }) divideBy 40)
    map ({screen: $ joinBy ""}))

var commands = (input1 splitBy '\n') 
        map ($ splitBy " ")

---
draw(execute(commands, 1))

