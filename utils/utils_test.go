package utils

import (
	"bytes"
	"compress/gzip"
	"testing"
)

var data = "[{\"name\":\"苹果\",\"category\":\"水果\",\"price\":5.6,\"rating\":4.5},{\"name\":\"牛奶\",\"category\":\"饮料\",\"price\":3.2,\"rating\":4.2},{\"name\":\"书包\",\"category\":\"文具\",\"price\":25.8,\"rating\":4.7},{\"name\":\"手机\",\"category\":\"电子产品\",\"price\":1999,\"rating\":4.8}]"

func BenchmarkGzipData(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = gzipData(data)
		}

	})
}

func BenchmarkGzipDataV2(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var buf bytes.Buffer
			zw := gzip.NewWriter(&buf)
			_, err := zw.Write([]byte(data))
			if err != nil {
				zw.Close()
			}
			zw.Close()
		}

	})
}
