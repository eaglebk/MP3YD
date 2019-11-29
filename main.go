package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"regexp"

	"github.com/atotto/clipboard"
)

func cleanup() {
	fmt.Println("cleanup")
}

func clearClipboard() {
	clipboard.WriteAll("")
}

func main() {

	clearClipboard()

	r, _ := regexp.Compile("^(https?\\:\\/\\/)?(www\\.youtube\\.com|youtu\\.?be)\\/.+$")
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()

	for {
		fmt.Println("sleeping...")
		time.Sleep(1 * time.Second)
		text, _ := clipboard.ReadAll()
		if text != "" {
			isValidURL := r.Match([]byte(text))
			if isValidURL {
				fmt.Println(text)
				clearClipboard()
				args := []string{"--extract-audio", "--audio-format=mp3", "--output=%(title)s.%(ext)s", fmt.Sprint(text)}
				cmd := exec.Command("youtube-dl", args...)

				cmd.Env = os.Environ()
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err != nil {
					log.Fatal(err)
				}
			}

		}

	}

}
