package main

import "math/rand/v2"

type Enemy struct {
	*Character
	move func(g *Game, e *Enemy)
}

func IsPositionValid(g *Game, position Vector) bool {
	return position.X >= 0 && position.X < 20 && position.Y < 15 && position.Y >= 0
}

func IsAnyEnemyAtPosition(g *Game, position Vector) bool {
	for _, enemy := range g.level.Enemies {
		if int(enemy.endPosition.X) == int(position.X) && int(enemy.endPosition.Y) == int(position.Y) {
			return true
		}
	}
	return false
}

func IsPositionWalkable(g *Game, position Vector) bool {
	if g.holdingTrolley {
		// if holding trolley, cannot walk into trolley position
		if VectorIs(VectorFloor(position), VectorFloor(VectorAdd(g.player.visiblePosition, g.player.facingDirection))) {
			return false
		}
	}
	return IsPositionValid(g, position) && g.level.Data[int(position.Y)][int(position.X)] == Floor && !IsAnyEnemyAtPosition(g, position)
}
func IsPositionWalkableNotIncludingTrolley(g *Game, position Vector) bool {
	return IsPositionValid(g, position) && g.level.Data[int(position.Y)][int(position.X)] == Floor && !IsAnyEnemyAtPosition(g, position)
}

func (e *Enemy) Die(g *Game) {
	e.deathPhase = 1
	e.startLerpT = -1000
	e.endPosition = e.visiblePosition
	sePlayer := g.rootGame.audioContext.NewPlayerFromBytes(enemyDeathSound)
	sePlayer.SetVolume(
		0.11,
	)
	sePlayer.Play()
}

func createEnemy1(position Vector) *Enemy {
	enemy := &Enemy{
		Character: &Character{
			startPositon:    position,
			endPosition:     position,
			startLerpT:      -1000,
			facingDirection: Vector{1, 0},
			walkSpeed:       .07,
			spriteIndex:     1,
		},
		move: func(g *Game, enemy *Enemy) {
			possibleDirection := []Vector{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

			// find the closest psoition to player if in range 10
			minDistance := 999999.0
			var bestDirection Vector
			for _, direction := range possibleDirection {
				newPosition := VectorAdd(enemy.endPosition, direction)
				if IsPositionWalkable(g, newPosition) {
					distanceToPlayer := VectorMagnitude(VectorAdd(g.player.endPosition, VectorScale(newPosition, -1)))
					if distanceToPlayer < minDistance {
						minDistance = distanceToPlayer
						bestDirection = direction
					}
				}
			}
			if minDistance > 10 {
				// do a random direction in possibleDirection
				tries := 0
				for {
					randomDirection := possibleDirection[rand.IntN(len(possibleDirection))]
					newPosition := VectorAdd(enemy.endPosition, randomDirection)
					if IsPositionValid(g, newPosition) && g.level.Data[int(newPosition.Y)][int(newPosition.X)] == Floor {
						bestDirection = randomDirection
						break
					}
					tries++
					if tries > 10 {
						bestDirection = Vector{0, 1}
						break
					}
				}
			}
			enemy.startPositon = enemy.endPosition
			enemy.endPosition = VectorAdd(enemy.endPosition, bestDirection)
			enemy.startLerpT = g.t
			enemy.facingDirection = bestDirection
		},
	}

	return enemy
}

func createEnemy2(position Vector) *Enemy {
	possibleDirection := []Vector{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

	enemy := &Enemy{
		Character: &Character{
			startPositon:    position,
			endPosition:     position,
			startLerpT:      -1000,
			facingDirection: Vector{1, 0},
			walkSpeed:       .14,
			spriteIndex:     2,
		},
		move: func(g *Game, enemy *Enemy) {
			// move forwards if possible, else turn around
			newPosition := VectorAdd(enemy.endPosition, enemy.facingDirection)
			if IsPositionWalkable(g, newPosition) {
				// move forwards
				enemy.startPositon = enemy.endPosition
				enemy.endPosition = newPosition
				enemy.startLerpT = g.t
			} else {
				// turn around
				enemy.facingDirection = VectorScale(enemy.facingDirection, -1)
			}
		},
	}

	enemy.facingDirection = possibleDirection[rand.IntN(len(possibleDirection))]

	return enemy
}
func createEnemy3(position Vector) *Enemy {
	possibleDirection := []Vector{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

	enemy := &Enemy{
		Character: &Character{
			startPositon:    position,
			endPosition:     position,
			startLerpT:      -1000,
			facingDirection: Vector{1, 0},
			walkSpeed:       .13,
			spriteIndex:     3,
		},
		move: func(g *Game, enemy *Enemy) {
			possibleDirection := []Vector{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

			// find the closest psoition to player if in range 10
			minDistance := 999999.0
			var bestDirection Vector
			for _, direction := range possibleDirection {
				newPosition := VectorAdd(enemy.endPosition, direction)
				if IsPositionWalkable(g, newPosition) {
					distanceToPlayer := VectorMagnitude(VectorAdd(g.player.endPosition, VectorScale(newPosition, -1)))
					if distanceToPlayer < minDistance {
						minDistance = distanceToPlayer
						bestDirection = direction
					}
				}
			}
			enemy.startPositon = enemy.endPosition
			enemy.endPosition = VectorAdd(enemy.endPosition, bestDirection)
			enemy.startLerpT = g.t
			enemy.facingDirection = bestDirection
		},
	}

	enemy.facingDirection = possibleDirection[rand.IntN(len(possibleDirection))]

	return enemy
}

func createEnemy4(position Vector) *Enemy {
	possibleDirection := []Vector{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

	enemy := &Enemy{
		Character: &Character{
			startPositon:    position,
			endPosition:     position,
			startLerpT:      -1000,
			facingDirection: Vector{1, 0},
			walkSpeed:       .2,
			spriteIndex:     4,
		},
		move: func(g *Game, enemy *Enemy) {
			var direction Vector
			if IsPositionWalkable(g, VectorAdd(enemy.endPosition, enemy.facingDirection)) && rand.Float64() < 0.7 {
				direction = enemy.facingDirection
			} else {
				// do a random direction in possibleDirection
				tries := 0
				for {
					randomDirection := possibleDirection[rand.IntN(len(possibleDirection))]
					newPosition := VectorAdd(enemy.endPosition, randomDirection)
					if IsPositionValid(g, newPosition) && g.level.Data[int(newPosition.Y)][int(newPosition.X)] == Floor {
						direction = randomDirection
						break
					}
					tries++
					if tries > 10 {
						direction = Vector{0, 1}
						break
					}
				}
			}
			enemy.startPositon = enemy.endPosition
			enemy.endPosition = VectorAdd(enemy.endPosition, direction)
			enemy.startLerpT = g.t
			enemy.facingDirection = direction
		},
	}

	enemy.facingDirection = possibleDirection[rand.IntN(len(possibleDirection))]

	return enemy
}
