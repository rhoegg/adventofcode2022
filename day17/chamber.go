package main

type Chamber struct {
	Width int
	Landed map[Point]bool
	Rocks []Rock
	nextRock int
	JetPattern []int
	nextJet int
	peak int
}

func NewChamber(width int, rocks []Rock, jetPattern []int) *Chamber {
	return &Chamber{Width: width, Landed: make(map[Point]bool), Rocks: rocks, JetPattern: jetPattern}
}

func (c *Chamber) Peak() int {
	return c.peak
}

func (c *Chamber) DropRock() int {

	deltaX := c.JetPattern[c.nextJet] + c.JetPattern[(c.nextJet + 1) % len(c.JetPattern)]
	// two below start position
	position := &Point{X: 2 + deltaX, Y: c.Peak() + 2}

	c.nextJet = (c.nextJet + 2) % len(c.JetPattern)

	var landed bool
	// loop
	for ! landed {
		//   jet
		position.X = position.X + c.JetPattern[c.nextJet]
		if c.Collides(c.CurrentRock(), *position) {
			position.X = position.X - c.JetPattern[c.nextJet]
		}
		c.nextJet = (c.nextJet + 1) % len(c.JetPattern)
		//   fall
		position.Y = position.Y - 1
		if position.Y == 0 || c.Collides(c.CurrentRock(), *position) {
			position.Y = position.Y + 1
			// land
			c.Land(c.CurrentRock(), *position)
			landed = true
		}
	}
	c.nextRock = (c.nextRock + 1) % len(c.Rocks)
	return c.Peak()
}

func (c *Chamber) Land(rock Rock, position Point) {
	for p := range rock.Shape {
		chamberPoint := Point{X: p.X + position.X, Y: p.Y + position.Y}
		c.Landed[chamberPoint] = true
	}
	for p := range c.Landed {
		if p.Y > c.peak {
			c.peak = p.Y
		}
	}
	if c.Peak() % 15 == 0 {
		c.PruneLanded()
	}
}

func (c *Chamber) PruneLanded() {
	for p := range c.Landed {
		if p.Y < c.Peak() - 20 {
			delete(c.Landed, p)
		}
	}
}

func (c Chamber) Collides(rock Rock, position Point) bool {
	for p := range rock.Shape {
		if p.X + position.X < 0 { return true }
		if p.X + position.X >= c.Width { return true }
	}
	if position.Y > c.Peak() {
		return false
	}
	return c.CollidesLanded(rock, position)
}

func (c Chamber) CollidesLanded(rock Rock, position Point) bool {
	for p := range rock.Shape {
		chamberPoint := Point{X: p.X + position.X, Y: p.Y + position.Y}
		if c.Landed[chamberPoint] {
			return true
		}
	}
	return false
}

func (c Chamber) CurrentRock() Rock {
	return c.Rocks[c.nextRock]
}

func (c Chamber) Draw() string {
	minY := c.Peak()
	for p := range c.Landed {
		if p.Y < minY {
			minY = p.Y
		}
	}
	canvas := ""
	for y := c.Peak(); y > minY; y-- {
		var line string
		for x := -1; x <= c.Width + 1; x++ {
			switch x {
			case -1, c.Width + 1: line += "|"
				break
			default:
				if c.Landed[Point{X: x, Y: y}] {
					line += "#"
				} else {
					line += "."
				}
			}
		}
		canvas = canvas + line + "\n"
	}
	return canvas
}