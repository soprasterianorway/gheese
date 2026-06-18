/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/dev-soprasteriano/gheese/internal/github"
	"github.com/spf13/cobra"
)

// costcenterCmd represents the costcenter command
var costcenterCmd = &cobra.Command{
	Use:   "costcenter",
	Short: "List users and their cost center assignments",
	Long: `List all users in the enterprise and their cost center assignments.

By default, shows all users. Use --only-none to show only users without a cost center,
or --cost-center to filter by a specific cost center name.`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := github.NewClient()
		if err != nil {
			fmt.Println(err)
			return
		}

		onlyNone, _ := cmd.Flags().GetBool("only-none")
		filterCC, _ := cmd.Flags().GetString("cost-center")

		users, err := github.GetUsersMissingCC(c, "soprasteriasca", onlyNone, filterCC)
		if err != nil {
			fmt.Println(err)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "USER\tCOST CENTER")
		for _, user := range users {
			fmt.Fprintf(w, "%s\t%s\n", user.Name, user.CostCenter)
		}
		_ = w.Flush()
	},
}

func init() {
	enterpriseCmd.AddCommand(costcenterCmd)

	costcenterCmd.Flags().Bool("only-none", false, "Show only users without a cost center assigned")
	costcenterCmd.Flags().String("cost-center", "", "Filter users by specific cost center name")
}
