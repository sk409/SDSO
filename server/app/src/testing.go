package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/sk409/gogit"

	"github.com/gorilla/websocket"
	"github.com/sk409/gotype"

	"github.com/go-yaml/yaml"

	"github.com/google/uuid"
	"github.com/sk409/gofile"
)

func runTest(userName, projectName, clonePath string) (bool, error) {
	user := user{}
	db.Where("name = ?", userName).First(&user)
	if db.Error != nil {
		// log.Print(db.Error)
		return false, db.Error
	}
	project := project{}
	db.Where("name = ? AND user_id = ?", projectName, user.ID).First(&project)
	if db.Error != nil {
		//log.Println(db.Error)
		return false, db.Error
	}
	// clonePath := filepath.Join(gitClones.RootDirectoryPath, filepath.Join(userName, projectName))
	uuid, err := uuid.NewUUID()
	if err != nil {
		return false, err
	}
	testPath := filepath.Join(cwd, "..", "testing", uuid.String())
	testAppPath := filepath.Join(testPath, "app")
	err = os.MkdirAll(testAppPath, 0755)
	if err != nil {
		return false, err
	}
	defer os.RemoveAll(testPath)
	git := gogit.NewGit(testAppPath, gitBinPath)
	err = git.Clone(clonePath, ".")
	if err != nil {
		// log.Println(err)
		return false, err
	}
	configFilePath := filepath.Join(testAppPath, ".sdso", "config.yml")
	if !gofile.IsExist(configFilePath) {
		return true, nil
	}
	configFile, err := os.Open(configFilePath)
	if err != nil {
		// log.Println(err.Error())
		return false, err
	}
	defer configFile.Close()
	configFileBytes, err := ioutil.ReadAll(configFile)
	if err != nil {
		// log.Println(err.Error())
		return false, err
	}
	config := config{}
	err = yaml.Unmarshal(configFileBytes, &config)
	if err != nil {
		// log.Println(err.Error())
		return false, err
	}
	workDir := "/app"
	serviceNames := []string{}
	serviceNameRegex := regexp.MustCompile("(.+):.+")
	for index, docker := range config.Jobs.Build.Docker {
		isPrimary := index == 0
		dockerDirectoryPath := testPath
		matches := serviceNameRegex.FindStringSubmatch(docker.Image)
		imageNameComponents := strings.Split(matches[1], "/")
		serviceName := imageNameComponents[len(imageNameComponents)-1]
		serviceNames = append(serviceNames, serviceName)
		if !isPrimary {
			dockerDirectoryPath = filepath.Join(testPath, serviceName)
			err = os.Mkdir(dockerDirectoryPath, 0755)
			if err != nil {
				// log.Println(err.Error())
				return false, err
			}
		}
		dockerfile, err := os.Create(filepath.Join(dockerDirectoryPath, "Dockerfile"))
		if err != nil {
			// log.Println(err.Error())
			return false, err
		}
		defer dockerfile.Close()
		dockerfileText := fmt.Sprintf("FROM %s\n", docker.Image)
		if isPrimary {
			dockerfileText += "COPY app " + workDir + "\n"
			dockerfileText += "WORKDIR " + workDir + "\n"
		}
		dockerfile.Write([]byte(dockerfileText))
	}
	dockerComposeText := "version: \"3.3\"\nservices:\n"
	for index, serviceName := range serviceNames {
		isPrimary := index == 0
		dockerDirectoryPath := testPath
		if !isPrimary {
			dockerDirectoryPath = filepath.Join(testPath, serviceName)
		}
		dockerComposeText += "  " + serviceName + ":\n    build: " + dockerDirectoryPath + "\n" + "    tty: true\n"
	}
	dockercomposeFile, err := os.Create(filepath.Join(testPath, "docker-compose.yml"))
	if err != nil {
		// log.Println(err.Error())
		return false, err
	}
	defer dockercomposeFile.Close()
	dockercomposeFile.Write([]byte(dockerComposeText))
	upCommand := exec.Command("docker-compose", "up", "-d", "--build")
	upCommand.Dir = testPath
	err = upCommand.Run()
	if err != nil {
		// log.Println(err.Error())
		return false, err
	}
	defer func() {
		downCommand := exec.Command("docker-compose", "down")
		downCommand.Dir = testPath
		downCommand.Run()
	}()
	test := test{
		Steps:     len(config.Jobs.Build.Steps),
		ProjectID: project.ID,
	}
	db.Save(&test)
	if db.Error != nil {
		// log.Print(db.Error)
		return false, db.Error
	}
	socketTest, exist := websocketsTest[user.ID]
	if exist {
		jsonBytes, err := json.Marshal(test)
		if err == nil {
			socketTest.WriteMessage(websocket.TextMessage, jsonBytes)
		}
	}
	for _, step := range config.Jobs.Build.Steps {
		if gotype.IsMap(step) {
			m := step.(map[interface{}]interface{})
			for key, value := range m {
				keyString := key.(string)
				valueString := value.(string)
				if keyString == "run" {
					sendToClient := func(testResult *testResult) {
						socket, exist := websocketsTestResult[user.ID]
						if !exist {
							return
						}
						jsonBytes, err := json.Marshal(testResult)
						if err != nil {
							// log.Println(err)
							return
						}
						socket.WriteMessage(websocket.TextMessage, []byte(jsonBytes))
					}
					compltedTestResult := func(status, output string, testResult *testResult) {
						var testStatus testStatus
						db.Where("text = ?", status).First(&testStatus)
						if db.Error != nil {
							// log.Println(db.Error)
							return
						}
						now := time.Now()
						testResult.Output = output
						testResult.TestStatusID = testStatus.ID
						testResult.CompletedAt = &now
						db.Save(&testResult)
						if db.Error != nil {
							// log.Println(db.Error)
							return
						}
					}
					var testStatusRunning testStatus
					db.Where("text = 'running'").First(&testStatusRunning)
					testResult := testResult{
						Command:      valueString,
						TestID:       test.ID,
						TestStatusID: testStatusRunning.ID,
					}
					db.Save(&testResult)
					if db.Error != nil {
						// log.Println(db.Error)
						return false, db.Error
					}
					sendToClient(&testResult)
					log.Println("ID = ", testResult.ID)
					args := []string{"exec", "-T", serviceNames[0]}
					args = append(args, strings.Split(valueString, " ")...)
					execCommand := exec.Command("docker-compose", args...)
					execCommand.Dir = testPath
					var output bytes.Buffer
					execCommand.Stderr = &output
					execCommand.Stdout = &output
					err = execCommand.Run()
					if err != nil {
						// log.Println(err.Error())
						compltedTestResult("failed", output.String(), &testResult)
						sendToClient(&testResult)
						return false, nil
					}
					compltedTestResult("success", output.String(), &testResult)
					sendToClient(&testResult)
				}
			}
		}
	}
	return true, nil
}
