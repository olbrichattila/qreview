package cmd

import (
	"fmt"
	"os"

	"github.com/olbrichattila/qreview/internal/format"
	"github.com/olbrichattila/qreview/internal/git"
	"github.com/olbrichattila/qreview/internal/review"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "qreview-go",
	Short: "AI-powered code review CLI using Amazon Q",
	Long: `qreview-go is a CLI tool that integrates with Amazon Q Developer CLI
to review your staged code changes for quality, security, and performance.`,
	Run: func(cmd *cobra.Command, args []string) {
		files, err := git.GetStagedFiles()
		if err != nil {
			fmt.Println("❌ Failed to get Git diff:", err)
			os.Exit(1)
		}

		if len(files) == 0 {
			fmt.Println("✅ No staged files to review.")
			return
		}

		hadIssues := false

		for _, file := range files {
			content, err := os.ReadFile(file)
			if err != nil {
				fmt.Printf("⚠️ Could not read %s: %v\n", file, err)
				continue
			}

			fmt.Printf("🔍 Reviewing %s...\n", file)
			result := review.AnalyzeCode(string(content), file)
			format.PrintToTerminal(file, result)

			// crude check to see if we should fail for Git hook
			if review.ContainsCritical(result) {
				hadIssues = true
			}
		}

		if hadIssues {
			fmt.Println("❌ Critical issues found. Commit aborted.")
			os.Exit(1)
		}
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
