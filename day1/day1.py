def total_calories(elfData):
    return sum(map(int, elfData.split("\n")))

with open("input.txt") as f:
    input = f.read()
    elves = input.split("\n\n")
    elf_calories = list(map(total_calories, elves))
    print("part 1: " + str(max(elf_calories)))
    elf_calories.sort(reverse=True)
    top3 = elf_calories[0:3]
    print("part 2: " + str(sum(top3)))