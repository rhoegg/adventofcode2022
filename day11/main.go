package main

import (
	"fmt"
	"math"
	"math/big"
	"os"
	"sort"
	"strconv"
	"strings"
)

var bigZero = big.NewInt(0)

type Operation func(old *big.Int) *big.Int
type Test func(worryLevel *big.Int) bool

type Monkey struct {
	Items []*big.Int
	Inspections int
	Operation Operation
	TestFactor int
	ThrowTest Test
	YesMonkey, NoMonkey int
}

func parseMonkeys(filename string) []*Monkey {
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	var monkeys []*Monkey
	for _, monkeyData := range strings.Split(string(data), "\n\n") {
		lines := strings.Split(monkeyData, "\n")
		monkey := &Monkey{}

		itemsTokens := strings.SplitN(strings.TrimLeft(lines[1], " "), " ", 3)
		items := strings.Split(itemsTokens[2], ", ")
		for _, token := range items {
			worryLevel, _ := strconv.ParseInt(token, 10, 64)
			monkey.Items = append(monkey.Items, big.NewInt(worryLevel))
		}

		monkey.Operation = parseOperation(lines[2])
		monkey.ThrowTest, monkey.TestFactor = parseTest(lines[3])
		monkey.YesMonkey = parseTrueMonkey(lines[4])
		monkey.NoMonkey = parseFalseMonkey(lines[5])

		monkeys = append(monkeys, monkey)
	}

	return monkeys
}

func parseOperation(line string) Operation {
	relevantParts := strings.Split(strings.TrimPrefix(line, "  Operation: new = old "), " ")
	if relevantParts[1] == "old" {
		if relevantParts[0] == "*" {
			return func(old *big.Int) *big.Int {
				//fmt.Printf("squaring %s", old.String())
				old.Mul(old, old)
				//fmt.Printf(" = %s\n", old.String())
				return old
			}
		} else {
			return func(old *big.Int) *big.Int {
				//fmt.Printf("doubling %s\n", old.String())
				return old.Add(old, old)
			}
		}
	} else {
		constant, _ := strconv.ParseInt(relevantParts[1], 10, 64)
		if relevantParts[0] == "*" {
			return func(old *big.Int) *big.Int {
				//fmt.Printf("%s * %d", old.String(), constant)
				old.Mul(old, big.NewInt(constant))
				//fmt.Printf(" = %s\n", old.String())
				return old
			}
		} else {
			return func(old *big.Int) *big.Int {
				//fmt.Printf("%s + %d\n", old.String(), constant)
				return old.Add(old, big.NewInt(constant))
			}
		}
	}
}

func parseTest(line string) (Test, int) {
	parameter := strings.TrimPrefix(line,"  Test: divisible by ")
	dividend, _ := strconv.ParseInt(parameter, 10, 64)
	return func(worryLevel *big.Int) bool {
		check := new(big.Int).Set(worryLevel)
		return check.Mod(worryLevel, big.NewInt(dividend)).Cmp(bigZero) == 0
	}, int(dividend)
}

func parseTrueMonkey(line string) int {
	token := strings.TrimPrefix(line, "    If true: throw to monkey ")
	num, _ := strconv.Atoi(token)
	return num
}

func parseFalseMonkey(line string) int {
	token := strings.TrimPrefix(line, "    If false: throw to monkey ")
	num, _ := strconv.Atoi(token)
	return num
}

func PrintItems(monkeys []*Monkey) {
	for i, monkey := range monkeys {
		fmt.Printf(" * %d: %v\n", i, monkey.Items)
	}
	fmt.Println()
}

