# FlexChannel

FlexChannel is a Go package that provides a non-blocking flexible channel implementation. It behaves like a channel with infinite size buffer, and it doesn't cause a panic even if you attempt to send after it has been closed.


## Installation

To install FlexChannel, use `go get`:

```sh
$ go get github.com/mochiya98/go-flexchannel
```

## Usage

Here is an example of how to use FlexChannel:

```go
package main

import (
	"fmt"
	"github.com/mochiya98/go-flexchannel"
)


func main() {
	fc := flexchannel.NewFlexChannel()

	// Start a goroutine to handle messages
	fc.Start(func(data interface{}) {
		// Receive message in order they were sent
		fmt.Println(data.(string))
	})

	start := time.Now()
	fc.Send("hello")
	fc.Send("works")
	assert.True(
		t,
		time.Since(start) < 1*time.Millisecond,
		"non-blocking send",
	)

	fc.Close()

	fc.Send("doesn't cause a panic after closed")
}

```

## API

### NewFlexChannel

```go
func NewFlexChannel() *FlexChannel
```

Creates a new FlexChannel instance.

### Start

```go
func (fc *FlexChannel) Start(h func(interface{})) error
```

Starts the FlexChannel with a handler function. The handler function will be called with each message in the order they were sent.

### Send

```go
func (fc *FlexChannel) Send(data interface{})
```

Sends a message to the FlexChannel.
It does not cause panic even if you attempt to send after the channel has been closed.

### Close

```go
func (fc *FlexChannel) Close()
```

Closes the FlexChannel.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
