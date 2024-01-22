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

// Utility types
type uploadthingFileStatus int

const (
	DeletionPending uploadthingFileStatus = iota
	Failed
	Uploaded
	Uploading
)

func createUploadthingFileStatus(status string) uploadthingFileStatus {
	switch status {
	case "Deletion Pending":
		return DeletionPending
	case "Failed":
		return Failed
	case "Uploaded":
		return Uploaded
	case "Uploading":
		return Uploading
	default:
		return Failed
	}
}

func (s uploadthingFileStatus) String() string {
	return [...]string{"Deletion Pending", "Failed", "Uploaded", "Uploading"}[s]
}

func (s uploadthingFileStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// Types for the uploadthing api to consume

type uploadthingConfig struct {
	Host    string
	ApiKey  string
	Version string
}

type uploadthingHeaders struct {
	ContentType  string
	ApiKey       string
	SdkVersion   string
	CacheControl string
}

type fileKeysPayload struct {
	FileKeys []string `json:"fileKeys"`
}

type SingleFileRename struct {
	FileKey string `json:"fileKey"`
	NewName string `json:"newName"`
}

// Arguments for the list files endpoint
type ListFilesOpts struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// Arguments for the rename files endpoint
type RenameFilesOpts struct {
	Files []SingleFileRename `json:"files"`
}

// Arguments for the presigned url endpoint
type PresignedUrlOpts struct {
	FileKey   string `json:"fileKey"`
	ExpiresIn int    `json:"expiresIn"`
}

// Types to decode from the uploadthing api responses

// Full response object for a delete file action
type DeleteFileResponse struct {
	Success bool `json:"success"`
}

func parseDeleteFileResponse(resp *http.Response) (DeleteFileResponse, error) {
	if resp == nil {
		return DeleteFileResponse{}, fmt.Errorf("response is nil")
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return DeleteFileResponse{}, fmt.Errorf("error reading response body: %v", err)
	}

	var response DeleteFileResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return DeleteFileResponse{}, fmt.Errorf("error unmarshaling response: %v", err)
	}

	return response, nil
}

// Represents a single uploadthing file
type uploadthingFile struct {
	Key    string                `json:"key"`
	Id     string                `json:"id"`
	Status uploadthingFileStatus `json:"status"`
}

// Represents a full response struct for a list of files
type UploadthingFileResponse struct {
	Files []uploadthingFile `json:"files"`
}

func parseUploadthingFileResponse(resp *http.Response) (UploadthingFileResponse, error) {
	if resp == nil {
		return UploadthingFileResponse{}, fmt.Errorf("response is nil")
	}
	defer resp.Body.Close()
	fmt.Println(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return UploadthingFileResponse{}, fmt.Errorf("error reading response body: %v", err)
	}
	var response UploadthingFileResponse
	json.Unmarshal(body, &response)
	if err != nil {
		return UploadthingFileResponse{}, fmt.Errorf("error unmarshaling response: %v", err)
	}
	return response, nil
}

// Represents a single uploadthing url
type UploadthingUrl struct {
	Url string `json:"url"`
	Key string `json:"key"`
}

// Represents a full response struct for a list of urls
type UploadthingUrlsResponse struct {
    Data []UploadthingUrl `json:"data"`
}

func parseUploadthingUrlsResponse(resp *http.Response) (UploadthingUrlsResponse, error) {
	if resp == nil {
		return UploadthingUrlsResponse{}, fmt.Errorf("response is nil")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return UploadthingUrlsResponse{}, fmt.Errorf("error reading response body: %v", err)
	}
	var response UploadthingUrlsResponse
	json.Unmarshal(body, &response)
	if err != nil {
		return UploadthingUrlsResponse{}, fmt.Errorf("error unmarshaling response: %v", err)
	}
	return response, nil
}

// Represents a full response struct for uploadthing usage info
type UploadthingUsageInfo struct {
	TotalBytes       int     `json:"totalBytes"`
	TotalReadable    string  `json:"totalReadable"`
	AppTotalBytes    float32 `json:"appTotalBytes"`
	AppTotalReadable string  `json:"appTotalReadable"`
	FilesUploaded    int     `json:"filesUploaded"`
	LimitBytes       float32 `json:"limitBytes"`
	LimitReadable    string  `json:"limitReadable"`
}

func parseUploadthingUsageInfoResponse(resp *http.Response) (UploadthingUsageInfo, error) {
	if resp == nil {
		return UploadthingUsageInfo{}, fmt.Errorf("response is nil")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return UploadthingUsageInfo{}, fmt.Errorf("error reading response body: %v", err)
	}
	var response UploadthingUsageInfo
	json.Unmarshal(body, &response)
	if err != nil {
		return UploadthingUsageInfo{}, fmt.Errorf("error unmarshaling response: %v", err)
	}
	return response, nil
}

func (u *UploadthingUsageInfo) AsString() string {
    return fmt.Sprintf("Total Bytes: %d\nTotal Readable: %s\nApp Total Bytes: %f\nApp Total Readable: %s\nFiles Uploaded: %d\nLimit Bytes: %f\nLimit Readable: %s", u.TotalBytes, u.TotalReadable, u.AppTotalBytes, u.AppTotalReadable, u.FilesUploaded, u.LimitBytes, u.LimitReadable)
}

// Represents a full response struct for a presigned url
type PresignedUrlResponse struct {
	Url string `json:"url"`
}

