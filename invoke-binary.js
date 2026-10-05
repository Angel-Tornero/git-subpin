// invoke-binary.js runs the platform-specific git-subpin binary bundled on the
// release branch. The action pretends to be a Node action; this shim maps the
// Node runtime to the matching GOOS/GOARCH build, fails clearly on unsupported
// platforms, and forwards arguments, environment, stdio and exit code.
const { spawnSync } = require("node:child_process");
const path = require("node:path");

// Pinned by the release workflow to the main commit that produced the bundled
// binaries, so the shim always runs the binaries it shipped with.
const SHA = "d58ae888aee3cebe79f276db27fd449b1b8add0e";

const platformToGOOS = {
  linux: "linux",
  darwin: "darwin",
  win32: "windows",
};

const archToGOARCH = {
  x64: "amd64",
  arm64: "arm64",
};

// The GOOS/GOARCH pairs built and committed to the release branch.
const supported = ["linux/amd64", "linux/arm64", "darwin/arm64", "windows/amd64"];

function chooseBinary() {
  const goos = platformToGOOS[process.platform];
  const goarch = archToGOARCH[process.arch];
  const target = `${goos}/${goarch}`;

  if (!goos || !goarch || !supported.includes(target)) {
    const detected = `${process.platform}/${process.arch}`;
    console.error(
      `git-subpin: unsupported platform ${detected}; supported: ${supported.join(", ")}`,
    );
    process.exit(1);
  }

  const ext = goos === "windows" ? ".exe" : "";
  return path.join(__dirname, `main-${goos}-${goarch}-${SHA}${ext}`);
}

const binary = chooseBinary();
const result = spawnSync(binary, process.argv.slice(2), {
  stdio: "inherit",
  env: process.env,
});

if (result.error) {
  console.error(`git-subpin: failed to run ${binary}: ${result.error.message}`);
  process.exit(1);
}

// status is null when the binary was killed by a signal; treat that as failure.
process.exit(result.status === null ? 1 : result.status);
