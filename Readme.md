# sappend
It is a simple command line utility to sync azure blob storage files to local files by appending only new data.

## Installation
```bash
go install github.com/vanshitkumar/sappend@latest
```

## Usage
```bash
sappend -file-url <azure-blob-url> [-local-file <local-file-path>]
```

## Authorization
`sappend` uses [Azure Default Credential](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#readme-defaultazurecredential) for authentication. Make sure to set up your environment accordingly. You can use environment variables, managed identities, or Azure CLI for authentication.

## Command line options
- `-file-url` (required): The URL of the Azure Blob Storage file to sync from.
- `-local-file` (optional): The path to the local file to sync to. If not provided, it defaults to the blob url path, relative to current directory.
