package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sk409/goconst"
	"github.com/sk409/goproxy"
)

func entrypointLogin() {
	mkDirIfNotExist(makeFilePath("auth"))
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	name := fs.String("name", "", "user name")
	password := fs.String("password", "", "password")
	fs.Parse(os.Args[2:])
	hashedPassword := fmt.Sprintf("%x", sha512.Sum512([]byte(*password)))
	fetchUsersResponse, err := sendRequest(
		http.MethodGet,
		route(pathUsers),
		map[string]string{
			"name":     *name,
			"password": hashedPassword,
		},
	)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer fetchUsersResponse.Body.Close()
	usersJSONBytes, err := ioutil.ReadAll(fetchUsersResponse.Body)
	if err != nil {
		log.Println(err.Error())
		return
	}
	users := []user{}
	err = json.Unmarshal(usersJSONBytes, &users)
	if len(users) == 0 {
		log.Println(string(usersJSONBytes))
		log.Println("The specified user does not exist")
		return
	}
	user := users[0]
	userJSONBytes, err := json.Marshal(user)
	if err != nil {
		return
	}
	userJSONFilePath := makeFilePath(filepath.Join("auth", "user"))
	userJSONFile, err := os.Create(userJSONFilePath)
	if err != nil {
		return
	}
	defer userJSONFile.Close()
	userJSONFile.Write(userJSONBytes)
}

func entrypointRecord() {
	mkDirIfNotExist(makeFilePath("requests"))
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	host := fs.String("host", "abc", "Vulnerability scan target host")
	// crtFile := fs.String("crt", makeFilePath(filepath.Join("ca", "server.crt")), "Root CA certificae")
	// keyFile := fs.String("key", makeFilePath(filepath.Join("ca", "server.key")), "Root CA private key")
	fs.Parse(os.Args[2:])
	targetHost = *host
	if len(targetHost) == 0 {
		fmt.Println("ホストを指定してください")
		return
	}
	files, err := ioutil.ReadDir("requests")
	if err != nil {
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			jsonFile, err := os.Open(makeFilePath(filepath.Join("requests", file.Name())))
			if err != nil {
				continue
			}
			jsonBytes, err := ioutil.ReadAll(jsonFile)
			if err != nil {
				continue
			}
			var request request
			err = json.Unmarshal(jsonBytes, &request)
			if err != nil {
				continue
			}
			httpRequest, err := request.toHTTPRequestWithSignature("")
			if err != nil {
				continue
			}
			storeRequestSignature(newRequestSignature(httpRequest))
		}
	}
	crtFilePath := makeFilePath(filepath.Join("ca", "server.crt"))
	keyFilePath := makeFilePath(filepath.Join("ca", "server.key"))
	p, err := goproxy.NewHTTPProxy(crtFilePath, keyFilePath)
	if err != nil {
		panic(err)
	}
	p.Hooks.Request = hookRequest
	// p.Hooks.Response = hookResponse
	http.ListenAndServe("0.0.0.0:4080", p)

}

func entrypointRequest() {
	if len(os.Args) <= 2 {
		return
	}
	subcommand := os.Args[2]
	switch subcommand {
	case "send":
		entrypointRequestSend()
	}
}

