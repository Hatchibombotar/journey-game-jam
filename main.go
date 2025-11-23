package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type BigGame struct {
	g            *Game
	audioContext *audio.Context
}

type Game struct {
	rootGame *BigGame
	player   *Character
	t        int
	level    *Level
	theMap   *Map

	mapX, mapY int

	projectiles []*Projectile

	holdingTrolley bool

	trolleyPosition  Vector
	trolleyDirection Vector

	viewedLevelData [MAP_SIZE][MAP_SIZE]*Level

	health        int
	timeRemaining float64

	// gameOver bool
	// mainMenu bool

	score int

	droppedTrolleyLevel [2]int

	specialMenu int

	prevLevel               *Level
	prevLevelTransitionTime int

	prevLevelDirection Vector
}

type SpecialMenu = int

const (
	NoMenu SpecialMenu = iota
	MainMenu
	Controls
	GameOver
	Win
)

type Character struct {
	startLerpT      int
	startPositon    Vector
	endPosition     Vector
	visiblePosition Vector
	facingDirection Vector
	spriteIndex     int
	walkSpeed       float64
	deathPhase      int
}

func (c *Character) GetLerpProgress(g *Game) float64 {
	return min((float64(g.t)-float64(c.startLerpT))*c.walkSpeed, 1)
}

func (c *Character) Draw(screen *ebiten.Image, g *Game, offsetX, offsetY float64, isTransition bool) {
	op := &ebiten.DrawImageOptions{}

	c.visiblePosition = LerpVectors(
		c.startPositon,
		c.endPosition,
		c.GetLerpProgress(g),
	)

	spriteDirection := 0

	frame := 1
	if c.GetLerpProgress(g) < 1 {
		frame = (g.t / 6) % 4
	}

	if VectorIs(c.facingDirection, Vector{0, 1}) {
		spriteDirection = 4 + frame
	} else if VectorIs(c.facingDirection, Vector{0, -1}) {
		spriteDirection = 8 + frame
	} else if VectorIs(c.facingDirection, Vector{-1, 0}) {
		spriteDirection = 4 + frame
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(16, 0)
	} else if VectorIs(c.facingDirection, Vector{1, 0}) {
		spriteDirection = 4 + frame
	}

	if c.deathPhase > 0 {
		spriteDirection = 12 + c.deathPhase - 1
	}

	op.GeoM.Translate(float64(c.visiblePosition.X*16), float64(c.visiblePosition.Y*16))
	op.GeoM.Translate(offsetX, offsetY)

	if isTransition {
		op.ColorScale.ScaleWithColor(color.Gray{50})
	}

	img := spritesheet.SubImage(image.Rect(spriteDirection*16, c.spriteIndex*16, (1+spriteDirection)*16, (1+c.spriteIndex)*16)).(*ebiten.Image)
	screen.DrawImage(img, op)
}

type Projectile struct {
	position  Vector
	image     *ebiten.Image
	direction Vector
	lifetime  int
}

