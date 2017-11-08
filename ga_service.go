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
	"os"
	"strings"
	"time"

	"github.com/golang/glog"
)

var (
	resyncPeriod     = 30 * time.Second
	resyncInProgress = false
)

// GA Server
type gaServer struct {
	// Command line / environment supplied configuration values
	cfg     *config
	repoURL string
	stopCh  chan struct{}
}

func newGAServer(cfg *config) *gaServer {
	return &gaServer{
		stopCh:  make(chan struct{}),
		cfg:     cfg,
		repoURL: GitProtocolHTTPS + *cfg.username + ":" + *cfg.password + "@" + *cfg.server + "/" + *cfg.account + "/" + *cfg.repository + ".git",
	}
}

func (as *gaServer) setDirectory() bool {

	err := os.Chdir(*as.cfg.directory)
	if err != nil {
		glog.Errorf("error: executing chdir %s, returned: %v", *as.cfg.directory, err)
		return false
	}

	if glog.V(2) {
		dir, _ := os.Getwd()
		glog.Infof("current working directory: %s", dir)
	}
	return true
}

func (as *gaServer) clone() bool {
	if *as.cfg.initialize == false {
		glog.V(2).Info("no initial repository clone requested")
		return true
	}
	glog.V(2).Infof("clone initial repository: %s", *as.cfg.repository)

	var cloneArgs = []string{GitClone, GitArgDepth, "1", as.repoURL}
	if *as.cfg.directory != "" {
		cloneArgs = append(cloneArgs, *as.cfg.directory)
	}
	_, err := Execute(GitCmd, cloneArgs)
	if err != nil {
		glog.Errorf("error: executing %s %s, returned: %v", GitCmd, GitClone, err)
		return false
	}
	return true
}

func (as *gaServer) pushUpdates() {
	glog.V(2).Infof("push updates at: %v", time.Now())
	_, err := Execute(GitCmd, []string{GitPush, as.repoURL, GitBranchMaster})
	if err != nil {
		glog.Warningf("error: executing %s %s %s %s, returned: %v",
			GitCmd, GitPush, GitRemoteOrigin, GitBranchMaster, err)
	}
}

func (as *gaServer) commitUpdates() {
	glog.V(2).Infof("commit updates at: %v", time.Now())
	configAuthor := GitConfigUserName + "=\"" + *as.cfg.username + "\""
	configEmail := GitConfigUserEmail + "=\"" + *as.cfg.email + "\""
	commitAuthor := "\"" + *as.cfg.username + " <" + *as.cfg.email + ">\""
	commitMessage := "\"" + "git-archivist: auto update: " + time.Now().String() + "\""
	_, err := Execute(GitCmd, []string{GitArgC, configAuthor, GitArgC, configEmail, GitCommit, GitArgAuthor, commitAuthor, GitArgAM, commitMessage})
	if err != nil {
		glog.Warningf("error: executing %s %s %s, returned: %v",
			GitCmd, GitCommit, commitMessage, err)
		return
	}
	as.pushUpdates()
}

func (as *gaServer) checkForUpdates() {
	// resyncInProgress is not a mutext, but a trivial way to avoid obvious failure modes
	if resyncInProgress {
		glog.V(2).Infof("resync currently in progress, skipping check for local updates at: %v", time.Now())
		return
	}

	glog.V(2).Infof("checking for local updates at: %v", time.Now())

	cmdOutBytes, err := Execute(GitCmd, []string{GitStatus, GitArgShort, GitArgNoUntracked})
	if err != nil {
		glog.Warningf("error: executing %s %s %s %s, returned: %v",
			GitCmd, GitStatus, GitArgShort, GitArgNoUntracked, err)
		return
	}
	mods := false
	outLines := strings.Split(string(cmdOutBytes), "\n")
	if len(outLines) > 0 {
		for _, line := range outLines {
			// Check to see if the 2nd character of the output line a liternal `M`, e.g.
			// $ git status --short --untracked-files=no
			//  M config.go
			//  ^
			if len(line) > 2 && line[1] == 77 {
				mods = true
				break
			}
		}
	}
	if mods {
		resyncInProgress = true
		defer func() {
			resyncInProgress = false
		}()
		as.commitUpdates()
	}
}

func (as *gaServer) checkForRemoteUpdates() {
	// resyncInProgress is not a mutext, but a trivial way to avoid obvious failure modes
	if resyncInProgress {
		glog.V(2).Infof("resync currently in progress, skipping check for remote updates at: %v", time.Now())
		return
	}

	glog.V(2).Infof("checking for remote updates at: %v", time.Now())

	resyncInProgress = true
	defer func() {
		resyncInProgress = false
	}()

	_, err := Execute(GitCmd, []string{GitPull, GitRemoteOrigin, GitBranchMaster})
	if err != nil {
		glog.Warningf("error: executing %s %s %s %s, returned: %v",
			GitCmd, GitStatus, GitArgShort, GitArgNoUntracked, err)
		glog.Info("info: are there uncommitted changes blocking the pull's merge?")
	}
}

func (as *gaServer) run() {
	glog.V(2).Infof("starting run at: %v", time.Now())
	if as.clone() == false {
		return
	}

	if *as.cfg.initonly == true {
		glog.Infof("exiting after initialization only, at: %v", time.Now())
		return
	}

	if as.setDirectory() == false {
		return
	}

	if *as.cfg.direction == "remote" || *as.cfg.direction == "both" {
		// sync (push) all changes locall to the remote origin master branch
		go Until(as.checkForUpdates, time.Duration(*as.cfg.frequency)*time.Second, as.stopCh)
		if *as.cfg.direction == "both" {
			time.Sleep(time.Duration(*as.cfg.frequency) * time.Second / 2)
		}
	}

	if *as.cfg.direction == "local" || *as.cfg.direction == "both" {
		// sync (pull) all changes from the remote origin master to the local master branch
		go Until(as.checkForRemoteUpdates, time.Duration(*as.cfg.frequency)*time.Second, as.stopCh)
	}

	for {
		time.Sleep(3600 * time.Second)
	}
}
