package main

import (
	"os"
	"bufio"
	"fmt"
	"io/ioutil"
	"io"
	"strings"

	"github.com/nlopes/slack"
	"github.com/alexflint/go-arg"
)

const (
	APP_VERSION = "0.1"
	APP_NAME = "SlackPipe"
	APP_LEGAL = "Copyright (C) 2016 Jeffrey Clark\nLicense GPLv3+: GNU GPL version 3 or later <http://gnu.org/licenses/gpl.html>.\n\nThis is free software: you are free to change and redistribute it.\nThere is NO WARRANTY, to the extent permitted by law."
	ENV_TOKEN_DISPLAY = "ENV['SLACK_TOKEN']"
	DEFAULT_EMOJI = ":gear:"
	DEFAULT_NAME = "pipe"
	DEFAULT_CHANNEL = "#general"
	DEFAULT_SLACK_TOKEN = ""
	DEFAULT_FILEMODE = false
)

type args struct {
	Message   string   `arg:"positional"`
	Username  string   `arg:"-u,help:Username"`
	Emoji     string   `arg:"-e,help:Emoji tag"`
	Channel   string   `arg:"-c,help:Channel name or ID"`
	Token     string   `arg:"-t,help:Slack Token"`
	FileMode  bool     `arg:"-f,help:Upload file (Message is pipe Title)"`
}

func (args) Version() string {
	return fmt.Sprintf("%s version %s\n%s\n", APP_NAME, APP_VERSION, APP_LEGAL)
}

func FatalCheck(e error) {
	if e != nil {
		fmt.Printf("Error: %s\n", e.Error())
		os.Exit(1)
	}
}

func main() {
	var args args
	args.Username = DEFAULT_NAME
	args.Emoji = DEFAULT_EMOJI
	args.Channel = DEFAULT_CHANNEL
	args.Token = ENV_TOKEN_DISPLAY
	args.FileMode = DEFAULT_FILEMODE

	env_token := os.Getenv("SLACK_TOKEN")
	if len(env_token) < 1 && len(DEFAULT_SLACK_TOKEN) > 0 {
		args.Token = "builtin"
	}

	parsedargs := arg.MustParse(&args)

	stat, _ := os.Stdin.Stat()
	have_stdin := (stat.Mode() & os.ModeNamedPipe != 0)
	if ! have_stdin && args.Message == "" {
		parsedargs.Fail("MESSAGE must be an argument or | to STDIN")
	}

	if args.Token == "builtin" {
		args.Token = DEFAULT_SLACK_TOKEN
	} else if args.Token == ENV_TOKEN_DISPLAY {
		args.Token = env_token
	}

	if args.Token == "" {
		parsedargs.Fail(fmt.Sprintf("TOKEN must be supplied or set %s", ENV_TOKEN_DISPLAY))
	}

	api := slack.New(args.Token)

	if have_stdin && args.FileMode {
		tempfile, err := ioutil.TempFile(os.TempDir(), APP_NAME)
		FatalCheck(err)
		defer os.Remove(tempfile.Name())

		data := make([]byte, 4096)

		input := bufio.NewReader(os.Stdin)
		output := bufio.NewWriter(tempfile)

		for {
			data = data[:cap(data)]
			n, err := input.Read(data)
			if err == io.EOF {
				break
			}
			data = data[:n]
			n, err = output.Write(data)
			FatalCheck(err)
		}
		output.Flush()

		params := slack.FileUploadParameters{}
		params.Title = args.Message
		params.File = tempfile.Name()
		params.Filetype = "auto"
		params.Channels = strings.Split(args.Channel, ",")
		tempfile.Close()

		_, uerr := api.UploadFile(params)
		FatalCheck(uerr)
	} else {
		params := slack.PostMessageParameters{}
		params.IconEmoji = args.Emoji
		params.Username = args.Username

		if have_stdin {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				stdin_line := scanner.Text()
				if len(stdin_line) > 0 {
					_, _, err := api.PostMessage(args.Channel, stdin_line, params)
					FatalCheck(err)
				}
			}
			FatalCheck(scanner.Err())
		} else {
			_, _, err := api.PostMessage(args.Channel, args.Message, params)
			FatalCheck(err)
		}
	}
}
