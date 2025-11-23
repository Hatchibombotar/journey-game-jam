package main

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func DrawLevel(screen *ebiten.Image, g *Game, level *Level, mapX, mapY int, offsetX, offsetY float64, isTransition bool) {
	bigOp := &ebiten.DrawImageOptions{}
	bigOp.GeoM.Translate(offsetX, offsetY)
	screen.DrawImage(bg, bigOp)

	for _, loot := range level.Loot {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(loot.position.X*16, loot.position.Y*16)
		op.GeoM.Translate(offsetX, offsetY)
		screen.DrawImage(loot.image, op)
	}

	isExit := VectorIs(g.theMap.WinPosition, Vector{float64(mapX), float64(mapY)})

	if isExit {
		screen.DrawImage(exitBack, bigOp)
	}

	for _, projectile := range g.projectiles {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-24, -24)
		op.GeoM.Rotate(VectorAngle(projectile.direction))

		op.GeoM.Translate(projectile.position.X*16, projectile.position.Y*16)
		op.GeoM.Translate(offsetX, offsetY)
		screen.DrawImage(projectile.image, op)
	}

	for _, enemy := range level.Enemies {
		enemy.Draw(screen, g, offsetX, offsetY, isTransition)
	}

	g.player.Draw(screen, g, offsetX, offsetY, isTransition)

	op := &ebiten.DrawImageOptions{}

	isDroppedTrolleyInThisLevel := g.droppedTrolleyLevel[0] == g.mapX && g.droppedTrolleyLevel[1] == g.mapY

	if g.holdingTrolley || isDroppedTrolleyInThisLevel {
		trolleyPosition := g.trolleyPosition
		if g.holdingTrolley {
			trolleyPosition = g.player.visiblePosition
		}
		if VectorIs(g.trolleyDirection, Vector{1, 0}) {
			op.GeoM.Translate(trolleyPosition.X*16, trolleyPosition.Y*16)
			op.GeoM.Translate(offsetX, offsetY)
			screen.DrawImage(trolley, op)
		} else if VectorIs(g.trolleyDirection, Vector{-1, 0}) {
			op.GeoM.Scale(-1, 1)
			op.GeoM.Translate(trolleyPosition.X*16, trolleyPosition.Y*16)
			op.GeoM.Translate(16, 0)
			op.GeoM.Translate(offsetX, offsetY)
			screen.DrawImage(trolley, op)
		} else if VectorIs(g.trolleyDirection, Vector{0, 1}) {
			op.GeoM.Translate(trolleyPosition.X*16, trolleyPosition.Y*16)
			op.GeoM.Translate(offsetX, offsetY)
			screen.DrawImage(trolley_down, op)
		} else if VectorIs(g.trolleyDirection, Vector{0, -1}) {
			op.GeoM.Translate(trolleyPosition.X*16, trolleyPosition.Y*16)
			op.GeoM.Translate(0, -16)
			op.GeoM.Translate(offsetX, offsetY)
			screen.DrawImage(trolley_up, op)
		}
	}

	for y := range level.Data {
		for x, tile := range level.Data[y] {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x)*16, float64(y)*16)
			op.GeoM.Translate(offsetX, offsetY)
			switch tile {
			case Floor:
			case WallTop:
				shelf_index := (x * 3) % 2
				shelf := bookshelf_top.SubImage(image.Rect(shelf_index*16, 0, (shelf_index+1)*16, 16)).(*ebiten.Image)
				screen.DrawImage(shelf, op)
			case Wall:
				screen.DrawImage(bookshelf, op)
			}
		}
	}

	if isExit && g.specialMenu == Win {
		screen.DrawImage(exitFront, bigOp)
	}

	if !g.holdingTrolley && isDroppedTrolleyInThisLevel {
		distanceFromTrolley := VectorMagnitude(VectorSubtract(g.player.visiblePosition, g.trolleyPosition))

		trolleyCenter := g.trolleyPosition
		if g.trolleyDirection.X == 1 {
			trolleyCenter = VectorAdd(trolleyCenter, Vector{1.5, 0.5})
		} else if g.trolleyDirection.X == -1 {
			trolleyCenter = VectorAdd(trolleyCenter, Vector{-0.5, 0.5})
		} else if g.trolleyDirection.Y == 1 {
			trolleyCenter = VectorAdd(trolleyCenter, Vector{0.5, 1.5})
		} else if g.trolleyDirection.Y == -1 {
			trolleyCenter = VectorAdd(trolleyCenter, Vector{0.5, -0.5})
		}

		if !g.holdingTrolley && distanceFromTrolley < 2.0 {
			vector.StrokeCircle(screen, float32(trolleyCenter.X*16)+float32(offsetX), float32(trolleyCenter.Y*16)+float32(offsetY), 16, 2, color.RGBA{255, 255, 255, 1}, false)
		}
	}

	screen.DrawImage(vignette, bigOp)
}
