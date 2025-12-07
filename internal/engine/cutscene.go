package engine

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nibrahim/eye-of-the-gopher/internal/formats"
)

type PixelIterator func() (image.Point, bool)

// Framework here
type CutSceneManager struct {
	scene    int
	assets   *formats.Assets
	subtitle *ebiten.Image
	scene0   *Scene0
	scene1   *Scene1
}

func NewCutSceneManager(assets *formats.Assets) (*CutSceneManager, error) {
	csm := &CutSceneManager{scene: 0,
		assets: assets,
	}
	sm0, err := NewScene0(csm)
	sm1, err := NewScene1(csm)

	if err != nil {
		return nil, err
	}
	csm.scene0 = sm0
	csm.scene1 = sm1
	return csm, nil

}

func (c *CutSceneManager) Update(game *Game) (bool, error) {
	switch c.scene {
	case 0:
		next, _ := c.Scene0Update(game)
		if next {
			c.scene = 1
		}
	case 1:
		next, _ := c.Scene1Update(game)
		if next {
			c.scene = 1 // Go to next
		}
	default:
		EngineLogger.Warn("Scene not implemented yet", "scene", c.scene)
	}
	return false, nil
}

func (c *CutSceneManager) Draw(screen *ebiten.Image, game *Game) {
	switch c.scene {
	case 0:
		c.Scene0Draw(screen, game)
	case 1:
		c.Scene1Draw(screen, game)
	default:
		EngineLogger.Warn("Scene not implemented yet", "scene", c.scene)
	}

}

// Helpers

func fadeGridGen(x0, y0, edge int) PixelIterator {
	i := 0
	dir := 0
	pt := image.Point{X: x0, Y: y0}

	return func() (image.Point, bool) {
		if edge <= 1 {
			EngineLogger.Debug("We're done. Stopping now")
			return image.Point{}, false
		}
		// EngineLogger.Debug("PixelIterator ", "i", i, "dir", dir, "point", pt)
		if i >= edge {
			dir += 1
			i = 0
			// EngineLogger.Debug("Resetting i to 0", "dir", dir)
		}

		if dir == 4 {
			edge -= 1
			i = 0
			dir = 0
			pt = image.Point{X: x0, Y: y0}
			// EngineLogger.Debug("Reducing edge", "edge", edge)
		}

		switch dir {
		case 0:
			// EngineLogger.Debug("Incrementing X")
			pt.X += 1
		case 1:
			// EngineLogger.Debug("Incrementing Y")
			pt.Y += 1
		case 2:
			// EngineLogger.Debug("Decrementing X")
			pt.X -= 1
		case 3:
			// EngineLogger.Debug("Decrementing Y")
			pt.Y -= 1
		}
		i++
		return pt, true // (value, has_more)

	}

}
