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
	"fmt"
	"strings"

	versioner "github.com/mudler/luet/pkg/versioner"

	"github.com/pkg/errors"
)

type GentooRepository struct{}

func (gr *GentooRepository) GetPackages(packageReq PackageRequest) ([]PackageResult, error) {
	var Packages []PackageResult

	client := GitHubClient()
	_, dir, _, err := client.Repositories.GetContents(context.TODO(), packageReq.Owner, packageReq.Repo, strings.Join([]string{packageReq.Category, packageReq.Name}, "/"), nil)
	if err != nil {
		return Packages, errors.Wrap(err, "Failed contacting github")
	}

	for _, file := range dir {
		if !strings.Contains(file.GetName(), ".ebuild") {
			continue
		}
		filename := strings.ReplaceAll(file.GetName(), ".ebuild", "")
		version := strings.ReplaceAll(filename, fmt.Sprintf("%s-", packageReq.Name), "")
		Packages = append(Packages, PackageResult{Path: file.GetPath(), Name: packageReq.Name, Category: packageReq.Category, Version: version})
	}

	return Packages, nil
}

func (gr *GentooRepository) AllPackages(packageReq PackageRequest) ([]PackageResult, error) {
	var Packages []PackageResult

	client := GitHubClient()
	_, dir, _, err := client.Repositories.GetContents(context.TODO(), packageReq.Owner, packageReq.Repo, "", nil)
	if err != nil {
		return Packages, errors.Wrap(err, "Failed contacting github")
	}

	for _, file := range dir {
		category := file.GetName()
		_, atoms, _, err := client.Repositories.GetContents(context.TODO(), packageReq.Owner, packageReq.Repo, strings.Join([]string{category}, "/"), nil)
		if err != nil {
			return Packages, errors.Wrap(err, "Failed contacting github")
		}
		for _, pack := range atoms {
			Packages = append(Packages, PackageResult{Path: pack.GetPath(), Name: pack.GetName(), Category: category, Version: ""})
		}
	}

	return Packages, nil
}

func (gr *GentooRepository) GetLatestPackage(packageReq PackageRequest) (PackageResult, error) {
	results, err := gr.GetPackages(packageReq)
	if err != nil {
		return PackageResult{}, err
	}

	versionMap := make(map[string]PackageResult)
	var versions []string
	for _, p := range results {
		versionMap[p.Version] = p
		versions = append(versions, p.Version)
	}

	ver := versioner.WrappedVersioner{}

	sorted := ver.Sort(versions)
	best := sorted[len(sorted)-1]

	return PackageResult{Path: versionMap[best].Path, Name: packageReq.Name, Category: packageReq.Category, Version: best}, nil
}
