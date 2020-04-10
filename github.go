// Copyright © 2020 Ettore Di Giacinto <mudler@gentoo.org>
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program; if not, see <http://www.gnu.org/licenses/>.

package main

import (
	"context"
	"os"

	"github.com/google/go-github/v31/github"
	"golang.org/x/oauth2"
)

var GithubToken = os.Getenv("GITHUB_TOKEN")

func GitHubClient() *github.Client {
	if len(GithubToken) != 0 {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: GithubToken},
		)
		tc := oauth2.NewClient(ctx, ts)

		return github.NewClient(tc)
	}

	return github.NewClient(nil)
}
