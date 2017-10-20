# Common logging across Virtuosys go binaries

This package uses github.com/sirupsen/logrus for logging and github.com/kelseyhightower/envconfig to read the environment

## Example

```bash

import "vsys/mas/pkg/maslog"

// Logger is used to send logging messages to stdout.
var logger = maslog.NewLogger()

func main() {

    logger.WithFields(logrus.Fields{
        "animal": "walrus",
        "number": 1,
        "size":   10,
    }).Info("A walrus appears")

    logger.Debug("An debug message")
    logger.Info("An information message")
    logger.Warn("An warning message")
    logger.Error("An error message")

    logger.Printf("Testing Version %s, BuildTime %s\n", config.Version, config.BuildTime)
}
```

## Guidelines.
1. Dont include newline statements on Print()/Println() but ok to include on Printf() functions
1. Dont pass parameters - these should be withFields
1. Logrus has global StandardLogger() or create your own (per package)
1. A common pattern is to re-use fields between logging statements by re-using
   the logrus.Entry returned from WithFields()
```bash
   contextLogger := log.WithFields(log.Fields{
     "common": "this is a common field",
     "other": "I also should be logged always",
   })
   contextLogger.Info("I'll be logged with common and other field")
```

## Glide dependency

Include this in your glide.yaml to add the dependency to your project (rather than git cloning the project). 
```bash
- package: vsys/golog
  repo: ssh://git@scm.vsys.net:7999/sw/golog.git
  vcs: git
  version: a63307f54b8f2a695f7e40d97e9022dc38f08d1e
  subpackages:
  - pkg/log
```
You also need to mount your $HOME/.ssh directory into the builder for permission to git pull. 

## Tags
This library is not built (yet) as it has no tests. The repository is manually tagged. Whenever changed this tag mneed to be manually bumped. A glide update is then required to pull the new version into your vendor directory 
```bash
git add README.md
git commit
git tag 0.0.2
git push --tags
git push
```


## Environment variables

The logging can be controlled by some environment variables

1. LOG_LEVEL

   Set the debug log level severity or above. "debug", "info", "warning", "error", "fatal", "panic"
   
   Default is log info for use in production. Development may want to use debug
   
1. LOG_DISABLE_TIMESTAMP

   Disable timestamps (useful when output is redirected to logging system that already adds timestamps)
  
   Default is no timestamps for use with docker logs which already include a timestamp

1. LOG_JSON_FORMAT  

   Log in JSON format (useful when output is redirected to logging system with JSON support).
   
   Default is text format for use with docker logs
   
1. LOG_CALLER_LEVELS
 
   List of log levels which should include the caller (file/line/function)
  
   Default is all levels: "debug", "info", "warning", "error", "fatal", "panic"
   
1. LOG_STACK_LEVELS

   List of log levels which should include a stack trace

   Default is "panic" level



## Unit Test
'''
    glide 
    docker exec -ti mas_mas-go-builder_1 sh
    cd /go/src/vsys/golog/pkg/log
'''

