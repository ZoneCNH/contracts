package scripts_test

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestReleasePreflightRejectsMissingOrInvalidVersion(t *testing.T) {
	for _, tc := range []struct {
		name      string
		version   string
		wantError string
	}{
		{
			name:      "missing version",
			version:   "",
			wantError: "ERROR: set VERSION=vX.Y.Z when running release preflight",
		},
		{
			name:      "invalid version",
			version:   "release-1",
			wantError: "ERROR: release version must look like vX.Y.Z: release-1",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result := runReleasePreflight(t, releasePreflightFixture{
				version: tc.version,
			})

			if result.code == 0 {
				t.Fatalf("expected release preflight to fail, got exit 0\noutput:\n%s", result.output)
			}
			if !strings.Contains(result.output, tc.wantError) {
				t.Fatalf("expected output to contain %q\noutput:\n%s", tc.wantError, result.output)
			}
		})
	}
}

func TestReleasePreflightRejectsWrongBranchAndDirtyTree(t *testing.T) {
	for _, tc := range []struct {
		name      string
		setup     func(*testing.T, string)
		wantError string
	}{
		{
			name: "wrong branch",
			setup: func(t *testing.T, repo string) {
				runGit(t, repo, "switch", "-c", "feature")
			},
			wantError: "ERROR: release preflight must run on main; current branch is feature",
		},
		{
			name: "dirty tree",
			setup: func(t *testing.T, repo string) {
				writeFixtureFile(t, repo, "dirty.txt")
			},
			wantError: "ERROR: release preflight requires a clean git worktree",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := newReleasePreflightRepo(t, false, false)
			tc.setup(t, repo)

			result := runReleasePreflightWithRepo(t, repo, "v1.2.3")
			if result.code == 0 {
				t.Fatalf("expected release preflight to fail, got exit 0\noutput:\n%s", result.output)
			}
			if !strings.Contains(result.output, tc.wantError) {
				t.Fatalf("expected output to contain %q\noutput:\n%s", tc.wantError, result.output)
			}
		})
	}
}

func TestReleasePreflightRejectsMissingChangelogHeading(t *testing.T) {
	repo := newReleasePreflightRepo(t, true, false)

	result := runReleasePreflightWithRepo(t, repo, "v1.2.3")
	if result.code == 0 {
		t.Fatalf("expected release preflight to fail, got exit 0\noutput:\n%s", result.output)
	}

	want := "ERROR: CHANGELOG.md must contain a release heading for v1.2.3"
	if !strings.Contains(result.output, want) {
		t.Fatalf("expected output to contain %q\noutput:\n%s", want, result.output)
	}
}

func TestReleasePreflightPassesWithRequiredMetadata(t *testing.T) {
	repo := newReleasePreflightRepo(t, true, true)
	binDir := filepath.Join(t.TempDir(), "bin")
	mkdirAll(t, binDir)
	fakeTool(t, binDir, "golangci-lint")
	fakeTool(t, binDir, "govulncheck")

	cmd := exec.Command("bash", filepath.Join("scripts", "check_release_preflight.sh"))
	cmd.Dir = repo
	cmd.Env = append(os.Environ(),
		"VERSION=v1.2.3",
		"PATH="+binDir+string(os.PathListSeparator)+os.Getenv("PATH"),
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("release preflight failed: %v\noutput:\n%s", err, output)
	}

	want := "release preflight metadata checks passed for v1.2.3"
	if !strings.Contains(string(output), want) {
		t.Fatalf("expected output to contain %q\noutput:\n%s", want, output)
	}
}

type releasePreflightFixture struct {
	version string
}

type releasePreflightRunResult struct {
	code   int
	output string
}

func runReleasePreflight(t *testing.T, fixture releasePreflightFixture) releasePreflightRunResult {
	t.Helper()

	repo := newReleasePreflightRepo(t, false, false)
	return runReleasePreflightWithRepoAndVersion(t, repo, fixture.version)
}

func runReleasePreflightWithRepo(t *testing.T, repo, version string) releasePreflightRunResult {
	t.Helper()

	return runReleasePreflightWithRepoAndVersion(t, repo, version)
}

func runReleasePreflightWithRepoAndVersion(t *testing.T, repo, version string) releasePreflightRunResult {
	t.Helper()

	cmd := exec.Command("bash", filepath.Join("scripts", "check_release_preflight.sh"))
	cmd.Dir = repo
	cmd.Env = append(os.Environ(), "VERSION="+version)

	output, err := cmd.CombinedOutput()
	if err == nil {
		return releasePreflightRunResult{code: 0, output: string(output)}
	}

	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("release preflight failed without exit status: %v\noutput:\n%s", err, output)
	}
	return releasePreflightRunResult{code: exitErr.ExitCode(), output: string(output)}
}

func newReleasePreflightRepo(t *testing.T, withOrigin bool, withChangelog bool) string {
	t.Helper()

	tempDir := t.TempDir()
	if withOrigin {
		origin := filepath.Join(tempDir, "origin.git")
		runGit(t, tempDir, "init", "--bare", origin)

		repo := filepath.Join(tempDir, "repo")
		runGit(t, tempDir, "clone", origin, repo)
		runGit(t, repo, "config", "user.name", "Release Preflight Test")
		runGit(t, repo, "config", "user.email", "release-preflight@example.com")

		writeRepoFile(t, repo, "README.md", "test fixture\n")
		if withChangelog {
			writeRepoFile(t, repo, "CHANGELOG.md", "## [v1.2.3]\n\n- release notes\n")
		} else {
			writeRepoFile(t, repo, "CHANGELOG.md", "## [v1.2.2]\n\n- previous release\n")
		}
		runGit(t, repo, "add", ".")
		runGit(t, repo, "commit", "-m", "initial release preflight fixture")
		runGit(t, repo, "branch", "-M", "main")
		runGit(t, repo, "push", "-u", "origin", "main")
		return repo
	}

	repo := filepath.Join(tempDir, "repo")
	mkdirAll(t, repo)
	runGit(t, repo, "init", "-b", "main")
	runGit(t, repo, "config", "user.name", "Release Preflight Test")
	runGit(t, repo, "config", "user.email", "release-preflight@example.com")
	writeRepoFile(t, repo, "README.md", "test fixture\n")
	if withChangelog {
		writeRepoFile(t, repo, "CHANGELOG.md", "## [v1.2.3]\n\n- release notes\n")
	}
	runGit(t, repo, "add", ".")
	runGit(t, repo, "commit", "-m", "initial release preflight fixture")
	return repo
}

func fakeTool(t *testing.T, binDir, name string) {
	t.Helper()

	path := filepath.Join(binDir, name)
	content := "#!/usr/bin/env bash\nexit 0\n"
	if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
		t.Fatalf("write fake tool %s: %v", path, err)
	}
}

func writeRepoFile(t *testing.T, repoDir, file, contents string) {
	t.Helper()

	path := filepath.Join(repoDir, filepath.FromSlash(file))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create parent dir for %s: %v", file, err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("write fixture file %s: %v", file, err)
	}
}
