package main

import (
	"log"
	"strings"
)

func init() {
	addSolutions(2, problem2)
}

func problem2(ctx *problemContext) {
	var games []*cubeGame
	scanner := ctx.scanner()
	for scanner.scan() {
		g, ok := parseCubeGame(scanner.text())
		if !ok {
			log.Fatalf("Bad game %q", scanner.text())
		}
		games = append(games, g)
	}
	ctx.reportLoad()

	var part1Score int64
	for _, g := range games {
		if g.part1OK() {
			part1Score += g.id
		}
	}

	ctx.reportPart1(part1Score)

	var part2Score int64
	for _, g := range games {
		part2Score += g.power()
	}
	ctx.reportPart2(part2Score)
}

type cubeGame struct {
	id   int64
	sets []cubeSet
}

func (g *cubeGame) part1OK() bool {
	for _, set := range g.sets {
		if set.red > 12 || set.green > 13 || set.blue > 14 {
			return false
		}
	}
	return true
}

func (g *cubeGame) power() int64 {
	var mins cubeSet
	for _, set := range g.sets {
		mins.red = max(mins.red, set.red)
		mins.green = max(mins.green, set.green)
		mins.blue = max(mins.blue, set.blue)
	}
	return int64(mins.red) * int64(mins.green) * int64(mins.blue)
}

func parseCubeGame(s string) (*cubeGame, bool) {
	var g cubeGame
	gameField, rest, ok := strings.Cut(s, ": ")
	if !ok {
		return nil, false
	}
	_, idStr, ok := strings.Cut(gameField, " ")
	if !ok {
		return nil, false
	}
	g.id = parseInt(idStr)
	for _, field := range strings.Split(rest, ";") {
		set, ok := parseCubeSet(strings.TrimSpace(field))
		if !ok {
			return nil, false
		}
		g.sets = append(g.sets, set)
	}
	return &g, true
}

type cubeSet struct {
	red   int
	green int
	blue  int
}

func parseCubeSet(s string) (cubeSet, bool) {
	var set cubeSet
	for _, field := range strings.Split(s, ",") {
		field = strings.TrimSpace(field)
		n, color, ok := strings.Cut(field, " ")
		if !ok {
			return set, false
		}
		num := int(parseInt(n))
		switch color {
		case "red":
			set.red += num
		case "green":
			set.green += num
		case "blue":
			set.blue += num
		default:
			return set, false
		}
	}
	return set, true
}
