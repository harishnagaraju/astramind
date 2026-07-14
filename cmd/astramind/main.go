package main

import (
	"fmt"
	"os"

	"github.com/harishnagaraju/astramind/internal/app"
)

func main() {
	if err := app.New().Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
