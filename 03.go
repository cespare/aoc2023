package main

func init() {
	addSolutions(3, problem3)
}

func problem3(ctx *problemContext) {
	var g grid[byte]
	var nums []schematicNum
	scanner := ctx.scanner()
	var y int64
	for scanner.scan() {
		line := scanner.text()
		g.addRow([]byte(line))
		for x := int64(0); x < int64(len(line)); x++ {
			r := line[x]
			start := x
			for r >= '0' && r <= '9' {
				x++
				if x == int64(len(line)) {
					break
				}
				r = line[x]
			}
			if x == start {
				continue
			}
			num := schematicNum{
				val:   parseInt(line[start:x]),
				start: vec2{start, y},
				len:   x - start,
			}
			nums = append(nums, num)
			x--
		}
		y++
	}
	ctx.reportLoad()

	var part1Sum int64
	gears := make(map[vec2][]schematicNum)
	for _, num := range nums {
		perim := []vec2{
			{num.start.x - 1, num.start.y},
			{num.start.x + num.len, num.start.y},
		}
		for x := num.start.x - 1; x <= num.start.x+num.len; x++ {
			perim = append(perim, vec2{x, num.start.y - 1})
			perim = append(perim, vec2{x, num.start.y + 1})
		}
		found := false
		for _, v := range perim {
			if !g.contains(v) {
				continue
			}
			r := g.at(v)
			if r == '.' || (r >= '0' && r <= '9') {
				continue
			}
			found = true
			if r == '*' {
				gears[v] = append(gears[v], num)
			}
		}
		if found {
			part1Sum += num.val
		}
	}

	ctx.reportPart1(part1Sum)

	var part2Sum int64
	for _, gearNums := range gears {
		if len(gearNums) != 2 {
			continue
		}
		part2Sum += gearNums[0].val * gearNums[1].val
	}

	ctx.reportPart2(part2Sum)
}

type schematicNum struct {
	val   int64
	start vec2
	len   int64
}
