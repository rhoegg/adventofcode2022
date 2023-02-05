%dw 2.0
output application/json
fun parse(text) = {(
    text splitBy "\n" map ($ splitBy ": " then {($[0]): $[1]})
)}

fun eval(monkeys, name) = monkeys[name] match {
    case monkey matches /\d/ -> monkey[0] as Number
    case monkey matches /([a-z]{4})\s(.)\s([a-z]{4})/ -> 
        monkey[2] match {
            case "+" -> eval(monkeys, monkey[1]) + eval(monkeys, monkey[3])
            case "-" -> eval(monkeys, monkey[1]) - eval(monkeys, monkey[3])
            case "*" -> eval(monkeys, monkey[1]) * eval(monkeys, monkey[3])
            case "/" -> eval(monkeys, monkey[1]) / eval(monkeys, monkey[3])
            else -> 0
        }
    case monkey if(true) -> monkey
}
---
eval(parse(input1), "root")
