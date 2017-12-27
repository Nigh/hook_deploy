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

	pipe := make(chan string)

	go auto_pull(pipe)
	go auto_build(pipe)
	go auto_run(pipe)
	pull(pipe)
	pipe <- "start"

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func pull(pipe chan string) {
	echo("pull")
	cmd := exec.Command("git", "pull", remote, branch)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	if !strings.Contains(string(out), "up-to-date") { // 如果不是最新
		echo("update")
		pipe <- "close"
		pipe <- "rebuild"
	}
}

func auto_pull(pipe chan string) {
	timer := time.NewTicker(30 * time.Minute)
	for {
		select {
		case <-timer.C:
			pull(pipe)
		case str := <-pipe:
			if str == "pull" {
				pull(pipe)
			}
		}
	}
}

func auto_build(pipe chan string) {
	for {
		select {
		case str := <-pipe:
			if str == "rebuild" {
				echo("rebuild")
				cmd := exec.Command("go", "build", "-v")
				out, err := cmd.CombinedOutput()
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(string(out))
				pipe <- "start"
			}
		}
	}
}

func auto_run(pipe chan string) {
	ctx, exit := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, bin)
	// cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Stdout = os.Stdout
	for {
		select {
		case <-ctx.Done():
			echo("done")
			pipe <- "rebuild"
		case str := <-pipe:
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
