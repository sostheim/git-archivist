# Git Archivist

## Running

The git-archivist application has a number of command line options that define how it operates.
```
$ ./git-archivist --help
Usage of ./git-archivist:
      --account string                   git account/owner/organization for repository to clone (default "samsung-cnct")
      --alsologtostderr                  log to standard error as well as files
      --directory string                 Required: The name of a new / existing repository directory to clone into / work in
      --initial-clone                    intialize the state of the repository by cloning the remote (default true)
      --log_backtrace_at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                   If non-empty, write log files in this directory
      --logtostderr                      log to standard error instead of files
      --password string                  Required: git remote login password
      --repository string                git repository to manage for archiving local updates (default "cluster-manifests")
      --server string                    git repository host (default "github.com")
      --stderrthreshold severity         logs at or above this threshold go to stderr (default 2)
      --sync-interval int                number of seconds between upstream sync's when changes are present (default 300)
      --username string                  git remote login userername (default "api-robot")
  -v, --v Level                          log level for V logs
      --version                          display version info and exit
      --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging
```

### Environment Variables
The git-archivist application is configurable through command line configuration flags, and through a subset of environment variables. Any configuration value set on the command line takes precedence over the same value from the environment.

The format of the environment variable for flag for flag is composed of the prefix `GA_` and the reamining text of the flag in all uppper case with all hyphens replaced by underscores.  Fore example, `--example-flag` would map to `GA_EXAMPLE_FLAG`. 

Not every flag can be set via an environment variable.  This is due to the fact that the total set of flags supported by the application is an aggregate of those that belong to git-archivist and 3rd party Go packages.  The set of flags that do have corresponding environment variable support are listed below:
* --account
* --directory
* --initial-clone
* --password
* --repository 
* --server
* --sync-interval
* --username
