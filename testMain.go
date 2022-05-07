package main

import (
	"log"
	"math"
	"testing"
)

//https://github.com/jsantore/FirstGameDemo/blob/master/game_test.go
//tried to do waht you did but could not get it to work on my end

func testCarSquishFrog(t *testing.T, name string) {
	pictNames, err := embeddedAssets.ReadDir("graphics")
	if err != nil {
		log.Fatal("failed to read embedded dir ", pictNames, " ", err)
	}
	embeddedFile, err := embeddedAssets.Open("graphics/" + name)
	if err != nil {
		log.Fatal("failed to load embedded image ", embeddedFile, err)
	}
	tests := []struct {
		frogSprite Sprite
		carSprite  Sprite
		crash      bool
	}{
		{Sprite{
			picts: froggerPict,
			xloc:  700,
			yloc:  400,
		},
			Sprite{picts: froggerPict,
				xloc: 150,
				yloc: 250,
				true},
		},
		t.Error("pass frog got squished"),
	}
}
func testNewGame(t *testing.T)bool{
	car Sprite{
		carWidth: 40,
		carHieght: 40
	}
	frog Sprite{
		frogWidth: 40,
		frogHeight: 50,
	}
	if carYellowSquishesFrog(frog, car Sprite)bool{
		if frog.xloc < car.xloc+carWidth &&
		frog.xloc+frogWidth > car.xloc &&
		frog.yloc < car.yloc+carHieght &&
		frog.yloc+frogHeight > car.yloc {
		return true
	}
		return false
	}
	t.Error("failed frog didnt collide")
}

//was trying to test if cars went across the screen at  90 degrees instead of curving
func testRotation(t *testing.T){
	angleHitsFrog(maxAngle)
	angle := (0.5 * math.Pi * float64(g.angle) / maxAngle)
	t.Error("car did not rotate properly")

}