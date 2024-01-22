# Examples

Below is example usage for each of the available functions.

At the bottom of the example you can find the construction of the utApi struct, which is used in the rest of the examples.

```go
package main

import (
	"fmt"
	"github.com/jesses-code-adventures/utapi-go"
	"os"
)

// Example - deleting a file
// Keys is a slice containing the keys returned by uploadthing on upload success
func deleteFileExample(utApi *utapi.UtApi, keys []string) {
	resp, err := utApi.DeleteFiles(keys)
	if err != nil {
		fmt.Println("Error deleting files")
		fmt.Println(fmt.Errorf("%s", err))
	} else {
		fmt.Println("Successfully deleted file")
		fmt.Println(resp.Success)
	}
}

// Example - Getting file urls
// Keys is a slice containing the keys returned by uploadthing on upload success
func getFileUrlsExample(utApi *utapi.UtApi, keys []string) {
	resp, err := utApi.GetFileUrls(keys)
	if err != nil {
		fmt.Println("Error getting file urls")
		fmt.Println(fmt.Errorf("%s", err))
	} else {
		fmt.Println("Successfully got file urls")
		fmt.Println(resp)
	}
}

// Example - Listing your files
func listFilesExample(utApi *utapi.UtApi) {
	opts := utapi.ListFilesOpts{Limit: 10, Offset: 0}
	resp, err := utApi.ListFiles(opts)
	if err != nil {
		fmt.Println("Error listing files")
		fmt.Println(fmt.Errorf("%s", err))
	} else {
		fmt.Println("Successfully listed files")
		fmt.Println(resp)
	}
}

// Example - Getting your usage info
func getUsageInfoExample(utApi *utapi.UtApi) {
	resp, err := utApi.GetUsageInfo()
	if err != nil {
		fmt.Println("Error getting usage info")
		fmt.Println(fmt.Errorf("%s", err))
	} else {
		fmt.Println("Successfully got usage info")
		fmt.Println(resp)
	}
}

// Example - Getting a presigned url
func getSignedUrlExample(utApi *utapi.UtApi, fileKey string, expiresIn int) {
    opts := utapi.PresignedUrlOpts{FileKey: fileKey, ExpiresIn: expiresIn}
    resp, err := utApi.GetPresignedUrl(opts)
    if err != nil {
        fmt.Println("Error getting presigned url")
        fmt.Println(fmt.Errorf("%s", err))
    } else {
        fmt.Println("Successfully got presigned url")
        fmt.Println(resp)
    }
}

// Example - Renaming files.
// This example takes a single old name and a single new name.
// You could aso construct an array of SingleFileRename structs and pass that.
func renameFilesExample(utApi *utapi.UtApi, oldFileName string, newFileName string) {
    singleRename := utapi.SingleFileRename{FileKey: oldFileName, NewName: newFileName}
    opts := utapi.RenameFilesOpts{Updates: []utapi.SingleFileRename{singleRename}}
    err := utApi.RenameFiles(opts)
    if err != nil {
        fmt.Println("Error renaming files")
        fmt.Println(fmt.Errorf("%s", err))
    } else {
        fmt.Println("Successfully renamed files")
    }
}

func main() {
	// Create api handler
	utApi, err := utapi.NewUtApi()
	if err != nil {
		fmt.Println("Error creating uploadthing api handler")
		fmt.Println(fmt.Errorf("%s", err))
		os.Exit(1)
	}
	listFilesExample(utApi)
}
```
