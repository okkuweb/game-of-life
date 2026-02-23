#!/bin/bash
cp "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" ./play-wasm/
GOOS=js GOARCH=wasm go build -o ./play-wasm/game-of-life.wasm .
