%dw 2.0
import * from dw::core::Arrays
import * from dw::core::Objects
import * from dw::core::Strings
output application/json

fun parseBlueprint(line) = do {
    var parts = line splitBy ". "
    ---
    {
        number: parts[0] substringAfter "Blueprint " substringBeforeLast ":",
        costs: {
            oreRobot: parts[0] substringAfterLast "costs " substringBefore " ore" then {ore: ($ as Number)},
            clayRobot: parts[1] substringAfterLast  "costs " substringBefore  " ore" then {ore: ($ as Number)},
            obsidianRobot: parts[2] substringAfter "costs " splitBy " and " map (
                $ splitBy " " then $[0] as Number
            ) then {ore: $[0], clay: $[1]},
            geodeRobot: parts[3] substringAfter "costs " splitBy " and " map (
                $ splitBy " " then $[0] as Number
            ) then {ore: $[0], obsidian: $[1]}
        }
    }
}

fun inventoryKey(inventory) = inventory pluck ($ as String) joinBy " "

fun runFactory(state, blueprint, inventory, minutes) = do {
    var cacheKey = (blueprint.number ++ " " ++ inventoryKey(inventory))
    var options = state[cacheKey] default choices(blueprint, inventory) map (transaction) -> do {
        var afterSpend = spendResources(inventory, transaction)
        ---
        {
            ore: afterSpend.ore + inventory.oreRobot,
            clay: afterSpend.clay + inventory.clayRobot,
            obsidian: afterSpend.obsidian + inventory.obsidianRobot,
            geode: afterSpend.geode + inventory.geodeRobot,
            oreRobot: afterSpend.oreRobot,
            clayRobot: afterSpend.clayRobot,
            obsidianRobot: afterSpend.obsidianRobot,
            geodeRobot: afterSpend.geodeRobot
        }
    }
    var worthyOptions = if (options some ($.geodeRobot > 0))
                    options filter ($.geodeRobot == max(options map ($.geodeRobot)))
                else options filter ($.obsidianRobot == max(options map ($.obsidianRobot)))
    var nextState = state ++ {
       (cacheKey) : options
    }
    ---
    if (minutes == 1) { state: {}, options: [options[0]] } // doesn't matter
    else worthyOptions reduce (afterCollection, acc = {state: nextState, options: []}) -> do {
        var optionResult = runFactory(acc.state, blueprint, afterCollection, minutes - 1)
        var aggregateResult = (acc.options ++ optionResult.options) distinctBy $ 
        ---
        {
            state: optionResult.state,
            options: aggregateResult
        }
    }
}

fun spendResources(inventory, transaction) = do {
    inventory mapObject ((balance, resource) -> 
        { (resource): balance + (transaction[resource] default 0)})
}

fun choices(blueprint, inventory) =
    if (afford(blueprint.costs.geodeRobot, inventory))
        ([spend(blueprint.costs.geodeRobot) ++ {geodeRobot: 1}])
    else if (afford(blueprint.costs.obsidianRobot, inventory))
        ([spend(blueprint.costs.obsidianRobot) ++ {obsidianRobot: 1}])
    else if (afford(blueprint.costs.clayRobot, inventory) and (inventory.clayRobot == 0))
        ([spend(blueprint.costs.clayRobot) ++ {clayRobot: 1}])
    else
        [{}]
        ++ (
            if (afford(blueprint.costs.clayRobot, inventory) and (inventory.clayRobot < blueprint.costs.obsidianRobot.clay))
            [spend(blueprint.costs.clayRobot) ++ {clayRobot: 1}] else [])
        ++ (
            if (afford(blueprint.costs.oreRobot, inventory) and (inventory.oreRobot < max(blueprint.costs pluck $.ore))) [spend(blueprint.costs.oreRobot) ++ { oreRobot: 1 }] else [])

fun spend(price) =
    price mapObject (price, resource) -> {(resource): -price}

fun afford(price, resources) =
    (price pluck (itemPrice, resource) -> resources[resource] - itemPrice)
    every (balance) -> balance >= 0


var lines = input1 splitBy "\n"
var blueprints = lines map parseBlueprint($)
var startInventory = {
    ore: 0,
    clay: 0,
    obsidian: 0,
    geode: 0,
    oreRobot: 1,
    clayRobot: 0,
    obsidianRobot: 0,
    geodeRobot: 0
}
---
{
    // inventory: startInventory,
    // blueprints: blueprints,
    run: blueprints map (blueprint) -> log(max(runFactory({}, blueprint, startInventory, 24).options map $.geode))
}