var bg = LoadImageFromPath("assets/testbg.png")
var spritesheet = LoadImageFromPath("assets/character.png")
var trolley = LoadImageFromPath("assets/trolley.png")
var trolley_down = LoadImageFromPath("assets/trolley_down.png")
var trolley_up = LoadImageFromPath("assets/trolley_up.png")
var bullet = LoadImageFromPath("assets/bullet.png")
var books = LoadImageFromPath("assets/books.png")
var clockLoot = LoadImageFromPath("assets/clock.png")
var mapLoot = LoadImageFromPath("assets/map.png")
var heart = LoadImageFromPath("assets/heart.png")
var bookshelf = LoadImageFromPath("assets/bookshelf.png")
var bookshelf_top = LoadImageFromPath("assets/bookshelf_top.png")
var vignette = LoadImageFromPath("assets/vignette.png")
var opening = LoadImageFromPath("assets/opening.png")
var gameover = LoadImageFromPath("assets/gameover.png")
var controls = LoadImageFromPath("assets/controls.png")
var winScreen = LoadImageFromPath("assets/end_screen.png")
var exitBack = LoadImageFromPath("assets/exit_back.png")
var exitFront = LoadImageFromPath("assets/exit_front.png")
var slash = LoadImageFromPath("assets/slash.png")
var title = LoadImageFromPath("assets/title.png")
var healthPowerup = LoadImageFromPath("assets/health.png")
var footstepSounds [][]byte = [][]byte{
	ReadOggBytesFromPath("assets/sounds/footstep00.ogg"),
	ReadOggBytesFromPath("assets/sounds/footstep01.ogg"),
	ReadOggBytesFromPath("assets/sounds/footstep02.ogg"),
	ReadOggBytesFromPath("assets/sounds/footstep03.ogg"),
	ReadOggBytesFromPath("assets/sounds/footstep04.ogg"),
	ReadOggBytesFromPath("assets/sounds/footstep05.ogg"),
	ReadOggBytesFromPath("assets/sounds/footstep06.ogg"),
	ReadOggBytesFromPath("assets/sounds/footstep07.ogg"),
	ReadOggBytesFromPath("assets/sounds/footstep08.ogg"),
	ReadOggBytesFromPath("assets/sounds/footstep09.ogg"),
}
var bookOpenSound = ReadOggBytesFromPath("assets/sounds/bookOpen.ogg")
var dropLeatherSound = ReadOggBytesFromPath("assets/sounds/dropLeather.ogg")
var impactSounds [][]byte = [][]byte{
	ReadOggBytesFromPath("assets/sounds/impactPlank_medium_000.ogg"),
	ReadOggBytesFromPath("assets/sounds/impactPlank_medium_001.ogg"),
	ReadOggBytesFromPath("assets/sounds/impactPlank_medium_002.ogg"),
	ReadOggBytesFromPath("assets/sounds/impactPlank_medium_003.ogg"),
	ReadOggBytesFromPath("assets/sounds/impactPlank_medium_004.ogg"),
}
var enemyDeathSound = ReadOggBytesFromPath("assets/sounds/slime_000.ogg")
var loseSound = ReadOggBytesFromPath("assets/sounds/jingles_NES11.ogg")
var winSound = ReadOggBytesFromPath("assets/sounds/jingles_NES12.ogg")
var hurtSound = ReadOggBytesFromPath("assets/sounds/jingles_HIT00_modified.ogg")

// var shootSound = ReadOggBytesFromPath("assets/sounds/laserRetro_002.ogg")
var shootSounds [][]byte = [][]byte{
	ReadOggBytesFromPath("assets/sounds/laserRetro_001.ogg"),
	// ReadOggBytesFromPath("assets/sounds/laserRetro_002.ogg"),
	// ReadOggBytesFromPath("assets/sounds/laserRetro_003.ogg"),
	ReadOggBytesFromPath("assets/sounds/laserRetro_004.ogg"),
}

const PREVLEVELTRANSITIONTIME = 32

