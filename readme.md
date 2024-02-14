# utapi-go

A thin wrapper for the uploadthing api.

If you'd like to add something you need, feel free to contribute in line with [contributing.md](contributing.md).


## why?

You have uploaded a large file to uploadthing and you'd like to process that big boy in go.

## setup

You will need a .env file with your uploadthing secret key.

```.env
UPLOADTHING_SECRET=sk_*************************
```

## usage

After adding your import statement as below, run go mod tidy.

```go
package main

import (
    "github.com/jesses-code-adventures/utapi-go"
    "os"
    "fmt"
)

func main() {
    // Create api handler
    utApi, err := utapi.NewUtApi()
    if err != nil {
        fmt.Println("Error creating uploadthing api handler")
        fmt.Println(fmt.Errorf("%s", err))
        os.Exit(1)
    }

    // Example - deleting a file
    // This is the key returned by uploadthing when you create a file
    keys := []string{"fc8d296b-20f6-4173-bfa5-5d6c32fc9f6b-geat9r.csv"}
	resp, err := utApi.DeleteFiles(keys)
	if err != nil {
		fmt.Println("Error deleting files")
		fmt.Println(fmt.Errorf("%s", err))
	} else {
		fmt.Println("Successfully deleted file")
		fmt.Println(resp.Success)
	}
}
```

