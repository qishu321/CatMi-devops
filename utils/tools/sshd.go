package tools

import (
	"CatMi-devops/request"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

func SshCommand(conf *request.SSHClientConfigReq, command string) (string, error) {
	config := &ssh.ClientConfig{
		User:            conf.UserName,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 忽略know_hosts检查
	}
	switch conf.AuthModel {
	case "PASSWORD":
		config.Auth = []ssh.AuthMethod{ssh.Password(conf.Password)}
	case "PrivateKey":
		logrus.Info("PrivateKey", conf.PrivateKey)
		signer, err := ssh.ParsePrivateKey([]byte(conf.PrivateKey))
		if err != nil {
			return "", err
		}
		config.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", conf.PublicIP, conf.Port), config)
	if err != nil {
		return "失败", err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	output, err := session.CombinedOutput(command)
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// CreateFileOnRemoteServer 在远程服务器上创建文件
func CreateFileOnRemoteServer(sshConfig *request.SSHClientConfigReq, filename, typename, content string) (string, error) {
	absoluteFilePath := "/tmp/" + filename + "." + typename

	// Escape special characters and format the content for bash script
	escapedContent := strings.ReplaceAll(content, "'", `'\''`)

	// Construct the bash script content
	scriptContent := fmt.Sprintf("echo '%s' > %s", escapedContent, absoluteFilePath)

	// Execute the combined script as a single SSH command
	command := fmt.Sprintf("%s && %s %s", scriptContent, typename, absoluteFilePath)

	output, err := SshCommand(sshConfig, command)
	if err != nil {
		fmt.Println("SSH Error:", err)
		return "", fmt.Errorf("Failed to execute SSH command: %v", err)
	}

	fmt.Println("SSH Output:", output) // Optional: Print the SSH output

	return output, nil
}
