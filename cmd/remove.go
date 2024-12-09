/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/jolt9dev/j9d/pkg/deployments"
	"github.com/spf13/cobra"
)

type removeOptions struct {
	project string
	target  string
	file    string
}

func registerRemoveCmd(rootCmd *cobra.Command) {
	// removeCmd represents the remove command

	removeArgs := removeOptions{}
	removeArgs.target = "default"

	// removeCmd represents the remove command
	var removeCmd = &cobra.Command{
		Use:   "remove",
		Short: "removes a deployment",
		Long:  `The remove command removes a deployment by using a j9d.yaml file`,
		RunE: func(cmd *cobra.Command, args []string) error {
			params := deployments.RemoveParams{}

			params.File = removeArgs.file
			params.Project = removeArgs.project
			params.Target = removeArgs.target

			return deployments.Remove(params)
		},
	}

	removeCmd.Flags().StringVarP(&removeArgs.project, "project", "p", "", "The project to remove. Projects should be in the @workspace/project format -e.g. @org/traefik.")
	removeCmd.Flags().StringVarP(&removeArgs.target, "target", "t", "", "The project target to remove. The target is generally used to specify the environment - e.g. dev, staging, prod.")
	removeCmd.Flags().StringVarP(&removeArgs.file, "file", "f", "", "The j9d file to use to remove the deployment. Supercedes the project and target flags.")

	rootCmd.AddCommand(removeCmd)
}

func init() {
	registerRemoveCmd(rootCmd)
}
