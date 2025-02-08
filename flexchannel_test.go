package flexchannel_test

import (
	"sync"
	"testing"
	"time"

	flexchannel "github.com/mochiya98/go-flexchannel"

	"github.com/stretchr/testify/assert"
)

const RepeatCount = 100

func TestFlexChannel(t *testing.T) {
	fc := flexchannel.NewFlexChannel()
	var wg sync.WaitGroup
	var mu sync.Mutex
	results := make([]*int, 0)

	h := func(data interface{}) {
		mu.Lock()
		defer mu.Unlock()
		defer wg.Done()
		results = append(results, data.(*int))
		// 50us * 100 = 5ms
		time.Sleep(200 * time.Microsecond)
	}

	wg.Add(100)
	err := fc.Start(h)
	assert.NoError(t, err)

	start := time.Now()
	for i := 0; i < RepeatCount; i++ {
		n := i
		fc.Send(&n)
	}
	elapsed := time.Since(start)
	assert.True(t, elapsed < 1*time.Millisecond, "should be fast")

	wg.Wait()

	assert.Equal(t, RepeatCount, len(results), "should all be received")
	for i := 0; i < RepeatCount; i++ {
		assert.Equal(t, i, *results[i], "messages should be in order")
	}

	fc.Close()
}

func TestFlexChannelClosedWrite(t *testing.T) {
	fc := flexchannel.NewFlexChannel()
	results := make([]interface{}, 0)

	h := func(data interface{}) {
		assert.Fail(t, "should not be received")
	}

	err := fc.Start(h)
	assert.NoError(t, err)

	fc.Close()

	func() {
		defer func() {
			if r := recover(); r != nil {
				assert.Fail(t, "panic should not happen")
			}
		}()
		fc.Send("test")
	}()

	assert.Equal(t, 0, len(results))
}
