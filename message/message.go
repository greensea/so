package message

import (
	"fmt"
	"strings"

	"github.com/greensea/so/sys"
)

type GPTRequest struct {
	Model     string              `json:"model"`
	Messages  []GPTRequestMessage `json:"messages"`
	Stream    bool                `json:"stream"`
	MaxTokens int                 `json:"max_tokens"`
	// Temperature float64 `json:"temperature"`
}

type GPTRequestMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Get command generate messages
func GetCommand(q string, sysInfo *sys.Sys) []GPTRequestMessage {

	prompt1 := []string{
		`你是一名资深工程师，能够专业而清晰地写出 shell 命令。

## 目标

你需要帮助我编写 shell 命令，完成我的需求

## 工作流

1. 我会告诉你当前的系统信息，以及相关的工作区情况；
2. 紧接着我会向你问一个问题，请你帮助我写出一条清晰且有用的 shell 命令；
3. 你需要将 shell 命令放在代码块中，并且只输出 shell 命令，不要输出任何无关的内容，避免影响到我的自动解析程序。

`,
	}

	if sysInfo != nil && sysInfo.Text() != "" {
		prompt1 = append(prompt1, fmt.Sprintf("## 当前工作区情况\n\n```\n%s\n```", sysInfo.Text()))
	}

	m := []GPTRequestMessage{
		{
			Role:    "user",
			Content: strings.Join(prompt1, "\n"),
		},
		{
			Role:    "assistant",
			Content: fmt.Sprintf("好的，我准备好了，请告诉我你的需求，我来帮助你编写 shell 命令"),
		},
		{
			Role:    "user",
			Content: fmt.Sprintf("列出当前目录的文件，按时间排序"),
		},
		{
			Role: "assistant",
			Content: fmt.Sprintf("```" + "\n" +
				"ls -lt" + "\n" +
				"```"),
		},
		{
			Role:    "user",
			Content: fmt.Sprintf("```\n%s\n```", q),
		},
	}

	return m
}

// Get command describe messages
// q is command
func GetExplain(q string) []GPTRequestMessage {
	// 	m := []GPTRequestMessage{
	// 		{
	// 			Role: "user",
	// 			Content: fmt.Sprintf(`你是一名资深工程师，能够专业而清晰地解释 shell 命令的参数。
	// 现在，请你帮我解读 shell 的参数的作用，你的工作流程如下：
	// 1. 我会告诉你一条 shell 命令，你需要解释这条命令的作用，以及每一个参数的作用；
	// 2. 你需要解释说明放在代码块中，并且只输出代码块，不要输出任何无关的内容，以便我使用程序自动解析，
	// 3. 我将把你提供的解释说明直接输出到终端中，所以请你稍注意使用缩进和空格进行排版，使其美观。
	// 如果你准备好了，就告诉我。`),
	// 		},
	// 		{
	// 			Role:    "assistant",
	// 			Content: fmt.Sprintf("好的，我准备好了。我将帮助你解释 shell 命令的参数"),
	// 		},
	// 		{
	// 			Role:    "user",
	// 			Content: fmt.Sprintf("```\n%s\n```", q),
	// 		},
	// 		{
	// 			Role: "assistant",
	// 			Content: `\x60\x60\x60
	// ls 是 unix 上最常用的命令之一，用于列出目录下的文件信息
	// -l 以长格式列出
	// -t 按照时间排序
	// \x60\x60\x60`,
	// 		},
	// 		{
	// 			Role:    "user",
	// 			Content: fmt.Sprintf("```\n%s\n```", q),
	// 		},
	// 	}
	m := []GPTRequestMessage{
		{
			Role: "user",
			Content: fmt.Sprintf(`你是一名资深系统工程师，能够专业而清晰地解释 shell 命令的参数。
现在，请你帮我解读 shell 的参数的作用，你的工作流程如下：
1. 我会告诉你一条 shell 命令，你需要简要解释这条命令的作用，以及每一个参数的作用；
2. 请你直接进行解释，不要输出无关的内容；
3. 我将把你提供的解释直接输出到终端(terminal)中，所以请你一些缩进和空格进行排版，使其输出美观。
现在，请你解读一下这条命令:
`+"```"+`
%s
`+"```"+`
`, q),
		},
	}

	return m
}

// Get command describe messages
// q is command
func GetExplainEN(q string) []GPTRequestMessage {
	// 	m := []GPTRequestMessage{
	// 		{
	// 			Role: "user",
	// 			Content: fmt.Sprintf(`You are a helpfule assistant.
	// Now, please help me explain the shell command. Your workflow is as follows:
	// 1. I will give you a shell command, and you need to explain the command and each of the parameters;
	// 2. Just do explain, do not output any unrelated content;
	// 3. Your explaination will be direct display in terminal, so please add some spaces and tabs to make it print pretty.
	// `),
	// 		},
	// 		{
	// 			Role:    "assistant",
	// 			Content: fmt.Sprintf("Great, I am ready. I will help you explain the parameters of shell commands."),
	// 		},
	// 		{
	// 			Role:    "user",
	// 			Content: fmt.Sprintf("```\n%s\n```", q),
	// 		},
	// 		{
	// 			Role: "assistant",
	// 			Content: `
	// ls is one of the most commonly used commands on Unix, used to list file information in a directory.
	// -l: Lists in long format.
	// -t: Sorts by time.`,
	// 		},
	// 		{
	// 			Role:    "user",
	// 			Content: fmt.Sprintf("```\n%s\n```", q),
	// 		},
	// 	}
	m := []GPTRequestMessage{
		{
			Role: "user",
			Content: fmt.Sprintf(`You are a helpfule assistant.
Now, please help me explain the shell command. Your workflow is as follows:
1. I will give you a shell command, and you need to explain the command and each of the parameters;
2. Just do explain, do not output any unrelated content;
3. Your explaination will be direct display in pure text format, please add some indent to make it more readable.
4. Keep your explaination simple and useful.
Now please explain the following command:
`+"```"+`
%s
`+"```"+`
	`, q),
		},
	}

	return m
}
