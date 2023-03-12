package apttransport

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type AptClient interface {
	Get(url string) (*http.Response, error)
}

type AptMethod struct {
	Client    AptClient
	AptString string
	Main      func()
}

func (a *AptMethod) GetAptString() string {
	if a.AptString == "" {
		return "http://"
	}
	return a.AptString
}

func (a *AptMethod) GetClient() AptClient {
	if a.Client == nil {
		return http.DefaultClient
	}
	return a.Client
}

func (a *AptMethod) Run() {
	c := make(chan *Message)
	go a.output(c)
	a.sendCapabilities(c)

	stdin := bufio.NewScanner(os.Stdin)
	for stdin.Scan() {
		line := stdin.Text()
		if line == "" {
			continue
		}
		fmt.Println(line)
		s := strings.SplitN(line, " ", 2)
		code, err := strconv.Atoi(s[0])
		if err != nil {
			fmt.Println("Malformed message!")
			os.Exit(100)
		}
		request := &Message{
			Status:     line,
			StatusCode: code,
			Header:     Header{},
		}

		for stdin.Scan() {
			line := stdin.Text()

			if line == "" {
				a.process(c, request)
				break
			}
			s := strings.SplitN(line, ": ", 2)
			request.Header.Add(s[0], s[1])
		}
	}
}

func (a *AptMethod) output(c <-chan *Message) {
	for {
		m := <-c
		os.Stdout.Write([]byte(m.String()))
		if m.Exit != 0 {
			os.Exit(m.Exit)
		}
	}
}

func (a *AptMethod) sendCapabilities(c chan<- *Message) {
	caps := &Message{
		Status: "100 Capabilities",
		Header: Header{},
	}

	caps.Header.Add("Version", "1.2")
	caps.Header.Add("Pipeline", "true")
	caps.Header.Add("Send-Config", "true")

	c <- caps
}

func (a *AptMethod) process(c chan<- *Message, m *Message) {
	switch m.StatusCode {
	case 600:
		go a.fetch(c, m)
	case 601:
		// TODO: parse config?
	default:
		fail := &Message{
			Status: "401 General Failure",
			Header: Header{},
			Exit:   100,
		}
		fail.Header.Add("Message", "Status code not implemented")

		c <- fail
	}
}

func (a *AptMethod) fetch(c chan<- *Message, m *Message) {
	uri := m.Header.Get("URI")
	filename := m.Header.Get("Filename")

	// TODO: If-Modified-Since

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		c <- &Message{
			Status: "400 URI Failure",
			Header: Header{
				"Message": []string{"Could not open file: " + err.Error()},
				"URI":     []string{uri},
			},
		}
		return
	}
	defer file.Close()

	// TODO: Fix bug with appending to existing files
	// TODO: implement range requests if file already exists

	realURI := strings.TrimPrefix(uri, a.AptString)
	log.Println("Get: ", realURI)
	resp, err := a.GetClient().Get(realURI)
	if err != nil {
		c <- &Message{
			Status: "400 URI Failure",
			Header: Header{
				"Message": []string{"Could not fetch URI: " + err.Error()},
				"URI":     []string{uri},
			},
		}
		return
	}
	defer resp.Body.Close()

	started := &Message{
		Status: "200 URI Start",
		Header: Header{
			"URI": []string{uri},
		},
	}
	// TODO: add Last-Modified header

	c <- started

	md5Hash := md5.New()
	sha1Hash := sha1.New()
	sha256Hash := sha256.New()
	sha512Hash := sha512.New()

	if _, err = io.Copy(io.MultiWriter(file, md5Hash, sha1Hash, sha256Hash, sha512Hash), resp.Body); err != nil {
		c <- &Message{
			Status: "400 URI Failure",
			Header: Header{
				"Message": []string{"Could not write file: " + err.Error()},
				"URI":     []string{uri},
			},
		}
		return
	}

	success := &Message{
		Status: "201 URI Done",
		Header: Header{},
	}
	success.Header.Add("URI", uri)
	success.Header.Add("Filename", filename)
	// TODO Size, Last-Modified
	md5Hex := hex.EncodeToString(md5Hash.Sum(nil)[:])
	success.Header.Add("MD5-Hash", md5Hex)
	success.Header.Add("MD5Sum-Hash", md5Hex)
	success.Header.Add("SHA1-Hash", hex.EncodeToString(sha1Hash.Sum(nil)[:]))
	success.Header.Add("SHA256-Hash", hex.EncodeToString(sha256Hash.Sum(nil)[:]))
	success.Header.Add("SHA512-Hash", hex.EncodeToString(sha512Hash.Sum(nil)[:]))

	c <- success
}

func (a *AptMethod) DefaultMain() {
	c := make(chan *Message)
	go a.output(c)
	a.sendCapabilities(c)

	stdin := bufio.NewScanner(os.Stdin)
	for stdin.Scan() {
		line := stdin.Text()
		if line == "" {
			continue
		}
		s := strings.SplitN(line, " ", 2)
		code, err := strconv.Atoi(s[0])
		if err != nil {
			fmt.Println("Malformed message!")
			os.Exit(100)
		}
		request := &Message{
			Status:     line,
			StatusCode: code,
			Header:     Header{},
		}

		for stdin.Scan() {
			line := stdin.Text()

			if line == "" {
				a.process(c, request)
				break
			}
			s := strings.SplitN(line, ": ", 2)
			request.Header.Add(s[0], s[1])
		}
	}
}
