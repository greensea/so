package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/carlmjohnson/versioninfo"
	"github.com/chzyer/readline"
	"github.com/eiannone/keyboard"
	"github.com/gdamore/tcell/v2"
	"github.com/go-color-term/go-color-term/coloring"
	"github.com/go-resty/resty/v2"
	"github.com/greensea/so/common"
	"github.com/greensea/so/config"
	"github.com/greensea/so/sys"
	"github.com/rivo/tview"
)

func main() {
	if len(os.Args) <= 1 {
		Usage()
		os.Exit(0)
	}

	go Umami(os.Args[1])

	switch os.Args[1] {
	case "help":
		Usage()
		os.Exit(0)
	case "version":
		Version()
		os.Exit(0)
	case "config":
		Config()
		os.Exit(0)
	case "sys":
		Sys()
		os.Exit(0)
	case "quota":
		Quota()
		os.Exit(0)
	default:
		Run()
		os.Exit(0)
	}
}

func Usage() {
	fmt.Println(`Usage: so <help | version | config | quota>
       so <Your question>

Examples:
	so How to list directories?
	`)
}

func Config() {
	conf := config.Get()
	refEnabled := conf.CompatiableOpenAIEnabled
	refAPIKey := conf.CompatiableOpenAIKey
	refAPIEndpoint := conf.CompatiableOpenAIEndpoint
	refModel := conf.CompatiableOpenAIModelName

	app := tview.NewApplication()

	onSave := func() {
		app.Stop()

		conf.CompatiableOpenAIKey = refAPIKey
		conf.CompatiableOpenAIEndpoint = refAPIEndpoint
		conf.CompatiableOpenAIModelName = refModel
		conf.CompatiableOpenAIEnabled = refEnabled

		if refAPIEndpoint != "" && strings.HasPrefix(refAPIEndpoint, "http://") == false && strings.HasPrefix(refAPIEndpoint, "https://") == false {
			fmt.Fprintf(os.Stderr, coloring.Red("API endpoint must start with http:// or https://\nConfig not saved."))
			return
		}

		err := config.Save(conf)
		if err != nil {
			fmt.Fprintf(os.Stderr, coloring.Red("Unable to save config: %v\n"), err)
		}
	}

	dropdownInitValue := 0
	if refEnabled {
		dropdownInitValue = 0
	} else {
		dropdownInitValue = 1
	}

	form := tview.NewForm().
		AddDropDown("Custom Endpoint", []string{"Enabled", "Disabled"}, dropdownInitValue, func(option string, index int) {
			refEnabled = index == 0
		}).
		AddTextView("", "You can use OpenAI compatiable API endpoint to handle your request.", 0, 3, true, false).
		AddInputField("API Endpoint", refAPIEndpoint, 0, nil, func(t string) { refAPIEndpoint = t }).
		AddTextView("", `e.g. https://api.openai.com/v1/chat/completions`, 0, 2, true, false).
		AddInputField("API Key", refAPIKey, 0, nil, func(t string) { refAPIKey = t }).
		AddTextView("", "API key for endpoint\ne.g. Bearer xxxxx", 0, 3, true, false).
		AddInputField("Model", refModel, 0, nil, func(t string) { refModel = t }).
		AddTextView("", "e.g. gpt-4o-mini", 0, 1, true, false).
		AddButton("Save", onSave).
		AddButton("Cancel", func() {
			app.Stop()
		})

	form.SetItemPadding(0)

	form.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			app.Stop()
		}
		return event
	})

	form.SetBorder(true).SetTitle(" Use <tab> or mouse to navigate ").SetTitleAlign(tview.AlignCenter)
	if err := app.SetRoot(form, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
		panic(err)
	}

}

func Quota() {
	// 1. Load quota info
	conf := config.Get()

	stop := displaySpin()
	res, err := resty.New().R().
		SetHeader("X-SO-Client-ID", conf.ClientID).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{
			"PremiumCode": conf.PremiumCode,
			"Lang":        common.Lang(),
		}).
		Post(common.SoEndpoint("/quota"))
	stop()

	if err != nil {
		fmt.Fprintf(os.Stderr, coloring.Red("Failed to get quota: %v\n"), err)
		return
	}
	statusCode := res.StatusCode()
	if statusCode != 200 {
		fmt.Fprintln(os.Stderr, coloring.Extras.BrightBlack(res.String()))
		fmt.Fprintf(os.Stderr, coloring.Red("Failed to get quota, got HTTP status %d\n"), statusCode)
		return
	}

	// 2. Display options and wait user input
	fmt.Println(coloring.Extras.BrightBlue(res.String()))
	fmt.Println(coloring.Extras.BrightBlack("<p> to Enter Premium Code; [Other key] to quit"))

	char, _, err := keyboard.GetSingleKey()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read input: %c\n", err)
	}
	if char != 'p' {
		os.Exit(0)
	}

	// 3. To input premium code
	rl, err := readline.New(coloring.Extras.BrightBlack("Enter Premium Code: "))
	if err != nil {
		fmt.Fprintf(os.Stderr, coloring.Red("Unable to read input: %v\n"), err)
		os.Exit(0)
	}

	code, err := rl.Readline()
	code = strings.TrimSpace(code)

	conf.PremiumCode = code
	err = config.Save(conf)
	if err != nil {
		fmt.Fprintf(os.Stderr, coloring.Red("Unable to save premium code: %v\n"), err)
	} else {
		fmt.Println(coloring.Green("Premium code saved"))
	}
}

func Version() {
	fmt.Printf("so %s\n", versioninfo.Version)
	fmt.Printf("Commit: %s\n", versioninfo.LastCommit)
}

func Sys() {
	fmt.Printf(sys.Get().Text())
}
