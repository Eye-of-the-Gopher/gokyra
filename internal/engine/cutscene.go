package engine

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/nibrahim/eye-of-the-gopher/internal/formats"
)

type PixelIterator func() (image.Point, bool)

// Framework here
type CutSceneManager struct {
	scene    int
	assets   *formats.Assets
	track    *audio.Player
	subtitle *ebiten.Image
	scene0   *Scene0
	scene1   *Scene1
}

func NewCutSceneManager(assets *formats.Assets) (*CutSceneManager, error) {
	csm := &CutSceneManager{scene: 0,
		assets: assets,
		track:  nil,
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
		// EngineLogger.Debug("Update: Cut scene 0")
		return c.Scene0Update(game)
	case 1:
		// EngineLogger.Debug("Update: Cut scene 1")
		return c.Scene1Update(game)
	default:
		EngineLogger.Warn("Scene not implemented yet", "scene", c.scene)
	}
	return false, nil
}

func (c *CutSceneManager) Draw(screen *ebiten.Image, game *Game) {
	switch c.scene {
	case 0:
		// EngineLogger.Debug("Draw: Cut scene 0")
		c.Scene0Draw(screen, game)
	case 1:
		// EngineLogger.Debug("Draw: Cut scene 1")
		c.Scene1Draw(screen, game)
	default:
		EngineLogger.Warn("Scene not implemented yet", "scene", c.scene)
	}

}

// actual scene 0 here. This is just a holding screen to fade out
type Scene0 struct {
	titleCard *ebiten.Image
	clearing  bool
	lineImg   *ebiten.Image
	fader     PixelIterator
}

func NewScene0(c *CutSceneManager) (*Scene0, error) {
	palette, err := c.assets.GetPalette("EOBPAL.COL")
	if err != nil {
		EngineLogger.Error("Couldn't load palette for title card ", "palette", "EOBPAL.COL")
		return nil, err
	}
	titleCard, err := c.assets.GetSprite("INTRO.CPS", palette, 320, 200, "")
	if err != nil {
		EngineLogger.Error("Couldn't load  title card sprite", "sprite", "INTRO.CPS")
		return nil, err
	}
	titleCardImage, err := titleCard.GetEbitenImage()
	if err != nil {
		EngineLogger.Error("Couldn't convert title card sprite into image", "image", "intro.cps")
		return nil, err
	}

	return &Scene0{
		titleCard: titleCardImage,
		clearing:  false,
	}, nil
}

func (c *CutSceneManager) Scene0Update(game *Game) (bool, error) {
	if c.scene0.lineImg == nil { // Create the fadeout line the first time
		EngineLogger.Debug("I'm initting the lineImg")
		c.scene0.lineImg = ebiten.NewImage(200, 200)
		c.scene0.clearing = true
		c.scene0.fader = fadeGridGen(0, 0, 20)
	}

	if c.scene0.clearing {
		if pt, hasMore := c.scene0.fader(); hasMore {
			// EngineLogger.Debug("Setting pixel", "x", pt.X, "y", pt.Y, "hasMore", hasMore)
			c.scene0.lineImg.Set(pt.X, pt.Y, color.Black)
		} else {
			EngineLogger.Debug("We're done drawing the fading square")
			c.scene0.clearing = false
		}
	}

	return false, nil
}

func (c *CutSceneManager) Scene0Draw(screen *ebiten.Image, game *Game) {
	screen.DrawImage(c.scene0.titleCard, nil)
	if c.scene0.clearing {
		for x := 0; x < 1000; x += 20 {
			for y := 0; y < 500; y += 20 {
				fop := &ebiten.DrawImageOptions{}
				fop.GeoM.Translate(float64(x), float64(y))
				screen.DrawImage(c.scene0.lineImg, fop)
			}
		}
	}
}

// Actual scene 1 here
type Scene1 struct {
	towerSprite *formats.Sprite
	textSprite  *formats.Sprite
	text1       image.Image
}

func NewScene1(c *CutSceneManager) (*Scene1, error) {
	textPalette, err := c.assets.GetPalette("TOWRMAGE.COL")
	if err != nil {
		EngineLogger.Error("Couldn't load palette for cutscene ", "palette", "TOWRMAGE.COL")
		return nil, err
	}
	textSprite, err := c.assets.GetSprite("TEXT.CMP", textPalette, 320, 200, "")
	if err != nil {
		EngineLogger.Error("Couldn't load Text sprite", "sprite", "TEXT.CMP")
		return nil, err
	}

	srcRect := image.Rect(0, 0, 1280, 131)
	textSpriteImage, err := textSprite.GetEbitenImage()
	if err != nil {
		return nil, err
	}

	return &Scene1{
		towerSprite: nil,
		textSprite:  textSprite,
		text1:       textSpriteImage.SubImage(srcRect).(*ebiten.Image),
	}, nil

}

func (c *CutSceneManager) Scene1Update(game *Game) (bool, error) {
	if c.subtitle == nil {
		c.subtitle = c.scene1.text1.(*ebiten.Image)
	}

	return false, nil
}

func (c *CutSceneManager) Scene1Draw(screen *ebiten.Image, game *Game) {
	// if c.scene1.clearing {
	// 	for x := 0; x < 1000; x += 20 {
	// 		for y := 0; y < 500; y += 20 {
	// 			fop := &ebiten.DrawImageOptions{}
	// 			fop.GeoM.Translate(float64(x), float64(y))
	// 			screen.DrawImage(c.scene1.lineImg, fop)
	// 		}
	// 	}
	// } else {
	// 	op := &ebiten.DrawImageOptions{}
	// 	op.GeoM.Translate(0, 669)
	// 	if c.subtitle != nil {
	// 		screen.DrawImage(c.subtitle, op)
	// 	}

	// }
}

// Helpers

func fadeGridGen(x0, y0, edge int) PixelIterator {
	i := 0
	dir := 0
	pt := image.Point{X: x0, Y: y0}

	return func() (image.Point, bool) {
		// EngineLogger.Debug("PixelIterator ", "i", i, "dir", dir, "point", pt)
		if i >= edge {
			dir += 1
			i = 0
			EngineLogger.Debug("Resetting i to 0", "dir", dir)
		}

		if dir == 5 {
			EngineLogger.Debug("Returning here", "dir", dir)
			return image.Point{}, false
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

// func fadeGridGen(x0, y0, edge int) PixelIterator {
// 	i := 0
// 	dir := 0 // 0R, 1D, 2L, 3U
// 	cEdge := edge

// 	lx := x0
// 	ly := y0
// 	return func() (image.Point, bool) {
// 		if edge <= 1 {
// 			return image.Point{}, false // done
// 		}
// 		if i == cEdge {
// 			i = 0
// 			if dir == 3 {
// 				cEdge -= 1
// 			}
// 			dir += 1
// 			dir %= 4
// 		}

// 		pt := image.Point{}
// 		switch dir {
// 		case 0:
// 			pt.X = lx + i
// 			// pt.Y = y0 + i
// 		case 1:
// 			pt.Y = ly + i
// 			// pt.Y = y0 + i
// 		case 2:
// 			pt.X = lx - i
// 			// pt.Y = y0 + i
// 		case 3:
// 			pt.Y = ly - i
// 			// pt.Y = y0 + i
// 		}
// 		lx = pt.X
// 		ly = pt.Y
// 		i += 1

// 		return pt, true // (value, has_more)
// 	}
// }
