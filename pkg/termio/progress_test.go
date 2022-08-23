package termio

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func ExampleProgressBar() {
	max := 100
	pb := NewProgressBar("example", int64(max))

	for i := 0; i < max+20; i++ {
		pb.Inc()
		pb.Add(23)
		pb.Set(42)
		time.Sleep(150 * time.Millisecond)
	}

	time.Sleep(5 * time.Second)
	pb.Done()
}

func TestProgress(t *testing.T) { //nolint:paralleltest
	max := 2
	pb := NewProgressBar("test", int64(max))
	pb.Hidden = true
	pb.Inc()
	assert.Equal(t, int64(1), pb.current())
}
