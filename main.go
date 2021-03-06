package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/dags-/promo-bot/github"
	"github.com/dags-/promo-bot/server"
)

func main() {
	ghtoken := flag.String("ghtoken", "", "Github auth token")
	owner := flag.String("owner", "", "Github repo owner")
	repo := flag.String("repo", "", "Github repo name")
	clientId := flag.String("clientId", "", "Discord bot client id")
	clientSecret := flag.String("clientSecret", "", "Discord bot client secret")
	redirectUri := flag.String("redirect", "", "Discord bot redirect uri")
	port := flag.Int("port", 8181, "The port to run the bot on")
	flag.Parse()

	flag.VisitAll(func(f *flag.Flag) {
		if f.Value.String() == "" {
			panic(fmt.Errorf("flag not set: %s (%s)", f.Name, f.Usage))
		}
	})

	session := github.NewSession(*ghtoken)
	rep := session.NewRepo(*owner, *repo)
	s := server.NewServer(session, rep, *clientId, *clientSecret, *redirectUri)
	go s.Start(*port)
	handleStop()
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
