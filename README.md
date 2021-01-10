# auralog
A Go logging package based off Go's Log package with some differences. This has about as much overhead as the stock Log package because it doesn't use thirdparty packages like the stock package. The differences it has over the stock Go log library are below:
* Different logging levels ``INFO, WARN, ERROR, FATAL, PANIC``
* Config struct for cleaner initialization.
* Removed helpers so it has to be manually.
* Built in Rotate Writer io.Writer interface for basic log file rotation.

## Getting Started
``go get github.com/saintwish/auralog`` Will get the latest from master branch.

## Documentation
Nothing here yet.

## Example
This is a basic example of how you will initialize the logger and use it with the RotateWriter.
```go
package main

import(
  "os"
  "io"
  "time"

  "github.com/saintwish/auralog"
)

var (
  flags = auralog.Ldate | auralog.Ltime
  wflags = auralog.Ldate | auralog.Ltime
  eflags = auralog.Ldate | auralog.Ltime | auralog.Lshortfile
  dflags = auralog.Ldate | auralog.Ltime | auralog.Lshortfile

  log *auralog.Logger
)

func main() {
  file := &auralog.RotateWriter{
    Dir: "./runtime/logs/", //Required
    Filename: "log.log", //Required
    ExTime: 24 * time.Hour, //Required if you want daily log rotation.
    MaxSize: 5 * auralog.Megabyte, //Can also use auralog.Kilobyte, Not required
  }

  log = auralog.New(auralog.Config{
    Output: io.MultiWriter(file, os.Stdout), //Required
    Prefix: "[PREFIX] ", //Not Required
    Level: auralog.LogLevelInfo, //Required
    Flag: flags, //Required
    WarnFlag: wflags, //Required
    ErrorFlag: eflags, //Required
    DebugFlag: dflags, //Required
  })
  log.Println("Test")
  log.Warnln("TEST WARN")
  log.Errorln("TEST ERROR")
  log.Debugln("TEST DEBUG")
}
```

## Contributing
If you have any additions or contributions you would like to make please do. Just keep the below in mind.
* Try to match current naming conventions as closely as possible.
* Create a Pull Request with your changes against the master branch.

## License
It's licensed under BSD-2-Clause License to stay compatible with Go's stock log package, since it's created from that.
