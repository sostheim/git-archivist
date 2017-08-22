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
	"fmt"
	"os"
	"strings"

	"github.com/golang/glog"
	flag "github.com/spf13/pflag"
)

type config struct {
	flagSet    *flag.FlagSet
	repository *string
	username   *string
	password   *string
	version    *bool
	frequency  *int
}

func newConfig() *config {
	return &config{
		repository: flag.String("repository", "", "git repository to manage for archiving local updates"),
		username:   flag.String("username", "", "git remote userername"),
		password:   flag.String("password", "", "git remote user's password"),
		version:    flag.Bool("version", false, "display version info and exit"),
		frequency:  flag.Int("sync-interval", 300, "number of seconds between upstream sync's when changes are present"),
	}
}

func (cfg *config) String() string {
	return fmt.Sprintf("repository: %s, username: %s, password: %s, frequency: %d, version: %t",
		*cfg.repository, *cfg.username, *cfg.password, *cfg.frequency, *cfg.version)
}

var envSupport = map[string]bool{
	"repository": true,
	"username":   true,
	"password":   true,
	"version":    false,
	"frequency":  true,
}

func variableName(name string) string {
	return "GA_" + strings.ToUpper(strings.Replace(name, "-", "_", -1))
}

// Just like Flags.Parse() except we try to get recognized values for the valid
// set of flags from environment variables.  We choose to use the environment
// value if 1) the value hasen't already been set as command line flags and the
// flas is a member of the supported set (see map defined above).
func (cfg *config) envParse() error {
	var err error

	alreadySet := make(map[string]bool)
	cfg.flagSet.Visit(func(f *flag.Flag) {
		if envSupport[f.Name] {
			alreadySet[variableName(f.Name)] = true
		}
	})

	usedEnvKey := make(map[string]bool)
	cfg.flagSet.VisitAll(func(f *flag.Flag) {
		if envSupport[f.Name] {
			key := variableName(f.Name)
			if !alreadySet[key] {
				val := os.Getenv(key)
				if val != "" {
					usedEnvKey[key] = true
					if serr := cfg.flagSet.Set(f.Name, val); serr != nil {
						err = fmt.Errorf("invalid value %q for %s: %v", val, key, serr)
					}
					glog.V(3).Infof("recognized and used environment variable %s=%s", key, val)
				}
			}
		}
	})

	for _, env := range os.Environ() {
		kv := strings.SplitN(env, "=", 2)
		if len(kv) != 2 {
			glog.Warningf("found invalid env %s", env)
		}
		if usedEnvKey[kv[0]] {
			continue
		}
		if alreadySet[kv[0]] {
			glog.V(3).Infof("recognized environment variable %s, but unused: superseeded by command line flag ", kv[0])
			continue
		}
		if strings.HasPrefix(env, "GA_") {
			glog.Warningf("unrecognized environment variable %s", env)
		}
	}

	return err
}