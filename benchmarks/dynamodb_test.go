package benchmarks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkDynamoBatchWrite(b *testing.B) {
	b.ReportAllocs()
	svc, err := setupDynamo()
	assert.NoError(b, err)

	createTable(svc)

	for i := 0; i < b.N; i++ {
		DynamoBatchWrite(svc)
	}

	deleteTable(svc)
}

func BenchmarkDynamoUpdate(b *testing.B) {
	b.ReportAllocs()
	svc, err := setupDynamo()
	assert.NoError(b, err)

	createTable(svc)

	for i := 0; i < b.N; i++ {
		DynamoUpdate(svc)
	}

	deleteTable(svc)
}
