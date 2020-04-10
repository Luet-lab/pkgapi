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
	"github.com/go-macaron/binding"
	macaron "gopkg.in/macaron.v1"
)

type RepositoryType interface {
	GetPackages(PackageRequest) ([]PackageResult, error)
	GetLatestPackage(PackageRequest) (PackageResult, error)
}

func NewRepositoryType(req PackageRequest) RepositoryType {
	switch req.RepositoryType {
	case "gentoo":
		return &GentooRepository{}
	default:
		return nil
	}

}

type Result struct {
	Packages []PackageResult
	Error    string
}
type PackageResult struct {
	Name     string
	Category string
	Version  string
}

type PackageRequest struct {
	Name           string
	Category       string
	Owner          string
	Repo           string
	Version        string
	RepositoryType string
}

func main() {
	m := macaron.Classic()

	m.Use(macaron.Renderer())

	m.Post("/api/latest/", binding.Bind(PackageRequest{}), LatestPackageVersion)
	m.Post("/api/versions/", binding.Bind(PackageRequest{}), PackageVersions)
	//m.Post("/api/versions/", binding.Bind(pkg.DefaultPackage{}), PackageVersions)

	m.Run()
}

func PackageVersions(ctx *macaron.Context, packageReq PackageRequest) {

	repo := NewRepositoryType(packageReq)
	Packages, err := repo.GetPackages(packageReq)
	if err != nil {
		handleErr(ctx, err)
		return
	}
	ctx.JSON(200, Result{Packages: Packages})
}

func handleErr(ctx *macaron.Context, err error) { ctx.JSON(500, Result{Error: err.Error()}) }

func LatestPackageVersion(ctx *macaron.Context, packageReq PackageRequest) {
	repo := NewRepositoryType(packageReq)
	Package, err := repo.GetLatestPackage(packageReq)
	if err != nil {
		handleErr(ctx, err)
		return
	}

	ctx.JSON(200, Result{Packages: []PackageResult{Package}})
}
