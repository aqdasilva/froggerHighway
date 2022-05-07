package main

import (
	"embed"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/rand"
	"time"
)

type Mode int

var (
	arcadeFont font.Face
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
	frogSpeed    = 2
	maxAngle     = 180
	//
)

var (
	emptyImage    = ebiten.NewImage(3, 3)
	emptySubImage = emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	emptyImage.Fill(color.Black)
}

var froggerPict *ebiten.Image
var blackRoadPict *ebiten.Image
var road1 *ebiten.Image
var road2 *ebiten.Image

func init() {
	rand.Seed(time.Now().UnixNano())
}

//frog sprite
type Sprite struct {
	picts *ebiten.Image
	//frog cords
	xloc  int
	yloc  int
	dx    int
	dy    int
	angle int
	//car cords
}

//go:embed graphics/frogger.png graphics/redCar.png graphics/yellowCar.png graphics/turtle.png graphics/wood.png graphics/turtle2.png graphics/wood2.png graphics/road.png
var embeddedAssets embed.FS
var road = ebiten.NewImage(ScreenWidth, ScreenHeight)

type Game struct {
	mode             Mode
	modeGameover     int
	froggerSprite    Sprite
	redcarSprites    Sprite
	yellowCarSprites Sprite
	turtleSprite     Sprite
	turtle2Sprite    Sprite
	woodSprite       Sprite
	wood2Sprite      Sprite
	roadSprite       Sprite
	froggy           []frogPlacement
	redCar           frogPlacement
	drawOps          ebiten.DrawImageOptions
	frogCollides     bool
	score            int
	carHighway       []Sprite
	timer            int
	moveTime         int
	angle            int
}
type frogPlacement struct {
	XSpot int
	YSpot int
}

func rectangleCollision(r1X, r1Y, r1Width, r1Height, r2X, r2Y, r2Width, r2Height float64) bool {
	return r1X-r1Width/2-r2Width/2 <= r2X && r2X <= r1X+r1Width/2+r2Width/2 && r1Y-r1Height/2-r2Height/2 <= r2Y && r2Y <= r1Y+r1Height/2+r2Height/2
}

func (g *Game) reset() {
	g.froggerSprite = Sprite{
		picts: loadPNGImageFromEmbedded("frogger.png"),
		xloc:  400,
		yloc:  700,
		dx:    0,
		dy:    0,
	}
	g.redcarSprites = Sprite{
		picts: loadPNGImageFromEmbedded("redCar.png"),
		xloc:  700,
		yloc:  300,
		dx:    0,
		dy:    0,
	}
	g.yellowCarSprites = Sprite{
		picts: loadPNGImageFromEmbedded("yellowCar.png"),
		xloc:  600,
		yloc:  160,
		dx:    0,
		dy:    0,
	}
	g.turtleSprite = Sprite{
		picts: loadPNGImageFromEmbedded("turtle.png"),
		xloc:  500,
		yloc:  100,
		dx:    0,
		dy:    0,
	}
	g.turtle2Sprite = Sprite{
		picts: loadPNGImageFromEmbedded("turtle2.png"),
		xloc:  300,
		yloc:  200,
		dx:    0,
		dy:    0,
	}
	g.woodSprite = Sprite{
		picts: loadPNGImageFromEmbedded("wood.png"),
		xloc:  150,
		yloc:  100,
		dx:    0,
		dy:    0,
	}
	g.wood2Sprite = Sprite{
		picts: loadPNGImageFromEmbedded("wood2.png"),
		xloc:  20,
		yloc:  30,
		dx:    0,
		dy:    0,
	}
	g.score = 0

}

//frog hits car
func carRedSquishesFrog(frog, car Sprite) bool {
	carWidth, carHieght := car.picts.Size()
	frogWidth, frogHeight := frog.picts.Size()
	if frog.xloc < car.xloc+carWidth &&
		frog.xloc+frogWidth > car.xloc &&
		frog.yloc < car.yloc+carHieght &&
		frog.yloc+frogHeight > car.yloc {
		return true
	}
	return false
}
func angleHitsFrog(frog, car Sprite) bool {
	carAngle, carHieght := car.picts.Size()
	frogWidth, frogHeight := frog.picts.Size()
	if frog.xloc < car.xloc+carAngle &&
		frog.xloc+frogWidth > car.xloc &&
		frog.yloc < car.yloc+carHieght &&
		frog.yloc+frogHeight > car.yloc {
		return true
	}
	return false
}

