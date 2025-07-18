package log

import (
	"context"
	"testing"
)

func BenchmarkInfo(b *testing.B) {
	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		Info(ctx).
			Str("key", "value").
			Msg("Hello world")
	}
}
