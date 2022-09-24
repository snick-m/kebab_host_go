package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/karalabe/hid"
)

const (
	FRAME_SIZE = 31
	DEBUG      = false
)

var (
	TOTAL_PRINTED = 0
	FPS           = 30
	PATH          = ""
	INVERT        = false

	frame_idx = 0
	dev       *hid.Device
	err       error
	frames    [][512]uint8
)

func init() {
	flag.IntVar(&FPS, "fps", 30, "Frames per second")
	flag.StringVar(&PATH, "path", "", "Path to folder with images")
	flag.BoolVar(&INVERT, "invert", false, "Invert colors")
	flag.Parse()

	if PATH == "" {
		fmt.Println("NO_PATH")
		os.Exit(1)
	}

	frames = readImages(PATH, INVERT)
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	// signal.Notify(c, syscall.SIGINT)

	// Find the device with the given VID/PID.
	dev, err = hid.Enumerate(0xFEED, 0x0000)[0].Open()
	if err != nil {
		panic(err)
	}
	defer dev.Close()

	go func() {
		for {
			sendFrame(frame_idx)
			frame_idx = (frame_idx + 1) % len(frames)

			// Control with FPS
			time.Sleep(time.Duration(1000/FPS) * time.Millisecond)
		}
	}()

	<-c
	// fmt.Println("exiting")
	dev.Close()
	os.Exit(0)

	// Read some data from the device.

	// Print out the returned buffer.
	// for i := 0; i < len(buf); i++ {
	// 	fmt.Printf("%d ", buf[i])
	// }
}

func sendFrame(idx int) {
	if DEBUG {
		fmt.Println(len(frames[idx]))
	}

	// Initial frame command id is 0x15
	writeToDev(dev, append([]uint8{0x15}, frames[idx][:FRAME_SIZE]...))

	send := frames[idx][FRAME_SIZE:]
	for len(send) > 0 {
		var send_snip []uint8

		// Continue frame command id is 0x16
		if len(send) > FRAME_SIZE {
			send_snip = append([]uint8{0x16}, send[:FRAME_SIZE]...)
			send = send[FRAME_SIZE:]
		} else {
			send_snip = append([]uint8{0x16}, send...)
			send = []uint8{}
		}

		writeToDev(dev, send_snip)
	}
}

func writeToDev(dev *hid.Device, data []uint8) {

	buf := make([]uint8, FRAME_SIZE+1)

	if DEBUG {
		fmt.Println(len(data))
	}
	_, err := dev.Write(data)
	if err != nil {
		panic(err)
	}

	dev.Read(buf)

	if DEBUG {
		for i := 1; i < len(buf); i++ {
			fmt.Print(buf[i], " ")
			TOTAL_PRINTED += 1

			if TOTAL_PRINTED%31 == 0 {
				fmt.Println("")
			}
		}
	}

}
