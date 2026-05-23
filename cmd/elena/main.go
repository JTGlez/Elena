package main

import (
	"elena/internal/app"
	"elena/internal/infrastructure/entrypoints"
)

func main() {
	a := app.Wire()
	if err := entrypoints.StartTUI(a); err != nil {
		panic(err)
	}
}
