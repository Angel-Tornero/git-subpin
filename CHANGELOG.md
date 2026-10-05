# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2026-10-05

### Added

- Resolve submodule pointer conflicts with the `keep-source` strategy, staging
  the source-side pointer so the merge can be committed.
- Declarative YAML policy: wildcard matching on source branch, target branch and
  submodule name, with specificity-based rule selection and ambiguity detection.
- Rule strategies `keep-source`, `keep-target` and `manual` are recognized;
  `keep-target` and `manual` are reported as unsupported in this version, and
  `fast-forward` is rejected when the policy is parsed.
- Distribution as a GitHub Action: a Node 24 shim dispatches to static Go
  binaries built for linux/amd64, linux/arm64, darwin/arm64 and windows/amd64,
  published to a `release` branch by CI.

[Unreleased]: https://github.com/Angel-Tornero/git-subpin/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/Angel-Tornero/git-subpin/releases/tag/v1.0.0
