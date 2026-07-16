package main

import (
	"fmt"
	"os"

	"github.com/harishnagaraju/astramind/internal/engine"
)

func main() {
	if err := engine.New().Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
