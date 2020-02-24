# update-go
**An automated update library for Go applications**

This library provides the functionality necessary to implement automated updates
for your applications.

## Example

```go

package main

import (
    "os"
    "log"
    "path/filepath"

    "github.com/SierraSoftworks/update-go"
)

// Update this on each new release (ideally using the "-X main.version=1.1.7" compiler flag)
var version = "1.0.0"

func main() {
    mgr := update.Manager{
        Application: os.Args[0],
        UpgradeApplication: filepath.Join(os.TempDir(), filepath.Base(os.Args[0])),

        Variant: update.MyPlatform(),
        Source: update.NewGitHubSource("SierraSoftworks/git-tool", "v", "git-tool-"),
    }

    // Resume an ongoing update operation (this may terminate the application as
    // part of the upgrade process).
    err := mgr.Continue()
    if err != nil {
        log.Fatalf("Unable to apply updates: %s", err)
    }

    // Your application code
    log.Println("Foo Bar!")

    rs, err := mgr.Source.Releases()
    if err != nil {
        log.Fatalf("Unable to fetch the list of available updates", err)
    }

    availableUpdate := update.LatestUpdate(rs, version)
    if availableUpdate != nil {
        log.Infof("Update available: %s", availableUpdate.ID)
        log.Infof("Changes: %s", availableUpdate.Changelog)

        err := mgr.Update(availableUpdate)
        if err != nil {
            log.Fatalf("Unable to start update: %s", err)
        }
    } else {
        log.Infof("No updates available")
    }
}

```