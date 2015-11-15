package main

import (
	"io/ioutil"
	"log"
	"os/exec"
)

func runMinioServer(automatedTestDir string) error {
	minioServerCmd := exec.Command("minio", "--anonymous", "server", automatedTestDir)
	return minioServerCmd.Run()
}

func main() {
	automatedTestDir, err := ioutil.TempDir(".", "automated-tests")
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Installing minio server.")
	minioInstallCmd := exec.Command("go", "get", "-u", "github.com/minio/minio")
	err = minioInstallCmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Success.")

	if err := runMinioServer(automatedTestDir); err != nil {
		log.Fatalln(err)
	}
}
