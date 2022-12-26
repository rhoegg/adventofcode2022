%dw 2.0
import * from dw::core::Strings
output application/json

fun snafuParse(digits) = do {
    var magnitude = sizeOf(digits) - 1
    var next = digits[0] match {
        case "1" -> 1
        case "2" -> 2
        case "0" -> 0
        case "-" -> -1
        case "=" -> -2
    }
    ---
    if (sizeOf(digits) == 1) next
    else next * (5 pow magnitude) + snafuParse(digits last magnitude)
}

fun snafuPrint(num) = do {
    var b5 = num match {
        case is Number -> base5(num)
        case is String -> num
    }
    var this = b5[sizeOf(b5) - 1]
    var thisDigit = this match {
        case "5" -> "0"
        case "4" -> "-"
        case "3" -> "="
        else -> this
    }
    var carry = if (this as Number >= 3) 1 else 0
    var nextDigit = 
        if (sizeOf(b5) > 1)
            (b5[sizeOf(b5) - 2] as Number + carry) as String
        else if (carry == 1) "1"
        else ""
    var rest = b5 first (sizeOf(b5) - 2) ++ nextDigit
    ---
    {
        b5: b5,
        this: this,
        nextDigit: nextDigit,
        carry: carry,
        rest: rest,
        almost: if (isEmpty(rest)) thisDigit 
            else snafuPrint(rest) ++ thisDigit
    }.almost
}

fun base5(num) = 
    (if (num > 4) base5(floor(num / 5)) else "")
    ++ (num mod 5) as String
---
{
    example: example2 splitBy("\n") map snafuParse($) map snafuPrint($),
    part1: snafuPrint(sum(input1 splitBy "\n" map snafuParse($)))
}

