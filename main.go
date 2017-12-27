// pull_and_build project main.go
package main

import (
	"fmt"
	"os/exec"
	"time"

	"golang.org/x/net/context"

	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/viper"
)

func echo(str string) {
	fmt.Println("@@@@@ECHO@@@@@: " + str)
}

var remote, branch, bin string
var pull_pipe, build_pipe, run_pipe chan string

func main() {
	viper.AddConfigPath(".")
	viper.SetConfigName("git")
	viper.SetConfigType("json")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	remote = viper.Get("remote").(string)
	branch = viper.Get("branch").(string)
	bin = viper.Get("exe").(string)

	pull_pipe = make(chan string)
	build_pipe = make(chan string)
	run_pipe = make(chan string)

	go auto_pull()
	go auto_build()
	go auto_run()
	pull_pipe <- "pull"
	run_pipe <- "start"

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func pull() {
	echo("pull")
	cmd := exec.Command("git", "pull", remote, branch)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	// if !strings.Contains(string(out), "up-to-date") { // 如果不是最新
	if strings.Contains(string(out), "changed") { // 如果不是最新
		echo("update")
		build_pipe <- "rebuild"
	} else {
		echo("no update")
	}
}

func auto_pull() {
	timer := time.NewTicker(30 * time.Minute)
	for {
		select {
		case <-timer.C:
			pull()
		case str := <-pull_pipe:
			if str == "pull" {
				pull()
			}
		}
	}
}

func auto_build() {
	for {
		select {
		case str := <-build_pipe:
			if str == "rebuild" {
				echo("rebuild")
				run_pipe <- "close"
				cmd := exec.Command("go", "build", "-v")
				out, err := cmd.CombinedOutput()
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(string(out))
				run_pipe <- "start"
			}
		}
	}
}

func auto_run() {
	ctx, exit := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, bin)
	// cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Stdout = os.Stdout
	for {
		select {
		case <-ctx.Done():
			echo("done")
			build_pipe <- "rebuild"
		case str := <-run_pipe:
			switch str {
			case "start":
				echo("start")
				cmd.Start()
			case "close":
				echo("close")
				exit()
			}
		}
	}
}
