# Contributing to bank-downloader

## Development

1. clone the repo
2. run `./setup.sh`
3. make a branch
4. make changes
5. make a pull request
   1. title MUST follow conventional-commit spec
   2. body MUST fill out the pull request template

### Support for new banks

1. each bank will be supported by providing the following three files 
   1. `processors/bankname.go`
   3. `processors/bankname_test.go`

### `processors/bankname.go`

This does all the work of downloading transactions from the bank.

1. It should provide a `NewBanknameProcessor` function that returns a `BanknameProcessor` struct.
2. `BanknameProcessor` should implement the `IProcessor` interface.
3. It should provide a struct to define the page objects for the bank's website. see `processors/anz.go` for an example.
4. Use all the `automation.Click`, `automation.Fill`, `automation.Find`, etc functions to automate the browser. see `processors/anz.go` for an example.
5. When downloading the exported transaction file, create a `FilenameTemplateContext` and use it to create a filename that gets passed to  `automation.DownloadFile`. see `processors/anz.go` for an example.

## Release

A release is controlled by repo maintainers.

When a PR is merged to `main`, a pending release PR will be created automatically or updated with new PR contents.

Mergin this pending release PR will generate an updated CHANGELOG and a build and publish a new version to the releases page.

