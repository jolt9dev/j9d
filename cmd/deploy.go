/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/jolt9dev/j9d/pkg/deployments"
	"github.com/spf13/cobra"
)

type deployOptions struct {
	project string
	target  string
	file    string
}

func registerDeployCmd(rootCmd *cobra.Command) {
	// deployCmd represents the deploy command

	deployArgs := deployOptions{}
	deployArgs.target = "default"

	var deployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "deploys an application or set of services",
		Long:  `The deploy command deploys an application or set of services by using a j9d.yaml file`,
		RunE: func(cmd *cobra.Command, args []string) error {
			params := deployments.DeployParams{
				File: deployArgs.file,
			}

			return deployments.Deploy(params)
		},
	}

	deployCmd.Flags().StringVarP(&deployArgs.project, "project", "p", "", "Project to deploy")
	deployCmd.Flags().StringVarP(&deployArgs.target, "target", "t", "", "Target to deploy to")
	deployCmd.Flags().StringVarP(&deployArgs.file, "file", "f", "", "Files to deploy")

	rootCmd.AddCommand(deployCmd)
}

func init() {
	registerDeployCmd(rootCmd)
}
