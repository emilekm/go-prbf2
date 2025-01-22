# go-prbf2

This repository serves libraries written in Golang for Battlefield 2: Reality mod.

## bf2demo

BF2Demo files are produced by the game server and contain binary representation of in-game events that are reproducible in the game client.

The library only reads the header of the whole file.

## prdemo

PRDemo files are binary files composed of structured messages produced by game server.

The library allows for iterating over messages and unmarshaling them into static types.

## prism

PRISM is a protocol used by GUI tool of the same name to manage a running server and it's players without the need of being logged in-game.

The library exposes complete set of funtionality provided by PRISM.