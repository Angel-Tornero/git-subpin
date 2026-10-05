# git-subpin

[![CI](https://github.com/Angel-Tornero/git-subpin/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/Angel-Tornero/git-subpin/actions/workflows/ci.yml)

A GitHub Action that resolves **Git submodule pointer conflicts** left by an
automated merge between environment branches, following a declarative policy.

When you merge one branch into another (for example `dev` into `prod`) and both
sides moved a submodule to a different commit, Git leaves the submodule gitlink
in conflict. `git-subpin` reads a policy file, decides which pointer to keep for
each conflicted submodule, and stages the result so the merge can be committed.

## What it does

Given a merge that has already produced conflicts in the working tree, the
action:

1. Reads the conflicted submodule pointers from `git status --porcelain=v2`.
2. Maps each conflicted path to its submodule name via `.gitmodules`.
3. Matches the name (and the source/target branches) against the policy rules,
   picking the most specific rule.
4. Applies the rule's strategy and **stages** the chosen pointer.

It does **not** run the merge or create the merge commit — your workflow does
that around it. If it cannot resolve a conflict, it exits non-zero so a human
can step in.

## Usage

```yaml
jobs:
  promote:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
          submodules: true

      - name: Attempt the merge
        run: |
          git config user.name "ci"
          git config user.email "ci@example.com"
          git checkout prod
          git merge --no-commit --no-ff dev || true

      - name: Resolve submodule pointers
        uses: Angel-Tornero/git-subpin@v1
        with:
          from: dev
          to: prod
          config: .subpin.yml

      - name: Commit the merge
        run: git commit --no-edit
```

### Inputs

| Input    | Required | Default       | Description                                        |
| -------- | -------- | ------------- | -------------------------------------------------- |
| `from`   | yes      | —             | Source branch of the merge (the branch merged in). |
| `to`     | yes      | —             | Target branch of the merge (merged into).          |
| `config` | no       | `.subpin.yml` | Path to the policy file.                           |

### Outputs

None in v1.

## Policy file

```yaml
version: "1"
rules:
  - from: "*"
    to: "*"
    submodule: "*"
    strategy: keep-source
    description: default to the pointer from the branch being merged in
  - from: dev
    to: prod
    submodule: keycloak
    strategy: keep-source
```

- `from`, `to`, `submodule` are glob-style patterns (`*` matches any run of
  characters). `submodule` matches the submodule **name** as declared in
  `.gitmodules`, not its path.
- When several rules match, the most specific one wins (a literal beats a
  partial wildcard, which beats `*`). If two equally specific rules disagree on
  the strategy, resolution stops with an "ambiguous rules" error.

## Behavior and limitations (v1)

This is a deliberately small first version.

- **Only `keep-source` is applied.** It keeps the pointer from the source branch
  (the branch being merged in, recorded by Git as "theirs").
- `keep-target` and `manual` are valid in a policy but **not yet supported**: if
  a matched rule selects them, the action stops with a clear error.
- `fast-forward` is **not a valid strategy** and is rejected when the policy is
  parsed.
- Only **submodule pointer** conflicts are handled. Any other conflict is left
  untouched for you or the rest of your workflow to deal with.
- A conflicted submodule with no matching rule, an ambiguous rule set, or a path
  missing from `.gitmodules` stops resolution with a non-zero exit.

## Development

```sh
go test ./...          # unit tests + a real-git integration test
go test -short ./...   # skip the integration test (no git required)
go build ./...
go vet ./...
gofmt -l .
staticcheck ./...      # go install honnef.co/go/tools/cmd/staticcheck@latest
```

The action is distributed following the "GitHub Actions in Go" pattern: source
lives on `main`, and CI publishes a `release` branch holding `action.yml`, a
small Node shim (`invoke-binary.js`) and static binaries for linux/amd64,
linux/arm64, darwin/arm64 and windows/amd64. The shim picks the binary matching
the runner, so no Go toolchain is needed at action run time.

## License

[MIT](LICENSE).
