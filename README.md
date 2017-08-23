# Git Archivist

## Running

The git-archivist application has a number of command line options that define how it operates.
```
$ ./git-archivist --help
Usage of ./git-archivist:
      --account string                   git account/owner/organization for repository to clone (default "samsung-cnct")
      --alsologtostderr                  log to standard error as well as files
      --directory string                 Required: The name of a new / existing repository directory to clone into / work in
      --email string                     Required: git user's email address (default "cnct.api.robot@gmail.com")
      --initial-clone                    initialize the state of the repository by cloning the remote (default true)
      --log_backtrace_at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                   If non-empty, write log files in this directory
      --logtostderr                      log to standard error instead of files
      --password string                  Required: git remote login password
      --repository string                git repository to manage for archiving local updates (default "cluster-manifests")
      --server string                    git repository host (default "github.com")
      --stderrthreshold severity         logs at or above this threshold go to stderr (default 2)
      --sync-interval int                number of seconds between upstream sync's when changes are present (default 300)
      --username string                  git remote login username (default "api-robot")
  -v, --v Level                          log level for V logs
      --version                          display version info and exit
      --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging
```
**Note 1** Most of the arguments have default values that are probably not useful in general, but just for their default application.

**Noe 2** There are two non-defaulted required parameters that must be supplied by the user, `--password` and `--directory`.  If either of these values are not specified the application will fail to start and print an appropriate error message.

### Environment Variables
The git-archivist application is configurable through command line configuration flags, and through a subset of environment variables. Any configuration value set on the command line takes precedence over the same value from the environment.

The format of the environment variable for flag for flag is composed of the prefix `GA_` and the remaining text of the flag in all uppercase with all hyphens replaced by underscores.  Fore example, `--example-flag` would map to `GA_EXAMPLE_FLAG`. 

Not every flag can be set via an environment variable.  This is due to the fact that the total set of flags supported by the application is an aggregate of those that belong to git-archivist and 3rd party Go packages.  The set of flags that do have corresponding environment variable support are listed below:
* --account
* --directory
* --initial-clone
* --password
* --repository 
* --server
* --sync-interval
* --username

### Example Invocation

The following example show how to 1) clone an existing repository (the default private repository `cluster-manifests`), to the desired location.  From there the default user `api-robot` with the robot's default email address `cnct.api.robot...` is used for all `git clone`, `git commit`, and `git push` operations.

```
$ git-archivist --v=4 --alsologtostderr --password **redacted** --directory /Users/sostheim/clusters
```
