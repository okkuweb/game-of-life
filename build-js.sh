#!/bin/bash
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" ./play-wasm/
GOOS=js GOARCH=wasm go build -o ./play-wasm/game-of-life.wasm .