func parsePresignedUrlResponse(resp *http.Response) (PresignedUrlResponse, error) {
	if resp == nil {
		return PresignedUrlResponse{Url: ""}, fmt.Errorf("response is nil")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return PresignedUrlResponse{Url: ""}, fmt.Errorf("error reading response body: %v", err)
	}

	var response PresignedUrlResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return PresignedUrlResponse{Url: ""}, fmt.Errorf("error unmarshaling response: %v", err)
	}
	return response, err
}

// manage environment

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

// functionality

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

func getUploadthingHeaders(apiKey string, version string) *uploadthingHeaders {
	return &uploadthingHeaders{ContentType: "application/json", ApiKey: apiKey, SdkVersion: version, CacheControl: "no-store"}
}

func setHeaders(req *http.Request, headers *uploadthingHeaders) {
	req.Header.Set("Content-Type", headers.ContentType)
	req.Header.Set("x-uploadthing-api-key", headers.ApiKey)
	req.Header.Set("x-uploadthing-version", headers.SdkVersion)
	req.Header.Set("Cache-Control", "no-store")
}

// UtApi - Interact with the uploadthing api.
// This struct is designed to replicate UTApi from the uploadthing typescript sdk.
// Please note that responses are encoded into structs that mirror the current json.
// Any errors are returned as-is.
type UtApi struct {
	config     *uploadthingConfig
	httpClient *http.Client
}

// Construct an instance of the UtApi struct.
// This will read the UPLOADTHING_SECRET environment variable from the .env file.
// If you don't have UPLOADTHING_SECRET set, the function will throw.
func NewUtApi() (*UtApi, error) {
	config, err := getUploadthingConfig()
	if err != nil {
		return nil, err
	}
	return &UtApi{config: config, httpClient: &http.Client{}}, nil
}

// Handler for all uploadthing requests.
// If this handler isn't used, you'll need to ensure that the correct headers are set in your separate solution.
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
	if resp.StatusCode != http.StatusOK {
		resp_body := bytes.NewBuffer([]byte{})
		_, err := io.Copy(resp_body, resp.Body)
		if err != nil {
			return nil, fmt.Errorf("uploadthing request failed, status code: %d, body: %s, req: %s", resp.StatusCode, resp_body, getDebugMessage(url, headers, body))
		} else {
			return nil, fmt.Errorf("uploadthing request failed, status code: %d, req: %s", resp.StatusCode, getDebugMessage(url, headers, body))
		}
	}
	return resp, nil
}

// Delete files from uploadthing.
func (ut *UtApi) DeleteFiles(fileKeys []string) (*DeleteFileResponse, error) {
	payload := fileKeysPayload{FileKeys: fileKeys}
	fileKeysJson, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	body := bytes.NewBuffer(fileKeysJson)
	utResponse, err := ut.requestUploadthing("/api/deleteFile", body)
	if err != nil {
		return nil, err
	}
	response, err := parseDeleteFileResponse(utResponse)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// Given an array of file keys, get the corresponding urls.
func (ut *UtApi) GetFileUrls(fileKeys []string) (*UploadthingUrlsResponse, error) {
	payload := fileKeysPayload{FileKeys: fileKeys}
	fileKeysJson, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	body := bytes.NewBuffer(fileKeysJson)
    fmt.Println(body)
	utResponse, err := ut.requestUploadthing("/api/getFileUrl", body)
	if err != nil {
		return nil, err
	}
	response, err := parseUploadthingUrlsResponse(utResponse)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// List files stored in uploadthing.
func (ut *UtApi) ListFiles(opts ListFilesOpts) (*UploadthingFileResponse, error) {
	optsJson, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}
	body := bytes.NewBuffer(optsJson)
	utResponse, err := ut.requestUploadthing("/api/listFiles", body)
	if err != nil {
		return nil, err
	}
	response, err := parseUploadthingFileResponse(utResponse)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// Rename files in uploadthing.
// No response is returned, but an error is returned if the request fails.
// This is in line with the behaviour of the uploadthing typescript sdk.
func (ut *UtApi) RenameFiles(files RenameFilesOpts) error {
	optsJson, err := json.Marshal(files)
	if err != nil {
		return err
	}
	body := bytes.NewBuffer(optsJson)
	_, err = ut.requestUploadthing("/api/renameFiles", body)
	if err != nil {
		return err
	}
	return nil
}

// Get usage info for the current uploadthing account.
func (ut *UtApi) GetUsageInfo() (*UploadthingUsageInfo, error) {
	utResponse, err := ut.requestUploadthing("/api/getUsageInfo", bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, err
	}
	response, err := parseUploadthingUsageInfoResponse(utResponse)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// Generate a presigned url for a file.
// expiresIn should be a duration in seconds.
// The maximum value for expiresIn is 604800 (7 days).
// You must accept overrides on the UploadThing dashboard for expiresIn to be accepted.
func (ut *UtApi) GetPresignedUrl(opts PresignedUrlOpts) (string, error) {
	if opts.ExpiresIn > 604800 {
		return "", errors.New("expiresIn must be less than 604800")
	}
	optsJson, err := json.Marshal(opts)
	if err != nil {
		return "", err
	}
	utResponse, err := ut.requestUploadthing("/api/requestFileAccess", bytes.NewBuffer(optsJson))
	if err != nil {
		return "", err
	}
	response, err := parsePresignedUrlResponse(utResponse)
	if err != nil {
		return "", err
	}
	return response.Url, nil
}
