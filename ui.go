package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
)

var (
	pid = getPid()
	wd  string
)

func setupUi() {
	mainWindow := ui.NewWindow("OLED Animation Stream", 500, 300, false)
	mainWindow.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})

	// Make a Grid UI with Text box, start, and stop button

	// box := ui.NewVerticalBox()
	// box.SetPadded(true)

	parentGrid := ui.NewGrid()
	parentGrid.SetPadded(true)

	grid := ui.NewGrid()
	grid.SetPadded(true)

	parentGrid.Append(grid,
		0, 0, 1, 1,
		true, ui.AlignCenter, true, ui.AlignCenter)

	// box.Append(grid, true)

	// Text box

	pathField := ui.NewEntry()
	fpsField := ui.NewEntry()
	inverCheck := ui.NewCheckbox("Invert")

	inverCheck.SetChecked(true)

	fpsField.SetText("30")

	openDir := ui.NewButton("Open")

	paramsForm := ui.NewForm()

	paramsForm.Append("Frames Directory ", pathField, true)
	paramsForm.Append("Open Dir", openDir, true)
	paramsForm.Append("FPS ", fpsField, true)
	paramsForm.Append("", inverCheck, true)

	grid.Append(paramsForm,
		0, 0, 2, 1,
		true, ui.AlignCenter, false, ui.AlignCenter)

	// Start button
	startButton := ui.NewButton("Start")

	grid.InsertAt(startButton,
		pathField, ui.Bottom, 1, 1,
		true, ui.AlignFill, true, ui.AlignCenter)

	// Stop button
	stopButton := ui.NewButton("Stop")

	grid.InsertAt(stopButton,
		startButton, ui.Trailing, 1, 1,
		true, ui.AlignFill, true, ui.AlignCenter)

	responseLabel := ui.NewLabel("")
	grid.InsertAt(responseLabel,
		startButton, ui.Bottom, 2, 1,
		true, ui.AlignStart, false, ui.AlignCenter)

	// Button events

	openDir.OnClicked(func(*ui.Button) {
		path := ui.OpenFile(mainWindow)
		if path != "" {
			pathField.SetText(path[:strings.LastIndex(path, "\\")])
		}
	})

	startButton.OnClicked(func(*ui.Button) {
		path := pathField.Text()
		fps := fpsField.Text()
		invert := inverCheck.Checked()
		pid = startProcess(path, fps, invert)
		if pid != -1 {
			responseLabel.SetText(fmt.Sprintf("Stream Started with PID %d", pid))
		} else {
			responseLabel.SetText("Please provide a valid path")
		}
	})

	stopButton.OnClicked(func(*ui.Button) {
		responseLabel.SetText("Stopping Stream...")
		resp := stopProcess()
		responseLabel.SetText(resp)
	})

	mainWindow.SetMargined(true)

	mainWindow.SetChild(parentGrid)
	mainWindow.Show()
}

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	wd = dir
	// fmt.Println(dir)
	ui.Main(setupUi)
}

func startProcess(path string, fps string, invert bool) int {
	if path == "" {
		return -1
	}

	if getPid() != -1 {
		stopProcess()
	}

	inv := ""
	if invert {
		inv = "-invert"
	}

	cmd := exec.Command(fmt.Sprintf("%s\\test_oled_frame.exe", wd), "-path", path, "-fps", fps, inv)
	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	pid = cmd.Process.Pid
	err = os.WriteFile("_pid", []byte(fmt.Sprintf("%d", pid)), 0644)
	if err != nil {
		panic(err)
	}

	return pid
}

func stopProcess() string {
	pid = getPid()
	if pid != -1 {
		// fmt.Println("Stopping process...")
		process, err := os.FindProcess(pid)
		if err != nil && strings.HasPrefix(err.Error(), "OpenProcess") {
			return "Stream not found"
		} else if err != nil {
			panic(err)
		}

		if err = process.Signal(syscall.SIGKILL); errors.Is(err, fs.ErrPermission) {
			return "Stream already stopped"
		} else if err != nil {
			panic(err)
		}

		os.WriteFile("_pid", []byte(""), 0644)

		return "Stream Stopped"
	} else {
		// fmt.Println("No process running.")
		return "No Stream running"
	}
}
