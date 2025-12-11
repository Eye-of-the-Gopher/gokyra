package engine

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Scene 2 is the orb
type Scene2 struct {
	orb   *ebiten.Image
	mages *ebiten.Image
}

func NewScene2(c *CutSceneManager) (*Scene2, error) {
	return &Scene2{}, nil
}

func (c *CutSceneManager) Scene2Update(game *Game) (bool, error) {
	EngineLogger.Debug("Running Scene 2")
	c.subtitle = c.subtitles[1]

	return false, nil
}

func (c *CutSceneManager) Scene2Draw(screen *ebiten.Image, game *Game) {

}
