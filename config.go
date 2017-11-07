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
	server     *string
	account    *string
	repository *string
	directory  *string
	username   *string
	password   *string
	email      *string
	version    *bool
	initialize *bool
	initonly   *bool
	frequency  *int
	direction  *string
}

func newConfig() *config {
	return &config{
		server:     flag.String("server", GitDefaultServer, "git repository host"),
		account:    flag.String("account", GitDefaultAccount, "git account/owner/organization for repository to clone"),
		repository: flag.String("repository", GitDefaultRepo, "git repository to manage for archiving local updates"),
		directory:  flag.String("directory", "", "Required: The name of a new / existing repository directory to clone into / work in"),
		username:   flag.String("username", GitDefaultUser, "git remote login username"),
		password:   flag.String("password", "", "Required: git remote login password"),
		email:      flag.String("email", GitDefaultEmail, "Required: git user's email address"),
		version:    flag.Bool("version", false, "display version info and exit"),
		initialize: flag.Bool("initial-clone", true, "initialize the state of the repository by cloning the remote"),
		initonly:   flag.Bool("init-only", false, "initialize the state of the repository only, then exit"),
		frequency:  flag.Int("sync-interval", 60, "number of seconds between upstream sync's when changes are present"),
		direction:  flag.String("sync-to", "remote", "must be one of `remote`, `local`, `both`.  Sync changes to/from: `remote` push local commits to remote, `local` pull remote changes to local, `both` manage bi-directional updates."),
	}
}

func (cfg *config) String() string {
	return fmt.Sprintf("server: %s, account: %s, repository: %s, directory: %s, "+
		"username: %s, password: %s, email: %s, frequency: %d, initialize: %t, "+
		"initonly: %t, sync-to: %s, version: %t",
		*cfg.server, *cfg.account, *cfg.repository, *cfg.directory,
		*cfg.username, *cfg.password, *cfg.email, *cfg.frequency, *cfg.initialize,
		*cfg.initonly, *cfg.direction, *cfg.version)
}

var envSupport = map[string]bool{
	"server":     true,
	"account":    true,
	"repository": true,
	"directory":  true,
	"username":   true,
	"password":   true,
	"email":      true,
	"version":    false,
	"initialize": true,
	"initonly":   true,
	"frequency":  true,
	"direction":  true,
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

func (cfg *config) validate() bool {
	if *cfg.directory == "" {
		glog.Error("Configuraiton flag: `--directory` can not be empty, a valid value is required.")
		return false
	}
	if *cfg.password == "" {
		glog.Error("Configuraiton flag: `--password` can not be empty, a valid value is required.")
		return false
	}
	if *cfg.initialize == false && *cfg.initonly == true {
		glog.Error("Configuraiton flags: `--initialize=false` and `--init-only=true` conflict.")
		return false
	}
	if *cfg.direction == "" ||
		(*cfg.direction != "remote" && *cfg.direction != "local" && *cfg.direction != "both") {
		glog.Error("Configuraiton flag: `--direction` must be one of: `remote`, `local`, or `both`")
		return false
	}
	return true
}
