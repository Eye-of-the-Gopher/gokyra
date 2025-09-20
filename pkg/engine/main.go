package engine

import (
	"log/slog"

	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/nibrahim/eye-of-the-gopher/internal/utils"
	"github.com/nibrahim/eye-of-the-gopher/pkg/formats"
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
	GamePlaying
	GamePaused
)

type Game struct {
	introManager *IntroManager
	state        GameState

	assets        formats.Assets
	audioContext  *audio.Context
	currentPlayer *audio.Player
}

func NewGame(assetDir string, extraAssetDir string) Game {
	EngineLogger.Debug("Creating game")
	assets := formats.LoadAssets(assetDir, extraAssetDir)
	audioContext, err := audio.NewContext(44100)

	assets.DumpAssets()

	type SceneConfig struct {
		name, asset, palette string
		duration1, duration2 int
	}

	configs := []SceneConfig{
		{"westwood", "ENHANCED/WESTWOOD.PNG", "WESTWOOD.COL", 4, 3},
		{"westwood And", "AND.CMP", "WESTWOOD.COL", 3, 2},
		{"ssi", "SSI.CMP", "WESTWOOD.COL", 4, 3},
		{"present", "PRESENT.CMP", "WESTWOOD.COL", 3, 2},
		{"dand", "DAND.CMP", "WESTWOOD.COL", 3, 2},
		{"dand", "ENHANCED/WESTWOOD.PNG", "WESTWOOD.COL", 3, 2},
	}
	var scenes []ImageStage
	for _, c := range configs {
		scene, err := NewImageStage(assets, c.name, c.asset, c.palette, c.duration1, c.duration2)
		if err != nil {
			EngineLogger.Error("Couldn't load asset ", "asset", c.asset, "error", err)
			panic("Asset loading failed")
		}
		scenes = append(scenes, *scene)
	}

	introManager := NewIntroManager(scenes)
	if err != nil {
		return Game{
			introManager: introManager,
			state:        GameIntro,
			assets:       *assets,
			audioContext: audioContext,
		}

	} else {
		return Game{
			introManager: introManager,
			state:        GameIntro,
			assets:       *assets,
			audioContext: nil,
		}

	}
}

func (g *Game) Update() error {
	switch g.state {
	case GameIntro:
		g.introManager.Update()

	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case GameIntro:
		g.introManager.Draw(screen)

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