func carYellowSquishesFrog(frog, car Sprite) bool {
	carWidth, carHieght := car.picts.Size()
	frogWidth, frogHeight := frog.picts.Size()
	if frog.xloc < car.xloc+carWidth &&
		frog.xloc+frogWidth > car.xloc &&
		frog.yloc < car.yloc+carHieght &&
		frog.yloc+frogHeight > car.yloc {
		return true
	}
	return false
}

//frog gotta move move
func (g *Game) frogGottaJump() bool {
	return g.timer%g.moveTime == 0
}

func (g *Game) Update() error {
	movements(g)
	g.froggerSprite.yloc += g.froggerSprite.dy
	g.froggerSprite.xloc += g.froggerSprite.dx
	if g.frogCollides == true {
		g.reset()
	}
	if g.frogCollides == false {
		g.frogCollides = carRedSquishesFrog(g.froggerSprite, g.redcarSprites)
	}
	if g.frogCollides == false {
		g.frogCollides = carYellowSquishesFrog(g.froggerSprite, g.yellowCarSprites)
	}
	if g.frogCollides == false {
		g.frogCollides = angleHitsFrog(g.froggerSprite, g.redcarSprites)
	}

	g.angle++
	if g.angle == maxAngle {
		g.angle = 5
	}
	return nil
}

func movements(g *Game) {
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		g.froggerSprite.dy = -frogSpeed
	} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		g.froggerSprite.dy = frogSpeed
	} else if inpututil.IsKeyJustReleased(ebiten.KeyUp) || inpututil.IsKeyJustReleased(ebiten.KeyDown) {
		g.froggerSprite.dy = 0
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		g.froggerSprite.dx = -frogSpeed
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		g.froggerSprite.dx = frogSpeed
	} else if inpututil.IsKeyJustReleased(ebiten.KeyArrowLeft) || inpututil.IsKeyJustReleased(ebiten.KeyArrowRight) {
		g.froggerSprite.dx = 0
	}
	g.froggerSprite.yloc += g.froggerSprite.dy
	if g.froggerSprite.yloc <= 0 {
		g.froggerSprite.dy = 0
		g.froggerSprite.yloc = 0
	} else if g.froggerSprite.yloc+g.froggerSprite.picts.Bounds().Size().Y > ScreenHeight {
		g.froggerSprite.dy = 0
		g.froggerSprite.yloc = ScreenHeight - g.froggerSprite.picts.Bounds().Size().Y
	}
	g.froggerSprite.xloc += g.froggerSprite.dx
	if g.froggerSprite.xloc <= 0 {
		g.froggerSprite.dx = 0
		g.froggerSprite.xloc = 0
	} else if g.froggerSprite.xloc+g.froggerSprite.picts.Bounds().Size().X > ScreenWidth {
		g.froggerSprite.dx = 0
		g.froggerSprite.xloc = ScreenWidth - g.froggerSprite.picts.Bounds().Size().X
	}

}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x80, 0xa0, 0xc0, 0xff})
	ebitenutil.DebugPrint(screen, "Score: %d")
	//first road
	road1 = ebiten.NewImage(800, 20)
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(64, 150)
	road1.Fill(color.Black)
	screen.DrawImage(road1, opts)
	g.drawOps.GeoM.Reset()
	//2nd road
	road2 = ebiten.NewImage(800, 20)
	opts2 := &ebiten.DrawImageOptions{}
	opts2.GeoM.Translate(64, 300)
	road2.Fill(color.Black)
	screen.DrawImage(road2, opts2)
	g.drawOps.GeoM.Reset()

	//frog image
	g.drawOps.GeoM.Translate(float64(g.froggerSprite.xloc), float64(g.froggerSprite.yloc))
	screen.DrawImage(g.froggerSprite.picts, &g.drawOps)
	g.drawOps.GeoM.Reset()
	//redcar img
	if !g.frogCollides {
		g.drawOps.GeoM.Translate(float64(g.redcarSprites.xloc), float64(g.redcarSprites.yloc))
		g.drawOps.GeoM.Rotate(0.5 * math.Pi * float64(g.angle) / maxAngle)
		screen.DrawImage(g.redcarSprites.picts, &g.drawOps)
		g.drawOps.GeoM.Reset()
	}
	//yellow car img
	g.drawOps.GeoM.Translate(float64(g.yellowCarSprites.xloc), float64(g.yellowCarSprites.yloc))
	g.drawOps.GeoM.Rotate(0.5 * math.Pi * float64(g.angle) / maxAngle)
	screen.DrawImage(g.yellowCarSprites.picts, &g.drawOps)
	g.drawOps.GeoM.Reset()
	//turtle img
	g.drawOps.GeoM.Translate(float64(g.turtleSprite.xloc), float64(g.turtleSprite.yloc))
	g.drawOps.GeoM.Rotate(0.08 * math.Pi * float64(g.angle) / maxAngle)
	screen.DrawImage(g.turtleSprite.picts, &g.drawOps)
	g.drawOps.GeoM.Reset()
	//turtle2
	g.drawOps.GeoM.Translate(float64(g.turtle2Sprite.xloc), float64(g.turtle2Sprite.yloc))
	g.drawOps.GeoM.Rotate(0.08 * math.Pi * float64(g.angle) / maxAngle)
	screen.DrawImage(g.turtle2Sprite.picts, &g.drawOps)
	g.drawOps.GeoM.Reset()
	//wood img
	g.drawOps.GeoM.Translate(float64(g.woodSprite.xloc), float64(g.woodSprite.yloc))
	g.drawOps.GeoM.Rotate(0.2 * math.Pi * float64(g.angle) / maxAngle)
	screen.DrawImage(g.woodSprite.picts, &g.drawOps)
	g.drawOps.GeoM.Reset()
	//wood2 img
	g.drawOps.GeoM.Translate(float64(g.wood2Sprite.xloc), float64(g.wood2Sprite.yloc))
	g.drawOps.GeoM.Rotate(0.2 * math.Pi * float64(g.angle) / maxAngle)
	screen.DrawImage(g.wood2Sprite.picts, &g.drawOps)
	g.drawOps.GeoM.Reset()

	if !g.frogCollides {
		g.drawOps.GeoM.Reset()
		g.drawOps.GeoM.Translate(float64(g.redcarSprites.xloc), float64(g.redcarSprites.yloc))
		screen.DrawImage(g.redcarSprites.picts, &g.drawOps)
	}

	if !g.frogCollides {
		g.drawOps.GeoM.Reset()
		g.drawOps.GeoM.Apply(400, 100)
		screen.DrawImage(g.redcarSprites.picts, &g.drawOps)
	}
	if !g.frogCollides {
		g.drawOps.GeoM.Reset()
		g.drawOps.GeoM.Translate(float64(g.yellowCarSprites.xloc), float64(g.yellowCarSprites.yloc))
		screen.DrawImage(g.yellowCarSprites.picts, &g.drawOps)
	}

}

