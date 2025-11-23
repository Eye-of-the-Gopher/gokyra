package engine

import (
	"fmt"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/nibrahim/eye-of-the-gopher/internal/formats"
	"github.com/nibrahim/eye-of-the-gopher/internal/utils"
)

var (
	EngineLogger *slog.Logger
)

const (
	ScreenWidth  = 320 * 4
	ScreenHeight = 200 * 4
)

type GameState int

const (
	GameIntro GameState = iota
	GameCutScene
	GameMenu
	GamePlaying
	GamePaused
)

type Game struct {
	introManager    *IntroManager
	cutSceneManager *CutSceneManager
	state           GameState

	assets       formats.Assets
	audioContext *audio.Context
	currentTrack *audio.Player
}

func NewGame(assetDir string, extraAssetDir string, enhanced bool) Game {
	EngineLogger.Debug("Creating game")
	assets := formats.LoadAssets(assetDir, extraAssetDir)
	audioContext := audio.NewContext(44100)

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
			{"present", "PRESENT.CMP", "WESTWOOD.COL", "", 2, 1},
			{"dand", "DAND.CMP", "WESTWOOD.COL", "", 8, 2},
			{"intro", "INTRO.CPS", "EOBPAL.COL", "", 8, 2},
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

	introManager := NewIntroManager(scenes)
	cutsceneManager := NewCutSceneManager()
	return Game{
		introManager:    introManager,
		cutSceneManager: cutsceneManager,
		state:           GameIntro,
		assets:          *assets,
		audioContext:    audioContext,
	}

}

func (g *Game) Update() error {
	switch g.state {
	case GameIntro:
		next, _ := g.introManager.Update(g)
		if next {
			g.state = GameCutScene
		}
	case GameCutScene:
		next, _ := g.cutSceneManager.Update(g)
		if next {
			g.state = GameCutScene
		}
	case GameMenu:
		fmt.Println("Menu : Not implemented")

	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case GameIntro:
		g.introManager.Draw(screen, g)
	case GameCutScene:

	}
	// screen.DrawImage(g.image, nil)
	ebitenutil.DebugPrint(screen, "Eye of the Gopher\nHello, Dungeon!")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func InitLogger(engineLevel slog.Level) {
	EngineLogger = utils.InitLogger("engine", engineLevel)

}
