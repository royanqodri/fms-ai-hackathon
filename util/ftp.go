package util

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jlaffaye/ftp"
	"github.com/royanqodri/fms-ai-hackathon/config"
)

// UploadFileToFTP uploads a file to FTP server
func UploadFileToFTP(ctx *gin.Context, filePath string, fileName string, fileLocal string) (string, error) {
	defer func() {
		err := os.Remove(fileLocal)
		if err != nil {
		}
	}()

	ftpHost := config.Get().FTP.Host
	ftpPort := config.Get().FTP.Port
	ftpUser := config.Get().FTP.User
	ftpPassword := config.Get().FTP.Password
	// ftpPath := config.Get().FTP.Path

	// Combine remote path and filename
	finalRemotePath := path.Join(filePath, fileName)
	finalRemotePath = strings.ReplaceAll(finalRemotePath, "\\", "/")

	// Open the local file for reading
	localFile, err := os.Open(fileLocal)
	if err != nil {
		return "", err
	}
	defer localFile.Close()

	// Convert ftpPort to string
	ftpPortStr := strconv.Itoa(int(ftpPort))

	// Connect to FTP server
	ftpAddress := fmt.Sprintf("%s:%s", ftpHost, ftpPortStr)
	conn, err := ftp.Dial(ftpAddress)
	if err != nil {
		return "", fmt.Errorf("failed to dial FTP server: %v", err)
	}
	defer conn.Quit()

	// Login to FTP server
	err = conn.Login(ftpUser, ftpPassword)
	if err != nil {
		return "", fmt.Errorf("failed to login to FTP server: %v", err)
	}

	// Upload the file to the FTP server
	err = conn.Stor(finalRemotePath, localFile)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to FTP server: %v", err)
	}

	// Construct the final URL
	ftpURL := fmt.Sprintf("%s", finalRemotePath)
	ftpURL = strings.ReplaceAll(ftpURL, "//", "/")

	return ftpURL, nil
}

// DeleteFileFromFTP deletes a file from FTP server
func DeleteFileFromFTP(filePath string) error {
	ftpHost := config.Get().FTP.Host
	ftpPort := config.Get().FTP.Port
	ftpUser := config.Get().FTP.User
	ftpPassword := config.Get().FTP.Password

	// Convert ftpPort to string
	ftpPortStr := strconv.Itoa(int(ftpPort))

	// Connect to FTP server
	ftpAddress := fmt.Sprintf("%s:%s", ftpHost, ftpPortStr)
	conn, err := ftp.Dial(ftpAddress)
	if err != nil {
		return fmt.Errorf("failed to dial FTP server: %v", err)
	}
	defer conn.Quit()

	// Login to FTP server
	err = conn.Login(ftpUser, ftpPassword)
	if err != nil {
		return fmt.Errorf("failed to login to FTP server: %v", err)
	}

	// Delete the file from the FTP server
	err = conn.Delete(filePath)
	if err != nil {
		return fmt.Errorf("failed to delete file from FTP server: %v", err)
	}

	return nil
}

// maxRetries is the number of attempts to retry the connection
const maxRetries = 3

// retryDelay is the time to wait between retries
const retryDelay = 5 * time.Second

// connectToFtp tries to connect to the FTP server with retries
func connectToFtp(ftpHost string, ftpPort int, ftpUser, ftpPassword string) (*ftp.ServerConn, error) {
	ftpAddress := fmt.Sprintf("%s:%d", ftpHost, ftpPort)

	for i := 0; i < maxRetries; i++ {
		conn, err := ftp.Dial(ftpAddress, ftp.DialWithTimeout(30*time.Second))
		if err == nil {
			err = conn.Login(ftpUser, ftpPassword)
			if err == nil {
				// Successful connection and login
				return conn, nil
			}
			conn.Quit()
		}

		// Log and retry
		fmt.Printf("Attempt %d to connect to FTP server failed: %v\n", i+1, err)
		time.Sleep(retryDelay)
	}

	return nil, fmt.Errorf("failed to connect to FTP server after %d attempts", maxRetries)
}

