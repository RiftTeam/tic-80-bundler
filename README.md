# TIC-80 Bundler

A Go program to bundle required files into one looong source file.

This is to ensure `require`d files work. Lua pulls these in as it runs. This bundler wraps the contents of those files into functions and calls those at run-time instead.

Necessary to pull the many individual RiFT libs into one file for submitting to a competition, or putting online on TIC80.com.

## Building

There are typical PC, Mac and Linux builds in the `build` directory.

This program is written in Go.

If you have Go installed, `go build` in this directory should create you an executable. You may need to run `go get` first.

If you want to refresh the `build` directory for others, `build.bat` will work on Windows.

## Running

```tic-80-bundler.exe --src=mainsourcefile.lua --dest=outfile.lua```

Both src and dest should be full or relative paths. *DO* take a copy / check your code first in if an accidental overwrite would ruin your day.

## Warning

This will be super brittle, but I (jtruk) got it to work on a complex, nested set of requires and it did good. Drop me a line if it fails to work and I'll try to either fix it or suggest a workaround for the code.

## Credit

jtruk/RiFT. Happy for this to be a liberal license, but please drop a mention / pop me a hello if you use it.
