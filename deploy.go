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
	r := regexp.MustCompile(`(\S+)`)
	matchs := r.FindAllString(cmd, -1)
	var err error = nil
	fmt.Println("try to run cmd:", matchs, "at ["+dir+"]")
	if len(matchs) > 0 {
		cmdobj := exec.Command(matchs[0], matchs[1:]...)
		cmdobj.Stdout = os.Stdout
		cmdobj.Stderr = os.Stderr
		cmdobj.Dir = dir
		err = cmdobj.Run()
	} else {
		err = errors.New("runCommand no matchs")
	}
	if err != nil {
		return err
	} else {
		fmt.Println("cmd run Finished")

	}
	return nil
}

func giteeHandler(j webhookGiteeJSON) {
	fmt.Println("Repo.Path = " + j.Repo.Path)
	for _, v := range config.Apps {
		fmt.Println("v.GitName = " + v.GitName)
		if v.Type == "gitee" && j.Repo.Path == v.GitName {
			absProjectPath, _ := filepath.Abs(v.ProjectDIR)
			// git reset
			err := runCommandAt("git reset --hard", absProjectPath)
			if err != nil {
				return
			}
			// pull
			err = runCommandAt("git fetch origin "+v.Branch, absProjectPath)
			if err != nil {
				return
			}
			// reset
			err = runCommandAt("git reset --hard origin/"+v.Branch, absProjectPath)
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
