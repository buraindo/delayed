# Delayed

A small and simple delayed tasks management library

[![go report card](https://goreportcard.com/badge/github.com/buraindo/delayed "go report card")](https://goreportcard.com/report/github.com/buraindo/delayed)
[![test status](https://github.com/buraindo/delayed/workflows/Go/badge.svg?branch=master "test status")](https://github.com/buraindo/delayed/actions)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)

## Overview

* No external dependencies
* Developer friendly

## Install

```
go get -u github.com/buraindo/delayed
```

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