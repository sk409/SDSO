package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"

	"github.com/sk409/goproxy"
)

func entrypointInit() {
	os.Mkdir(directoryProject, 0755)
	os.Mkdir(directoryRequests, 0755)
	os.Mkdir(directoryVulnerabilities, 0755)
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	teamname := fs.String("team", "", "team name")
	projectname := fs.String("project", "", "project name")
	fs.Parse(os.Args[2:])
	c := config{Teamname: *teamname, Projectname: *projectname}
	saveJSON(filepathConfig, c)
}

func entrypointLogin() {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	name := fs.String("name", "", "user name")
	password := fs.String("password", "", "password")
	fs.Parse(os.Args[2:])
	users := []user{}
	err := fetch(pathUsers, map[string]interface{}{"name": *name}, &users)
	if err != nil {
		return
	}
	if len(users) == 0 {
		return
	}
	u := users[0]
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(*password))
	if err != nil {
		return
	}
	err = saveJSON(userJSONFilePath, u)
	if err != nil {
		return
	}
}

func entrypointPush() {
	// fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	command("git", "push", os.Args[2], os.Args[3])
}

func entrypointRecord() {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	host := fs.String("host", "", "")
	fs.Parse(os.Args[2:])
	targetHost = *host
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

func entrypointSetup() {
	directories := []string{directoryApp, directoryAuth, directoryCA, directoryRequests, directoryVulnerabilities}
	for _, directory := range directories {
		os.Mkdir(directory, 0755)
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
	err = store(pathScans, map[string]interface{}{"commitSHA1": string(commitSHA1), "username": u.Name, "password": u.Password, "teamname": c.Teamname, "projectname": c.Projectname}, &s)
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
			"teamname":    c.Teamname,
			"projectname": c.Projectname,
			"username":    u.Name,
			"password":    u.Password,
		}
		err = store(pathVulnerabilities, data, &v)
		if err != nil {
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
