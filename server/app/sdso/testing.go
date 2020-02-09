package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

type tester struct {
}

func (t *tester) checkout(clonePath, testPath, branchname string) error {
	git := gogit.NewGit(testPath, gitBinPath)
	err := git.Clone(clonePath, ".", "-b", branchname)
	return err
}

func (t *tester) config(testPath string) (*config, error) {
	configFilePath := filepath.Join(testPath, ".sdso", "config.yml")
	if !gofile.IsExist(configFilePath) {
		return nil, errNotExist
	}
	configFile, err := os.Open(configFilePath)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()
	configFileBytes, err := ioutil.ReadAll(configFile)
	if err != nil {
		return nil, err
	}
	c := config{}
	err = yaml.Unmarshal(configFileBytes, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (t *tester) execTestCommand(test test, testPath, primaryServicename, command string) error {
	testStatusRunning := testStatus{}
	err := first(map[string]interface{}{"text": testStatusRunningText}, &testStatusRunning)
	if err != nil {
		return err
	}
	testResult := testResult{
		Command:  command,
		TestID:   test.ID,
		StatusID: testStatusRunning.ID,
	}
	gormDB.Save(&testResult)
	if gormDB.Error != nil {
		return err
	}
	t.sendTest(test.ID)
	args := []string{"exec", "-T", primaryServicename}
	args = append(args, strings.Split(command, " ")...)
	cmd := exec.Command("docker-compose", args...)
	cmd.Dir = testPath
	var output bytes.Buffer
	cmd.Stderr = &output
	cmd.Stdout = &output
	failed := cmd.Run()
	status := testStatusSuccessText
	if failed != nil {
		status = testStatusFailedText
	}
	testStatus := testStatus{}
	err = first(map[string]interface{}{"text": status}, &testStatus)
	if err != nil {
		return err
	}
	now := time.Now()
	testResult.Output = output.String()
	testResult.StatusID = testStatus.ID
	testResult.CompletedAt = &now
	gormDB.Save(&testResult)
	if gormDB.Error != nil {
		return err
	}
	t.sendTest(test.ID)
	return failed
}

func (t *tester) makeDockerfiles(tmpPath string, c *config) ([]string, error) {
	workDir := "/app"
	servicenames := []string{}
	servicenameRegex := regexp.MustCompile("(.+):.+")
	for index, docker := range c.Jobs.Build.Docker {
		isPrimary := index == 0
		dockerDirectoryPath := tmpPath
		matches := servicenameRegex.FindStringSubmatch(docker.Image)
		imageNameComponents := strings.Split(matches[1], "/")
		servicename := imageNameComponents[len(imageNameComponents)-1]
		servicenames = append(servicenames, servicename)
		if !isPrimary {
			dockerDirectoryPath = filepath.Join(tmpPath, servicename)
			err := os.Mkdir(dockerDirectoryPath, 0755)
			if err != nil {
				return nil, err
			}
		}
		dockerfile, err := os.Create(filepath.Join(dockerDirectoryPath, "Dockerfile"))
		if err != nil {
			return nil, err
		}
		defer dockerfile.Close()
		dockerfileText := fmt.Sprintf("FROM %s\n", docker.Image)
		if isPrimary {
			dockerfileText += "COPY app " + workDir + "\n"
			dockerfileText += "WORKDIR " + workDir + "\n"
		}
		dockerfile.Write([]byte(dockerfileText))
	}
	return servicenames, nil
}

func (t *tester) makeTestDirectory() (string, string, error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return "", "", err
	}
	tmpPath := filepath.Join(cwd, "..", "tmp", uuid.String())
	testPath := filepath.Join(tmpPath, "app")
	err = os.MkdirAll(testPath, 0755)
	if err != nil {
		return "", "", err
	}
	return tmpPath, testPath, nil
}

func (t *tester) project(teamname, projectname string) (*project, error) {
	team := team{}
	err := first(map[string]interface{}{"name": teamname}, &team)
	if err != nil {
		return nil, err
	}
	p := project{}
	err = first(map[string]interface{}{"name": projectname, "team_id": team.ID}, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (t *tester) run(teamname, projectname, clonePath, branchname, commitSHA1 string) (bool, error) {
	p, err := t.project(teamname, projectname)
	if err != nil {
		return false, err
	}
	tmpPath, testPath, err := t.makeTestDirectory()
	if err != nil {
		return false, err
	}
	defer os.RemoveAll(tmpPath)
	err = t.checkout(clonePath, testPath, branchname)
	if err != nil {
		return false, err
	}
	c, err := t.config(testPath)
	if err != nil {
		return false, err
	}
	servicenames, err := t.makeDockerfiles(tmpPath, c)
	if err != nil {
		return false, err
	}
	err = t.upDockerCompose(tmpPath, servicenames)
	if err != nil {
		return false, err
	}
	defer func() {
		downCommand := exec.Command("docker-compose", "down")
		downCommand.Dir = testPath
		downCommand.Run()
	}()
	//************
	// test := test{
	// 	Steps:      len(c.Jobs.Build.Steps),
	// 	Branchname: branchname,
	// 	CommitSHA1: commitSHA1,
	// 	ProjectID:  p.ID,
	// }
	// gormDB.Save(&test)
	// if gormDB.Error != nil {
	// 	return false, gormDB.Error
	// }
	//************
	test := test{}
	err = save(map[string]interface{}{"Steps": len(c.Jobs.Build.Steps), "Branchname": branchname, "CommitSHA1": commitSHA1, "ProjectID": p.ID}, &test)
	if err != nil {
		return false, err
	}
	//************
	err = t.sendTest(test.ID)
	if err != nil {
		return false, err
	}
	t.runSteps(test, c, testPath, servicenames[0])
	return true, nil
}

func (t *tester) runSteps(test test, c *config, testPath, primaryServicename string) error {
	for _, step := range c.Jobs.Build.Steps {
		if gotype.IsMap(step) {
			m := step.(map[interface{}]interface{})
			for key, value := range m {
				keyString := key.(string)
				valueString := value.(string)
				if keyString == "run" {
					err := t.execTestCommand(test, testPath, primaryServicename, valueString)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func (t *tester) sendTest(id uint) error {
	test, err := testRepository.findByID(id, loadAllRelation)
	if err != nil {
		return err
	}
	testBytes, err := json.Marshal(test)
	if err != nil {
		return err
	}
	p := project{}
	err = first(map[string]interface{}{"id": test.ProjectID}, &p)
	if err != nil {
		return err
	}
	team := team{}
	err = first(map[string]interface{}{"id": p.TeamID}, &team)
	if err != nil {
		return err
	}
	teamUsers := []teamUser{}
	err = find(map[string]interface{}{"teamID": team.ID}, &teamUsers)
	if err != nil {
		return err
	}
	for _, teamUser := range teamUsers {
		u := user{}
		err = first(map[string]interface{}{"id": teamUser.UserID}, &u)
		if err != nil {
			continue
		}
		socket, exist := testWebsockets[u.ID]
		if !exist {
			continue
		}
		socket.WriteMessage(websocket.TextMessage, testBytes)
	}
	return nil
}

func (t *tester) upDockerCompose(tmpPath string, servicenames []string) error {
	dockerComposeText := "version: \"3.3\"\nservices:\n"
	for index, servicename := range servicenames {
		isPrimary := index == 0
		dockerDirectoryPath := tmpPath
		if !isPrimary {
			dockerDirectoryPath = filepath.Join(tmpPath, servicename)
		}
		dockerComposeText += "  " + servicename + ":\n    build: " + dockerDirectoryPath + "\n" + "    tty: true\n"
	}
	dockerComposeFile, err := os.Create(filepath.Join(tmpPath, "docker-compose.yml"))
	if err != nil {
		return err
	}
	defer dockerComposeFile.Close()
	dockerComposeFile.Write([]byte(dockerComposeText))
	upCommand := exec.Command("docker-compose", "up", "-d", "--build")
	upCommand.Dir = tmpPath
	return upCommand.Run()
}
