# go-ebiten-suika

A Suika Game clone built with Go and the Ebiten game engine. Drop balls into the bin, watch them collide and stack with realistic physics.

![](./.github/capture1.gif)

## Features

- Physics-based ball collision and stacking
- Balls of varying sizes with realistic momentum transfer
- Simple and responsive controls

## Requirements

- Go 1.21 or later

## Installation

```sh
go install github.com/conao3/go-ebiten-suika@latest
```

Or clone and build from source:

```sh
git clone https://github.com/conao3/go-ebiten-suika.git
cd go-ebiten-suika
go build
```

## Usage

Run the game:

```sh
go run .
```

### Controls

| Key | Action |
|-----|--------|
| Left/Right Arrow | Move ball horizontally |
| Space | Drop ball |

## License

Apache License 2.0. See [LICENSE](LICENSE) for details.
