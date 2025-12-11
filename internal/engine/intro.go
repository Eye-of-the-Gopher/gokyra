package engine

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nibrahim/eye-of-the-gopher/internal/formats"
)

type ImageStage struct { // This will later become an interface
	startedAt   time.Time
	running     bool
	name        string
	image       *formats.Sprite
	fadeStart   time.Duration
	displayTime time.Duration
	// track        *formats.AudioTrack
	track        string
	trackStarted bool
}

func NewImageStage(assets *formats.Assets, name string, assetName string, paletteName string, trackName string, displayDuration int, fadeDuration int) (*ImageStage, error) {
	image, err := assets.GetSprite(assetName, paletteName, 320, 200, "")
	if err != nil {
		return nil, err
	}
	ret := ImageStage{
		name:         name,
		image:        image,
		displayTime:  time.Duration(displayDuration) * time.Second,
		fadeStart:    time.Duration(fadeDuration) * time.Second,
		trackStarted: false,
		track:        trackName,
	}
	return &ret, nil
}

type IntroManager struct {
	scale      float64
	stageIndex int
	stages     []ImageStage
	fading     bool
	fadeStart  time.Time
	fadeAlpha  float64
}

func (i *IntroManager) Update(game *Game) (bool, error) {
	// EngineLogger.Debug("Game is playing Intro")

	// If there are 6 stages, we move to next the next Game state
	// after the 5th stage. Playing the 5th stage will crossfade
	// into the 6th. If we check only for 6th, the crossfade will
	// attempt to play the 7th and crash. This probably needs to
	// be done better. TBD
	if i.stageIndex >= len(i.stages) {
		EngineLogger.Debug("Ending intro")
		return true, nil
	}
	stage := &i.stages[i.stageIndex]
	// EngineLogger.Debug("Time now", "startedAt", stage.startedAt, "time since", time.Since(stage.startedAt), "duration", stage.displayTime, "running", stage.running)
	if !stage.running {
		stage.running = true
		stage.startedAt = time.Now()
		EngineLogger.Debug("Starting stage", "name", stage.name, "at", stage.startedAt)
		game.EnsureTrackPlaying(stage.track)
	} else {
		if (time.Since(stage.startedAt) > stage.fadeStart) && !i.fading { // Otherwise, start cross fade if appropriate
			EngineLogger.Debug("Starting fade", "name", stage.name, "at", stage.fadeStart)
			i.fading = true
			i.fadeStart = time.Now()
			i.fadeAlpha = 0
		}
		if i.fading { //If fading is running, adjust screen alpha
			i.fadeAlpha = time.Since(i.fadeStart).Seconds() / (float64(stage.displayTime.Seconds()) - float64(stage.fadeStart.Seconds()))
		}

		if time.Since(stage.startedAt) > stage.displayTime { //If time is up for this stage, then we need to decide what to do
			if i.stageIndex == len(i.stages)-1 { // If it's the last stage, then tell the framework to move on
				EngineLogger.Debug("On last page")
				return true, nil
			} else { // Else, go to the next scene
				i.stageIndex++
				i.fading = false
				EngineLogger.Debug("Going to next page")
			}
		}
	}
	return false, nil
}

func (i *IntroManager) Draw(screen *ebiten.Image, game *Game) {
	stage := i.stages[i.stageIndex]
	img, _ := stage.image.GetEbitenImage()
	screen.DrawImage(img, nil)
	if i.fading && i.stageIndex+1 < len(i.stages) {
		nStage := i.stages[i.stageIndex+1]
		nImg, _ := nStage.image.GetEbitenImage()
		op := &ebiten.DrawImageOptions{}
		op.ColorScale.ScaleAlpha(float32(i.fadeAlpha)) // 0.0 to 1.0
		screen.DrawImage(nImg, op)
	}
}

func NewIntroManager(assets *formats.Assets, enhanced bool) *IntroManager {
	type SceneConfig struct {
		name, asset, palette, trackname string
		duration1, duration2            int
	}

	var configs []SceneConfig

	if enhanced {
		EngineLogger.Debug("Using enhanced assets")
		configs = []SceneConfig{
			{"westwood", "ENHANCED/WESTWOOD.PNG", "WESTWOOD.COL", "", 4, 3},
			{"westwood And", "ENHANCED/AND.PNG", "WESTWOOD.COL", "", 3, 2},
			{"ssi", "ENHANCED/SSI.PNG", "WESTWOOD.COL", "", 4, 3},
			{"present", "ENHANCED/PRESENT.PNG", "WESTWOOD.COL", "", 3, 2},
			{"dand", "ENHANCED/DAND.PNG", "WESTWOOD.COL", "", 3, 2},
			{"dand", "ENHANCED/WESTWOOD.PNG", "WESTWOOD.COL", "", 3, 2},
		}
	} else {
		EngineLogger.Debug("Using classic assets")
		configs = []SceneConfig{
			{"westwood", "WESTWOOD.CMP", "WESTWOOD.COL", "ENHANCED/INTRO.WAV", 8, 3},
			{"westwood And", "AND.CMP", "WESTWOOD.COL", "", 3, 2},
			{"ssi", "SSI.CMP", "WESTWOOD.COL", "", 5, 3},
			{"present", "PRESENT.CMP", "WESTWOOD.COL", "", 3, 2},
			{"dand", "DAND.CMP", "WESTWOOD.COL", "", 7, 2},
			{"intro", "INTRO.CPS", "EOBPAL.COL", "", 2, 0}, //ENHANCED/CUTSCENE.WAV
		}
	}
	var scenes []ImageStage
	for _, c := range configs {
		scene, err := NewImageStage(assets, c.name, c.asset, c.palette, c.trackname, c.duration1, c.duration2)
		if err != nil {
			EngineLogger.Error("Couldn't load asset ", "asset", c.asset, "error", err)
			panic("Asset loading failed")
		}
		scenes = append(scenes, *scene)
	}
	ret := IntroManager{
		fadeAlpha:  0,
		fadeStart:  time.Time{},
		fading:     false,
		stages:     scenes,
		stageIndex: 0,
	}
	return &ret
}
