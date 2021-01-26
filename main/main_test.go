package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

const (
	TestFile   = "tmp_test.txt"
	ticketName = "ABC-123"
)

func TestHook(t *testing.T) {
	commitMsg := []byte("Added new feature")
	branchName := "feature/" + ticketName
	exec.Command("git", "checkout", "-b", branchName).Run()

	ioutil.WriteFile(TestFile, commitMsg, 0644)
	exec.Command("go", "run", "main.go", TestFile).Run()

	msg, err := ioutil.ReadFile(TestFile)

	if err != nil {
		t.Fail()
	}

	if string(msg) != string(formatMsg(ticketName, string(commitMsg))) {
		t.Fail()
	}

	exec.Command("git", "checkout", "master").Run()
	exec.Command("git", "branch", "-d", branchName).Run()
	os.Remove(TestFile)
}
