package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// rootCmd 作为 CLI 入口
var rootCmd = &cobra.Command{
	Use:   "OneInfer",
	Short: "A CLI tool for managing AI models",
	Long:  `OneInfer is a portable CLI tool for managing AI models like LLMs, embeddings, SD, and speech models.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to OneInfer!")
	},
}

// Execute 运行 rootCmd
func Execute() error {
	return rootCmd.Execute()
}
