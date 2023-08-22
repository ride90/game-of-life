# Multiverse Game of Life

- Create multiple universes.
- Static & empty universes are deleted automatically.
- Merge all universes into one.
- Configurable fps.
- Full reset.
- Stream updates to clients via websockets.
- Render updates in the browser as canvas.
- Concurrent evolution of each universe (spawn a virtual thread per universe).

  [Demo video]([https://www.google.com](https://raw.githubusercontent.com/ride90/game-of-life/master/static/demo.mp4))

## Installation
`go mod tidy`

## Run
`go run cmd/main.go`

## TODO
- Rendering in the browser is inefficient -> generate a video on the server side and stream it to the browser.
