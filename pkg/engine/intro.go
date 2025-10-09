package engine

import (
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
	stageIndex    int
	stages        []ImageStage
	fading        bool
	fadingCounter float32
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
			if game.currentTrack == nil { // But nothing is playing
				audioPlayer, err := stage.track.GetEbintenPlayer(game.audioContext) // Get ready to play
				if err == nil {
					game.currentTrack = audioPlayer
				}
			} else if game.currentTrack.IsPlaying() { // There's something already playing
				err := game.currentTrack.Close() // Stop it
				if err != nil {
					EngineLogger.Warn("Couldn't stop current track ", "reason", err)
				}
				audioPlayer, err := stage.track.GetEbintenPlayer(game.audioContext) // and create a new player
				if err == nil {
					game.currentTrack = audioPlayer
				}
			}
		} else { // The stage doesn't have a track
			if game.currentTrack != nil { // If there's something playing
				err := game.currentTrack.Close() // Stop it
				if err != nil {
					EngineLogger.Warn("Couldn't stop current track ", "reason", err)
				}
			}
		}
	} else {
		if time.Since(stage.startedAt) > stage.fadeStart {
			EngineLogger.Debug("Starting fade", "name", stage.name, "at", stage.fadeStart)
			i.fading = true
			i.fadingCounter = 0
		}

		if time.Since(stage.startedAt) > stage.displayTime {
			i.stageIndex++
			EngineLogger.Debug("Going to next page")
		}
	}
	return nil
}

func (i *IntroManager) Draw(screen *ebiten.Image, game *Game) {
	// EngineLogger.Debug("Game is Drawing Intro")
	stage := i.stages[i.stageIndex]
	img, _ := stage.image.GetEbitenImage()
	screen.DrawImage(img, nil)
	if game.currentTrack != nil && !game.currentTrack.IsPlaying() {
		game.currentTrack.Play()
	}
	if i.fading {
		nStage := i.stages[i.stageIndex+1]
		nImg, _ := nStage.image.GetEbitenImage()
		op := &ebiten.DrawImageOptions{}
		op.ColorScale.ScaleAlpha(float32(i.fadingCounter) / 10) // 0.0 to 1.0
		i.fadingCounter = +0.05
		screen.DrawImage(nImg, op)
		if i.fadingCounter > 10 {
			i.fading = false
		}
	}

}

func NewIntroManager(scenes []ImageStage) *IntroManager {
	ret := IntroManager{
		fadingCounter: 0,
		fading:        false,
		stages:        scenes,
		stageIndex:    0,
	}
	return &ret
}
