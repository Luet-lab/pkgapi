# :package: :mag: PkgAPI

PkgAPI is a simple API that retrieves informations about package available inside Github repositories.

## Example

Retrieve latest version of a package:
`> curl -d "owner=gentoo&repo=gentoo&category=app-text&name=tree&repository_type=gentoo" -X POST http://127.0.0.1:4000/api/latest`

Retrieve all versions available of a package
`> curl -d "owner=gentoo&repo=gentoo&category=app-text&name=tree&repository_type=gentoo" -X POST http://127.0.0.1:4000/api/versions`

## Options

- **owner** repository owner (e.g. gentoo, Sabayon, etc. )
- **repo** repository name (e.g. sabayon-distro, gentoo, etc. )
- **category** package category (e.g. app-text, sys-libs, etc. )
- **name** package name (e.g. tree, vim, etc. )
- **repository_type** repository type (e.g. gentoo )

## Supported repository trees

- Gentoo

## Run

Simply start it. You can specify the `HOST`, `PORT` and `GITHUB_TOKEN` environment variables