func TakeTurn(monkeys []*Monkey, i int, relief bool, round int) []*Monkey {
	lcm := int64(1)
	for _, monkey := range monkeys {
		lcm *= int64(monkey.TestFactor)
	}
	for _, itemWorryLevel := range monkeys[i].Items {
		monkeys[i].Inspections += 1
		//fmt.Printf("inspect item with worry level is %d\n", itemWorryLevel)
		// operation
		old := itemWorryLevel.Int64()
		itemWorryLevel := monkeys[i].Operation(itemWorryLevel)
		//fmt.Printf("operation result %s\n", itemWorryLevel.String())
		//fmt.Printf("Worry level is %d\n", itemWorryLevel)
		// get bored
		if relief {
			itemWorryLevel = itemWorryLevel.Div(itemWorryLevel, big.NewInt(3))
			//fmt.Printf("Bored! Worry level is %d\n", itemWorryLevel)
		} else {
			//if itemWorryLevel > lcm && (itemWorryLevel % lcm == 0) {
			//	itemWorryLevel = itemWorryLevel / lcm
			//}
			if big.NewInt(old).Cmp(itemWorryLevel) > 0 {
				fmt.Printf("Max int64 %d\n", math.MaxInt64)
				panic(fmt.Sprintf("monkey %d overflow round %d: %d -> %d\n", i, round, old, itemWorryLevel.Int64()))
			}
			bigLcm := big.NewInt(lcm)
			if itemWorryLevel.Cmp(bigLcm) > 0 {
				itemWorryLevel.Mod(itemWorryLevel, bigLcm)
				//fmt.Printf("fixed worry level %s\n", itemWorryLevel.String())
			}
		}
		// test worry level
		destMonkey := -1
		if monkeys[i].ThrowTest(itemWorryLevel) {
			destMonkey = monkeys[i].YesMonkey
		} else {
			destMonkey = monkeys[i].NoMonkey
		}
		// throw
		//fmt.Printf("Item with worry level %d is thrown to monkey %d\n", itemWorryLevel, destMonkey)
		monkeys[destMonkey].Items = append(monkeys[destMonkey].Items, itemWorryLevel)
		//fmt.Printf("sent %s to %d: %v\n", itemWorryLevel.String(), destMonkey, monkeys[destMonkey].Items)
	}
	monkeys[i].Items = nil
	return monkeys
}

func PlayRound(monkeys []*Monkey, relief bool, round int) []*Monkey {
	for i := 0; i < len(monkeys); i++ {
		monkeys = TakeTurn(monkeys, i, relief, round)
	}
	return monkeys
}

func Part1() {
	monkeys := parseMonkeys("input.txt")
	for i := 0; i < 20; i++ {
		PlayRound(monkeys, true, i)
	}
	var inspections []int
	for _, monkey := range monkeys {
		inspections = append(inspections, monkey.Inspections)
	}
	sort.Ints(inspections)
	most, secondMost := inspections[len(inspections) - 1], inspections[len(inspections) - 2]
	fmt.Printf("inspections %v\n", inspections)
	fmt.Printf("monkey business: %d\n", most * secondMost)
}

func Part2() {
	monkeys := parseMonkeys("input.txt")
	printRounds := []int {1, 20, 1000, 2000, 3000, 4000, 5000, 6000, 7000, 8000, 9000, 10000}
	for i := 0; i < 10000; i++ {
		fmt.Printf("Round %d...\n", i)
		PlayRound(monkeys, false, i)
		for _, r := range printRounds {
			if i == (r - 1) {
				var inspections []int
				for _, monkey := range monkeys {
					inspections = append(inspections, monkey.Inspections)
				}
				fmt.Printf("Round %d\n%v\n\n", r, inspections)
				PrintItems(monkeys)
			}
		}
	}
	var inspections []int
	for _, monkey := range monkeys {
		inspections = append(inspections, monkey.Inspections)
	}
	sort.Ints(inspections)
	most, secondMost := inspections[len(inspections) - 1], inspections[len(inspections) - 2]
	fmt.Printf("inspections %v\n", inspections)
	fmt.Printf("monkey business: %d\n", most * secondMost)
}

func main() {
	Part2()
}
