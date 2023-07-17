package main

import "fmt"

func PrintHelp() {
	fmt.Println("Usage: flow [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --init --py\t\t\tInitialize flow structure for Python including GitLab-File creation")
	fmt.Println("  --init --ts\t\t\tInitialize flow structure for TypeScript including GitLab-File creation")
	fmt.Println("  --init --py\t\t\tInitialize flow structure for Go including GitLab-File creation")
	fmt.Println("  --create-gitlab-file --py\tCreate a GitLab-File for Python based Pipelines")
	fmt.Println("  --create-gitlab-file --ts\tCreate a GitLab-File for TypeScript based Pipelines")
	fmt.Println("  --create-gitlab-file --go\tCreate a GitLab-File for Go based Pipelines")
	fmt.Println("  --help\t\t\tShow Help")
}
