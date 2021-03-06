
<a name="0.0.4"></a>
## [0.0.4](https://github.com/prologic/shops/compare/0.0.3...0.0.4) (2021-03-06)

### Features

* Add display of ascii poo on non-zero target errors and exit with exit status 3
* Add uptime to testdata/ping.yml spec


<a name="0.0.3"></a>
## [0.0.3](https://github.com/prologic/shops/compare/0.0.2...0.0.3) (2021-03-06)

### Bug Fixes

* Fix SSHRunner to capture both stdout/stderr
* Fix bad call to log.Debugf()
* Fix bug with exitStatus.Error() and improve error hanadling when actions fail

### Features

* Add toc to README
* Add example of ensuring and installing node_exporter on Linux hosts
* Add additinoal debug logging until we iron out subtle bugs/issues
* Add docs on authentication and reference [#9](https://github.com/prologic/shops/issues/9)
* Add star button to README

### Updates

* Update CHANGELOG for 0.0.3
* Update README.md
* Update README
* Update README.md


<a name="0.0.2"></a>
## [0.0.2](https://github.com/prologic/shops/compare/0.0.1...0.0.2) (2021-03-05)

### Bug Fixes

* Fix error handling for copying files/dirs
* Fix exit error/status handling
* Fix versioning the binaries
* Fix bugs in local and ssh runners

### Features

* Add another test case for ssh://host:port
* Add better help on targets ([#8](https://github.com/prologic/shops/issues/8))
* Add support for functions ([#7](https://github.com/prologic/shops/issues/7))
* Add GHA workflows for CI
* Add -c/--continue-on-error option with default fail fast ([#5](https://github.com/prologic/shops/issues/5))

### Updates

* Update CHANGELOG for 0.0.2
* Update README with latest demo example


<a name="0.0.1"></a>
## 0.0.1 (2021-03-01)

### Bug Fixes

* Fix GoReleeaser config (dependent library does not support freebsd :/)
* Fix GoReleaser cofig
* Fix race condition in GroupRunner
* Fix printing last result
* Fix arg handling
* Fix typos

### Features

* Add /dist to .gitignore
* Add empty CHANGELOG
* Add support for local runner (Closes [#2](https://github.com/prologic/shops/issues/2))
* Add Makefile and release tools
* Add todo list
* Add support for copying files and directories
* Add config structs

### Updates

* Update CHANGELOG for 0.0.1
* Update CHANGELOG for 0.0.1
* Update TODO
* Update sample config
* Update module path

