# delayed

A small and simple delayed tasks management library

## Overview

* No external dependencies
* Developer friendly

## Usage

```go
package main

import (
	"fmt"
	"time"

	"github.com/buraindo/delayed"
)

func main() {
	m, err := delayed.NewTaskManager(time.Second) // err = nil
	future, err := m.Run(func() (any, error) { // err = nil
		return "hello", nil
	}, 10*time.Second)
	if future.HasError() { // blocks for 10 sec, returns `false` after
		fmt.Printf("error: %v", future.Error())
	}
	fmt.Printf("result: %v", future.Get()) // "hello"
}
```

## License

© buraindo, 2023~time.Now

Released under the [MIT License](https://github.com/buraindo/delayed/blob/master/License)