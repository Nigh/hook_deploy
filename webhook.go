package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type giteeRepo struct {
	Name string `json:"name"`
	Path string `json:"path"`
}
type giteeSender struct {
	Name      string `json:"name"`
	AvatarUrl string `json:"avatar_url"`
}
type webhookGiteeJSON struct {
	Repo     giteeRepo   `json:"repository,omitempty"`
	Sender   giteeSender `json:"sender,omitempty"`
	Password string      `json:"password,omitempty"`
}

// TODO: add gitee password
type appsJSON struct {
	Type       string `json:"type"`
	GitName    string `json:"git_name"`
	GitURL     string `json:"git_url"`
	Branch     string `json:"branch"`
	ProjectDIR string `json:"project_dir,omitempty"`
	BuildCMD   string `json:"build,omitempty"`
	DeployCMD  string `json:"deploy,omitempty"`
}
type configJSON struct {
	Port int        `json:"port"`
	Apps []appsJSON `json:"apps"`
}

var config configJSON

func getConfig(path string) error {
	absPath, _ := filepath.Abs(path)
	jsonFile, err := os.Open(absPath)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &config)
	fmt.Printf("%v\n", config)
	return nil
}
func main() {
	fmt.Println("Hook Deploy Start")
	if err := getConfig("config.json"); err == nil {
		Listen(config.Port)
	} else {
		log.Panic(err)
	}
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("headers: %v\n", r.Header)

	if r.Header.Get("Content-Type") != "application/json" {
		return
	}
	decoder := json.NewDecoder(r.Body)
	switch r.Header.Get("User-Agent") {
	case "git-oschina-hook":
		var t webhookGiteeJSON
		decoder.Decode(&t)
		go giteeHandler(t)
	}
}

func Listen(port int) {
	fmt.Println("Listen on port:", port)
	http.HandleFunc("/", handleWebhook)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
