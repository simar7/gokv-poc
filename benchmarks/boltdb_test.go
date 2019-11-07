package benchmarks

import "testing"

func BenchmarkBoltUpdate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BoltUpdate()
	}
}

func BenchmarkBoltBatch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BoltBatch()
	}
}
