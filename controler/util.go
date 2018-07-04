package main

import (
	"fmt"
	"math/rand"
)

func UUID() string {
	return fmt.Sprintf("%x", rand.Uint64())
}
