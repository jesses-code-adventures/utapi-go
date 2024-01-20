package utapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"net/http"
	"os"
	"strings"
)

type uploadthingConfig struct {
	Host    string
	ApiKey  string
	Version string
}

type fileKeysPayload struct {
	FileKeys []string `json:"fileKeys"`
}

type uploadthingHeaders struct {
	ContentType  string
	ApiKey       string
	SdkVersion   string
	CacheControl string
}

func setEnvironmentVariablesFromFile() error {
	return godotenv.Load(".env")
}

func handleSetEnvironmentVariables() error {
	err := setEnvironmentVariablesFromFile()
	if err != nil {
		fmt.Printf("Couldn't set environment variables from file")
		return err
	}
	return nil
}

func validateEnvironmentVariables(keys []string) error {
	for _, key := range keys {
		if os.Getenv(key) == "" {
			return errors.New(fmt.Sprintf("%s not set", key))
		}
	}
	return nil
}

func getUploadthingConfig() (*uploadthingConfig, error) {
	err := handleSetEnvironmentVariables()
	if err != nil {
		return nil, err
	}
	err = validateEnvironmentVariables([]string{"UPLOADTHING_SECRET"})
	if err != nil {
		return nil, err
	}
	return &uploadthingConfig{Host: "https://uploadthing.com", ApiKey: os.Getenv("UPLOADTHING_SECRET"), Version: "6.2.0"}, nil
}

func getUploadthingHeaders(apiKey string, version string) *uploadthingHeaders {
	return &uploadthingHeaders{ContentType: "application/json", ApiKey: apiKey, SdkVersion: version, CacheControl: "no-store"}
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

func getDebugMessage(url string, headers *uploadthingHeaders, body *bytes.Buffer) string {
	return fmt.Sprintf("url: %s, headers: %s, body: %s", url, headers, body.String())
}

func setHeaders(req *http.Request, headers *uploadthingHeaders) {
	req.Header.Set("Content-Type", headers.ContentType)
	req.Header.Set("x-uploadthing-api-key", headers.ApiKey)
	req.Header.Set("x-uploadthing-version", headers.SdkVersion)
	req.Header.Set("Cache-Control", "no-store")
}

type UtApi struct {
	config     *uploadthingConfig
	httpClient *http.Client
}

func NewUtApi() (*UtApi, error) {
	config, err := getUploadthingConfig()
	if err != nil {
		return nil, err
	}
	return &UtApi{config: config, httpClient: &http.Client{}}, nil
}

func (ut *UtApi) requestUploadthing(pathname string, body *bytes.Buffer) (*http.Response, error) {
	url := getUploadthingUrl(pathname, ut.config.Host)
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	headers := getUploadthingHeaders(ut.config.ApiKey, ut.config.Version)
	setHeaders(req, headers)
	resp, err := ut.httpClient.Do(req)
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
	return resp, nil
}

func (ut *UtApi) DeleteFiles(ids []string) (*http.Response, error) {
	payload := fileKeysPayload{FileKeys: ids}
	idsJson, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	body := bytes.NewBuffer(idsJson)
    response, err := ut.requestUploadthing("/api/deleteFile", body)
	if err != nil {
		return nil, err
	}
	return response, nil
}
