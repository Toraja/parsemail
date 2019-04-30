package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/mail"
	"os"
	"runtime"
	"strings"
	"sync"

	"gopkg.in/gomail.v2"
)

func main() {
	// "singlemail.eml"
	fmt.Printf("Start #go routine: %d\n", runtime.NumGoroutine())
	defer fmt.Printf("End #go routine: %d\n", runtime.NumGoroutine())
	to := []string{"bob@example.com", "cora@example.com", "david@example.com", "ema@example.com", "fox@example.com", "giraffe@example.com", "hen@example.com", "ice@example.com"}
	// to := []string{"bob@example.com"}

	// if err := send(to...); err != nil {
	// if err := send2(to...); err != nil {
	if err := send3(to...); err != nil {
		panic(err)
	}
}

func newDialer() gomail.Dialer {
	return gomail.Dialer{Host: "localhost", Port: 1025}
}

func send1(to ...string) error {
	d := newDialer()
	sc, err := d.Dial()
	if err != nil {
		return err
	}
	defer sc.Close()

	errCh := make(chan error)
	limitCh := make(chan struct{}, 1) // parallel work causes error "nested MAIL command"
	go func() {
		var wg sync.WaitGroup
		for i := range to {
			wg.Add(1)
			limitCh <- struct{}{}
			go func(to1 string) {
				defer func() {
					<-limitCh
					wg.Done()
				}()
				m := newMail(to1)
				if err := gomail.Send(sc, m); err != nil {
					errCh <- err
					return
				}
			}(to[i])
		}
		wg.Wait()
		close(errCh)
	}()
	for err := range errCh {
		if err != nil {
			return err
		}
	}

	return nil
}

func send2(to ...string) error {
	d := newDialer()
	sc, err := d.Dial()
	if err != nil {
		return err
	}
	defer sc.Close()

	msgs := make([]*gomail.Message, len(to))
	for i := range to {
		msgs[i] = newMail(to[i])
	}

	chunk := 5
	for len(msgs) != 0 {
		if len(msgs) < chunk {
			chunk = len(msgs)
		}
		if err := gomail.Send(sc, msgs[:chunk]...); err != nil {
			return err
		}
		msgs = msgs[chunk:]
	}

	return nil
}

func send3(to ...string) error {
	d := newDialer()
	sc, err := d.Dial()
	if err != nil {
		return err
	}
	defer sc.Close()

	msgs := make([]*gomail.Message, len(to))
	for i := range to {
		msgs[i] = newMail(to[i])
	}

	if err := gomail.Send(sc, msgs...); err != nil {
		return err
	}

	return nil
}
func dialAndSend1(to ...string) error {
	d := newDialer()

	errCh := make(chan error)
	limitCh := make(chan struct{}, 5)
	go func() {
		var wg sync.WaitGroup
		for i := range to {
			wg.Add(1)
			limitCh <- struct{}{}
			go func(to1 string) {
				defer func() {
					<-limitCh
					wg.Done()
				}()
				m := newMail(to1)
				if err := d.DialAndSend(m); err != nil {
					errCh <- err
					return
				}
			}(to[i])
		}
		wg.Wait()
		close(errCh)
	}()
	for err := range errCh {
		if err != nil {
			return err
		}
	}

	return nil
}

func dialAndSend2(to ...string) error {
	msgs := make([]*gomail.Message, len(to))
	for i := range to {
		msgs[i] = newMail(to[i])
	}

	d := newDialer()
	if err := d.DialAndSend(msgs...); err != nil {
		return err
	}

	return nil
}

func newMail(to string) *gomail.Message {
	m := gomail.NewMessage()
	m.SetHeader("From", "alex@example.com")
	m.SetHeader("To", to)
	// m.SetAddressHeader("Cc", "dan@example.com")
	m.SetHeader("Cc", "dan@example.com")
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/plain", "Dear gopher,\nNever give up! Tomorrow always comes.\nBest Regards,\n")
	m.Attach("attachment/xyz.csv")
	return m
}

func getSingleMail(path string) (string, error) {
	return "", nil
}

func parseMail(path string) (*mail.Message, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	msg, err := mail.ReadMessage(file)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func parseMultipart(m *mail.Message) (emlBody, error) {
	mediaType, params, err := mime.ParseMediaType(m.Header.Get("Content-Type"))
	if err != nil {
		return emlBody{}, err
	}

	var eb emlBody
	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(m.Body, params["boundary"])
		for {
			part, err := mr.NextPart()
			if err != nil {
				if err == io.EOF {
					break
				} else {
					return emlBody{}, err
				}
			}
			partBody, err := ioutil.ReadAll(part)
			if part.Header.Get("Content-Disposition") != "" {
				dec, err := base64.StdEncoding.DecodeString(string(partBody))
				// dec := make([]byte, len(partBody))
				// _, err := base64.StdEncoding.Decode(dec, partBody) <- this will have extra bytes and disturbs assertion
				if err != nil {
					return emlBody{}, err
				}
				eb.Attachment = dec
			} else {
				eb.Text = partBody
			}
		}
	}
	return eb, nil
}

type eml struct {
	From       string
	To         string
	Cc         string
	Subject    string
	Body       string
	Attachment string
}

type emlBody struct {
	Text       []byte
	Attachment []byte
}
