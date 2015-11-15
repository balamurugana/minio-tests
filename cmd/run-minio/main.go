/*
 * Minio Tests (C) 2015 Minio, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