func (g *Game) Update(rootGame *BigGame) error {
	g.t += 1
	// inpututil.KeyPressDuration()

	if g.specialMenu == MainMenu {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.specialMenu = Controls
		}
		return nil
	}

	if g.specialMenu == Controls {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.specialMenu = NoMenu
		}
		return nil
	}

	if g.specialMenu == GameOver {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			rootGame.g = createGame(rootGame)
		}
		return nil
	}

	if g.specialMenu == Win {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			rootGame.g = createGame(rootGame)
		}
		return nil
	}

	if g.prevLevelTransitionTime > 0 {
		g.prevLevelTransitionTime -= 1
		return nil
	}

	moveLerpProgress := g.player.GetLerpProgress(g)

	movementVector := Vector{}
	moving := false
	if ebiten.IsKeyPressed(ebiten.KeyD) && moveLerpProgress == 1 { // || ebiten.IsKeyPressed(ebiten.KeyArrowRight) && (g.t%int(walkSpeed) == 0) {
		movementVector = Vector{X: 1}
		moving = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) && moveLerpProgress == 1 {
		movementVector = Vector{X: -1}
		moving = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) && moveLerpProgress == 1 {
		movementVector = Vector{Y: -1}
		moving = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) && moveLerpProgress == 1 {
		movementVector = Vector{Y: 1}
		moving = true
	}

	isDroppedTrolleyInThisLevel := g.droppedTrolleyLevel[0] == g.mapX && g.droppedTrolleyLevel[1] == g.mapY

	if moving {
		g.player.facingDirection = movementVector
		endPosition := VectorAdd(g.player.endPosition, (movementVector))

		movedLevel := false

		if endPosition.X >= 20 {
			g.player.startPositon = g.player.endPosition
			g.player.endPosition = endPosition

			g.player.startPositon.X = -1
			g.player.endPosition.X = 0

			g.player.startLerpT = g.t

			movedLevel = true
			g.mapX += 1
			g.prevLevelDirection = Vector{X: 1}
		} else if endPosition.X < 0 {
			g.player.startPositon = g.player.endPosition
			g.player.endPosition = endPosition

			g.player.startPositon.X = 20
			g.player.endPosition.X = 19

			g.player.startLerpT = g.t
			movedLevel = true
			g.mapX -= 1
			g.prevLevelDirection = Vector{X: -1}

		} else if endPosition.Y >= 15 {
			g.player.startPositon = g.player.endPosition
			g.player.endPosition = endPosition

			g.player.startPositon.Y = -1
			g.player.endPosition.Y = 0

			g.player.startLerpT = g.t
			movedLevel = true
			g.mapY += 1
			g.prevLevelDirection = Vector{Y: 1}

		} else if endPosition.Y < 0 {
			g.player.startPositon = g.player.endPosition
			g.player.endPosition = endPosition

			g.player.startPositon.Y = 15
			g.player.endPosition.Y = 14
			g.player.startLerpT = g.t

			movedLevel = true
			g.mapY -= 1
			g.prevLevelDirection = Vector{Y: -1}

		} else {

			if g.level.Data[int(endPosition.Y)][int(endPosition.X)] == Floor {
				g.player.startPositon = g.player.endPosition
				g.player.endPosition = endPosition
				g.player.startLerpT = g.t
			}
		}

		if movedLevel {
			if g.viewedLevelData[g.mapY][g.mapX] == nil {
				g.viewedLevelData[g.mapY][g.mapX] = generateLevelForMapPosition(g.theMap, g.mapX, g.mapY)
			}
			g.prevLevel = g.level
			g.prevLevelTransitionTime = PREVLEVELTRANSITIONTIME
			g.level = g.viewedLevelData[g.mapY][g.mapX]
			g.projectiles = []*Projectile{}
		}
	}

	for projIndex, projectile := range g.projectiles {
		projectile.lifetime -= 1
		if projectile.lifetime <= 0 {
			g.projectiles = append(g.projectiles[:projIndex], g.projectiles[projIndex+1:]...)
			continue
		}
		projectile.position = VectorAdd(projectile.position, VectorScale(projectile.direction, 0.5))
		// if in wall, remove projectile
		if !IsPositionValid(g, projectile.position) || g.level.Data[int(projectile.position.Y)][int(projectile.position.X)] != Floor {
			g.projectiles = append(g.projectiles[:projIndex], g.projectiles[projIndex+1:]...)
			break
		}

		// var nearestEnemy *Enemy
		// nearestEnemyIndex := -1
		// minDistance := 9999.9
		// for enemyIndex, enemy := range g.level.Enemies {

		// 	if enemy.deathPhase > 0 {
		// 		continue
		// 	}
		// 	distance := VectorMagnitude(VectorAdd(enemy.visiblePosition, VectorScale(projectile.position, -1)))
		// 	if distance < minDistance {
		// 		minDistance = distance
		// 		nearestEnemy = enemy
		// 		nearestEnemyIndex = enemyIndex
		// 	}
		// }
		// if nearestEnemyIndex == -1 {
		// 	continue
		// }
		// if minDistance < 0.5 {
		// 	g.score += 3
		// 	nearestEnemy.Die(g)
		// 	// g.level.Enemies = append(g.level.Enemies[:nearestEnemyIndex], g.level.Enemies[nearestEnemyIndex+1:]...)
		// 	g.projectiles = append(g.projectiles[:projIndex], g.projectiles[projIndex+1:]...)
		// 	break
		// } else if minDistance < 1.5 {
		// 	// home in on enemy
		// 	projectile.direction = VectorScale(VectorNormalise(VectorAdd(nearestEnemy.visiblePosition, VectorScale(projectile.position, -1))), 0.5)
		// }

		// if enemy within 1.5 of projectile, kill all enemies within 1.5 of projectile and remove projectile
		hasDeleted := false
		for _, enemy := range g.level.Enemies {
			if enemy.deathPhase > 0 {
				continue
			}
			distance := VectorMagnitude(VectorAdd(enemy.visiblePosition, VectorScale(projectile.position, -1)))
			if distance <= 1.5 {
				enemy.Die(g)
			}
		}

		if hasDeleted {
			g.projectiles = append(g.projectiles[:projIndex], g.projectiles[projIndex+1:]...)
			break
		}

	}

	distanceFromTrolley := VectorMagnitude(VectorSubtract(g.player.visiblePosition, g.trolleyPosition))

	canPickupTrolley := !g.holdingTrolley && isDroppedTrolleyInThisLevel && distanceFromTrolley < 2
	trolleyPosition := VectorAdd(g.player.visiblePosition, g.player.facingDirection)

	if g.holdingTrolley {
		g.trolleyDirection = g.player.facingDirection
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if canPickupTrolley {
			g.holdingTrolley = true
		} else {
			// position := VectorAdd
			p := &Projectile{
				position:  g.player.endPosition,
				image:     slash,
				direction: g.player.facingDirection,
				lifetime:  24,
			}
			g.timeRemaining -= 12
			g.projectiles = append(g.projectiles, p)
			sePlayer := rootGame.audioContext.NewPlayerFromBytes(shootSounds[rand.IntN(len(shootSounds))])
			sePlayer.SetVolume(
				0.1,
			)
			sePlayer.Play()

			if g.holdingTrolley {
				g.trolleyPosition = g.player.visiblePosition

				for range 5 {
					newPos := VectorAdd(g.trolleyPosition, g.player.facingDirection)
					if IsPositionWalkableNotIncludingTrolley(g, newPos) {
						g.trolleyPosition = newPos
					}
				}
				g.trolleyPosition = VectorSubtract(g.trolleyPosition, g.player.facingDirection)

				g.trolleyDirection = g.player.facingDirection
				// if inpututil.IsKeyJustPressed(ebiten.KeyM) {
				g.holdingTrolley = false
				g.droppedTrolleyLevel[0] = g.mapX
				g.droppedTrolleyLevel[1] = g.mapY
				// }
			}
		}
	}

	if g.holdingTrolley {
		for i, loot := range g.level.Loot {
			if VectorIs(VectorFloor(loot.position), VectorFloor(trolleyPosition)) {
				g.level.Loot = append(g.level.Loot[:i], g.level.Loot[i+1:]...)
				loot.onPickup(g, loot)
				break
			}
		}
	} else {
		for i, loot := range g.level.Loot {
			if !loot.canBePickedUpWithoutCart {
				continue
			}
			if VectorIs(VectorFloor(loot.position), VectorFloor(trolleyPosition)) {
				g.level.Loot = append(g.level.Loot[:i], g.level.Loot[i+1:]...)
				loot.onPickup(g, loot)
				break
			}
		}
	}

	for _, enemy := range g.level.Enemies {
		if enemy.deathPhase > 0 {
			continue
		}
		if enemy.GetLerpProgress(g) == 1 {
			enemy.move(g, enemy)
		}

		if g.holdingTrolley && VectorIs(VectorFloor(enemy.visiblePosition), VectorFloor(trolleyPosition)) {
			enemy.Die(g)
			g.score += 10
			break
		} else if VectorIs(VectorFloor(enemy.visiblePosition), VectorFloor(g.player.visiblePosition)) {
			enemy.Die(g)
			g.health -= 1
			sePlayer := rootGame.audioContext.NewPlayerFromBytes(impactSounds[rand.IntN(len(impactSounds))])
			sePlayer.SetVolume(
				1,
			)
			sePlayer.Play()
			break
		}
		// for _, projectile := range g.projectiles {
		// 	if VectorIs(VectorFloor(enemy.visiblePosition), VectorFloor(projectile.position)) {
		// 		g.level.Enemies = append(g.level.Enemies[:i], g.level.Enemies[i+1:]...)
		// 		break
		// 	}
		// }
	}

	// if inpututil.IsKeyJustPressed(ebiten.KeyF5) {
	// 	fillInMapWithLevels(g)
	// }

	g.timeRemaining -= 1.0 / 24
	if g.timeRemaining <= 0 || g.health <= 0 {
		g.specialMenu = GameOver
		sePlayer := g.rootGame.audioContext.NewPlayerFromBytes(loseSound)
		sePlayer.SetVolume(
			0.5,
		)
		sePlayer.Play()
	}

	// if inpututil.IsKeyJustPressed(ebiten.KeyF6) {
	// 	g.specialMenu = GameOver
	// }

	if g.player.GetLerpProgress(g) < 1 && g.t%20 == 0 {
		sePlayer := rootGame.audioContext.NewPlayerFromBytes(footstepSounds[rand.IntN(len(footstepSounds))])
		sePlayer.SetVolume(
			0.15,
		)
		sePlayer.Play()
	}

	for i, enemy := range g.level.Enemies {
		if enemy.deathPhase > 0 && g.t%10 == 0 {
			enemy.deathPhase += 1
		}
		if enemy.deathPhase > 4 {
			g.level.Enemies = append(g.level.Enemies[:i], g.level.Enemies[i+1:]...)
			break
		}
	}

	if g.theMap.Data[g.mapY][g.mapX] == Exit {
		// if player near middle of screen:
		if VectorMagnitude(VectorSubtract(g.player.visiblePosition, Vector{10, 7.5})) < 2 {
			g.specialMenu = Win
			g.player.startPositon = g.player.visiblePosition
			g.player.walkSpeed = 0.02
			g.player.endPosition = VectorAdd(g.player.visiblePosition, Vector{Y: 15})

			sePlayer := g.rootGame.audioContext.NewPlayerFromBytes(winSound)
			sePlayer.SetVolume(
				0.5,
			)
			sePlayer.Play()
			// g.player.facingDirection = Vector{0, 1}
		}
	}
	return nil
}

