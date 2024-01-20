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

// Populates the local environment with variables from .env file
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

type uploadthingConfig struct {
	Host    string
	ApiKey  string
	Version string
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

type uploadthingHeaders struct {
	ContentType  string
	ApiKey       string
	SdkVersion   string
	CacheControl string
}

func getDefaultUploadthingHeaders(apiKey string, version string) *uploadthingHeaders {
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

type UtApi struct {
	config *uploadthingConfig
}

func NewUtApi() (*UtApi, error) {
	config, err := getUploadthingConfig()
	if err != nil {
		return nil, err
	}
	return &UtApi{config: config}, nil
}

func getDebugMessage(url string, headers *uploadthingHeaders, body *bytes.Buffer) string {
	return fmt.Sprintf("url: %s, headers: %s, body: %s", url, headers, body.String())
}

type fileKeysPayload struct {
	FileKeys []string `json:"fileKeys"`
}

func setHeaders(req *http.Request, headers *uploadthingHeaders) {
	req.Header.Set("Content-Type", headers.ContentType)
	req.Header.Set("x-uploadthing-api-key", headers.ApiKey)
	req.Header.Set("x-uploadthing-version", headers.SdkVersion)
	req.Header.Set("Cache-Control", "no-store")
}

func (ut *UtApi) requestUploadthing(pathname string, body *bytes.Buffer) (*http.Response, error) {
	url := getUploadthingUrl(pathname, ut.config.Host)
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	headers := getDefaultUploadthingHeaders(ut.config.ApiKey, ut.config.Version)
	setHeaders(req, headers)
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

func (ut *UtApi) DeleteFiles(ids []string) error {
	payload := fileKeysPayload{FileKeys: ids}
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
