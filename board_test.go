package main

import (
	"reflect"
	"testing"
)

// Tests
func TestEmptyScreen(t *testing.T) {
	packets := packets{
		{0, 0, 0, 255, 255, 255, 255, 255, 255},
		{0, 0, 2, 255, 255, 255, 255, 255, 255},
		{0, 0, 4, 255, 255, 255, 255, 255, 255},
		{0, 0, 6, 255, 255, 255},
	}

	ledArray := LedArray{}
	checkConvertLedArray(t, &packets, &ledArray)
}

func TestFullScreen(t *testing.T) {
	packets := packets{
		{0, 0, 0, 224, 0, 0, 224, 0, 0},
		{0, 0, 2, 224, 0, 0, 224, 0, 0},
		{0, 0, 4, 224, 0, 0, 224, 0, 0},
		{0, 0, 6, 224, 0, 0},
	}

	ledArray := LedArray{}
	for x := range ledArray {
		for y := range ledArray[x] {
			ledArray[x][y] = true
		}
	}
	checkConvertLedArray(t, &packets, &ledArray)
}

func TestSinglePixelScreen(t *testing.T) {
	packets := packets{
		{0, 0, 0, 255, 255, 255, 255, 255, 255},
		{0, 0, 2, 255, 255, 255, 255, 255, 255},
		{0, 0, 4, 255, 255, 247, 255, 255, 255},
		{0, 0, 6, 255, 255, 255},
	}

	ledArray := LedArray{}
	ledArray[3][4] = true
	checkConvertLedArray(t, &packets, &ledArray)
}

func TestEdgesScreen(t *testing.T) {
	packets := packets{
		{0, 0, 0, 224, 0, 0, 239, 255, 254},
		{0, 0, 2, 239, 255, 254, 239, 255, 254},
		{0, 0, 4, 239, 255, 254, 239, 255, 254},
		{0, 0, 6, 224, 0, 0},
	}

	ledArray := LedArray{}
	for x := 0; x < BoardWidth; x++ {
		ledArray[x][0] = true
		ledArray[x][BoardHeight-1] = true
	}
	for y := 0; y < BoardHeight; y++ {
		ledArray[0][y] = true
		ledArray[BoardWidth-1][y] = true
	}
	checkConvertLedArray(t, &packets, &ledArray)
}

func checkConvertLedArray(t *testing.T, expectedPackets *packets, ledArray *LedArray) {
	packets := convertLedArray(ledArray)
	if !reflect.DeepEqual(packets, expectedPackets) {
		t.Fatal("Packets don't match", "Got:", packets, "Expected:", expectedPackets)
	}
}

// Benchmarks
func BenchmarkConvertLedArrayEmpty(b *testing.B) {
	var ledArray LedArray
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		convertLedArray(&ledArray)
	}
}

func BenchmarkConvertLedArrayFull(b *testing.B) {
	var ledArray LedArray
	for x := 0; x < BoardWidth; x++ {
		for y := 0; y < BoardHeight; y++ {
			ledArray[x][y] = true
		}
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		convertLedArray(&ledArray)
	}
}
