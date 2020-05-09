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
	"compress/bzip2"
	"io"
	"io/ioutil"
	"os"

	entropy "github.com/Sabayon/pkgs-checker/pkg/entropy"
	versioner "github.com/mudler/luet/pkg/versioner"
	"github.com/pkg/errors"
)

type SabayonRepository struct{}

func CurrentPackageList(url string) ([]*entropy.EntropyPackage, error) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "entropy-")
	if err != nil {
		return nil, err
	}

	decompressed, err := ioutil.TempFile(os.TempDir(), "entropy-")
	if err != nil {
		return nil, err
	}

	defer os.Remove(decompressed.Name())
	defer os.Remove(tmpFile.Name())
	if err := DownloadFile(tmpFile.Name(), url); err != nil {
		return nil, err
	}
	bz := bzip2.NewReader(tmpFile)
	out, err := os.Create(decompressed.Name())
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(out, bz)
	if err != nil {
		return nil, err
	}
	out.Close()

	return entropy.RetrieveRepoPackages(decompressed.Name())
}

func (gr *SabayonRepository) GetPackages(packageReq PackageRequest) ([]PackageResult, error) {
	var Packages []PackageResult

	packages, err := CurrentPackageList(packageReq.Repo)
	if err != nil {
		return Packages, errors.Wrap(err, "Failed getting Sabayon package lists")
	}

	for _, p := range packages {
		if p.Name == packageReq.Name && p.Category == packageReq.Category {
			Packages = append(Packages, PackageResult{Name: p.Name, Category: p.Category, Version: p.Version})
		}
	}

	return Packages, nil
}

func (gr *SabayonRepository) AllPackages(packageReq PackageRequest) ([]PackageResult, error) {
	var Packages []PackageResult

	packages, err := CurrentPackageList(packageReq.Repo)
	if err != nil {
		return Packages, errors.Wrap(err, "Failed getting Sabayon package lists")
	}

	for _, p := range packages {
		Packages = append(Packages,
			PackageResult{Name: p.Name, Category: p.Category, Version: p.Version})

	}
	return Packages, nil
}

func (gr *SabayonRepository) GetLatestPackage(packageReq PackageRequest) (PackageResult, error) {
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
