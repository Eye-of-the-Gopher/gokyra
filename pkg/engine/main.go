package engine

import (
	"fmt"
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

type Game struct {
	assets        formats.Assets
	image         *ebiten.Image
	audioContext  *audio.Context // ONE for entire game
	currentPlayer *audio.Player  // Changes per track
}

func NewGame(assetDir string, extraAssetDir string) Game {
	EngineLogger.Debug("Creating game")
	assets := formats.LoadAssets(assetDir, extraAssetDir)
	eg := "TITLE-V.CMP"
	pal := "WESTWOOD.COL"
	p, err := assets.GetPalette(pal)
	t, err := assets.GetSprite(eg, p, 320, 200, "")
	if err != nil {
		utils.ErrorAndExit("Couldn't load File %s: %v", eg, err)
	} else {
		fmt.Println(t)
	}
	audioContext, err := audio.NewContext(44100)
	//assets.DumpAssets()
	if err != nil {
		return Game{
			assets:       *assets,
			image:        ebiten.NewImageFromImage(t.Image),
			audioContext: audioContext,
		}

	} else {
		return Game{
			assets:       *assets,
			image:        ebiten.NewImageFromImage(t.Image),
			audioContext: nil,
		}

	}
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.image, nil)
	ebitenutil.DebugPrint(screen, "Eye of the Gopher\nHello, Dungeon!")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func InitLogger(engineLevel slog.Level) {
	EngineLogger = utils.InitLogger("engine", engineLevel)

}
