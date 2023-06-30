# TIC-80 Bundler

A Go program that helps manage large TIC projects.

Sometimes it's useful to split your code across multiple files. This command-line program does one job - it bundles together source files into one looong file, which can then be used for demo submissions, or loaded onto [the TIC-80 Website](https://tic80.com)

## How it works

Lua has a `require(filename)` function, which allows you to pull in other files at runtime. It also provides a `package.path` variable, which specifies which directories to look in for these files.

This bundler does text pattern matching to:
- Pick up the `package.path` variables throughout the code, using this to discover subsequent required files.
- Find `require(...)` declarations, pasting the required file into a function, and switching the `require` into a call to that function.

Required files are deduplicated, so they will only appear once in the output, and ought to be included in the order they're picked up in the lua source. Requires are picked up throughout the whole tree, so you should be able to require in required files as much as you like.

You *MUST* make sure your package path is set up in this format, otherwise the regex match (which is very specific) will fail:

`package.path=package.path..";C:\\path\\to\\your\\lua\\files\\?.lua"`

(This should be configured to match your path!)

## Examples

Here are some suggestions for how you might employ required files:

### Example 1: Including a single file into global scope
```
-- FILE1.lua
MY_GLOBAL_VAR="hello"

-- MAIN.lua
package.path=package.path..";C:\\???\\???\\lua-tests\\example1\\?.lua"

require("./file1")

function TIC()
    cls()
    print(MY_GLOBAL_VAR,106,60,12)
end
```

### Example 2: Static data
```
-- FILE1.lua
return {"Hello", "There"}

-- MAIN.lua
package.path=package.path..";C:\\???\\???\\lua-tests\\example2\\?.lua"

DATA=require("./file1")

function TIC()
    cls()
    print(DATA[1],106,54,12)
    print(DATA[2],106,66,12)
end
```

### Example 3: Namespaced Functions
```
-- FILE1.lua
return {
    doSomePrint=function(txt,x,y)
        print(txt,x,y,12)
    end,
}

-- MAIN.lua
package.path=package.path..";C:\\???\\???\\lua-tests\\example3\\?.lua"

FNS=require("./file1")

function TIC()
    cls(0)
    FNS.doSomePrint("hello",106,60)
end
```

### Example 4: Library aggregation
```
-- FILE1.lua
return function(L)
    L=L or {}
    L.calculate=function(op1,op2)
        return op1 + op2
    end
end

-- FILE2.lua
return function(L)
    L=L or {}
    L.print=function(txt)
        print(txt,116,60,12)
    end
end

-- MAIN.lua
package.path=package.path..";C:\\???\\???\\lua-tests\\example4\\?.lua"

L={}
require("./file1")(L)
require("./file2")(L)

function TIC()
    cls()
    local val=L.calculate(23,6)
    L.print(val)
end
```

### Example 5: Data factory
```
-- FILE1.lua
return function(reps)
    local text="HO"
    for i=1,reps do
        text=text.."HO"
    end
    return text
end

-- MAIN.lua
package.path=package.path..";C:\\???\\???\\lua-tests\\example5\\?.lua"

HOHOHO=require("./file1")(3)

function TIC()
    cls()

    print(HOHOHO,98,60,12)
end
```

## Building

There are typical PC, Mac and Linux builds in the `build` directory.

This program is written in Go.

If you have Go installed, `go build` in this directory should create you an executable. You may need to run `go get` first.

If you want to refresh the `build` directory for others, `build.bat` will work on Windows.

## Running

This tool takes two arguments, `--src` and `--dest`. Both can be relative or full paths.

```tic-80-bundler.exe --src=mainsourcefile.lua --dest=outfile.lua```

*DO* take a copy / check your code into source control first in if an accidental overwrite would ruin your day!

## Known Issues

This tool will be super brittle, but I ([jtruk](https://mastodon.social/@jtruk)) got it to work on a complex, nested set of requires and it did good.

- It doesn't interpret Lua code - it just does pattern matching, so exotic runtime require systems will fail, e.g.
    ```
    local path="./file1.lua"
    require(path)
    ```
    ...will not work
- It seems to have some issues overwriting the destination file, overwriting part, but not pulling in the whole source. If this happens then try deleting the destination file and run again.
- TIC will fail to parse if one of your source files includes meta data blocks (e.g. `<PALETTE>...</PALETTE>` declarations). These were intended to be included at the end of the source file. I don't consider it an issue for this tool to resolve that - you should delete that meta data from your source file. It should be clear where this is happening from the line number of the error in your bundled file.

Drop me a line if it fails to work and I'll try to either fix it or suggest a workaround for the code.

## Credit

jtruk/RiFT. Happy for this to be a liberal license - [Code Credit 1.1.0](https://codecreditlicense.com/license/1.1.0) - you are welcome to fork and modify as you like, but please credit and ideally drop me a mention / pop me a hello if you use it.
