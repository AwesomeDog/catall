# catall

`catall` is a smarter cross-platform alternative to the classic `cat` command
that recursively concatenates text files while respecting `.gitignore` and `.ignore` files.  
Perfect for preparing codebases for AI/LLM processing.

## Features

✨ Recursive Concatenation  
Automatically discovers and concatenates all text files in current directory and subdirectories

🔍 Ignore-Aware  
Respects `.gitignore`, `.ignore`, and custom patterns to exclude irrelevant files
(node_modules, build artifacts, temporary files, etc.)

🤖 AI-Friendly Output  
Produces clean concatenated text perfect for feeding to AI models like GitHub Copilot, ChatGPT, or Claude

⚡️ Lightning Fast  
Optimized for large codebases with thousands of files

## Ideal AI Use Cases

1. Feed entire projects to AI for:
    - Code analysis and refactoring
    - Documentation generation
    - Cross-file context understanding
    - Repository-level optimizations

2. Prepare training data for custom LLMs

3. Create context bundles for GPT-based tools

### Comparison

| Command  | Recursive | Ignore Support | Output Control | Ideal For      |
|----------|-----------|----------------|----------------|----------------|
| `cat *`  | ❌         | ❌              | ❌              | Single files   |
| `find`   | ✅         | ❌              | Manual         | File discovery |
| `catall` | ✅         | ✅              | ✅              | AI code input  |

## Usage

if you want to ignore some files, add a rule to `.gitignore` before running the cmd and everything will work like a
charm

```shell
cd <your code directory>

# Basic concatenation:
catall > all.txt  # that will print all to all.txt
catall --filter '\.(go|java)$'  # process only Go and Java files and print to stdout

# Preview files without concatenating:
catall --list
catall --list --filter '\.go$'  # list only Go files
```

## Installation

May require sudo if you want it in `/usr/local/bin`

Linux (x86_64/amd64)

```shell
sudo curl -L https://github.com/AwesomeDog/catall/releases/latest/download/catall_linux_amd64 -o /usr/local/bin/catall
sudo chmod +x /usr/local/bin/catall
catall --version
```

macOS (Apple Silicon)

```shell
sudo curl -L https://github.com/AwesomeDog/catall/releases/latest/download/catall_darwin_arm64 -o /usr/local/bin/catall
sudo chmod +x /usr/local/bin/catall
catall --version
```

Windows (amd64)

```shell
curl.exe -LO https://github.com/AwesomeDog/catall/releases/latest/download/catall_windows_amd64.exe
mv catall_windows_amd64.exe "$env:LOCALAPPDATA\Microsoft\WindowsApps\catall.exe"
catall --version
```

You may find executables for other platforms in [Releases](https://github.com/AwesomeDog/catall/releases)

## Manual Build (Requires Go)

```shell
go mod tidy
go build -o catall main.go
./catall > all.txt
```

## Prompts

```text
please make a cmdline app using golang:
1. it's like bash's cat, instead it will recursively cat text files under pwd.
2. it will respect .gitignore and .ignore files nested.

code structure:
1. first iterate pwd with github.com/boyter/gocodewalker v1.4.0
2. then list all files by alphabetical order
3. cat the files one by one, note that the file path starts with pwd should also be printed, like '`n==== ./src/main.go ====`n'
4. if the file is not of text format or malformed, print text 'BINARY OR BAD FORMAT', else the content itself
5. print all that to stdout, note that users may redirect the output to xxx.txt, you'd be careful not including the out file.

usage:
`catall > all.txt`
```
