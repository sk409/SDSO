package main

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func sendMail(username, subject, body string, html bool) error {
	u, err := userRepository.findByName(username)
	if err != nil {
		return err
	}
	client, err := newMailClient(u.Email)
	if err != nil {
		return err
	}
	err = client.send(u.Email, subject, body, html)
	if err != nil {
		return err
	}
	return nil
}

func sendMailMessage(sendername, place, message string, users []user) error {
	regex := regexp.MustCompile("@([a-zA-Z0-9]+)")
	matches := regex.FindAllStringSubmatch(message, -1)
	if 0 < len(matches) {
		file, err := os.Open(pathMailTemplateMessage)
		if err != nil {
			return err
		}
		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}
		text := string(bytes)
		text = strings.Replace(text, "${sender}", sendername, 1)
		text = strings.Replace(text, "${message}", message, 1)
		text = strings.Replace(text, "${place}", place, 1)
		for _, match := range matches {
			for _, u := range users {
				if match[1] != u.Name {
					continue
				}
				message := strings.ReplaceAll(text, "${receiver}", u.Name)
				err := sendMail(u.Name, "test", message, true)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
