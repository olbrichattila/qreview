*This is a submission for the [Amazon Q Developer "Quack The Code" Challenge](https://dev.to/challenges/aws-amazon-q-v2025-04-30): Crushing the Command Line*

## What I Built
I created an automation tool using golang, called QReview that enhances code review workflows using Amazon Q Developer CLI, and optionally integrates with GitHub, Ollama, and Amazon Bedrock. It helps developers automate reviews, generate documentation, and surface security or improvement suggestions — all from the command line.

**Key Features:**

- **Local or GitHub PR Code Review:**
Perform in-depth reviews using local changes or GitHub Pull Requests. Output can be shown directly in the terminal (with formatting and colors) or posted as inline comments in the PR.
- **Auto-Generated Documentation:**
Creates structured local HTML documentation for each review session — broken down by year, month/day, and hour/minute — and includes:
- **Review Documentation:** AI-driven comments and suggestions
- **Code Documentation:** Describes what the code does
- **Update Documentation:** Explains what changed and why
- **Custom Prompt Support:**

Extend reviews using your own YAML configuration. You can define prompts, input modes (diff, file, etc.), and outputs (HTML, Markdown, PR comments). Example:
```yaml
- prompt: "Summarize the differences."
  retrieverKind: diff
  commentOnPr: true
  reporters:
    - kind: html
      name: diff-summary
    - kind: markdown
      name: diff-summary
```

- **Flexible AI Client Integration:**
Supports Amazon Q Developer CLI by default, but can also run with:
- Amazon Bedrock
- Ollama (locally installed)


**GitHub Action Support:**

Fully automatable via GitHub workflows. Example:
```
name: Run code review on PR

on:
  pull_request:
    types: [opened]

jobs:
  run-qreview:
    runs-on: ubuntu-latest
    steps:
      - name: Pull Docker image
        run: docker pull aolb/qreview:latest

      - name: Run the container
        env:
          AI_CLIENT: ${{ secrets.AI_CLIENT }}
          PR_URL: ${{ github.event.pull_request.html_url }}
          GITHUB_TOKEN: ${{ secrets.GH_TOKEN }}
          AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }}
          AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          AWS_REGION: ${{ secrets.AWS_REGION }}
        run: |
          docker run \
            -e PR_URL="$PR_URL" \
            -e AI_CLIENT="$AI_CLIENT" \
            -e GITHUB_TOKEN="$GITHUB_TOKEN" \
            -e AWS_ACCESS_KEY_ID="$AWS_ACCESS_KEY_ID" \
            -e AWS_SECRET_ACCESS_KEY="$AWS_SECRET_ACCESS_KEY" \
            -e AWS_REGION="$AWS_REGION" \
            aolb/qreview
```
> Note: Amazon Q Developer CLI currently requires interactive login, so GitHub Action support is limited to Bedrock or other backends unless expanded.

---

## Demo
Here’s a quick demo of QReview in action:
TODO insert link, or video
You’ll see:

- Terminal reviews and summaries
- Markdown + HTML output generation
- PR inline comments via GitHub API
- Dynamic, timestamped documentation

## Code Repository
[https://github.com/olbrichattila/qreview](https://github.com/olbrichattila/qreview)

---

## How I Used Amazon Q Developer
Amazon Q Developer CLI is at the heart of this project. I used it to:
- Analyze code diffs and full files using custom prompts
- Extract improvement proposals, summaries, and documentation
- Output structured, reliable feedback in both markdown and HTML
- The Amazon Badrock implementation was done by Q developer, installed in vscode with a prompt like please use the existing interface like ollama and create a Badrock implementation.

---

Tips:
- Use the --code and --diff features of Amazon Q CLI for narrow, focused analysis
- Build reusable prompts with different retrieverKind values (e.g., file, diff)
- Since Q requires interactive login, local reviews are ideal — but the architecture supports extending it with Bedrock or Ollama for CI/CD compatibility

