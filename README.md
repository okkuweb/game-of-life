# Game of Life
Conway's Game of Life using [Gruid](https://codeberg.org/anaseto/gruid).

Runs in a terminal and on a browser with WASM.

## Try out the web version
A WASM build is hosted on [github pages](https://okkuweb.github.io/game-of-life/)

## Keybindings:
- q, -: reduce update interval
- e, +: increase update interval
- space, p, P: pause/unpause game updates
- left mouse: add life
- right mouse: remove life
- click and drag mouse: paints with one of the above
### Terminal specific bindings:
- Q, ESC: quit game
- w, a, s, d: increase and decrease window size


## Dependencies
- [go](https://go.dev/doc/install)
## To run in the terminal
1. `go get`
2. `go run ./...`
## To build wasm version
1. Copy your local wasm by running `./dev/copy-wasm.sh` in the project root dir
2. Run `./dev/build-js.sh`
3. Run `go run -tags serve ./dev/serve.go` to serve on http://localhost:8080
