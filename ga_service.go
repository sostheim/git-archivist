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

import "time"

var (
	resyncPeriod = 30 * time.Second
)

// GA Server
type gaServer struct {
	// Command line / environment supplied configuration values
	cfg    *config
	stopCh chan struct{}
}

func newGAServer(cfg *config) *gaServer {
	return &gaServer{
		stopCh: make(chan struct{}),
		cfg:    cfg,
	}
}

func (as *gaServer) run() {
	// run the controller and queue goroutines
	// go as.apiServer.Run(as.stopCh)
	// Allow time for the initial startup
	time.Sleep(5 * time.Second)
}
