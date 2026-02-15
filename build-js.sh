#!/bin/bash
GOOS=js GOARCH=wasm go build -o ./play-wasm/app.wasm ./*.go
