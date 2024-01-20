package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
    "errors"
)

func ValidateEnvironmentVariables(keys []string) error {
    for _, key := range keys {
        if os.Getenv(key) == "" {
            return errors.New(fmt.Sprintf("%s not set", key))
        }
    }
    return nil
}

type UploadthingConfig struct {
	Host    string
	ApiKey  string
	Version string
}

func getUploadthingConfig() (*UploadthingConfig, error) {
	err := ValidateEnvironmentVariables([]string{"UPLOADTHING_SECRET"})
	if err != nil {
		return nil, err
	}
    return &UploadthingConfig{Host: "https://uploadthing.com", ApiKey: os.Getenv("UPLOADTHING_SECRET"), Version: "6.2.0"}, nil
}

type UploadthingHeaders struct {
    ContentType string
	ApiKey     string
	SdkVersion string
    CacheControl string
}

func getDefaultUploadthingHeaders(apiKey string, version string) *UploadthingHeaders {
    return &UploadthingHeaders{ContentType: "application/json", ApiKey: apiKey, SdkVersion: version, CacheControl: "no-store"}
}

func getUploadthingUrl(pathname string, host string) string {
    if !strings.HasPrefix(pathname, "/") {
        pathname = fmt.Sprintf("/%s", pathname)
    }
    if !strings.HasPrefix(pathname, "/api") {
        pathname = fmt.Sprintf("/api%s", pathname)
    }
    url := fmt.Sprintf("%s%s", host, pathname)
    return url
}


type UtApi struct {
    config *UploadthingConfig
}

func NewUtApi() (*UtApi, error) {
    config, err := getUploadthingConfig()
    if err != nil {
        return nil, err
    }
    return &UtApi{config: config}, nil
}

func getDebugMessage(url string, headers *UploadthingHeaders, body *bytes.Buffer) string {
    return fmt.Sprintf("url: %s, headers: %s, body: %s", url, headers, body.String())
}


type FileKeysPayload struct {
    FileKeys []string `json:"fileKeys"`
}


func (ut *UtApi) requestUploadthing(pathname string, body *bytes.Buffer) (*http.Response, error) {
    url := getUploadthingUrl(pathname, ut.config.Host)
    headers := getDefaultUploadthingHeaders(ut.config.ApiKey, ut.config.Version)
    req, err := http.NewRequest(http.MethodPost, url, body)
    if err != nil {
        return nil, err
    }

    // Set headers
    req.Header.Set("Content-Type", headers.ContentType)
    req.Header.Set("x-uploadthing-api-key", headers.ApiKey)
    req.Header.Set("x-uploadthing-version", headers.SdkVersion)
    req.Header.Set("Cache-Control", "no-store")

    // Create client and send request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        resp_body := bytes.NewBuffer([]byte{})
        _, err := io.Copy(resp_body, resp.Body)
        if err != nil {
            return nil, fmt.Errorf("failed to delete files, status code: %d, body: %s, req: %s", resp.StatusCode, resp_body, getDebugMessage(url, headers, body))
        } else {
            return nil, fmt.Errorf("failed to delete files, status code: %d, req: %s", resp.StatusCode, getDebugMessage(url, headers, body))
        }
    }

    responseBytes, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    fmt.Println(string(responseBytes))

    return resp, nil
}


func (ut *UtApi) DeleteUploadthingFiles(ids []string) error {
    payload := FileKeysPayload{FileKeys: ids}
    idsJson, err := json.Marshal(payload)
    if err != nil {
        return err
    }
    body := bytes.NewBuffer(idsJson)
    _, err = ut.requestUploadthing("/api/deleteFile", body)
    if err != nil {
        return err
    }
    return nil
}

