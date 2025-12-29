package internal

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
	"github.com/schollz/progressbar/v3"
)

const (
	PrivateDirPerm  fs.FileMode = 0700
	PrivateFilePerm fs.FileMode = 0600
)

// if destination is empty string, it will parse url from source
func SyncFile(source string, destination string) error {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return fmt.Errorf("failed to obtain a credential: %w", err)
	}

	blobClient, err := blob.NewClient(source, cred, nil)
	if err != nil {
		return fmt.Errorf("failed to create blob client: %w", err)
	}

	if destination == "" {
		u, err := url.Parse(blobClient.URL())
		if err != nil {
			return fmt.Errorf("failed to parse blob url: %w", err)
		}
		destination, err = filepath.Localize(strings.TrimPrefix(u.Path, "/"))
		if err != nil {
			return fmt.Errorf("failed to get local path: %w", err)
		}
	}

	destFile, err := getFileClient(destination)
	if err != nil {
		return fmt.Errorf("failed to get file client: %w", err)
	}
	defer destFile.Close()

	err = syncFile(blobClient, destFile)
	if err != nil {
		return fmt.Errorf("failed to sync file: %w", err)
	}

	return nil
}

// returns file client with append only rights
// if file already exists does not overwrite it
func getFileClient(path string) (*os.File, error) {
	err := os.MkdirAll(filepath.Dir(path), PrivateDirPerm)
	if err != nil {
		return nil, fmt.Errorf("failed to create directories for %s: %w", path, err)
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, PrivateFilePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", path, err)
	}
	return file, nil
}

// syncs and also shows a progress bar
func syncFile(src *blob.Client, dest *os.File) error {
	info, err := dest.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat destination file: %w", err)
	}
	offset := info.Size()

	resp, err := src.DownloadStream(context.Background(), &blob.DownloadStreamOptions{
		Range: blob.HTTPRange{Offset: offset},
	})
	if err != nil {
		if bloberror.HasCode(err, bloberror.InvalidRange) {
			log.Printf("%s is already fully synced", info.Name())
			return nil
		}
		return fmt.Errorf("failed to download blob stream: %w", err)
	}
	defer resp.Body.Close()

	bar := progressbar.DefaultBytes(*resp.ContentLength, fmt.Sprintf("syncing %s", info.Name()))
	_, err = io.Copy(io.MultiWriter(dest, bar), resp.Body)
	if err != nil {
		return fmt.Errorf("failed to copy data to destination file: %w", err)
	}

	return nil
}
