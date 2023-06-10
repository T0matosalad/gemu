<p align="center">
    <img src="gemu.png" width="200">
</p>

# GEMU: Gameboy EMUlator
GEMU is a GAMEBOY emulator written in golang.

## Installation
### Building from source
```
$ cd gemu
$ make
$ ls dist
gemu
```

## Usage
```
$ gemu -h
Usage of gemu:

gemu [-vrd] ROM
    -v         display version
    -r int     magnification ratio of screen (default: 1)
    -l string  log level {verbose, debug, warn, error, fatal} (default: debug)
    -d         start debug mode
```

# Resources
- [The Ultimate Game Boy Talk (33c3)](https://youtu.be/HyzD8pNlpwI)
- [GB DEV](https://gbdev.io/)
- [Game Boy CPU (SM83) instruction set](https://gbdev.io/gb-opcodes/optables/)
- [Game boy Architecture](https://www.copetti.org/writings/consoles/game-boy/)

# ROMs
- [Tobu Tobu Girl](https://tangramgames.dk/tobutobugirl/)
- [gb-test-roms](https://github.com/retrio/gb-test-roms)
- [helloworld](https://github.com/gitendo/helloworld)