package main

import (
	"github.com/dags-/promo-bot/server"
	"github.com/dags-/promo-bot/github"
	"flag"
	"fmt"
	"bufio"
	"os"
)

func main() {
	ghtoken := flag.String("ghtoken", "",  "Github auth token")
	owner := flag.String("owner", "", "Github repo owner")
	repo := flag.String("repo", "", "Github repo name")
	clientId := flag.String("clientId", "",  "Discord bot client id")
	clientSecret := flag.String("clientSecret", "",  "Discord bot client secret")
	redirectUri := flag.String("redirect", "",  "Discord bot redirect uri")
	port := flag.Int("port", 8181, "The port to run the bot on")
	flag.Parse()

	if !checkFlag(ghtoken, "Github Token") {
		return
	}

	if !checkFlag(clientId, "Client ID") {
		return
	}

	if !checkFlag(clientSecret, "Client Secret") {
		return
	}

	if !checkFlag(redirectUri, "Redirect URI") {
		return
	}

	go handleStop()

	session := github.NewSession(*ghtoken)
	rep := session.NewRepo(*owner, *repo)
	s := server.NewServer(session, rep, *clientId, *clientSecret, *redirectUri)
	s.Start(*port)
}

func handleStop() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if scanner.Text() == "stop" {
			fmt.Println("Stopping...")
			os.Exit(0)
			break
		}
	}
}

func checkFlag(flag *string, name string) bool {
	if *flag == "" {
		fmt.Println("Flag ", name, " has not been provided!")
		return false
	}
	return true
}