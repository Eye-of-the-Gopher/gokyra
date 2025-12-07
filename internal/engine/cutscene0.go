package engine

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// actual scene 0 here. This is just a holding screen to fade out
type Scene0 struct {
	titleCard *ebiten.Image
	clearing  bool
	lineImg   *ebiten.Image
	fader     PixelIterator
	done      bool
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
	game.EnsureTrackPlaying("ENHANCED/CUTSCENE.WAV")
	if c.scene0.lineImg == nil { // Create the fadeout line the first time
		EngineLogger.Debug("I'm initting the lineImg")
		c.scene0.lineImg = ebiten.NewImage(200, 200)
		c.scene0.clearing = true
		c.scene0.fader = fadeGridGen(0, 0, 20)
	}
	if c.scene0.clearing {
		for i := 0; i < 15; i++ {
			if pt, hasMore := c.scene0.fader(); hasMore {
				c.scene0.lineImg.Set(pt.X, pt.Y, color.Black)
			} else {
				EngineLogger.Debug("We're done drawing the fading square")
				c.scene0.done = true
				c.scene0.clearing = false
				return true, nil
			}
		}
	}

	return false, nil
}

func (c *CutSceneManager) Scene0Draw(screen *ebiten.Image, game *Game) {
	if c.scene0.done == true {
		screen.Fill(color.Black)
	} else {
		screen.DrawImage(c.scene0.titleCard, nil)
	}

	if c.scene0.clearing {
		bounds := screen.Bounds()
		width := bounds.Dx()  // Delta X (max.X - min.X)
		height := bounds.Dy() // Delta Y (max.Y - min.Y)

		for x := 0; x < width; x += 20 {
			for y := 0; y < height; y += 20 {
				fop := &ebiten.DrawImageOptions{}
				fop.GeoM.Translate(float64(x), float64(y))
				screen.DrawImage(c.scene0.lineImg, fop)
			}
		}
	}
}
