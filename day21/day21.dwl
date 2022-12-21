%dw 2.0
import * from dw::core::Objects
output application/json
fun parse(text) = {(
    text splitBy "\n" map ($ splitBy ": " then {($[0]): $[1]})
)}

fun eval(monkeys, name) = monkeys[name] match {
    case monkey matches /^\d+$/ -> monkey[0] as Number
    case monkey matches /([a-z]{4})\s(.)\s([a-z]{4})/ -> 
        monkey[2] match {
            case "+" -> eval(monkeys, monkey[1]) + eval(monkeys, monkey[3])
            case "-" -> eval(monkeys, monkey[1]) - eval(monkeys, monkey[3])
            case "*" -> eval(monkeys, monkey[1]) * eval(monkeys, monkey[3])
            case "/" -> eval(monkeys, monkey[1]) / eval(monkeys, monkey[3])
        }
}

fun part2(monkeys, name = "root") = name match {
        case "root" -> do {
            var parsed = monkeys["root"] splitBy " "
            ---
            part2(monkeys, parsed[0]) - part2(monkeys, parsed[2])
        }
        else -> eval(monkeys mergeWith {humn: 3952673930912}, name)
        // zero in by changing values of humn
    }
---
{
    part1: eval(parse(input1), "root"),
    part2: part2(parse(input1))
}