func entrypointRequestSend() {
	if f, err := os.Stat(".sdso"); os.IsNotExist(err) || !f.IsDir() {
		return
	}
	requestDirectory := filepath.Join(".sdso", "request")
	files, err := filepath.Glob(filepath.Join(requestDirectory, "*.sh"))
	if err != nil {
		return
	}
	cwd, err := os.Getwd()
	if err != nil {
		return
	}
	for _, file := range files {
		// log.Print(file)
		command := exec.Command(filepath.Join(cwd, file))
		command.Dir = requestDirectory
		// command.Env = append(
		// 	os.Environ(),
		// 	"SDSO_CRT="+makeFilePath(filepath.Join("ca", "server.crt")),
		// 	"SDSO_KEY="+makeFilePath(filepath.Join("ca", "server.key")),
		// )
		err = command.Run()
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func entrypointScan() {
	mkDirIfNotExist(makeFilePath("vulnerabilities"))
	fs := flag.NewFlagSet("scan", flag.ExitOnError)
	projectName := fs.String("project", "", "project name")
	verbose := fs.Bool("verbose", false, "verbose")
	fs.Parse(os.Args[2:])
	if len(*projectName) == 0 {
		fmt.Println("--projectオプションを指定してください")
		return
	}
	user, err := loadUser()
	if err != nil {
		fmt.Println("ログインしてください")
		return
	}
	userIDString := strconv.Itoa(int(user.ID))
	exist, err := sendRequestExist(
		http.MethodGet,
		route(pathProjectsExist),
		map[string]string{
			"user_id": userIDString,
			"name":    *projectName,
		},
	)
	if err != nil {
		return
	}
	if !exist {
		fmt.Println("プロジェクトが存在していません")
		return
	}
	storeScanResponse, err := sendRequest(
		http.MethodPost,
		route(pathScans),
		map[string]string{
			"projectName": *projectName,
			"userName":    user.Name,
			"password":    user.Password,
		},
	)
	if err != nil {
		return
	}
	defer storeScanResponse.Body.Close()
	scanJSONBytes, err := ioutil.ReadAll(storeScanResponse.Body)
	if err != nil {
		return
	}
	var scan map[string]interface{}
	err = json.Unmarshal(scanJSONBytes, &scan)
	if err != nil {
		return
	}
	scanID, exist := scan["ID"]
	if !exist {
		return
	}
	scanIDUint := uint(scanID.(float64))
	signatures := []signature{
		newXSSSignatureAngleBrackets(),
		newOSCommandInjectionSleep(),
		newSQLInjectionSleep(),
	}
	requestFiles, err := ioutil.ReadDir(makeFilePath("requests"))
	if err != nil {
		return
	}
	fmt.Println("*****スキャンを開始します*****")
	for _, file := range requestFiles {
		if file.IsDir() {
			continue
		}
		jsonFile, err := os.Open(makeFilePath(filepath.Join("requests", file.Name())))
		if err != nil {
			continue
		}
		jsonBytes, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			continue
		}
		var request request
		json.Unmarshal(jsonBytes, &request)
		for _, sig := range signatures {
			detected, requestString, responseString, err := sig.diagnostician().diagnose(request)
			if err != nil {
				continue
			}
			if !detected {
				continue
			}
			vulnerability := vulnerability{
				Name:        sig.name(),
				Description: sig.description(),
				Path:        request.Path,
				Method:      request.Method,
				Request:     requestString,
				Response:    responseString,
				ScanID:      scanIDUint,
			}
			if *verbose {
				fmt.Println("-------------------------------------")
				fmt.Print("Type: ", vulnerability.Name)
				fmt.Println()
				fmt.Print("Description: ", vulnerability.Description)
				fmt.Println()
				fmt.Print("Path: ", vulnerability.Path)
				fmt.Println()
				fmt.Print("Method: ", vulnerability.Method)
				fmt.Println()
				fmt.Println("Request:")
				fmt.Print(vulnerability.Request)
				fmt.Println("Response:")
				fmt.Print(vulnerability.Response)
			}
			jsonBytes, err := json.Marshal(vulnerability)
			if err != nil {
				continue
			}
			id, err := uuid.NewUUID()
			if err != nil {
				continue
			}
			file, err := os.Create(makeFilePath(filepath.Join("vulnerabilities", id.String())))
			if err != nil {
				continue
			}
			defer file.Close()
			file.Write(jsonBytes)
		}
	}
	fmt.Println("*****スキャンが終了しました*****")
}

func entrypointUpload() {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	projectName := fs.String("project", "", "project name")
	fs.Parse(os.Args[2:])
	if len(*projectName) == 0 {
		fmt.Println("--projectオプションを指定してください")
		return
	}
	user, err := loadUser()
	if err != nil {
		fmt.Println("ログインしてください")
		return
	}
	userIDString := strconv.Itoa(int(user.ID))
	exist, err := sendRequestExist(
		http.MethodGet,
		route(pathProjectsExist),
		map[string]string{
			"user_id": userIDString,
			"name":    *projectName,
		},
	)
	if err != nil {
		return
	}
	if !exist {
		fmt.Println("プロジェクトが存在していません")
		return
	}
	fileInfos, err := ioutil.ReadDir(makeFilePath("vulnerabilities"))
	if err != nil {
		return
	}
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			continue
		}
		filePath := makeFilePath(filepath.Join("vulnerabilities", fileInfo.Name()))
		jsonBytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			continue
		}
		vulnerability := vulnerability{}
		json.Unmarshal(jsonBytes, &vulnerability)
		params := url.Values{}
		params.Set("name", vulnerability.Name)
		params.Set("description", vulnerability.Description)
		params.Set("path", vulnerability.Path)
		params.Set("method", vulnerability.Method)
		params.Set("request", vulnerability.Request)
		params.Set("response", vulnerability.Response)
		params.Set("scanID", strconv.Itoa(int(vulnerability.ScanID)))
		params.Set("userName", user.Name)
		params.Set("password", user.Password)
		params.Set("projectName", *projectName)
		url := fmt.Sprintf("%s/vulnerabilities", serverOrigin)
		request, err := http.NewRequest(http.MethodPost, url, strings.NewReader(params.Encode()))
		if err != nil {
			continue
		}
		request.Header.Set(goconst.HTTP_HEADER_CONTENT_TYPE, goconst.HTTP_HEADER_CONTENT_TYPE_URLENCODED)
		client := new(http.Client)
		response, err := client.Do(request)
		if err != nil {
			continue
		}
		defer response.Body.Close()
		os.Remove(filePath)
	}
}

func entrypointX509() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	publicKey := privateKey.Public()
	subject := pkix.Name{
		CommonName: "SDSO",
		Country:    []string{"JP"},
	}
	serialNumber, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return
	}
	now := time.Now()
	notBefore := now.Add(-time.Minute)
	notAfter := now.Add(24 * time.Hour * 365)
	tpl := x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               subject,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
	}
	certificate, err := x509.CreateCertificate(rand.Reader, &tpl, &tpl, publicKey, privateKey)
	if err != nil {
		return
	}
	mkDirIfNotExist(makeFilePath("ca"))
	crtFile, err := os.Create(makeFilePath(filepath.Join("ca", "server.crt")))
	if err != nil {
		return
	}
	defer crtFile.Close()
	err = pem.Encode(crtFile, &pem.Block{Type: "CERTIFICATE", Bytes: certificate})
	if err != nil {
		return
	}
	keyFile, err := os.Create(makeFilePath(filepath.Join("ca", "server.key")))
	if err != nil {
		return
	}
	defer keyFile.Close()
	derPrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
	err = pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: derPrivateKey})
	if err != nil {
		return
	}
}
