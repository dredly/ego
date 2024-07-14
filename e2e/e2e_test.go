package e2e

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
)

const testBinary = "../ego_test"
var testDir string

func TestMain(m *testing.M) {
	var err error
	testDir, err = os.MkdirTemp("", "e2e")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("created testDir in %s", testDir)
	m.Run()
	log.Printf("Removing testDir %s", testDir)
	os.RemoveAll(testDir)
}

func TestBasicHappyPath(t *testing.T) {
	dbFileName := uuid.New().String() + ".db"
	dbPath := filepath.Join(testDir, dbFileName)
	dbPathArg := fmt.Sprintf("-dbpath=%s", dbPath)

	execAndVerify(t, []string{"create", dbPathArg}, "Successfully initialised leaderboard", "")
	execAndVerify(t, []string{"add", "-name=Player1", dbPathArg}, "added new player Player1 with starting ELO 1000", "")
	execAndVerify(t, []string{"add", "-name=Player2", dbPathArg}, "added new player Player2 with starting ELO 1000", "")
	execAndVerify(
		t, 
		[]string{"record", "-p1=Player1", "-p2=Player2", "-score=11-6", dbPathArg}, 
		"Recorded result 11-6 between Player1 and Player2\nPlayer1 elo: 1000.00 -> 1010.00. Player2 elo: 1000.00 -> 990.00", 
		"",
	)
	execAndVerify(t, []string{"leaderboard", dbPathArg}, "Leaderboard\n1. Player1: 1010.00\n2. Player2: 990.00", "")
}

func TestRunWithoutSubcommand(t *testing.T) {
	execAndVerify(t, []string{}, "expected a subcommand", "exit status 1")
}

func TestRunWithUnrecognisedSubcommand(t *testing.T) {
	execAndVerify(t, []string{"foo"}, "unrecognised subcommand: foo", "exit status 1")
}

func TestWithWrongDbPath(t *testing.T) {
	execAndVerify(t, []string{"add", "-name=Bob", "-dbpath=/home/not/exists"}, 
		"failed to add player: unable to open database file: no such file or directory", 
		"exit status 1",
	)
}

func execAndVerify(t *testing.T, args []string, expectedOut, expectedErr string) {
	t.Helper()
	cmd := exec.Command(testBinary, args...)
	stdout, err := cmd.Output()
	actualOut := strings.TrimSpace(string(stdout))

	if expectedErr == "" {
		if err != nil {
			t.Errorf("Expected no error but got %s", err.Error())
		}
	} else {
		if err == nil {
			t.Fatal("Expected command to return an error but it did not")
		}
		if err.Error() != expectedErr {
			t.Errorf("Expected error to be '%s' but was '%s'", expectedErr, err.Error())
		}
	}

	if actualOut != expectedOut {
		t.Errorf("Expected output to be '%s' but got '%s'", expectedOut, actualOut)
	}
}