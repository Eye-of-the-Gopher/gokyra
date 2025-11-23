package engine

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type CutSceneManager struct {
	stage int
}

func (c *CutSceneManager) Update(game *Game) (bool, error) {
	fmt.Println("Well, I've reached here!")
	return false, nil
}

func (c *CutSceneManager) Draw(screen *ebiten.Image, game *Game) {

}

func NewCutSceneManager() *CutSceneManager {
	return &CutSceneManager{stage: 0}
}
