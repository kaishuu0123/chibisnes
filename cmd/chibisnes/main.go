package main

import (
	"encoding/binary"
	"flag"
	"image"
	"log"
	"math"
	"os"
	"unsafe"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/inkyblackness/imgui-go/v4"
	"github.com/kaishuu0123/chibisnes/chibisnes"
	"github.com/kaishuu0123/chibisnes/internal/gui"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	SCALE               int     = 2
	WINDOW_WIDTH        int     = 256 * SCALE
	WINDOW_HEIGHT       int     = 240 * SCALE
	AUDIO_MASTER_VOLUME float64 = 0.5 // min: 0.0 max: 1.0
)

var (
	windowFlags imgui.WindowFlags = imgui.WindowFlagsNoCollapse |
		imgui.WindowFlagsNoMove |
		imgui.WindowFlagsNoResize |
		imgui.WindowFlagsHorizontalScrollbar
	isRunning = false

	audioDevice sdl.AudioDeviceID
	audioBuffer [735 * 4]int16 // *2 for stereo, *2 for sizeof(int16)
)

// For pprof
// func init() {
// 	go func() {
// 		log.Println(http.ListenAndServe("localhost:6060", nil))
// 	}()
// }

func main() {
	flag.Parse()
	if len(flag.Args()) >= 1 {
		_, err := os.Stat(flag.Arg(0))
		if err != nil {
			log.Fatalln("no SNES ROM file specified or found")
		}
	}

	window := gui.NewMasterWindow("ChibiSNES", WINDOW_WIDTH, WINDOW_HEIGHT, -1)
	screenImage := image.NewRGBA(image.Rect(0, 0, WINDOW_WIDTH, WINDOW_HEIGHT))

	console := chibisnes.NewConsole()

	romFilePath := flag.Arg(0)
	data, err := readFile(romFilePath)
	if err != nil {
		log.Fatalf("readFile error: %s\n", err)
	}
	if err := console.LoadROM(romFilePath, data, len(data)); err != nil {
		log.Fatalf("%s\n", err)
	}

	initAudio()

	var texture imgui.TextureID
	for !window.Platform.ShouldStop() {
		window.Platform.ProcessEvents()

		if window.Platform.Window.GetKey(glfw.KeyL) == glfw.Press {
			console.Debug = !console.Debug
		}
		processInputController1(window.Platform.Window, console)
		// want to more keys
		// processInputController2(window.Platform.Window, console)

		console.RunFrame()

		// clear screen
		for i := 0; i < len(screenImage.Pix); i++ {
			screenImage.Pix[i] = 0
		}

		console.SetPixels(screenImage.Pix)

		texture, _ = window.Renderer.CreateImageTexture(screenImage)
		renderGUI(window, &texture)
		window.Renderer.ReleaseImage(texture)

		PlayAudio(console)
	}

	console.Close()
}

func readFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	data := make([]byte, stat.Size())
	if err := binary.Read(file, binary.LittleEndian, data); err != nil {
		return nil, err
	}

	return data, nil
}

func renderGUI(w *gui.MasterWindow, texture *imgui.TextureID) {
	w.Platform.NewFrame()
	imgui.NewFrame()

	imgui.BackgroundDrawList().
		AddImage(
			*texture,
			imgui.Vec2{X: 0, Y: 0},
			imgui.Vec2{X: float32(WINDOW_WIDTH), Y: float32(WINDOW_HEIGHT)},
		)

	imgui.Render()

	w.Renderer.PreRender(w.ClearColor)
	w.Renderer.Render(w.Platform.DisplaySize(), w.Platform.FramebufferSize(), imgui.RenderedDrawData())
	w.Platform.PostRender()
}

func initAudio() {
	err := sdl.InitSubSystem(sdl.INIT_AUDIO)
	if err != nil {
		log.Fatalf("Failed to init SDL: %s\n", err)
	}

	var want, have sdl.AudioSpec
	want.Freq = 44100
	want.Format = sdl.AUDIO_S16
	want.Channels = 2
	want.Samples = 2048
	want.Callback = nil // use queue
	audioDevice, err = sdl.OpenAudioDevice("", false, &want, &have, 0)
	if err != nil {
		log.Fatalf("SDL OpenAudioDevice error: %s\n", err)
	}
	sdl.PauseAudioDevice(audioDevice, false)
}

func PlayAudio(console *chibisnes.Console) {
	console.SetAudioSamples(audioBuffer[:], 735)
	if sdl.GetQueuedAudioSize(audioDevice) <= uint32(len(audioBuffer)*6) {
		src := (*[len(audioBuffer) * 2]uint8)(unsafe.Pointer(&audioBuffer[0]))
		dst := make([]uint8, len(audioBuffer)*2)
		volume := int(math.Floor(float64(sdl.MIX_MAXVOLUME) * AUDIO_MASTER_VOLUME))

		sdl.MixAudioFormat(&dst[0], &src[0], sdl.AUDIO_S16, uint32(len(audioBuffer)*2), volume)

		// don't queue audio if buffer is still filled
		sdl.QueueAudio(audioDevice, dst[:len(audioBuffer)])
	}
}

func processInputController1(window *glfw.Window, console *chibisnes.Console) {
	var result [12]bool
	result[chibisnes.ButtonB] = window.GetKey(glfw.KeyX) == glfw.Press
	result[chibisnes.ButtonY] = window.GetKey(glfw.KeyC) == glfw.Press
	result[chibisnes.ButtonSelect] = window.GetKey(glfw.KeyRightShift) == glfw.Press
	result[chibisnes.ButtonStart] = window.GetKey(glfw.KeyEnter) == glfw.Press
	result[chibisnes.ButtonUp] = window.GetKey(glfw.KeyUp) == glfw.Press
	result[chibisnes.ButtonDown] = window.GetKey(glfw.KeyDown) == glfw.Press
	result[chibisnes.ButtonLeft] = window.GetKey(glfw.KeyLeft) == glfw.Press
	result[chibisnes.ButtonRight] = window.GetKey(glfw.KeyRight) == glfw.Press
	result[chibisnes.ButtonA] = window.GetKey(glfw.KeyZ) == glfw.Press
	result[chibisnes.ButtonX] = window.GetKey(glfw.KeyV) == glfw.Press
	result[chibisnes.ButtonL] = window.GetKey(glfw.KeyA) == glfw.Press
	result[chibisnes.ButtonR] = window.GetKey(glfw.KeyF) == glfw.Press

	for i := 0; i < len(result); i++ {
		console.SetButtonState(1, i, result[i])
	}
}

// func processInputController2(window *glfw.Window) [8]bool {
// 	var result [8]bool
// 	result[chibines.ButtonA] = window.GetKey(glfw.KeyA) == glfw.Press
// 	result[chibines.ButtonB] = window.GetKey(glfw.KeyS) == glfw.Press
// 	result[chibines.ButtonSelect] = window.GetKey(glfw.KeyLeftShift) == glfw.Press
// 	result[chibines.ButtonStart] = window.GetKey(glfw.KeyE) == glfw.Press
// 	result[chibines.ButtonUp] = window.GetKey(glfw.KeyI) == glfw.Press
// 	result[chibines.ButtonDown] = window.GetKey(glfw.KeyK) == glfw.Press
// 	result[chibines.ButtonLeft] = window.GetKey(glfw.KeyJ) == glfw.Press
// 	result[chibines.ButtonRight] = window.GetKey(glfw.KeyL) == glfw.Press
// 	return result
// }
