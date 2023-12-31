# Bank Transaction Downloader

Simple command line tool that you'd use in a scheduled task to 
download bank transactions.


## Supported Banks

- ANZ

(that's it for now, but [feel free to add more!](#contributing))

## Getting Started

1. create a example config and download the playwright patched browsers with `bank-downloader init`
2. edit the config file
3. run `bank-downloader download`


### Prerequisites

`bank-downloader` makes use of chromedp to automate a real browser.

This means you need to have google chrome, or chromium installed on your system.

## Configuration

`bank-downloader` uses a YAML config file to describe the jobs to run.

Each job should describe a single institution and login details, and a list of accounts to download.
It is assumed that each account will be downloaded to a separate file. 

### `config.yaml`
The config file is located at `~/.bank-downloader/config.yaml` by default.

```yaml
---
sources:
  - name: "anz"
    exportFormat: "AgriBank (CSV)"
    outputTemplate: "mybank-firstuser-{{.Account.NameSlug}}-{{.Account.NumberSlug}}-{{.FromDateUnix}}-{{.ToDateUnix}}.csv"
    daysToFetch: 7
    credentials:
        type: file
        username: "myusername"
        password: "mypassword"
    accounts:
      - name: "My Everyday Cash Account Name"
        number: "123456789"
  - name: "anz"
    exportFormat: "Agrimaster(CSV)"
    outputTemplate: "mybank-seconduser-{{.Account.NameSlug}}-{{.Account.NumberSlug}}-{{.FromDateUnix}}-{{.ToDateUnix}}.csv"
    daysToFetch: 7
    credentials:
        type: file
        username: "my-other-username"
        password: "my-other-password"
    accounts:
      - name: "My Offset Account Name"
        number: "123456789"
      - name: "My Fixed Rate Mortage Account Name"
        number: "123456789"
        
```

#### `sources`

Sources where bankdownloader can download transactions from.

#### `source[].name`

The name of the bank to download from. Currently only `anz` is supported.

#### `source[].exportFormat`

The format of the downloaded file. Supported formats depends on the source.

For reference, the ANZ source supports: 

- `Microsoft Money(OFC)`
- `MYOB(OFX)`
- `MYOB(QIF)`
- `Quicken(OFX)`
- `Quicken(QIF)`
- `Microsoft Excel(CSV)`
- `Agrimaster(CSV)`
- `Phoenix Gateway(CSV)`

#### `source[].outputTemplate`

The template to use for the output file name.

_example_: `mybank-firstuser-{{.Account.NameSlug}}-{{.Account.NumberSlug}}-{{.FromDateUnix}}-{{.ToDateUnix}}.csv`

The following variables are available:

- `{{.Source}}` - the name of the source
- `{{.SourceSlug}}` - the name of the source, with spaces replaced with dashes
- `{{.Account.Name}}` - the name of the account
- `{{.Account.NameSlug}}` - the name of the account, with spaces replaced with dashes
- `{{.Account.Number}}` - the account number
- `{{.Account.NumberSlug}}` - the account number, with spaces replaced with dashes
- `{{.FromDate}}` - the date of the first transaction in the file
- `{{.FromDateSlug}}` - the date of the first transaction in the file, with spaces replaced with dashes
- `{{.FromDateUnix}}` - the date of the first transaction in the file, in unix time
- `{{.ToDate}}` - the date of the last transaction in the file
- `{{.ToDateSlug}}` - the date of the last transaction in the file, with spaces replaced with dashes
- `{{.ToDateUnix}}` - the date of the last transaction in the file, in unix time
- `{{.Now}}` - the current date

#### `source[].daysToFetch`

The number of days to fetch. Defaults to `7`.

#### `source[].credentials`

The credentials to use to log in to the bank.

#### `source[].credentials.type`

The type of credentials to use to log in to the bank. Currently only `file`, `gopass`, `env`, `gopass-otp`, `keychain` are supported.

If `file` is used, the username and password are read from the `username` and `password` fields.

If `gopass` is used, the username and password are read from gopass using the `usernameKey` and `passwordKey` fields, and the `gopass` binary must be installed.

If `env` is used, the username and password are read from environment variables using the `usernameKey` and `passwordKey` fields.

If `gopass-otp` is used, the username, password and otp are read from gopass using the `usernameKey`, `passwordKey` and `otpKey` fields, and the `gopass` binary must be installed.

If `keychain` is used, the username and password are read from the keychain using the `usernameKey` and `passwordKey` fields. `keychain` is only [supported on osx, linux, bsd or windows](https://pkg.go.dev/github.com/zalando/go-keyring@v0.2.3).


#### `source[].accounts`

A list of accounts to download.

#### `source[].accounts[].name`

The name of the account to download.

#### `source[].accounts[].number`

The number of the account to download.

## Running

### `bank-downloader init`

Creates an example config file at `~/.bank-downloader/config.yaml`.

### `bank-downloader download`

Downloads transactions for all jobs in the config file.

#### Arguments

##### `--config`

The path to the config file to use. Defaults to `~/.bank-downloader/config.yaml`.

##### `--headless`

Whether to run the browser in headless mode. Defaults to `true`.

##### `--debug`

Whether to run the browser in debug mode. Defaults to `false`.

##### `--range-strategy`

The date range mode to use. Defaults to `days`.

Options:

- `days-ago` - use the `daysToFetch` config to calculate the date range, but always download transactions from yesterday to the last `daysToFetch` days.
- `since-last-download` - uses the `daysToFetch` config to calculate the date range, but always download transactions from the last downloaded transaction date. If the last downloaded transaction date is not available, it will download transactions from yesterday to the last `daysToFetch` days. If the end date is in the future, it will be set to yesterday.

The idea here is that with a `daysToFetch` of `60`, you can set up a scheduled task to run `bank-downloader download` every week, and it will:
- on first run download transactions from the last 60 days
- on subsequent runs download transactions from the last 60 days, or since the last downloaded transaction date, whichever is more recent


## How it works

`bank-downloader` automates your installed instance of google chrome.

1. logs in to the bank, 
2. for each account:
   1. navigates to the accounts page,
   2. navigates to the transactions page,
   3. calculates the date range based on the `daysToFetch` config, and the last downloaded transaction date
   4. downloads the transactions for the date range
   5. saves the transactions to a file, using the `outputTemplate` config
   6. saves the last downloaded transaction date to a file, so that next time it can calculate the date range correctly

## Contributing


see [CONTRIBUTING.md](CONTRIBUTING.md)
