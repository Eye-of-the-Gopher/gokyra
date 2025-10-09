package engine

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nibrahim/eye-of-the-gopher/pkg/formats"
)

type ImageStage struct { // This will later become an interface
	startedAt    time.Time
	running      bool
	name         string
	image        *formats.Sprite
	fadeStart    time.Duration
	displayTime  time.Duration
	track        *formats.AudioTrack
	trackStarted bool
}

func NewImageStage(assets *formats.Assets, name string, assetName string, paletteName string, trackName string, displayDuration int, fadeDuration int) (*ImageStage, error) {
	p, err := assets.GetPalette(paletteName)
	if err != nil {
		return nil, err
	}

	image, err := assets.GetSprite(assetName, p, 320, 200, "")
	if err != nil {
		return nil, err
	}
	var track *formats.AudioTrack
	if trackName != "" {
		track, err = assets.GetAudioTrack(trackName)
		if err != nil {
			return nil, err
		}
	}

	ret := ImageStage{
		name:         name,
		image:        image,
		displayTime:  time.Duration(displayDuration) * time.Second,
		fadeStart:    time.Duration(fadeDuration) * time.Second,
		trackStarted: false,
		track:        track,
	}
	return &ret, nil
}

func (g *Game) PlayScene(scenes []ImageStage) {
	EngineLogger.Debug("Playing scene")
	for _, i := range scenes {
		EngineLogger.Debug(i.name)
	}
}

type IntroManager struct {
	stageIndex int
	stages     []ImageStage
	fading     bool
	fadeStart  time.Time
	fadeAlpha  float64
}

func (i *IntroManager) Update(game *Game) error {
	// EngineLogger.Debug("Game is playing Intro")
	stage := &i.stages[i.stageIndex]
	// EngineLogger.Debug("Time now", "startedAt", stage.startedAt, "time since", time.Since(stage.startedAt), "duration", stage.displayTime, "running", stage.running)
	if !stage.running {
		stage.running = true
		stage.startedAt = time.Now()
		EngineLogger.Debug("Starting stage", "name", stage.name, "at", stage.startedAt)
		if stage.track != nil { // Stage has a track
			EngineLogger.Debug("Staget track is not nil")
			if game.currentTrack == nil { // But nothing is playing
				audioPlayer, err := stage.track.GetEbintenPlayer(game.audioContext) // Get ready to play
				if err == nil {
					game.currentTrack = audioPlayer
					game.currentTrack.Play()
				}
			} else if game.currentTrack.IsPlaying() { // There's something already playing
				err := game.currentTrack.Close() // Stop it
				if err != nil {
					EngineLogger.Warn("Couldn't stop current track ", "reason", err)
				}
				audioPlayer, err := stage.track.GetEbintenPlayer(game.audioContext) // and create a new player
				if err == nil {
					game.currentTrack = audioPlayer
					game.currentTrack.Play()
				}
			}
		} // else { // The stage doesn't have a track
		// 	EngineLogger.Debug("Staget track is nil")
		// 	if game.currentTrack != nil { // If there's something playing
		// 		err := game.currentTrack.Close() // Stop it
		// 		if err != nil {
		// 			EngineLogger.Warn("Couldn't stop current track ", "reason", err)
		// 		}
		// 	}
		// }
	} else {
		if (time.Since(stage.startedAt) > stage.fadeStart) && i.fading == false {
			EngineLogger.Debug("Starting fade", "name", stage.name, "at", stage.fadeStart)
			i.fading = true
			i.fadeStart = time.Now()
			i.fadeAlpha = 0
		}
		if i.fading {
			i.fadeAlpha = time.Since(i.fadeStart).Seconds() / (float64(stage.displayTime.Seconds()) - float64(stage.fadeStart.Seconds()))
		}

		if time.Since(stage.startedAt) > stage.displayTime {
			i.stageIndex++
			i.fading = false
			EngineLogger.Debug("Going to next page")
		}
	}
	return nil
}

func (i *IntroManager) Draw(screen *ebiten.Image, game *Game) {
	stage := i.stages[i.stageIndex]
	img, _ := stage.image.GetEbitenImage()
	screen.DrawImage(img, nil)
	if i.fading {
		nStage := i.stages[i.stageIndex+1]
		nImg, _ := nStage.image.GetEbitenImage()
		op := &ebiten.DrawImageOptions{}
		fmt.Println(i.fadeAlpha)
		op.ColorScale.ScaleAlpha(float32(i.fadeAlpha)) // 0.0 to 1.0
		screen.DrawImage(nImg, op)
	}

}

func NewIntroManager(scenes []ImageStage) *IntroManager {
	ret := IntroManager{
		fadeAlpha:  0,
		fadeStart:  time.Time{},
		fading:     false,
		stages:     scenes,
		stageIndex: 0,
	}
	return &ret
}