// Example usage in your function
func GetFileFromFTP(remoteFilePath string, fileName string) (string, error) {
	ftpHost := config.Get().FTP.Host
	ftpPort := config.Get().FTP.Port
	ftpUser := config.Get().FTP.User
	ftpPassword := config.Get().FTP.Password

	// Connect to FTP server with retries
	conn, err := connectToFtp(ftpHost, int(ftpPort), ftpUser, ftpPassword)
	if err != nil {
		return "", fmt.Errorf("failed to connect to FTP server: %v", err)
	}
	defer conn.Quit()

	// Ensure the remoteFilePath uses the correct path separator
	remoteFilePath = strings.ReplaceAll(remoteFilePath, "\\", "/")
	dir, file := path.Split(remoteFilePath)

	// Change working directory
	err = conn.ChangeDir(dir)
	if err != nil {
		return "", fmt.Errorf("failed to change directory on FTP server: %v", err)
	}

	// Retrieve the file from FTP server
	retrievedFile, err := conn.Retr(file)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve file from FTP server: %v", err)
	}
	defer retrievedFile.Close()

	// Read the file content into a buffer
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, retrievedFile)
	if err != nil {
		return "", fmt.Errorf("failed to copy file content: %v", err)
	}

	// Encode the file content to base64
	base64Data := base64.StdEncoding.EncodeToString(buf.Bytes())

	return base64Data, nil
}

// GetFileFromFTPWithFileName retrieves a file from the FTP server, saves it locally, and returns its base64 representation, file name, and local file path.
func GetFileFromFTPWithFileName(remoteFilePath string, fileName string) (string, string, string, error) {
	ftpHost := config.Get().FTP.Host
	ftpPort := config.Get().FTP.Port
	ftpUser := config.Get().FTP.User
	ftpPassword := config.Get().FTP.Password

	// Convert ftpPort to string
	ftpPortStr := strconv.Itoa(int(ftpPort))

	// Combine host and port
	ftpAddress := fmt.Sprintf("%s:%s", ftpHost, ftpPortStr)

	// Connect to FTP server
	conn, err := ftp.Dial(ftpAddress)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to dial FTP server: %v", err)
	}
	defer conn.Quit()

	// Login to FTP server
	err = conn.Login(ftpUser, ftpPassword)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to login to FTP server: %v", err)
	}

	// Ensure the remoteFilePath uses the correct path separator
	remoteFilePath = strings.ReplaceAll(remoteFilePath, "\\", "/")

	// Separate directory and file name
	dir, file := path.Split(remoteFilePath)

	// Change working directory
	err = conn.ChangeDir(dir)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to change directory on FTP server: %v", err)
	}

	// Retrieve the file from FTP server
	retrievedFile, err := conn.Retr(file)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to retrieve file from FTP server: %v", err)
	}
	defer retrievedFile.Close()

	// Read the file content into a buffer
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, retrievedFile)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to copy file content: %v", err)
	}

	// Encode the file content to base64
	base64Data := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Save the file locally
	fileLocalPath := config.Get().FileLocalPath
	localFilePath := filepath.Join(fileLocalPath, fileName)
	err = ioutil.WriteFile(localFilePath, buf.Bytes(), 0644)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to write file locally: %v", err)
	}

	// Return base64 data, file name, and local file path
	return base64Data, fileName, localFilePath, nil
}

// GetFileName extracts the file name from a file path
func GetFileName(filePath string) string {
	_, file := path.Split(filePath)
	return file
}

// CheckFileExistsOnFTP checks if a file exists on the FTP server
func CheckFileExistsOnFTP(filePath string) (bool, error) {
	ftpHost := config.Get().FTP.Host
	ftpPort := config.Get().FTP.Port
	ftpUser := config.Get().FTP.User
	ftpPassword := config.Get().FTP.Password

	// Connect to FTP server
	conn, err := ftp.Dial(fmt.Sprintf("%s:%d", ftpHost, ftpPort))
	if err != nil {
		return false, fmt.Errorf("failed to connect to FTP server: %v", err)
	}
	defer conn.Quit()

	// Login to FTP server
	err = conn.Login(ftpUser, ftpPassword)
	if err != nil {
		return false, fmt.Errorf("failed to login to FTP server: %v", err)
	}

	// Ensure the file path uses the correct path separator
	filePath = strings.ReplaceAll(filePath, "\\", "/")
	dir, file := path.Split(filePath)

	// Change to the directory containing the file
	err = conn.ChangeDir(dir)
	if err != nil {
		return false, fmt.Errorf("failed to change directory on FTP server: %v", err)
	}

	// Try to retrieve the file listing to check if the file exists
	entries, err := conn.List(file)
	if err != nil {
		// If error, check if it's a not-found error or something else
		if strings.Contains(err.Error(), "550") {
			// 550 error is a common FTP code for "file not found"
			return false, nil
		}
		return false, fmt.Errorf("error while checking file: %v", err)
	}

	// If any entries are found, the file exists
	return len(entries) > 0, nil
}
