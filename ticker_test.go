package main

import "testing"

// Tests
func TestTickerEventNormal(t *testing.T) {
	eventData := `{"amount": 0.0035188699999999999, "buy_order_id": 808325081, "sell_order_id": 808325501, "amount_str": "0.00351887", "price_str": "11420.29", "timestamp": "1516557485", "price": 11420.290000000001, "type": 1, "id": 48598213}`
	testTickerEvent(t, eventData, 11420.29, false)
}

func TestTickerEventInvalidPrice(t *testing.T) {
	eventData := `{"amount": 0.0035188699999999999, "buy_order_id": 808325081, "sell_order_id": 808325501, "amount_str": "0.00351887", "price_str": "11420.29", "timestamp": "1516557485", "price": test, "type": 1, "id": 48598213}`
	testTickerEvent(t, eventData, 0, true)
}

func TestTickerEventPriceOnly(t *testing.T) {
	eventData := `{"price":3}`
	testTickerEvent(t, eventData, 3, false)
}

func TestTickerEventEmptyJson(t *testing.T) {
	eventData := `{}`
	testTickerEvent(t, eventData, 0, true)
}

func TestTickerEventInvalid(t *testing.T) {
	eventData := `test`
	testTickerEvent(t, eventData, 0, true)
}

func TestTickerEventEmpty(t *testing.T) {
	eventData := ``
	testTickerEvent(t, eventData, 0, true)
}

func testTickerEvent(t *testing.T, eventData string, expectedPrice float64, throwsErr bool) {
	price, err := getPriceForEventData(eventData)
	if err != nil && !throwsErr {
		t.Fatal("Function throws error but should not.")
	}
	if err == nil && throwsErr {
		t.Fatal("Function does not throw error but should.")
	}
	if price != expectedPrice {
		t.Fatal("Price is wrong:", "Got:", price, "Expected:", expectedPrice)
	}
}

// Benchmarks
func BenchmarkTickerEvent(b *testing.B) {
	eventData := `{"amount": 0.0035188699999999999, "buy_order_id": 808325081, "sell_order_id": 808325501, "amount_str": "0.00351887", "price_str": "11420.29", "timestamp": "1516557485", "price": 11420.290000000001, "type": 1, "id": 48598213}`
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		getPriceForEventData(eventData)
	}
}
