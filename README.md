# SlackPipe

Simple command line tool to send messages to a slack channel.

Reads standard input and last argument for messages.

```shell
usage: slackpipe [--username USERNAME] [--emoji EMOJI] [--channel CHANNEL] [--token TOKEN] MESSAGE

positional arguments:
  message

  options:
    --username USERNAME    Username [default: pipe]
    --emoji EMOJI          Emoji tag [default: :gear:]
    --channel CHANNEL      Channel name or ID [default: #general]
    --token TOKEN          Slack Token [default: ENV['SLACK_TOKEN']]
    --help, -h             display this help and exit
    --version              display version and exit
```

_Pipe input is split by newline and sent as multiple messages_

## Requirements

* Go
* Slack API Token
* Clever messages to pipe (or not)

## Install

```shell
$ go get github.com/h0tw1r3/slackpipe
```

## About

```shell
$ ./slackpipe -v
SlackPipe version 0.1
Copyright (C) 2016 Jeffrey Clark
License GPLv3+: GNU GPL version 3 or later <http://gnu.org/licenses/gpl.html>.

This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.
```
