package main

import (
	goflag "flag"

	flag "github.com/spf13/pflag"

	"testing"
)

func validateBoolFlag(flag string, value bool, b *bool, t *testing.T) bool {
	if b == nil {
		t.Errorf("validateBoolFlag(%s): want object, have nil", flag)
		return false
	} else if *b != value {
		t.Errorf("validateBoolFlag(%s): want %t, have %t", flag, value, *b)
		return false
	}
	return true
}

func validateStringFlag(flag, value string, s *string, t *testing.T) bool {
	if s == nil {
		t.Errorf("validateStringFlag(%s): want object, have nil", flag)
		return false
	} else if *s != value {
		t.Errorf("validateStringFlag(%s): want %s, have %s", flag, value, *s)
		return false
	}
	return true
}

func validateIntFlag(flag string, value int, i *int, t *testing.T) bool {
	if i == nil {
		t.Errorf("validateIntFlag(%s): want object, have nil", flag)
		return false
	} else if *i != value {
		t.Errorf("validateIntFlag(%s): want %d, have %d", flag, value, *i)
		return false
	}
	return true
}

func TestNewConfig(t *testing.T) {
	cfg := gaCfg
	f := flag.NewFlagSet("test", flag.ContinueOnError)
	f.AddGoFlagSet(goflag.CommandLine)

	if !validateStringFlag("server", GitDefaultServer, cfg.server, t) {
		t.Error("TestNewConfig() want valid server")
	}
	if !validateStringFlag("account", GitDefaultAccount, cfg.account, t) {
		t.Error("TestNewConfig() want valid account")
	}
	if !validateStringFlag("repository", GitDefaultRepo, cfg.repository, t) {
		t.Error("TestNewConfig() want valid repository")
	}
	if !validateStringFlag("email", GitDefaultEmail, cfg.email, t) {
		t.Error("TestNewConfig() want valid email")
	}
	if !validateStringFlag("password", "", cfg.password, t) {
		t.Error("TestNewConfig() want valid password")
	}
	if !validateStringFlag("directory", "", cfg.directory, t) {
		t.Error("TestNewConfig() want valid directory")
	}
	if !validateBoolFlag("version", false, cfg.version, t) {
		t.Error("TestConfigOverrides() want valid version")
	}
	if !validateBoolFlag("initial-clone", true, cfg.initialize, t) {
		t.Error("TestNewConfig() want valid initial-clone")
	}
	if !validateBoolFlag("init-only", false, cfg.initonly, t) {
		t.Error("TestNewConfig() want valid init-only")
	}
	if !validateIntFlag("frequency", 60, cfg.frequency, t) {
		t.Error("TestNewConfig() want valid sync-internval")
	}
	if !validateStringFlag("sync-to", "remote", cfg.direction, t) {
		t.Error("TestNewConfig() want valid sync-to")
	}
}

func TestConfigOverrides(t *testing.T) {
	cfg := gaCfg
	f := flag.NewFlagSet("test", flag.ContinueOnError)
	f.AddGoFlagSet(goflag.CommandLine)

	flag.Set("directory", "/Users/sostheim/work/source/github.com/git-archivist")
	flag.Set("password", "redacted")
	flag.Set("sync-interval", "3600")
	flag.Set("sync-to", "both")
	flag.Parse()

	if !validateStringFlag("server", GitDefaultServer, cfg.server, t) {
		t.Error("TestNewConfig() want valid server")
	}
	if !validateStringFlag("account", GitDefaultAccount, cfg.account, t) {
		t.Error("TestNewConfig() want valid account")
	}
	if !validateStringFlag("repository", GitDefaultRepo, cfg.repository, t) {
		t.Error("TestNewConfig() want valid repository")
	}
	if !validateStringFlag("email", GitDefaultEmail, cfg.email, t) {
		t.Error("TestNewConfig() want valid email")
	}
	if !validateStringFlag("password", "redacted", cfg.password, t) {
		t.Error("TestNewConfig() want valid password")
	}
	if !validateStringFlag("directory", "/Users/sostheim/work/source/github.com/git-archivist", cfg.directory, t) {
		t.Error("TestNewConfig() want valid directory")
	}
	if !validateBoolFlag("version", false, cfg.version, t) {
		t.Error("TestConfigOverrides() want valid version")
	}
	if !validateBoolFlag("initial-clone", true, cfg.initialize, t) {
		t.Error("TestNewConfig() want valid initial-clone")
	}
	if !validateBoolFlag("init-only", false, cfg.initonly, t) {
		t.Error("TestNewConfig() want valid init-only")
	}
	if !validateIntFlag("frequency", 3600, cfg.frequency, t) {
		t.Error("TestNewConfig() want valid sync-internval")
	}
	if !validateStringFlag("sync-to", "both", cfg.direction, t) {
		t.Error("TestNewConfig() want valid sync-to")
	}
}
