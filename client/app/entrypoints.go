package main

import (
	"crypto/rand"
	"crypto/rsa"
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
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/sk409/gofile"

	"golang.org/x/crypto/bcrypt"

	"github.com/sk409/goproxy"
)

func entrypointInit() {
	if !gofile.IsExist(projectDirectory) {
		err := os.Mkdir(projectDirectory, 0755)
		if err != nil {
			return
		}
	}
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	projectname := fs.String("project", "", "project name")
	fs.Parse(os.Args[2:])
	if emptyAny(projectname) {
		return
	}
	c := config{Projectname: *projectname}
	saveJSON(filepathConfig, c)
}

func entrypointLogin() {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	name := fs.String("name", "", "user name")
	password := fs.String("password", "", "password")
	fs.Parse(os.Args[2:])
	//hashedPassword := fmt.Sprintf("%x", sha512.Sum512([]byte(*password)))
	fetchUsersResponse, err := sendRequest(
		http.MethodGet,
		route(pathUsers),
		map[string]string{
			"name": *name,
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
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(*password))
	if err != nil {
		log.Println("Password does not match")
		return
	}
	userJSONBytes, err := json.Marshal(user)
	if err != nil {
		return
	}
	userJSONFilePath := filepath.Join(directoryAuth, "user")
	userJSONFile, err := os.Create(userJSONFilePath)
	if err != nil {
		return
	}
	defer userJSONFile.Close()
	userJSONFile.Write(userJSONBytes)
}

func entrypointRecord() {
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
			jsonFile, err := os.Open(filepath.Join(directoryRequests, file.Name()))
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
	crtFilePath := filepath.Join(directoryCA, "server.crt")
	keyFilePath := filepath.Join(directoryCA, "server.key")
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
	commitSHA1, err := command("git", "rev-parse", "HEAD")
	if err != nil {
		return
	}
	c := config{}
	err = readJSON(filepathConfig, &c)
	if err != nil {
		return
	}
	u := user{}
	err = readJSON(filepathUser, &u)
	s := scan{}
	err = store(pathScans, map[string]interface{}{"commitSHA1": string(commitSHA1), "projectname": c.Projectname, "username": u.Name}, &s)
	if err != nil {
		return
	}
	signatures := []signature{
		newXSSSignatureAngleBrackets(),
		newOSCommandInjectionSleep(),
		newSQLInjectionSleep(),
	}
	requestFiles, err := ioutil.ReadDir(directoryRequests)
	if err != nil {
		return
	}
	count := 0
	for _, requestFile := range requestFiles {
		if requestFile.IsDir() {
			continue
		}
		r := request{}
		err = readJSON(filepath.Join(directoryRequests, requestFile.Name()), &r)
		for _, sig := range signatures {
			detected, requestString, responseString, err := sig.diagnostician().diagnose(r)
			if err != nil {
				continue
			}
			if !detected {
				continue
			}
			vulnerability := vulnerability{
				Name:        sig.name(),
				Description: sig.description(),
				Path:        r.Path,
				Method:      r.Method,
				Request:     requestString,
				Response:    responseString,
				ScanID:      s.ID,
			}
			id, err := uuid.NewUUID()
			if err != nil {
				continue
			}
			saveJSON(filepath.Join(directoryVulnerabilities, id.String()), vulnerability)
			count++
		}
	}
	if count == 0 {
		fmt.Println("脆弱性は検出されませんでした")
	} else {
		fmt.Println(fmt.Sprintf("%d個の脆弱性が検出されました", count))
	}
}

func entrypointUpload() {
	c := config{}
	err := readJSON(filepathConfig, &c)
	if err != nil {
		return
	}
	u := user{}
	err = readJSON(filepathUser, &u)
	if err != nil {
		return
	}
	vulnerabilityFiles, err := ioutil.ReadDir(directoryVulnerabilities)
	if err != nil {
		return
	}
	for _, vulnerabilityFile := range vulnerabilityFiles {
		if vulnerabilityFile.IsDir() {
			continue
		}
		v := vulnerability{}
		vulnerabilityFilepath := filepath.Join(directoryVulnerabilities, vulnerabilityFile.Name())
		err := readJSON(vulnerabilityFilepath, &v)
		if err != nil {
			continue
		}
		data := map[string]interface{}{
			"name":        v.Name,
			"description": v.Description,
			"path":        v.Path,
			"method":      v.Method,
			"request":     v.Request,
			"response":    v.Response,
			"scanID":      strconv.Itoa(int(v.ScanID)),
			"username":    u.Name,
			"projectname": c.Projectname,
		}
		err = store(pathVulnerabilities, data, &v)
		if err != nil {
			log.Println(err)
			continue
		}
		os.Remove(vulnerabilityFilepath)
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
	crtFile, err := os.Create(filepath.Join(directoryCA, "server.crt"))
	if err != nil {
		return
	}
	defer crtFile.Close()
	err = pem.Encode(crtFile, &pem.Block{Type: "CERTIFICATE", Bytes: certificate})
	if err != nil {
		return
	}
	keyFile, err := os.Create(filepath.Join(directoryCA, "server.key"))
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
