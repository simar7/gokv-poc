package benchmarks

import "testing"

func BenchmarkBoltUpdate(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		BoltUpdate()
	}
}

func BenchmarkBoltUpdateGoRoutines(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		BoltUpdateGoRoutines()
	}
}

func BenchmarkBoltBatch(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		BoltBatch()
	}
}

func BenchmarkBoltBatchGoRoutines(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		BoltBatchGoRoutines()
	}
}
