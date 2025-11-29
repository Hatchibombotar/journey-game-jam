package main

import (
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
)

type Tile = int

const (
	Floor = iota
	Wall
	WallTop
)

type Level struct {
	Data    [15][20]Tile
	Enemies []*Enemy
	Loot    []Loot
}

type Loot struct {
	position                 Vector
	image                    *ebiten.Image
	onPickup                 func(g *Game, l Loot)
	canBePickedUpWithoutCart bool
}

func createWalledLevel() *Level {
	l := &Level{
		Data: [15][20]Tile{},
	}
	for i := range 20 {
		l.Data[0][i] = WallTop
	}

	for i := range 20 {
		l.Data[14][i] = Wall
	}

	for i := range 15 {
		l.Data[i][0] = Wall
	}

	for i := range 15 {
		l.Data[i][19] = Wall
	}

	return l
}

func generateLevelForMapPosition(m *Map, x, y int) *Level {
	level := createWalledLevel()
	targetLevelType := m.Data[y][x]

	var tileRight LevelType
	if x < MAP_SIZE-1 {
		tileRight = m.Data[y][x+1]
	} else {
		tileRight = Filled
	}

	var tileLeft LevelType
	if x > 0 {
		tileLeft = m.Data[y][x-1]
	} else {
		tileLeft = Filled
	}

	var tileUp LevelType
	if y > 0 {
		tileUp = m.Data[y-1][x]
	} else {
		tileUp = Filled
	}

	var tileDown LevelType
	if y < MAP_SIZE-1 {
		tileDown = m.Data[y+1][x]
	} else {
		tileDown = Filled
	}

	if tileLeft != Filled {
		level.Data[3][0] = WallTop
		for i := range 6 {
			level.Data[4+i][0] = Floor
		}
	}
	if tileRight != Filled {
		level.Data[3][19] = WallTop
		for i := range 6 {
			level.Data[4+i][19] = Floor
		}
	}
	if tileUp != Filled {
		for i := range 6 {
			level.Data[0][7+i] = Floor
		}
	}
	if tileDown != Filled {
		for i := range 6 {
			level.Data[14][7+i] = Floor
		}
	}

	switch targetLevelType {
	case Filled:
		return level
	case Exit:
		return level
	}

	patterns := []func(level *Level){
		createPattern0,
		createPattern1,
		createPattern2,
	}
	pattern := patterns[rand.IntN(len(patterns))]
	pattern(level)

	// levelDifficulty := (float64(y+x) / 4)

	levelDifficulty := 22 - VectorMagnitude(VectorSubtract(m.WinPosition, Vector{float64(x), float64(y)}))

	enemyAmount := rand.IntN(int(2 + (levelDifficulty / 3)))
	if x == 0 && y == 0 {
		enemyAmount = 0
	}

	for range enemyAmount {
		enemyType := rand.IntN(4)
		position := Vector{X: 1 + float64(rand.IntN(16)), Y: 1 + float64(rand.IntN(13))}
		switch enemyType {
		case 0:
			level.Enemies = append(level.Enemies, createEnemy1(position))
		case 1:
			level.Enemies = append(level.Enemies, createEnemy2(position))
			level.Enemies = append(level.Enemies, createEnemy2(position))
			level.Enemies = append(level.Enemies, createEnemy2(position))
		case 2:
			level.Enemies = append(level.Enemies, createEnemy3(position))
		case 3:
			level.Enemies = append(level.Enemies, createEnemy4(position))
		}
	}

	if levelDifficulty > 18 {
		for range 1 + rand.IntN(4) {
			if tileUp != Filled {
				level.Enemies = append(level.Enemies, createEnemy2(Vector{8, 1}))
			}
			if tileDown != Filled {
				level.Enemies = append(level.Enemies, createEnemy2(Vector{8, 13}))
			}
			if tileLeft != Filled {
				level.Enemies = append(level.Enemies, createEnemy2(Vector{1, 1}))
			}
			if tileRight != Filled {
				level.Enemies = append(level.Enemies, createEnemy2(Vector{8, 18}))
			}
		}
	}

	s2 := rand.NewPCG(uint64(x), uint64(y)*1000)
	r2 := rand.New(s2)

	// add loot randomly in level
	numLoot := r2.IntN(5) + 1
	for range numLoot {
		for {
			lootX := r2.IntN(20)
			lootY := r2.IntN(15)
			if level.Data[lootY][lootX] == Floor {
				lootType := r2.Float64()

				if lootType < 0.7 {
					level.Loot = append(level.Loot, Loot{
						image:    books,
						position: Vector{float64(lootX), float64(lootY)},
						onPickup: func(g *Game, l Loot) {
							g.score += 3
							sePlayer := g.rootGame.audioContext.NewPlayerFromBytes(bookOpenSound)
							sePlayer.SetVolume(
								1,
							)
							sePlayer.Play()
						},
					})
				} else if lootType < 0.85 {
					level.Loot = append(level.Loot, Loot{
						image:    healthPowerup,
						position: Vector{float64(lootX), float64(lootY)},
						onPickup: func(g *Game, l Loot) {
							g.health += 1
							sePlayer := g.rootGame.audioContext.NewPlayerFromBytes(bookOpenSound)
							sePlayer.SetVolume(
								1,
							)
							sePlayer.Play()
						},
						canBePickedUpWithoutCart: true,
					})
				} else if lootType < 0.95 {
					level.Loot = append(level.Loot, Loot{
						image:    clockLoot,
						position: Vector{float64(lootX), float64(lootY)},
						onPickup: func(g *Game, l Loot) {
							g.timeRemaining += 30
							sePlayer := g.rootGame.audioContext.NewPlayerFromBytes(bookOpenSound)
							sePlayer.SetVolume(
								1,
							)
							sePlayer.Play()
						},
						canBePickedUpWithoutCart: true,
					})
				} else {
					level.Loot = append(level.Loot, Loot{
						image:    mapLoot,
						position: Vector{float64(lootX), float64(lootY)},
						onPickup: func(g *Game, l Loot) {
							fillInMapWithLevels(g)
							sePlayer := g.rootGame.audioContext.NewPlayerFromBytes(bookOpenSound)
							sePlayer.SetVolume(
								1,
							)
							sePlayer.Play()
						},
						canBePickedUpWithoutCart: true,
					})
				}
				break
			}
		}
	}

	return level
}

func createPattern0(level *Level) {
}
func createPattern1(level *Level) {
	for x := range 10 {
		level.Data[4][5+x] = WallTop
		level.Data[7][5+x] = WallTop
		level.Data[10][5+x] = WallTop
	}
}
func createPattern2(level *Level) {
	for x := range 10 {
		level.Data[5][5+x] = Wall
		level.Data[6][5+x] = Wall
		level.Data[7][5+x] = Wall
		level.Data[8][5+x] = Wall
		level.Data[9][5+x] = WallTop
	}
}
