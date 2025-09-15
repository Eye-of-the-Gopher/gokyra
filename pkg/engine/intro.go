package engine

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nibrahim/eye-of-the-gopher/pkg/formats"
)

type ImageStage struct { // This will later become an interface
	name     string
	image    *formats.Sprite
	duration time.Duration
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
}

func (i *IntroManager) Update() error {
	EngineLogger.Debug("Game is playing Intro")
	return nil
}

func (i *IntroManager) Draw(screen *ebiten.Image) {
	EngineLogger.Debug("Game is Drawing Intro")
	stage := i.stages[i.stageIndex]
	img, _ := stage.image.GetEbitenImage()
	screen.DrawImage(img, nil)

}

func NewIntroManager(scenes []ImageStage) *IntroManager {
	ret := IntroManager{
		stages:     scenes,
		stageIndex: 0,
	}
	return &ret
}
