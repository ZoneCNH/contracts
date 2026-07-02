package scripts_test

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestBoundaryCheckAllowsBusinessTermsInContractsPackage(t *testing.T) {
	repo := newBoundaryFixture(t, false)

	result := runBoundaryCheck(t, repo)
	if result.code != 0 {
		t.Fatalf("expected boundary check to pass, got exit %d\noutput:\n%s", result.code, result.output)
	}

	if !strings.Contains(result.output, "boundary check passed") {
		t.Fatalf("expected output to contain boundary success message\noutput:\n%s", result.output)
	}
}

func TestBoundaryCheckRejectsBusinessTermsInRuntimePackages(t *testing.T) {
	repo := newBoundaryFixture(t, true)

	result := runBoundaryCheck(t, repo)
	if result.code == 0 {
		t.Fatalf("expected boundary check to fail, got exit 0\noutput:\n%s", result.output)
	}

	want := "ERROR: forbidden business term found: MacroRegime"
	if !strings.Contains(result.output, want) {
		t.Fatalf("expected output to contain %q\noutput:\n%s", want, result.output)
	}
}

func newBoundaryFixture(t *testing.T, withRuntimeTerm bool) string {
	t.Helper()

	repo := t.TempDir()
	mkdirAll(t, filepath.Join(repo, "scripts"))
	mkdirAll(t, filepath.Join(repo, "pkg", "contracts"))
	mkdirAll(t, filepath.Join(repo, "internal", "runtime"))
	copyFile(t, "check_boundary.sh", filepath.Join(repo, "scripts", "check_boundary.sh"), 0o755)
	writeFile(t, filepath.Join(repo, "go.mod"), "module example.com/boundary\n\ngo 1.23\n")
	writeFile(t, filepath.Join(repo, "pkg", "contracts", "contracts.go"), "package contracts\n\n// MacroRegime is allowed in public contracts.\ntype MacroRegime = string\n")

	runtimeComment := "// clean runtime package\n"
	if withRuntimeTerm {
		runtimeComment = "// MacroRegime must not appear in runtime code.\n"
	}
	writeFile(t, filepath.Join(repo, "internal", "runtime", "runtime.go"), "package runtime\n\n"+runtimeComment+"func Ready() {}\n")

	return repo
}

func runBoundaryCheck(t *testing.T, repo string) commandResult {
	t.Helper()

	cmd := exec.Command("bash", "scripts/check_boundary.sh")
	cmd.Dir = repo
	cmd.Env = append(os.Environ(), "GOWORK=off")
	output, err := cmd.CombinedOutput()
	if err == nil {
		return commandResult{code: 0, output: string(output)}
	}

	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("boundary check failed without exit status: %v\noutput:\n%s", err, output)
	}
	return commandResult{code: exitErr.ExitCode(), output: string(output)}
}
