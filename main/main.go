package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
)

const (
	jiraIssueFormat = "[A-Z]{2,}-\\d+"
)

func messageChan(messageFile string) <-chan []byte {
	messageChan := make(chan []byte)
	go func() {
		message, err := ioutil.ReadFile(messageFile)
		if err != nil {
			messageChan <- nil
		}
		messageChan <- message
		close(messageChan)
	}()

	return messageChan
}

func gitBranchChan() <-chan []byte {
	gitBranchChan := make(chan []byte)
	go func() {
		branch, err := exec.Command("git", "symbolic-ref", "--short", "HEAD").CombinedOutput()

		if err != nil {
			gitBranchChan <- nil
		}
		regex := regexp.MustCompile(jiraIssueFormat)
		name := regex.Find(branch)
		if len(name) > 0 {
			gitBranchChan <- name
		} else {
			gitBranchChan <- name
		}
		close(gitBranchChan)
	}()

	return gitBranchChan
}

func formatMsg(branch string, message string) []byte {
	return []byte(fmt.Sprintf("[%s] %s", string(branch), string(message)))
}

func main() {
	if len(os.Args[1]) == 0 {
		return
	}

	messageFile := os.Args[1]

	messageChan := messageChan(messageFile)
	branchChan := gitBranchChan()

	message, branch := <-messageChan, <-branchChan

	if message == nil || branch == nil {
		return
	}

	ioutil.WriteFile(messageFile, formatMsg(string(branch), string(message)), 0644)
}
