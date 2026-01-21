package coingecko

import (
	"bytes"
	"encoding/json"
	"io"
	"sync"
	"testing"
)

// Benchmark implementations for reading response body

func BenchmarkDecode_Stream(b *testing.B) {
	jsonBody := []byte(`{"bitcoin": {"usd": 50000.0, "usd_24h_change": 2.5, "usd_24h_vol": 100000}}`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(jsonBody)
		var data simplePriceResponse
		_ = json.NewDecoder(r).Decode(&data)
	}
}

func BenchmarkDecode_ReadAll(b *testing.B) {
	jsonBody := []byte(`{"bitcoin": {"usd": 50000.0, "usd_24h_change": 2.5, "usd_24h_vol": 100000}}`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(jsonBody)
		body, _ := io.ReadAll(r)
		var data simplePriceResponse
		_ = json.Unmarshal(body, &data)
	}
}

func BenchmarkDecode_SyncPool(b *testing.B) {
	jsonBody := []byte(`{"bitcoin": {"usd": 50000.0, "usd_24h_change": 2.5, "usd_24h_vol": 100000}}`)
	var pool = sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(jsonBody)
		buf := pool.Get().(*bytes.Buffer)
		buf.Reset()
		_, _ = io.Copy(buf, r)

		var data simplePriceResponse
		_ = json.Unmarshal(buf.Bytes(), &data)

		pool.Put(buf)
	}
}
