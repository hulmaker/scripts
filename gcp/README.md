# gcp - GCP CLI Extensions

`gcp` extends `gcloud` with productivity tools for managing configurations, builds, and console access.


## Installation

```bash
go install gcp
```

## Commands

- **`login`**: Authenticate and setup a `gcloud` configuration with application default credentials.
- **`builds`**: Utilities for Cloud Builds.
  - **`list`**: List and filter recent cloud builds.
  - **`watch`**: Monitor a build and trigger notifications or commands on completion.
- **`open`**: Open the Google Cloud Console with given page in your default browser.
  - **`console`**: Open the Google Cloud Console
  - **`storage`**: Open the storage with given location

Use `gcp [command] --help` for more information.

## Developer Instructions

### Building locally

```bash
go build -o gcp main.go
```

### Running

```bash
go run main.go [command]
```
