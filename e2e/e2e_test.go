package e2e

import (
	"log"
	"os/exec"
	"strings"
	"testing"
)

const testBinary = "../ego_test"

func TestMain(m *testing.M) {
	log.Println("Running end to end tests")
	m.Run()
	cleanup()
}

func TestRunWithoutSubcommand(t *testing.T) {
	// given
	cmd := exec.Command(testBinary)
	expectedOut := "expected a subcommand"
	expectedErr := "exit status 1"
	
	// when
	stdout, err := cmd.Output()
	actualOut := strings.TrimSpace(string(stdout))
	
	// then
	if err.Error() != expectedErr {
		t.Errorf("Expected error to be '%s' but was '%s'", expectedErr, err.Error())
	}
	if actualOut != expectedOut {
		t.Errorf("Expected output to be '%s' but got '%s'", expectedOut, actualOut)
	}
}

func TestRunWithUnrecognisedSubcommand(t *testing.T) {
	// given
	cmd := exec.Command(testBinary, "foo")
	expectedOut := "unrecognised subcommand: foo"
	expectedErr := "exit status 1"

	// when
	stdout, err := cmd.Output()
	actualOut := strings.TrimSpace(string(stdout))
	
	// then
	if err.Error() != expectedErr {
		t.Errorf("Expected error to be '%s' but was '%s'", expectedErr, err.Error())
	}
	if actualOut != expectedOut {
		t.Errorf("Expected output to be '%s' but got '%s'", expectedOut, actualOut)
	}
}

func TestWithWrongDbPath(t *testing.T) {
	// given
	cmd := exec.Command(testBinary, "add", "-name=Bob", "-dbpath=/home/not/exists")
	expectedOut := "failed to add player: unable to open database file: no such file or directory"
	expectedErr := "exit status 1"

	// when
	stdout, err := cmd.Output()
	actualOut := strings.TrimSpace(string(stdout))
	
	// then
	if err.Error() != expectedErr {
		t.Errorf("Expected error to be '%s' but was '%s'", expectedErr, err.Error())
	}
	if actualOut != expectedOut {
		t.Errorf("Expected output to be '%s' but got '%s'", expectedOut, actualOut)
	}
}

func TestBasicHappyPath(t *testing.T) {
	// TODO
}


func cleanup() {
	log.Println("Cleaning up")
}