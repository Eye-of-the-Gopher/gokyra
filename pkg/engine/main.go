package engine

import (
	"log/slog"
	"time"

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
	// eg := "TITLE-V.CMP"
	// pal := "WESTWOOD.COL"
	// p, err := assets.GetPalette(pal)
	// t, err := assets.GetSprite(eg, p, 320, 200, "")
	// if err != nil {
	// 	utils.ErrorAndExit("Couldn't load File %s: %v", eg, err)
	// } else {
	// 	fmt.Println(t)
	// }
	audioContext, err := audio.NewContext(44100)
	//assets.DumpAssets()
	eg := "TITLE-V.CMP"
	pal := "WESTWOOD.COL"
	p, err := assets.GetPalette(pal)
	t, err := assets.GetSprite(eg, p, 320, 200, "")
	scene1 := ImageStage{
		name:     "Title-v",
		image:    t,
		duration: time.Duration(1000),
	}

	introManager := NewIntroManager([]ImageStage{scene1})
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
