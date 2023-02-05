%dw 2.0
output application/json
var rounds = input1 splitBy '\n'
var opponent = {
    A: "rock",
    B: "paper",
    C: "scissors"
}
var recommended = {
    X: "rock",
    Y: "paper",
    Z: "scissors"
}

var desiredOutcome = {
    X: "lose",
    Y: "draw",
    Z: "win"
}

fun shapeScore(play) = play match {
    case "rock" -> 1
    case "paper" -> 2
    case "scissors" -> 3
}

fun chooseShape(round) = do {
    var opponentPos = shapeScore(round.opponent) - 1
    var delta = round.outcome match {
        case "lose" -> -1
        case "draw" -> 0
        case "win" -> 1
    }
    ---
    ["rock", "paper", "scissors"][(opponentPos + delta) mod 3]
}


fun wins(me, them) = (shapeScore(me) - 1) == (shapeScore(them) mod 3)
---
rounds map (
    ($ splitBy " ") then {
        opponent: opponent[$[0]],
        outcome: desiredOutcome[$[1]]
    }
) map {
    shape: chooseShape($),
    outcome: $.outcome match {
        case "lose" -> 0
        case "draw" -> 3
        case "win" -> 6
    }
} map (shapeScore($.shape) + $.outcome)
then sum($)
