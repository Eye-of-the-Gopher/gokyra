package engine

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/nibrahim/eye-of-the-gopher/internal/formats"
)

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
		EngineLogger.Debug("Update: Cut scene 0")
		return c.Scene0Update(game)
	case 1:
		EngineLogger.Debug("Update: Cut scene 1")
		return c.Scene1Update(game)
	default:
		EngineLogger.Warn("Scene not implemented yet", "scene", c.scene)
	}
	return false, nil
}

func (c *CutSceneManager) Draw(screen *ebiten.Image, game *Game) {
	switch c.scene {
	case 0:
		EngineLogger.Debug("Draw: Cut scene 0")
		c.Scene0Draw(screen, game)
	case 1:
		EngineLogger.Debug("Draw: Cut scene 1")
		c.Scene1Draw(screen, game)
	default:
		EngineLogger.Warn("Scene not implemented yet", "scene", c.scene)
	}

}

// actual scene 0 here. This is just a holding screen to fade out
type Scene0 struct {
	titleCard *ebiten.Image
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
	}, nil
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

func (c *CutSceneManager) Scene0Update(game *Game) (bool, error) {
	return false, nil
}

func (c *CutSceneManager) Scene0Draw(screen *ebiten.Image, game *Game) {
	screen.DrawImage(c.scene0.titleCard, nil)
}

func (c *CutSceneManager) Scene1Update(game *Game) (bool, error) {
	// start the audio track for the cutscene. This is for the whole scene
	// if c.track == nil { // No music. Start playing something here
	// 	track, err := c.assets.GetAudioTrack("ENHANCED/CUTSCENE.WAV")
	// 	if err == nil {
	// 		audioPlayer, err := track.GetEbintenPlayer(game.audioContext) // Get ready to play
	// 		if err == nil {
	// 			c.track = audioPlayer
	// 			c.track.Play()
	// 		} else {
	// 			EngineLogger.Debug("Couldn't get player for track - ENHANCED/CUTSCENE.WAV")
	// 		}
	// 	} else {
	// 		EngineLogger.Warn("Couldn't load track - ENHANCED/CUTSCENE.WAV")
	// 	}
	// }
	// Now paint the text
	if c.subtitle == nil {
		c.subtitle = c.scene1.text1.(*ebiten.Image)
	}

	return false, nil
}

func (c *CutSceneManager) Scene1Draw(screen *ebiten.Image, game *Game) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 669)
	if c.subtitle != nil {
		screen.DrawImage(c.subtitle, op)
	}

}
