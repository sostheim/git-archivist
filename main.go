/*
Copyright 2017 Samsung SDSA CNCT

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	goflag "flag"
	"fmt"

	"github.com/blang/semver"
	"github.com/golang/glog"
	flag "github.com/spf13/pflag"
)

// MajorMinorPatch - semantic version string
var MajorMinorPatch string

// ReleaseType - release type
var ReleaseType = "alpha"

// GitCommit - git commit sha-1 hash
var GitCommit string

var gaCfg *config

func init() {
	gaCfg = newConfig()
}

func displayVersion() {
	semVer, err := semver.Make(MajorMinorPatch + "-" + ReleaseType + "+git.sha." + GitCommit)
	if err != nil {
		panic(err)
	}
	fmt.Println(semVer.String())
}

func main() {
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()

	gaCfg.flagSet = flag.CommandLine
	gaCfg.envParse()

	// check for version flag, if present print veriosn and exit
	if *gaCfg.version {
		displayVersion()
		return
	}

	// Git Archivist Service
	glog.V(3).Infof("main(): staring Git Archivist Service")
	srv := newGAServer(gaCfg)

	srv.run()
}
