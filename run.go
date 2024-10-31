package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/go-color-term/go-color-term/coloring"

	"github.com/eiannone/keyboard"

	"github.com/go-resty/resty/v2"
	"github.com/greensea/so/common"
	"github.com/greensea/so/config"
	"github.com/greensea/so/message"
	"github.com/greensea/so/sys"
	jsoniter "github.com/json-iterator/go"
)

type SoJobType string

const (
	// Generate command
	SoJobTypeCommand SoJobType = "command"
	// Explain command
	SoJobTypeExplain SoJobType = "explain"
)

func Run() {
	stopSpin := displaySpin()

	// 1. Request AI
	q := strings.Join(os.Args[1:], " ")
	chatResponse, err := requestChat(q, SoJobTypeCommand)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch: %v\n", err)
		os.Exit(-1)
	}

	cmds, err := extractCodeBlocks(chatResponse)
	if err != nil || len(cmds) == 0 {
		fmt.Fprintf(os.Stderr, "Raw response:\n%s\n", coloring.Extras.BgBrightWhite(string(chatResponse)))
		fmt.Fprintf(os.Stderr, "Unable to extract suggested command: %v\n", err)
		os.Exit(-1)
	}

	stopSpin()

	// 2. Display Options
	cmdExplain := ""
	for {
		if cmdExplain != "" {
			fmt.Printf("%s\n", coloring.Extras.BrightBlack(cmdExplain))
		}

		if cmdExplain == "" {
			fmt.Printf(coloring.Extras.BrightBlack("<Enter> to run; <e> to explain; [Other key] to quit\n"))
		} else {
			fmt.Printf(coloring.Extras.BrightBlack("<Enter> to run; [Other key] to quit\n"))
		}

		fmt.Printf("%s\n", coloring.Green(cmds[0]))

		// 3. Wait for user input
		char, key, err := keyboard.GetSingleKey()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to read input: %c\n", err)
		}

		if key != 0 && key != keyboard.KeyEnter {
			os.Exit(0)
		}

		if char == '\x00' && key == keyboard.KeyEnter {
			// Enter key pressed, execute cmd, and concat Stdout and Stderr
			err := ExecCmd(cmds[0], os.Stdout, os.Stderr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to run command: %v\n", err)
				os.Exit(-1)
			}
			os.Exit(0)
		} else if cmdExplain == "" && strings.ToLower(string(char)) == "e" {

			cancel := displaySpin()
			cmdExplain, err = explainCommand(cmds[0])
			cancel()

			if err != nil {
				fmt.Fprintf(os.Stderr, coloring.Red("Unable to explain command: %v\n"), err)
			}
			continue

		} else {
			os.Exit(0)
		}
	}
}

func explainCommand(c string) (string, error) {
	resp, err := requestChat(c, SoJobTypeExplain)
	if err != nil {
		return "", err
	}

	// cmds, err := extractCodeBlocks(resp)
	// if err != nil || len(cmds) == 0 {
	// 	return "", fmt.Errorf("Unable to extract command explains from response: %v", err)
	// }

	// return cmds[0], nil

	return string(resp), nil
}

// Request AI for answer
// q is the question asked by user
func requestChat(q string, job SoJobType) ([]byte, error) {
	c := config.Get()

	if c.CompatiableOpenAIEnabled == true && c.PremiumCode == "" {
		raw, err := requestChatOpenAI(q, job)
		if err != nil {
			return nil, err
		}

		content := jsoniter.Get(raw, "choices", 0, "message", "content").ToString()
		if content == "" {
			return nil, fmt.Errorf("Invalid response:\n%s", coloring.Extras.BgBrightWhite(string(raw)))
		}

		return []byte(content), nil

	} else {
		return requestChatSo(q, job)
	}
}

func requestChatSo(q string, job SoJobType) ([]byte, error) {

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-SO-Client-ID", config.Get().ClientID).
		SetBody(map[string]any{
			"Q":           q,
			"Job":         string(job),
			"Lang":        common.Lang(),
			"Sys":         sys.Get(),
			"PremiumCode": config.Get().PremiumCode,
		}).
		Post(common.SoEndpoint("/run"))

	if err != nil {
		return nil, err
	}

	statusCode := resp.StatusCode()
	if statusCode == 403 {
		return nil, errors.New(string(resp.Body()))
	} else if statusCode != 200 {
		return nil, fmt.Errorf("Got HTTP status code %d\n", statusCode)
	}

	buf := resp.Body()
	ret := jsoniter.Get(buf).Get("response").ToString()

	if ret == "" {
		return nil, fmt.Errorf("Invalid response: %s", string(buf))
	}

	return []byte(ret), nil
}

func requestChatOpenAI(q string, job SoJobType) ([]byte, error) {
	conf := config.Get()
	model := conf.CompatiableOpenAIModelName
	authorization := conf.CompatiableOpenAIKey

	if job == SoJobTypeCommand {
		return requestGPT(message.GetCommand(q, sys.Get()), model, authorization)
	} else if job == SoJobTypeExplain {
		if common.Lang() == "zh" {
			return requestGPT(message.GetExplain(q), model, authorization)
		} else {
			return requestGPT(message.GetExplainEN(q), model, authorization)
		}
	}

	return nil, fmt.Errorf("Invalid job type: %s", job)
}

// Extract command from chat content
// Command is located in the last code block (```)
func extractCodeBlocks(raw []byte) ([]string, error) {

	//   \x60 == `
	re := regexp.MustCompile(`(?sU)\x60\x60\x60[a-z]*[\r|\n](.+)[\r|\n]\x60\x60\x60`)
	matches := re.FindAllStringSubmatch(string(raw), -1)
	if len(matches) < 1 {
		return nil, fmt.Errorf("Unable to extract command from response")
	}

	ret := []string{}
	for _, v := range matches {
		ret = append(ret, strings.TrimSpace(v[1]))
	}

	return ret, nil
}

// Display a loading spin
// Return a function to stop the spin
func displaySpin() func() {
	t := time.NewTicker(time.Millisecond * 100)

	stop := make(chan struct{})

	go func() {
		defer t.Stop()
		chars := []string{"|", "/", "-", "\\"}
		i := 0

		for {
			i++
			select {
			case <-t.C:
				fmt.Printf("%s\r", chars[i%len(chars)])
			case <-stop:
				return
			}
		}
	}()

	return func() {
		close(stop)
	}

}

func ExecCmd(c string, stdout io.Writer, stderr io.Writer) error {
	s := sys.Get()
	if s.Shell == "" {
		s.Shell = "sh"
	}

	cmd := exec.Command(s.Shell)
	cmd.Stdin = strings.NewReader(c)
	if stdout != nil {
		cmd.Stdout = stdout
	}
	if stderr != nil {
		cmd.Stderr = stderr
	}

	err := cmd.Run()

	return err
}
