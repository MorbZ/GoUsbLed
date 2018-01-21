package main

import (
	"github.com/karalabe/hid"
	"log"
	"sync"
	"time"
)

const (
	// BoardWidth Number of horizintal LEDs
	BoardWidth = 21

	// BoardHeight Number of vertical LEDs
	BoardHeight = 7

	hidVendorID          = 0x1D34
	hidProductID         = 0x0013
	bytesPerRow          = 3
	rowsPerPacket        = 2
	numPackets           = 4
	displayClearInterval = time.Millisecond * 400
)

// LedArray for each LED on the message board
type LedArray [BoardWidth][BoardHeight]bool

// Reset turn off all LEDs
func (l *LedArray) Reset() {
	for x := range l {
		for y := range l[x] {
			l[x][y] = false
		}
	}
}

// USB packets
type packets [numPackets][]byte

// Board connects and communicates with the USB message board
type Board struct {
	device     *hid.Device
	deviceLock sync.Mutex
	timer      *time.Timer
	Packets    *packets
}

// Start connection to the message board
func (b *Board) Start(ledChannel <-chan LedArray) {
	// Connect to display
	go b.connectDisplay()

	// Update display on new data or timer
	b.timer = time.NewTimer(displayClearInterval)
	for {
		select {
		case ledArray := <-ledChannel:
			b.Packets = convertLedArray(&ledArray)
			b.updateDisplay()

		case <-b.timer.C:
			b.updateDisplay()
		}
	}
}

func (b *Board) connectDisplay() {
	for {
		// Search for device
		log.Println("Searching for devices..")
		deviceInfos := hid.Enumerate(hidVendorID, hidProductID)
		if len(deviceInfos) > 0 {
			// Connect to device
			deviceInfo := deviceInfos[0]
			device, err := deviceInfo.Open()

			if err == nil {
				b.deviceLock.Lock()
				b.device = device
				b.deviceLock.Unlock()
				log.Println("Connected to device.")
				return
			}
		}
		time.Sleep(time.Second * 5)
	}
}

func (b *Board) updateDisplay() {
	// Reset timer
	b.timer.Reset(displayClearInterval)

	// Lock device
	b.deviceLock.Lock()
	defer b.deviceLock.Unlock()
	if b.Packets == nil || b.device == nil {
		return
	}

	// Write packets
	for _, packetBytes := range b.Packets {
		_, err := b.device.Write(packetBytes)

		// Reconnect to device on error
		if err != nil {
			b.device = nil
			log.Println("Device error: ", err)
			go b.connectDisplay()
			return
		}
	}
}

// Convert the LED array to USB packets
func convertLedArray(ledArray *LedArray) *packets {
	packets := packets{}
	for packetIndex := 0; packetIndex < numPackets; packetIndex++ {
		rowOffset := packetIndex * rowsPerPacket

		// Init packet {(0), brightness, starting row}
		packetBytes := make([]byte, 0, 9)
		if hidPacketHasZeroByte {
			packetBytes = append(packetBytes, 0)
		}
		packetBytes = append(packetBytes, 0, byte(rowOffset))

		for packetRowOffset := 0; packetRowOffset < rowsPerPacket; packetRowOffset++ {
			y := rowOffset + packetRowOffset
			if y >= BoardHeight {
				break
			}

			// Set single LEDs
			rowBytes := [bytesPerRow]byte{}
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
			for i := range rowBytes {
				rowBytes[i] ^= 0xFF
			}
			packetBytes = append(packetBytes, rowBytes[:]...)
		}
		packets[packetIndex] = packetBytes
	}
	return &packets
}
