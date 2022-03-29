package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

func runCommandAt(cmd string, dir string) error {
	r := regexp.MustCompile(`(.+?)\s+(.*)`)
	matchs := r.FindStringSubmatch(cmd)
	var err error = nil
	fmt.Println("try to run cmd:", "["+cmd+"]", "at ["+dir+"]")
	if len(matchs) > 0 {
		cmdobj := exec.Command(matchs[1], matchs[2])
		cmdobj.Stdout = os.Stdout
		cmdobj.Stderr = os.Stderr
		cmdobj.Dir = dir
		err = cmdobj.Run()
	} else {
		err = errors.New("runCommand no matchs")
	}
	if err != nil {
		return err
	}
	return nil
}

func giteeHandler(j webhookGiteeJSON) {
	for _, v := range config.Apps {
		if j.Repo.Path == v.GitName {
			absProjectPath, _ := filepath.Abs(v.ProjectDIR)
			// git reset
			err := runCommandAt("git reset --hard", absProjectPath)
			if err != nil {
				return
			}
			// git checkout
			err = runCommandAt("git checkout "+v.Branch, absProjectPath)
			if err != nil {
				return
			}
			// pull
			err = runCommandAt("git pull origin "+v.Branch+" -f", absProjectPath)
			if err != nil {
				return
			}
			// build
			if len(v.BuildCMD) > 0 {
				err = runCommandAt(v.BuildCMD, absProjectPath)
				if err != nil {
					return
				}
			}
			// deploy
			if len(v.DeployCMD) > 0 {
				err = runCommandAt(v.DeployCMD, absProjectPath)
				if err != nil {
					return
				}
			}
		}
	}
}
