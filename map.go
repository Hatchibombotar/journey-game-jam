package main

import (
	"math/rand/v2"
)

type Map struct {
	Data        [MAP_SIZE][MAP_SIZE]LevelType
	WinPosition Vector
}

type LevelType int

const (
	Filled LevelType = iota
	Empty
	Exit
)

const MAP_SIZE = 16

func GenerateMap() *Map {
	m := &Map{
		Data: [MAP_SIZE][MAP_SIZE]LevelType{},
	}
	for y := range MAP_SIZE {
		for x := range MAP_SIZE {
			m.Data[y][x] = Filled
		}
	}

	currentRow, currentColumn := 0, 0

	remainingTunnels := 32
	maxLength := 8

	directions := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	lastDirection := [2]int{0, 0}

	for remainingTunnels > 0 {
		randomDirection := directions[rand.IntN(len(directions))]
		if (randomDirection[0] == -lastDirection[0] && randomDirection[1] == -lastDirection[1]) ||
			(randomDirection[0] == lastDirection[0] && randomDirection[1] == lastDirection[1]) {
			continue
		}
		randomLength := rand.IntN(maxLength) + 1
		tunnelLength := 0

		for tunnelLength < randomLength {
			if (currentRow == 0 && randomDirection[0] == -1) ||
				(currentColumn == 0 && randomDirection[1] == -1) ||
				(currentRow == MAP_SIZE-1 && randomDirection[0] == 1) ||
				(currentColumn == MAP_SIZE-1 && randomDirection[1] == 1) {
				break
			}
			m.Data[currentRow][currentColumn] = Empty
			currentRow += randomDirection[0]
			currentColumn += randomDirection[1]
			tunnelLength++
		}

		if tunnelLength > 0 {
			lastDirection = randomDirection
			remainingTunnels--
		}

		if remainingTunnels == 0 {
		}
	}

	smallestWinPDistance := 99
	winX, winY := 0, 0
	for y, row := range m.Data {
		for x, sq := range row {
			if sq == Empty && (x+y) < smallestWinPDistance {
				winX, winY = x, y
			}
		}
	}
	m.Data[winY][winX] = Exit
	m.WinPosition = Vector{float64(winX), float64(winY)}

	return m
}
