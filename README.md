# tic-80-bundler

A Go program to bundle required files into one looong source file

This is to take of `require`d files that are loaded at run-time. This bundler pulls the contents of those files into functions and calls them at run-time instead.

Necessary to pull the many individual RiFT libs into one file for submitting to a competition, or putting online on TIC80.com.

## Building

This program is written in Go.

If you have Go installed, `go build` should create you an executable.

## Running

```tic-80-bundler.exe --src=mainsourcefile.lua --dest=outfile.lua```
