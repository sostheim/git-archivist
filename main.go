package main

import (
	goflag "flag"
	"fmt"
	"time"

	"github.com/blang/semver"
	"github.com/golang/glog"
	flag "github.com/spf13/pflag"
)

// MajorMinorPatch - semantic version string
var MajorMinorPatch string

// ReleaseType - release type
var ReleaseType = "alpha"

// GitCommitSha - git commit sha-1 hash
var GitCommitSha string

var gaCfg *config

func init() {
	go Until(glog.Flush, 10*time.Second, NeverStop)
	gaCfg = newConfig()
}

func displayVersion() {
	semVer, err := semver.Make(MajorMinorPatch + "-" + ReleaseType + "+git.sha." + GitCommitSha)
	if err != nil {
		panic(err)
	}
	fmt.Println(semVer.String())
}

func main() {
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()

	gaCfg.flagSet = flag.CommandLine

	// check for version flag, if present print veriosn and exit
	if *gaCfg.version {
		displayVersion()
		return
	}

	gaCfg.envParse()
	glog.V(2).Infof("main(): configuration: %v", gaCfg.String())

	if false == gaCfg.validate() {
		return
	}

	// Git Archivist Service
	glog.V(2).Infof("main(): staring Git Archivist Service")
	srv := newGAServer(gaCfg)

	srv.run()
}
