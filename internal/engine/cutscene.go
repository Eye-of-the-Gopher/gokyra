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
	scene1   *Scene1
}

func NewCutSceneManager(assets *formats.Assets) (*CutSceneManager, error) {
	csm := &CutSceneManager{scene: 1,
		assets: assets,
		track:  nil,
	}
	sm1, err := NewScene1(csm)
	if err != nil {
		return nil, err
	}
	csm.scene1 = sm1
	return csm, nil

}

func (c *CutSceneManager) Update(game *Game) (bool, error) {
	switch c.scene {
	case 1:
		return c.Scene1Update(game)
	default:
		EngineLogger.Warn("Scene not implemented yet", "scene", c.scene)
	}
	return false, nil
}

func (c *CutSceneManager) Draw(screen *ebiten.Image, game *Game) {
	switch c.scene {
	case 1:
		c.Scene1Draw(screen, game)
	default:
		EngineLogger.Warn("Scene not implemented yet", "scene", c.scene)
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
