%dw 2.0
import * from dw::core::Arrays
// import * from dw::core::Strings

output application/json
var pairInputs = input1 splitBy "\n\n" 
    map (p) -> p splitBy "\n"
var packetInputs = input1 splitBy "\n" filter ! isEmpty($)
var dividerPackets = [
    [[2]],
    [[6]]
]

fun parsePackets(pInput) = pInput map read($, "application/json")

fun cmp(l, r) =
    l match {
        case is Number -> r match {
            case is Number -> if (l != r) (if (l < r) -1 else 1)
                                else 0
            else -> cmp([l], r)
        }
        else -> r match {
            case is Number -> cmp(l, [r])
            else -> cmpArrays(l, r)
        }
    }

fun cmpArrays(l, r) =
    if (isEmpty(l))
        if (isEmpty(r)) 0
        else -1
    else if (isEmpty(r)) 1
        else cmp(l[0], r[0]) match {
                case -1 -> -1
                case 1 -> 1
                case 0 -> cmp(l drop 1, r drop 1)
            }

fun merge(a, b) =
    if (isEmpty(a)) b
    else if (isEmpty(b)) a
    else if (cmp(a[0], b[0]) == -1) [a[0]] ++ merge(a drop 1, b)
    else [b[0]] ++ merge(a, b drop 1)

fun mergesort(packets) =
    if (sizeOf(packets) < 2) packets
    else do {
        var halves = packets splitAt (sizeOf(packets) / 2)
        ---
        merge(mergesort(halves[0]), mergesort(halves[1]))
    }

@TailRec()
fun bubblesort(packets, sorted = []) = do {
    var nextSmallest = packets reduce (packet, smallest) -> 
        if (cmp(packet, smallest) == -1) packet else smallest
    ---
    bubblesort(packets -- nextSmallest, sorted << nextSmallest)
}
    

var pairs = pairInputs map (pairInput, i) -> do {
    var pair = parsePackets(pairInput)
    ---
    { 
        pair: pair,
        cmp: cmp(pair[0], pair[1]),
        index: i + 1
    }
}
fun countLessThan(packets, element) =
    packets countBy (packet) -> (cmp(packet, element) == -1)

// full sorting ended up blowing the stack (mergesort) or way too slow (bubblesort)
// var sortedPackets = bubblesort(parsePackets(packetInputs) ++ dividerPackets)
var p2packets = parsePackets(packetInputs)

---
{
    part1: sum(pairs filter $.cmp == -1 map $.index),
    part2: do {
        var divider1 = countLessThan(p2packets, [[2]]) + 1
        var divider2 = countLessThan(p2packets, [[6]]) + 2 // counting [[2]]
        ---
        {
            dividers: [divider1, divider2],
            product: divider1 * divider2
        }
    }
}