func fillInMapWithLevels(g *Game) {
	for y := range MAP_SIZE {
		for x := range MAP_SIZE {
			if g.theMap.Data[y][x] == Filled {
				continue
			}
			g.viewedLevelData[y][x] = generateLevelForMapPosition(g.theMap, x, y)
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {

	if g.prevLevelTransitionTime > 0 {
		progress := 1 - (float64(g.prevLevelTransitionTime) / float64(PREVLEVELTRANSITIONTIME))

		xDif := -320 * progress * g.prevLevelDirection.X
		yDif := -240 * progress * g.prevLevelDirection.Y

		xDif2 := 320 * g.prevLevelDirection.X
		yDif2 := 240 * g.prevLevelDirection.Y

		// prevLevel := g.viewedLevelData[g.mapY-int(g.prevLevelDirection.Y)][g.mapX-int(g.prevLevelDirection.X)]

		DrawLevel(screen, g, g.prevLevel, g.mapX+1, g.mapY+1, float64(xDif), float64(yDif), true)
		DrawLevel(screen, g, g.level, g.mapX, g.mapY, xDif2+float64(xDif), yDif2+float64(yDif), true)
	} else {
		DrawLevel(screen, g, g.level, g.mapX, g.mapY, 0, 0, false)
	}

	if g.specialMenu != NoMenu {
		screen.DrawImage(vignette, nil)
	}

	if g.specialMenu == MainMenu {
		screen.DrawImage(title, nil)
		screen.DrawImage(opening, nil)
		return
	}

	if g.specialMenu == Controls {
		screen.DrawImage(title, nil)
		screen.DrawImage(controls, nil)
		return
	}

	if g.specialMenu == Win {
		screen.DrawImage(winScreen, nil)

		mediumText := &text.GoTextFace{
			Source: fontFaceSource,
			Size:   16,
		}

		op := &text.DrawOptions{}
		op.GeoM.Translate(128, 180)
		text.Draw(screen, fmt.Sprintln("Score: ", g.score), mediumText, op)
		return
	}

	if g.specialMenu == GameOver {
		screen.DrawImage(gameover, nil)
		return
	}

	uiTransaparency := 175
	if g.player.visiblePosition.X < 5 && g.player.visiblePosition.Y < 5 {
		uiTransaparency = 50
	}

	for i := range g.health {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(44+i*9), 8)
		// make transparent with uiTransaparency
		op.ColorScale.ScaleAlpha(float32(uiTransaparency) / 255.0)

		screen.DrawImage(heart, op)
	}

	mapOffsetX, mapOffsetY := 8.0, 8.0
	vector.FillRect(screen, float32(mapOffsetX), float32(mapOffsetY), 32, 32, color.RGBA{0, 0, 0, uint8(uiTransaparency)}, false)
	for y := range g.viewedLevelData {
		for x := range g.viewedLevelData[y] {
			if g.viewedLevelData[y][x] != nil {
				tileColor := color.RGBA{180, 180, 180, uint8(uiTransaparency)}
				if g.mapX == x && g.mapY == y {
					tileColor = color.RGBA{255, 255, 255, uint8(uiTransaparency)}
				}
				if g.theMap.Data[y][x] == Exit {
					tileColor = color.RGBA{255, 215, 0, uint8(uiTransaparency)}
				}
				vector.FillRect(screen, float32(mapOffsetX)+float32(x*2), float32(mapOffsetY)+float32(y*2), 2, 2, tileColor, false)
			}
		}
	}

	// time remaining
	vector.FillRect(screen, 44, 16, float32(g.timeRemaining)*0.3, 5, color.RGBA{100, 100, 200, uint8(uiTransaparency)}, false)

	scoreText := &text.GoTextFace{
		Source: fontFaceSource,
		Size:   8,
	}
	op1 := &text.DrawOptions{}
	op1.GeoM.Translate(44, 22)
	op1.ColorScale.ScaleAlpha(float32(uiTransaparency) / 255.0)

	text.Draw(screen, fmt.Sprintln("Score: ", g.score), scoreText, op1)
}

func (big *BigGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func (big *BigGame) Update() error {
	return big.g.Update(big)
}
func (big *BigGame) Draw(screen *ebiten.Image) {
	big.g.Draw(screen)
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")

	root := &BigGame{
		audioContext: audio.NewContext(SAMPLE_RATE),
	}
	g := createGame(root)
	root.g = g
	if err := ebiten.RunGame(root); err != nil {
		log.Fatal(err)
	}
}

const SAMPLE_RATE = 48000

func createGame(rootgame *BigGame) *Game {
	g := &Game{
		theMap:          GenerateMap(),
		mapX:            0,
		mapY:            0,
		holdingTrolley:  true,
		viewedLevelData: [MAP_SIZE][MAP_SIZE]*Level{},
		health:          5,
		timeRemaining:   300,
		score:           0,
		rootGame:        rootgame,
		specialMenu:     MainMenu,
	}
	g.level = generateLevelForMapPosition(g.theMap, g.mapX, g.mapY)
	g.viewedLevelData[g.mapY][g.mapX] = g.level

	p := &Character{
		startPositon:    Vector{1, 7},
		endPosition:     Vector{1, 7},
		startLerpT:      -1000,
		facingDirection: Vector{1, 0},
		walkSpeed:       .15,
	}
	g.player = p
	return g
}
