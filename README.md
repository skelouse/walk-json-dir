# walk-json-dir

This tool parses and aggregates JSON files from a specified directory (and its subdirectories) into a single output JSON file.

## Installation

First, ensure you have Go installed. If not, go [here](https://golang.org/dl/).

Then run the following command:

```sh
go install github.com/skelouse/walk-json-dir@latest
```

## Usage

`walk-json-dir` provides the following options:

- `-root` or `-r`: Specifies the root directory to search for JSON files. Default is the current directory (`./`).
- `-output-file-path` or `-o`: Specifies the path to the output JSON file. Default is `output.json` in the current directory.


output all json files in current working directory and all subdirectories
```sh
walk-json-dir
```

output all json files in a provided directory
```sh
walk-json-dir -r <path/to/directory> -o <output/file/path>
```

## Example

Assuming you have a directory structure as follows:

```
data/
    config.json
    users/
        john.json
```

Contents:

- `config.json`

```json
{
  "appVersion": "1.0.0",
  "environment": "production"
}
```

- `users/john.json`

```json
{
  "name": "John Doe",
  "email": "john.doe@example.com"
}
```

And you want to aggregate these JSON files into a single file named `combined.json`, run:

```sh
walk-json-dir -r ./data -o combined.json
```

`combined.json`

```json
{
  "config": {
    "appVersion": "1.0.0",
    "environment": "production"
  },
  "users": {
    "john": {
      "email": "john.doe@example.com",
      "name": "John Doe"
    }
  }
}
```

## Notes

- Expects all JSON files to have valid JSON formats.
- If you get an out of memory error:
  - Increase swap
  - Close Chrome
  - Maybe there's too much data
  - [Download more RAM](https://downloadmoreram.com/)
