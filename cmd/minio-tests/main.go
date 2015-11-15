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
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var minGolangRuntimeVersion = "1.5.1"

// following code handles the current Golang release styles, we might have to update them in future
// if golang community divulges from the below formatting style.
const (
	betaRegexp = "beta[0-9]"
	rcRegexp   = "rc[0-9]"
)

func getNormalizedGolangVersion() string {
	version := strings.TrimPrefix(runtime.Version(), "go")
	br := regexp.MustCompile(betaRegexp)
	rr := regexp.MustCompile(rcRegexp)
	betaStr := br.FindString(version)
	version = strings.TrimRight(version, betaStr)
	rcStr := rr.FindString(version)
	version = strings.TrimRight(version, rcStr)
	return version
}

type version struct {
	major, minor, patch string
}

func (v1 version) String() string {
	return fmt.Sprintf("%s%s%s", v1.major, v1.minor, v1.patch)
}

func (v1 version) Version() int {
	ver, e := strconv.Atoi(v1.String())
	if e != nil {
		log.Fatal(e)
	}
	return ver
}

func (v1 version) LessThan(v2 version) bool {
	if v1.Version() < v2.Version() {
		return true
	}
	return false
}

func newVersion(v string) version {
	ver := version{}
	verSlice := strings.Split(v, ".")
	ver.major = verSlice[0]
	ver.minor = verSlice[1]
	if len(verSlice) == 3 {
		ver.patch = verSlice[2]
	} else {
		ver.patch = "0"
	}
	return ver
}

func checkGolangRuntimeVersion() {
	log.Println("Checking golang runtime version.")
	v1 := newVersion(getNormalizedGolangVersion())
	v2 := newVersion(minGolangRuntimeVersion)
	if v1.LessThan(v2) {
		log.Fatalln("Old Golang runtime version ‘" + v1.String() + "’ detected, ‘automated-tests’ requires minimum go1.5.1 or later.")
	}
	log.Println("Success.")
}

func checkGolangEnvironment() {
	log.Println("Checking golang GOPATH.")
	if goPath := os.Getenv("GOPATH"); strings.TrimSpace(goPath) == "" {
		log.Fatalln("GOPATH not set, cannot continue please follow https://github.com/minio/mc/blob/master/INSTALLGO.md.")
	} else {
		if !strings.Contains(os.Getenv("PATH"), goPath) {
			log.Fatalln("GOPATH not part of PATH, cannot continue please follow https://github.com/minio/mc/blob/master/INSTALLGO.md.")
		}
	}
	log.Println("Success.")
}

func checkIfServerRunning() {
	log.Println("Checking if server is running.")
	resp, err := http.Get("http://localhost:9000")
	if err != nil {
		log.Fatalln(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalln("Server replied back with %s.", resp.Status)
	}
	log.Println("Success.")
}

func verifyRuntime() {
	checkGolangRuntimeVersion()
	checkGolangEnvironment()
	checkIfServerRunning()
}

func updateMinioClient() {
	log.Println("Installing minio client.")
	mcInstallCmd := exec.Command("go", "get", "-u", "github.com/minio/mc")
	err := mcInstallCmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Success.")
}

func main() {
	verifyRuntime()
	if len(os.Args) < 1 {
		log.Fatal("Invalid number of arguments please make sure \n\t$ ./minio-tests <alias> <corpus-datadir>\n")
	}
	err := runTests(os.Args[1], os.Args[2])
	if err != nil {
		log.Fatalln(err)
	}
}

func mcCmd(args ...string) (err error) {
	log.Printf("Running test %s\n", args[0])
	newArgs := []string{"--json"}
	newArgs = append(newArgs, args...)
	cmd := exec.Command("mc", newArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed, %s.\n", output)
		return err
	}
	log.Println("Success.")
	return nil
}

const (
	recursive = "..."
)

func runTests(serverAlias, directory string) error {
	err := mcCmd("ls", serverAlias)
	if err != nil {
		return err
	}
	err = mcCmd("mb", filepath.Join(serverAlias, "testbucket"))
	if err != nil {
		return err
	}
	err = mcCmd("access", "set", "readonly", filepath.Join(serverAlias, "testbucket"))
	if err != nil {
		return err
	}
	err = mcCmd("cp", directory+recursive, filepath.Join(serverAlias, "testbucket"))
	if err != nil {
		return err
	}
	err = mcCmd("rm", "--force", filepath.Join(serverAlias, "testbucket")+recursive)
	if err != nil {
		return err
	}
	return nil
}
