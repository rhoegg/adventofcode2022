%dw 2.0
import * from dw::core::Arrays
import * from dw::core::Strings
output application/json

fun parseBlueprint(line) = do {
    var parts = line splitBy ". "
    ---
    {
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

fun runFactory(blueprint, inventory, minutes) = do {
    var options = choices(blueprint, inventory) map (transaction) -> do {
        var afterSpend = spendResources(blueprint, inventory, transaction)
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
    ---
    if (minutes == 1) [options[0]] // doesn't matter
    else flatten(options map (afterCollection) -> runFactory(blueprint, afterCollection, minutes - 1))
}

fun spendResources(blueprint, inventory, transaction) = do {
    inventory mapObject ((balance, resource) -> { (resource): balance + (transaction[resource] default 0)})
}

fun choices(blueprint, inventory) =
    if (afford(blueprint.geodeRobot, inventory))
        [spend(blueprint.geodeRobot) ++ {geodeRobot: 1}] 
    else
        [{}] // wait
        ++ if (afford(blueprint.obsidianRobot, inventory)) [spend(blueprint.obsidianRobot) ++ {obsidianRobot: 1}] else []
        ++ if (afford(blueprint.clayRobot, inventory)) [spend(blueprint.clayRobot) ++ {clayRobot: 1}] else []
        ++ if (afford(blueprint.oreRobot, inventory)) [spend(blueprint.oreRobot) ++ { oreRobot: 1 }] else []

fun spend(price) =
    price mapObject (price, resource) -> {(resource): -price}

fun afford(price, resources) =
    (price pluck (itemPrice, resource) -> resources[resource] - itemPrice)
    every (balance) -> balance >= 0


var lines = example splitBy "\n"
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
    blueprints: blueprints,
    run: runFactory(blueprints[1], startInventory, 24) maxBy $.geode
}

