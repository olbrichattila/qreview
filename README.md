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

### Possible Future Improvements
This project is already quite powerful, but there are several areas I plan to enhance further:

- **Improve Code Quality and Refactoring:**
The tool has grown quickly, and parts of the codebase would benefit from refactoring and better organization. As it continues to evolve, I plan to clean up internal abstractions, enforce stricter typing, and improve test coverage.

- **Additional Reporter Types:**
Expand the current HTML and Markdown reporters to include:

- Direct uploads to Amazon S3 for documentation hosting
- Integration with Confluence for team-wide visibility
- Custom API calls to push documentation or review results to internal systems

- **Dynamic Documentation Portal:**
Create a dedicated documentation viewer site with search and filtering by file, time, and type (review, code explanation, change summary). This would make it easier for teams to explore the review history over time.

- **Better GitHub Integration:**
Support for more GitHub workflows (e.g., triggered on push or comment), better diff visualization, and inline responses to review comments.

- **Slack or Other Chat Notifications for Reviews:**
Add integration with Slack (or similar chat platforms) to send real-time notifications when a new review is completed. This could include summaries, critical suggestions, and direct links to full reports or PR comments — helping teams stay informed and act faster.


- **Domain-Aware, Agentic AI Reviews:**
Extend the reviewer with agentic AI capabilities using Retrieval-Augmented Generation (RAG), where reviews can be enhanced by specific domain knowledge or internal documentation. This would allow the tool to not only identify code issues, but also evaluate business logic, validate domain-specific rules, and flag inconsistencies based on proprietary requirements or best practices.

These improvements aim to make the tool even more developer-friendly, team-scalable, and capable of reasoning beyond the code itself.

---

## Demo
Here’s a quick demo of QReview in action: Note: This video uses AI voice over.

[![Watch the video](https://img.youtube.com/vi/7C9rcr6esXU/default.jpg)](https://www.youtube.com/watch?v=7C9rcr6esXU)


You’ll see:
- At the beginning of the video I explain what the tool can do
- Terminal reviews and summaries
- Markdown + HTML output generation
- PR inline comments via GitHub API
- Dynamic, timestamped documentation
- GitGub action set-up

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

## Development Experience and Challenges
What started as a quick automation experiment quickly evolved into a much larger and more powerful tool as new ideas emerged during development — especially around making the system modular, with support for multiple and extendable reporters (e.g., HTML, Markdown, PR comments).

### Key Challenges:
- **Aligning Reviews to the Correct Lines in GitHub PRs**
One of the most complex problems was ensuring that AI-generated comments could be accurately mapped to the correct lines in a Pull Request. Since the AI (Amazon Q or others) may tokenize input, omit blank lines, or slightly reformat code internally, it often drifted from the original structure.
To solve this, I preprocess the full source file (removing blank lines) before feeding it into the model, then carefully map the cleaned lines back to their positions in the raw file. This also includes parsing and respecting the Git diff hunk headers `@@ -a,b +c,d @@` to ensure the review appears exactly where intended in the PR.

- **Handling Incomplete Context in Diffs:**
Reviewing just the diff often lacked enough context, leading to poor or irrelevant suggestions. I had to design a hybrid approach that uses both the full source file (for context-aware reviews) and the diff (for accurate line mapping and change detection).

- **Unexpected but Welcome Feedback — From Itself:**
Once the system was functional, I turned it loose on its own source code. The result? A cascade of inline suggestions and improvement ideas — effectively reviewing and improving itself. This was a rewarding moment that reinforced the value and practical utility of the tool

What began as a utility has grown into a versatile, self-improving, and extensible code review companion.

## Local installation guide:
Please use the command:
```
go install github.com/olbrichattila/qreview@latest
```

## GitHub automation installation guide:

1. Set Up GitHub Secrets
```
AI_CLIENT=bedrock
FILE_EXTENSIONS=php,go,js  # Add/extend file types to analyze
GH_TOKEN=<your_GitHub_token>  
AWS_ACCESS_KEY_ID=<your_AWS_access_key>  
AWS_SECRET_ACCESS_KEY=<your_AWS_secret_key>  
AWS_REGION=<your_AWS_region>  
```

> Note: Ensure your AWS IAM user has permissions for Amazon Bedrock (anthropic.claude-v2).

2. Create GitHub Workflow

Add this YAML to .github/workflows/code-review.yml:
```yaml
name: Automated Code Review on PR  

on:  
  pull_request:  
    types: [opened]  

jobs:  
  review:  
    runs-on: ubuntu-latest  
    steps:  
      - name: Pull QReview Docker Image  
        run: docker pull aolb/qreview:latest  

      - name: Run Code Review  
        env:  
          AI_CLIENT: ${{ secrets.AI_CLIENT }}  
          FILE_EXTENSIONS: ${{ secrets.FILE_EXTENSIONS }}  
          PR_URL: ${{ github.event.pull_request.html_url }}  
          GITHUB_TOKEN: ${{ secrets.GH_TOKEN }}  
          AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }}  
          AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_SECRET_ACCESS_KEY }}  
          AWS_REGION: ${{ secrets.AWS_REGION }}  
        run: |  
          docker run \  
            -e PR_URL="$PR_URL" \  
            -e AI_CLIENT="$AI_CLIENT" \  
            -e FILE_EXTENSIONS="$FILE_EXTENSIONS" \  
            -e GITHUB_TOKEN="$GITHUB_TOKEN" \  
            -e AWS_ACCESS_KEY_ID="$AWS_ACCESS_KEY_ID" \  
            -e AWS_SECRET_ACCESS_KEY="$AWS_SECRET_ACCESS_KEY" \  
            -e AWS_REGION="$AWS_REGION" \  
            aolb/qreview  
```

**Notes:**
Default model: anthropic.claude-v2 (ensure IAM permissions).
Extend FILE_EXTENSIONS to include other languages as needed.