func loadPNGImageFromEmbedded(name string) *ebiten.Image {
	pictNames, err := embeddedAssets.ReadDir("graphics")
	if err != nil {
		log.Fatal("failed to read embedded dir ", pictNames, " ", err)
	}
	embeddedFile, err := embeddedAssets.Open("graphics/" + name)
	if err != nil {
		log.Fatal("failed to load embedded image ", embeddedFile, err)
	}
	rawImage, err := png.Decode(embeddedFile)
	if err != nil {
		log.Fatal("failed to load embedded image ", name, err)
	}
	gameImage := ebiten.NewImageFromImage(rawImage)
	return gameImage
}

func (g Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) newGame() {

}

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Frogger")
	gameObject := Game{}
	gameObject.froggerSprite = Sprite{
		picts: loadPNGImageFromEmbedded("frogger.png"),
		xloc:  400,
		yloc:  700,
		dx:    0,
		dy:    0,
	}
	gameObject.redcarSprites = Sprite{
		picts: loadPNGImageFromEmbedded("redCar.png"),
		xloc:  700,
		yloc:  300,
		dx:    0,
		dy:    0,
	}
	gameObject.yellowCarSprites = Sprite{
		picts: loadPNGImageFromEmbedded("yellowCar.png"),
		xloc:  600,
		yloc:  150,
		dx:    0,
		dy:    0,
	}
	gameObject.turtleSprite = Sprite{
		picts: loadPNGImageFromEmbedded("turtle.png"),
		xloc:  500,
		yloc:  100,
		dx:    0,
		dy:    0,
	}
	gameObject.turtle2Sprite = Sprite{
		picts: loadPNGImageFromEmbedded("turtle2.png"),
		xloc:  300,
		yloc:  200,
		dx:    0,
		dy:    0,
	}
	gameObject.woodSprite = Sprite{
		picts: loadPNGImageFromEmbedded("wood.png"),
		xloc:  150,
		yloc:  100,
		dx:    0,
		dy:    0,
	}
	gameObject.wood2Sprite = Sprite{
		picts: loadPNGImageFromEmbedded("wood2.png"),
		xloc:  20,
		yloc:  30,
		dx:    0,
		dy:    0,
	}

	if err := ebiten.RunGame(&gameObject); err != nil {
		log.Fatal("Oh no! something terrible happened", err)
	}
}
