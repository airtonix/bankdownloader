# Bank Transaction Downloader

Simple command line tool that you'd use in a scheduled task to 
download bank transactions.


## Getting Started

1. create a example config `bank-downloader init`
2. edit the config file
3. run `bank-downloader download`


## Config File

The config file should describe a list of download jobs.

Each job should target a single institution and login details.

```yaml
jobs:
  - source: "anz"
    format: "AgriBank (CSV)"
    output: "mybank-{accountName}-{accountNumber}-{dateRange}.csv"
    daysToFetch: 7
    credentials:
        username: "myusername"
        password: "mypassword"
    accounts:
      - name: "My Everyday Cash Account Name"
        number: "123456789"
      - name: "My Offset Account Name"
        number: "123456789"
      - name: "My Fixed Rate Mortage Account Name"
        number: "123456789"
```


## Contributing

> ⚠️ pull request titles must follow conventional commits
