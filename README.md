[![Go Report Card](https://goreportcard.com/badge/github.com/samsung-cnct/git-archivist)](https://goreportcard.com/report/github.com/samsung-cnct/git-archivist)
[![Docker Repository on Quay](https://quay.io/repository/samsung_cnct/git-archivist/status "Docker Repository on Quay")](https://quay.io/repository/samsung_cnct/git-archivist)
[![maturity](https://img.shields.io/badge/status-alpha-red.svg)](https://github.com/github.com/samsung-cnct/git-archivist)

# Git Archivist

Automation for commiting git repository tracked file updates to an upstream master repository on a `--sync-interval`

## Running

The git-archivist application has a number of command line options that define how it operates.
```
$ ./git-archivist --help
Usage of ./git-archivist:
      --account string                   git account/owner/organization for repository to clone (default "samsung-cnct")
      --alsologtostderr                  log to standard error as well as files
      --directory string                 Required: The name of a new / existing repository directory to clone into / work in
      --email string                     Required: git user's email address (default "cnct.api.robot@gmail.com")
      --init-only                        initialize the state of the repository only, then exit
      --initial-clone                    initialize the state of the repository by cloning the remote (default true)
      --log_backtrace_at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                   If non-empty, write log files in this directory
      --logtostderr                      log to standard error instead of files
      --password string                  Required: git remote login password
      --repository string                git repository to manage for archiving local updates (default "cluster-manifests")
      --server string                    git repository host (default "github.com")
      --stderrthreshold severity         logs at or above this threshold go to stderr (default 2)
      --sync-interval int                number of seconds between upstream sync's when changes are present (default 60)
      --username string                  git remote login username (default "api-robot")
  -v, --v Level                          log level for V logs
      --version                          display version info and exit
      --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging
```
**Note 1** Most of the flags have default values.  The defaults are probably not useful in general, but just for their default application deployment.  In particular, the flags `--account` (default value "samsung-cnct"), and `--username` (default value "api-robot"), will need to be overridden for all but Samsung CNCT operations.  

**Note 2** There are two required parameters that must be supplied by the user, `--password` and `--directory`.  If either of these values are not specified the application will fail to start and print an appropriate error message.

**Note 3** It is an error to set `--initial-clone=false` and `--init-only=true`

### Arguments
A closer look at some of the application flags:
- **--account**: The github account/owner/organization that own's the repository that will, optionally be cloned, used for updates
- **--username**: A real, valid, credentialed account, that has at a minimum "repo" access.
- **--initialize**: default value: true.  If `--initialize=false`, then then the archivist expects there to be an existing valid git repository found at `--directory`.  Otherwise, the default behavior is to clone the repository in to the working `--directory`
- **--init-only**: default value: false.  If `--init-only=true`, then then the archivist will perform the inital clone of the specified repository and exit upon completion. Otherwise, the default behavior is to start executing the timed syncrhonization process.
- **--sync-interval**: default value: 60s.  The frequency with which working `--directory` is checked for tracked file updates.

### Environment Variables
The git-archivist application is configurable through command line configuration flags, and through a subset of environment variables. Any configuration value set on the command line takes precedence over the same value from the environment.

The format of the environment variable for flag is composed of the prefix `GA_` and the remaining text of the flag in all uppercase with all hyphens replaced by underscores.  Fore example, `--example-flag` would map to `GA_EXAMPLE_FLAG`. 

Not every flag can be set via an environment variable.  This is due to the fact that the total set of flags supported by the application is an aggregate of those that belong to git-archivist and 3rd party Go packages.  The set of flags that do have corresponding environment variable support are listed below:
* --server
* --account
* --username
* --password
* --repository 
* --directory
* --init-only
* --initial-clone
* --sync-interval
* --sync-to

## Example Invocation

### Simple Default Example

The following example show how to 1) clone an existing repository (the default private repository `cluster-manifests`), to the desired location.  From there the default user `api-robot` with the robot's default email address `cnct.api.robot...` is used for all `git clone`, `git commit`, and `git push` operations.

```
$ git-archivist --v=4 --alsologtostderr --password **redacted** --directory /Users/sostheim/clusters
```

### As an Init Container and Sidecar

The following example shows how to use git-archivist as an init container to ensure that the git repository is cloned before the same executeable begins life as a sidecar.  This example again uses default values for the user `api-robot` and default email address `cnct.api.robot...`

First, as an init container clone an existing repository (the default private repository `cluster-manifests`), to the desired location and exit.
```
$ git-archivist --v=4 --alsologtostderr=true --password=**redacted** --directory=/.kraken/ --init-only=true
```

Second, in the long running container, the `--init-only` assumes it's default value: false, but we must now set `--initial-clone` to false. 
```
$ git-archivist --v=4 --alsologtostderr=true --password=**redacted** --directory=/.kraken/ --initial-clone=false
```

#### Example Deployment Container Mainfest for Init Container and Sidecar
```
. . . Deployment Details Omitted ... 

 spec:
  containers:
  - name: git-archivist
    image: quay.io/samsung_cnct/git-archivist:latest
    imagePullPolicy: Always
    args:
    - --v=2
    - --logtostderr=true
    - --directory=/root/.kraken/
    - --initial-clone=false
    env:
    - name: GA_ACCOUNT
      valueFrom:
        secretKeyRef:
          key: account
          name: git-archivist-secret
    - name: GA_PASSWORD
      valueFrom:
        secretKeyRef:
          key: password
          name: git-archivist-secret
    - name: GA_EMAIL
      valueFrom:
        secretKeyRef:
          key: email
          name: git-archivist-secret
    - name: GA_REPOSITORY
      valueFrom:
        secretKeyRef:
          key: repository
          name: git-archivist-secret
    - name: GA_SERVER
      valueFrom:
        secretKeyRef:
          key: server
          name: git-archivist-secret
    - name: GA_USERNAME
      valueFrom:
        secretKeyRef:
          key: username
          name: git-archivist-secret
    resources:
      limits:
        cpu: 200m
        memory: 128Mi
      requests:
        cpu: 50m
        memory: 128Mi
  initContainers:
  - name: git-archivist-init
    image: quay.io/samsung_cnct/git-archivist:latest
    imagePullPolicy: Always
    args:
    - --v=4
    - --logtostderr=true
    - --directory=/.kraken/
    - --init-only=true
    env:
    - name: GA_ACCOUNT
      valueFrom:
        secretKeyRef:
          key: account
          name: git-archivist-secret
    - name: GA_EMAIL
      valueFrom:
        secretKeyRef:
          key: email
          name: git-archivist-secret
    - name: GA_REPOSITORY
      valueFrom:
        secretKeyRef:
          key: repository
          name: git-archivist-secret
    - name: GA_SERVER
      valueFrom:
        secretKeyRef:
          key: server
          name: git-archivist-secret
    - name: GA_USERNAME
      valueFrom:
        secretKeyRef:
          key: username
          name: git-archivist-secret
    - name: GA_PASSWORD
      valueFrom:
        secretKeyRef:
          key: password
          name: git-archivist-secret
```

$ dep init -no-examples
  Using ^3.5.1 as constraint for direct dep github.com/blang/semver
  Locking in v3.5.1 (2ee8785) for direct dep github.com/blang/semver
  Using master as constraint for direct dep github.com/golang/glog
  Locking in master (23def4e) for direct dep github.com/golang/glog
  Using ^1.0.0 as constraint for direct dep github.com/spf13/pflag
  Locking in v1.0.0 (e57e3ee) for direct dep github.com/spf13/pflag
## Building the Project
This project provides a `makefile` that includcs the basic functionality to build and package the project.

### Dependencies
A note on Go dependencies for this project.  The current dependency management strategy is to use the `dep` tool.  For more information on `dep` see the [GitHub project page readme](https://github.com/golang/dep).

So, why not `glide`?, you ask - we'll according to the [`glide` project readme](https://github.com/Masterminds/glide#golang-dep):
> The Go community now has the dep project to manage dependencies. Please consider trying to migrate from Glide to dep. If there is an issue preventing you from migrating please file an issue with dep so the problem can be corrected. Glide will continue to be supported for some time but is considered to be in a state of support rather than active feature development.

This *is not* a guide to using `dep`, refere to project page for advanced help.  You really should at least read the following sections to understand more about managing a projects dependencies: 
* [Adding a Dependency](https://github.com/golang/dep#adding-a-dependency)
* [Changing Dependencies](https://github.com/golang/dep#changing-dependencies)
* [Checking the Status of Dependencies](https://github.com/golang/dep#checking-the-status-of-dependencies)

For the sake of this maintaining this project's dependencies, please follow these steps to add a new package reference to the project.
1. Run the command: `$ dep ensure -add github.com/foo/bar`
1. Now, import and use the package in your code!

To modify an existing dependency, please follow these steps:
1. Manually edit the file: `Gopkg.toml`
1. Run the command: `$ dep ensure`

### Build and Install the Project
To build just the executeable locally:
```
$ make build
```
To build the executelable image locally and install it on the local system:
```
$ make install
```
### Build the Container
Build the container image (resides in local images only):
```
$ make conatiner
 
  . . . 

$ docker images
REPOSITORY                             TAG      IMAGE ID        CREATED           SIZE
quay.io/samsung_cnct/git-archivist     0.2.0    7e9416dd3acb    15 minutes ago    28.9MB

```
### Build and Push the Container
Build the container image and push the image to [Quay.io git-archivist repository](https://quay.io/repository/samsung_cnct/git-archivist):
```
$ make push
```

### Cutting a Release

1. Install `github-release` from https://github.com/c4milo/github-release
1. Create a github personal access token with repo read/write permissions and export it as GITHUB_TOKEN  
1. Adjust VERSION and TYPE variables in the [Makefile](Makefile) as needed  
1. Generate the cross compiled sources, tagged, inventoried, and make a release in GitHub:
```
make release
```
