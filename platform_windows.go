package main

import (
	"os"
	"os/exec"
	"path/filepath"
)

func GetPlatform() (Platform, error) {
	appDataDir := os.Getenv("APPDATA")
	dataDir := filepath.Join(appDataDir, "VocabMaster")

	err := os.MkdirAll(dataDir, os.ModePerm)
	if err != nil {
		return Platform{}, err
	}

	return Platform{
		DataDir: appDataDir,
		CertDir: filepath.Join(dataDir, "cert"),
		Font:    filepath.Join(dataDir, "font", "red_bean.ttf"),
	}, nil
}

func (receiver Platform) OpenCertDir() {
	cmd := exec.Command("explorer", receiver.CertDir)
	cmd.Start()
}
