# Changelog

We document all notable changes to this project in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.5] - 2022-03-29

### Fixed

* Adds logging parameters to the log function as an array instead of a map in
  the logging middleware handler. This will address inconsistent parameter
  ordering in the log output for HTTP handlers.

## [0.1.4] - 2021-01-06

### Build

* Updated dependency on Leg lifecycle package.

## [0.1.3] - 2020-12-04

### Fixed

* Removed unnecessary `replace` directives in Go module definition.

## [0.1.2] - 2020-12-04

### Fixed

* Removed unnecessary `replace` directives in Go module definition.

## [0.1.1] - 2020-12-04

### Fixed

* Removed unnecessary `replace` directives in Go module definition.

## [0.1.0] - 2020-12-04

### Changed

* Renamed project to Leg.

[Unreleased]: https://github.com/puppetlabs/leg/compare/httputil/v0.1.5...HEAD
[0.1.5]: https://github.com/puppetlabs/leg/compare/httputil/v0.1.4...httputil/v0.1.5
[0.1.4]: https://github.com/puppetlabs/leg/compare/httputil/v0.1.3...httputil/v0.1.4
[0.1.3]: https://github.com/puppetlabs/leg/compare/httputil/v0.1.2...httputil/v0.1.3
[0.1.2]: https://github.com/puppetlabs/leg/compare/httputil/v0.1.1...httputil/v0.1.2
[0.1.1]: https://github.com/puppetlabs/leg/compare/httputil/v0.1.0...httputil/v0.1.1
[0.1.0]: https://github.com/puppetlabs/leg/compare/d290e8e835c3fa3ea4e93073bfe19e1958493d47...httputil/v0.1.0
