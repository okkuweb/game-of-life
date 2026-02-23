#!/bin/bash
GOOS=js GOARCH=wasm go build -o ./play-wasm/game-of-life.wasm .
