# Guru WSL Bridge

This project is a little helper to achieve full
GO `guru` functionality on Windows systems.

## Problem:
When searching referrers of a GO type or function while
positioning the cursor on the definition rather than on a usage,
`guru` will not find all usages on Windows compared to what it
finds on Linux systems.

## Solution:
I've set up my Windows 10 WSL system to be the main development
environment and use the Windows side mostly for editing. For GO
code I use Visual Studio Code together with the GO plugin.

Since I have set up GO on both sides of the fence anyway,
I have created a little `guru` proxy that will read input from
VS Code, translate the GOPATH, call the Linux `guru` executable
and return the result. VS Code picks up the result and shows the
referrers as it would on a Linux system.

I hope this is a short lived workaround. It will be obsolete as
soon as the Windows version of `guru` fixes this problem. Until
then I hope it does work- and maybe helps a few others...

## Installation

First make sure you have GO installed on Windows and on WSL- I'm using Ubuntu 18.04.

The `guru` proxy calls `wsl.exe guru.sh` together with the parameters
handed over by VS Code.

When we call an executable this way it's not calling the users .bashrc, which makes the call way more performant, but misses out on
important environment settings like the GOROOT and GOPATH. These are
set within `guru.sh`, which in turn than just calls the Linux `guru`.

* Install `guru` in WSL (if it's not there already):

```
go get -u -v golang.org/x/tools/cmd/guru
```

* Please make sure you adapt `guru.sh` to your environment:

```
#!/bin/sh

export GOROOT=/ewu/go-1.10
export GOPATH=/entw/go

$GOPATH/bin/guru "$@"
```

* Put `guru.sh` somewhere in your (Windows) path. When we execute `wsl.exe guru.sh` wsl will append the Windows PATH to the global Linux path (no `.bashrc` is called), that's why an executable put into a Windows path will be
found on the Linux side. Try it our by calling `wsl echo $PATH` from a `cmd.exe` window. 

* Put `guru-wsl-bridge.json.example` into your `%USERPROFILE%` folder (`/Users/tuhlmann` for me), adapt it  and rename it to
`.guru-wsl-bridge.json`:

```
{
  "GOPATHOnWindows": "c:\\entw\\go",
  "GOPATHOnLinux": "/entw/go"
}
```

That's the config file read by the `guru` proxy so it knows how
to rewrite path names.

* run `go install` in this repository. This will create a
binary `guru-wsl-bridge.exe` inside your `%GOPATH%/bin` directory.
Go there, backup the original `guru.exe` and rename the new one into
`guru.exe` so VS Code will find it: `ren guru-wsl-bridge.exe guru.exe`


That should be it, let me know if I forgot something.

# License

MIT
