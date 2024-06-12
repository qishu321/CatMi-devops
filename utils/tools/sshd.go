package tools

import (
	"CatMi-devops/request"
	"fmt"

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
		return "", err
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
