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
	ScreenWidth  = 320
	ScreenHeight = 200
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
	// Internal
	// Audio related
	audioContext     *audio.Context
	currentTrack     *audio.Player
	currentTrackName string

	// This allows us from move to one state (e.g. intro, cutscenes etc. to the next)
	state GameState

	// Managers are similar to Game structs. We delegate Update
	// and Draw to them depending on stage
	introManager    *IntroManager
	cutSceneManager *CutSceneManager

	// Assets used in the game
	assets formats.Assets
}

func NewGame(assetDir string, extraAssetDir string, enhanced bool) Game {
	EngineLogger.Debug("Creating game")
	assets := formats.LoadAssets(assetDir, extraAssetDir)
	audioContext := audio.NewContext(44100)

	introManager := NewIntroManager(assets, enhanced)
	cutsceneManager, err := NewCutSceneManager(assets)
	if err != nil {
		panic(err)
	}

	return Game{
		introManager:    introManager,
		cutSceneManager: cutsceneManager,
		state:           GameIntro,
		// state:        GameCutScene,
		assets:       *assets,
		audioContext: audioContext,
	}

}

func (g *Game) EnsureTrackPlaying(trackname string) {
	if trackname == "" {
		EngineLogger.Debug("No track to play here")
	} else if g.currentTrackName == trackname {
		EngineLogger.Debug("Track is already playing. Not changing", "name", trackname)
	} else {
		track, err := g.assets.GetAudioTrack(trackname)
		if err != nil {
			EngineLogger.Debug("Can't get track ", "name", trackname)
		} else if g.currentTrack != nil && g.currentTrack.IsPlaying() { // There's something already playing
			err := g.currentTrack.Close() // Stop it
			if err != nil {
				EngineLogger.Warn("Couldn't stop current track ", "reason", err, "name", g.currentTrackName)
			}
		}
		audioPlayer, err := track.GetEbintenPlayer(g.audioContext) // and create a new player
		if err == nil {
			g.currentTrack = audioPlayer
			g.currentTrack.Play()
		}
		g.currentTrackName = trackname
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
			g.state = GameCutScene // Change this when done
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
		g.cutSceneManager.Draw(screen, g)
	}
	// screen.DrawImage(g.image, nil)
	fps := ebiten.ActualFPS()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f", fps))
	// ebitenutil.DebugPrint(screen, "Eye of the Gopher\nHello, Dungeon!")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func InitLogger(engineLevel slog.Level) {
	EngineLogger = utils.InitLogger("engine", engineLevel)
}
