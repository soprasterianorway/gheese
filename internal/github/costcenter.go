package github

import (
	"context"
	"fmt"
	"sort"

	"github.com/google/go-github/v84/github"
)

type UserCostCenter struct {
	Name       string
	CostCenter string
}

func GetUsersMissingCC(c *github.Client, ent string, onlyNone bool, filterCC string) ([]UserCostCenter, error) {
	ctx := context.Background()

	// Fetch all consumed licenses with pagination
	allUsers := []*github.EnterpriseLicensedUsers{}
	opts := &github.ListOptions{PerPage: 100}
	for {
		l, resp, err := c.Enterprise.GetConsumedLicenses(ctx, ent, opts)
		if err != nil {
			return nil, fmt.Errorf("get consumed licenses: %w", err)
		}
		allUsers = append(allUsers, l.Users...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	// Fetch all cost centers (no pagination for this API)
	cc, _, err := c.Enterprise.ListCostCenters(ctx, ent, nil)
	if err != nil {
		return nil, fmt.Errorf("list cost centers: %w", err)
	}

	// Build map of user login to cost center name
	resourceCostCenters := make(map[string]string)
	for _, center := range cc.CostCenters {
		for _, resource := range center.Resources {
			resourceCostCenters[resource.Name] = center.Name
		}
	}

	users := make([]UserCostCenter, 0, len(allUsers))
	for _, u := range allUsers {
		costcenter := resourceCostCenters[u.GithubComLogin]
		if costcenter == "" {
			costcenter = "none"
		}

		// Apply filters
		if onlyNone && costcenter != "none" {
			continue
		}
		if filterCC != "" && costcenter != filterCC {
			continue
		}

		users = append(users, UserCostCenter{
			Name:       u.GithubComLogin,
			CostCenter: costcenter,
		})
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].Name < users[j].Name
	})

	return users, nil
}
