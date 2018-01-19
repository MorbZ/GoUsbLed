package main

import "testing"

// Tests
type LedArrayByte [BoardHeight][BoardWidth]byte

func TestText3(t *testing.T) {
	ledArrayByte := LedArrayByte{
		{0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 1, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 1, 1, 0, 0, 1, 0, 0, 1, 0, 1, 0, 0, 1, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0},
	}
	CheckTextLedArray(t, "123", &ledArrayByte)
	CheckTextLedArray(t, "12a3", &ledArrayByte)
}

func TestText5(t *testing.T) {
	ledArrayByte := LedArrayByte{
		{0, 0, 1, 0, 0, 1, 1, 0, 0, 1, 0, 1, 0, 1, 1, 1, 0, 0, 1, 1, 0},
		{0, 1, 0, 1, 0, 0, 0, 1, 0, 1, 0, 1, 0, 1, 0, 0, 0, 1, 0, 0, 0},
		{0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 0, 1, 0, 1, 0, 0, 0, 1, 0, 0, 0},
		{0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 1, 1, 0, 1, 1, 0, 0, 1, 1, 0, 0},
		{0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 0, 1, 0},
		{0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 0, 1, 0},
		{0, 1, 1, 1, 0, 1, 1, 0, 0, 0, 0, 1, 0, 1, 1, 0, 0, 0, 1, 0, 0},
	}
	CheckTextLedArray(t, "23456", &ledArrayByte)
	CheckTextLedArray(t, "./-23test456@#''", &ledArrayByte)
}

func TestText8(t *testing.T) {
	ledArrayByte := LedArrayByte{
		{1, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 1, 1, 1, 0, 0, 1, 0, 0, 0, 1},
		{0, 1, 0, 1, 0, 1, 0, 0, 1, 1, 0, 0, 0, 1, 0, 1, 0, 1, 0, 1, 0},
		{0, 1, 0, 1, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 0, 1, 0, 1, 0},
		{1, 0, 0, 1, 0, 1, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1},
		{0, 1, 0, 1, 0, 1, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 1, 0, 0, 0},
		{0, 1, 0, 1, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0},
		{1, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1, 0, 0, 1, 1},
	}
	CheckTextLedArray(t, "23017894", &ledArrayByte)
}

func CheckTextLedArray(t *testing.T, str string, ledArrayByte *LedArrayByte) {
	text := NewText()
	ledArray := text.MakeText(str)
	ledArrayExpected := ConvertLedArrayByte(ledArrayByte)
	if *ledArray != *ledArrayExpected {
		t.Fatal("Resulting LED arrays don't match.")
	}
}

func ConvertLedArrayByte(ledArrayByte *LedArrayByte) *LedArray {
	ledArray := LedArray{}
	for x := 0; x < BoardWidth; x++ {
		for y := 0; y < BoardHeight; y++ {
			// x, y values are swapped on the LedArrayByte
			if ledArrayByte[y][x] != 0 {
				ledArray[x][y] = true
			}
		}
	}
	return &ledArray
}

// Benchmarks
func BenchmarkText3(b *testing.B)    { RunTextBenchmark(b, "123") }
func BenchmarkText4(b *testing.B)    { RunTextBenchmark(b, "1234") }
func BenchmarkText5(b *testing.B)    { RunTextBenchmark(b, "12345") }
func BenchmarkText10(b *testing.B)   { RunTextBenchmark(b, "0123456789") }
func BenchmarkTextLong(b *testing.B) { RunTextBenchmark(b, "asd4asd6ads4ads1asd7asd1") }

func RunTextBenchmark(b *testing.B, str string) {
	text := NewText()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		text.MakeText(str)
	}
}
