package main

import (
	"github.com/karalabe/hid"
	"log"
	"sync"
	"time"
)

const (
	HIDVendorId             = 0x1D34
	HIDProductId            = 0x0013
	BoardWidth, BoardHeight = 21, 7
	BytesPerRow             = 3
	RowsPerPacket           = 2
	NumPackets              = 4
	DisplayClearInterval    = time.Millisecond * 400
)

type LedArray [BoardWidth][BoardHeight]bool

func (l *LedArray) Reset() {
	for x, _ := range l {
		for y, _ := range l[x] {
			l[x][y] = false
		}
	}
}

// USB packets
type Packets [NumPackets][]byte

type Board struct {
	Device     *hid.Device
	DeviceLock sync.Mutex
	Timer      *time.Timer
	Packets    *Packets
}

func (b *Board) Start(ledChannel <-chan LedArray) {
	// Connect to display
	go b.ConnectDisplay()

	// Update display on new data or timer
	b.Timer = time.NewTimer(DisplayClearInterval)
	for {
		select {
		case ledArray := <-ledChannel:
			b.Packets = ConvertLedArray(&ledArray)
			b.UpdateDisplay()

		case <-b.Timer.C:
			b.UpdateDisplay()
		}
	}
}

func (b *Board) ConnectDisplay() {
	for {
		// Search for device
		log.Println("Searching for devices..")
		deviceInfos := hid.Enumerate(HIDVendorId, HIDProductId)
		if len(deviceInfos) > 0 {
			// Connect to device
			deviceInfo := deviceInfos[0]
			device, err := deviceInfo.Open()

			if err == nil {
				b.DeviceLock.Lock()
				b.Device = device
				b.DeviceLock.Unlock()
				log.Println("Connected to device.")
				return
			}
		}
		time.Sleep(time.Second * 5)
	}
}

func (b *Board) UpdateDisplay() {
	// Reset timer
	b.Timer.Reset(DisplayClearInterval)

	// Lock device
	b.DeviceLock.Lock()
	defer b.DeviceLock.Unlock()
	if b.Packets == nil || b.Device == nil {
		return
	}

	// Write packets
	for _, packetBytes := range b.Packets {
		_, err := b.Device.Write(packetBytes)

		// Reconnect to device on error
		if err != nil {
			b.Device = nil
			log.Println("Device error: ", err)
			go b.ConnectDisplay()
			return
		}
	}
}

// Convert the LED array to USB packets
func ConvertLedArray(ledArray *LedArray) *Packets {
	packets := Packets{}
	for packetIndex := 0; packetIndex < NumPackets; packetIndex++ {
		rowOffset := packetIndex * RowsPerPacket

		// Init packet {0, brightness, starting row}
		packetBytes := make([]byte, 0, 9)
		packetBytes = append(packetBytes, []byte{0, 0, byte(rowOffset)}...)

		for packetRowOffset := 0; packetRowOffset < RowsPerPacket; packetRowOffset++ {
			y := rowOffset + packetRowOffset
			if y >= BoardHeight {
				break
			}

			// Set single LEDs
			rowBytes := [BytesPerRow]byte{}
			byteCursor := 2
			bitCursor := byte(0)
			for x := 0; x < BoardWidth; x++ {
				// Turn LED on
				if ledArray[x][y] {
					rowBytes[byteCursor] |= 1 << bitCursor
				}

				// Next bit
				bitCursor++
				if bitCursor >= 8 {
					bitCursor = 0
					byteCursor--
				}
			}

			// Invert bits
			for i, _ := range rowBytes {
				rowBytes[i] ^= 0xFF
			}
			packetBytes = append(packetBytes, rowBytes[:]...)
		}
		packets[packetIndex] = packetBytes
	}
	return &packets
}
