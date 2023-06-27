# gitlab-extract-variable

A tool to extract GitLab CI/CD variables written in Golang and save it in csv file

## Installing

`go install github.com/diduk001/gitlab-extract-variable@latest`

## Usage examples

```gitlab-extract-variable -token=TOKEN -project=ProjectOwner/ProjectName -compact```

```gitlab-extract-variable -token=TOKEN -project=ProjectOwner/ProjectName -output=.env -format=env```

### Options

- `-token=` - required, your personal GitLab access token
- `-project=` - required, the format is `{project owner}/{project name}`
- `-output=` - output file, by default `output.txt`
- `-compact` - extract only key and value
- `-format=` - format, `csv` or `env` only

# TODO
- Add several output formats
- Add colors to output