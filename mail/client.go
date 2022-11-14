/**
 * Copyright © 2022, Staufi Tech - Switzerland
 * All rights reserved.
 *
 *  THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 *  AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 *  IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 *  ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
 *  LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 *  CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 *  SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 *  INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 *  CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 *  ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 *  POSSIBILITY OF SUCH DAMAGE.
 */

package mail

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/ChrIgiSta/swiss-qr-bill/qr"
	"github.com/ChrIgiSta/swiss-qr-bill/specs"
	"github.com/knadh/go-pop3"
	gomail "gopkg.in/mail.v2"
)

const (
	MIME_TYPE_TEXT = "text/plain"
	MIME_TYPE_HTML = "text/html"
	MIME_TYPE_JSON = "application/json"
	MIME_TYPE_PDF  = "application/pdf"

	MAX_DOWNLOAD_SIZE = 100000000
)

type Client struct {
	Username   string
	Password   string
	From       string
	SmtpSecure bool
	ImapSecure bool
	SmtpServer string
	SmtpPort   uint16
	Pop3Server string
	Pop3Port   uint16
	Token      string
}

type Attachments struct {
	FileName string
	MimeTyoe string
}
type Message struct {
	To           []string
	CC           []string
	BCC          []string
	Subject      string
	Body         string
	BodyMimeType string
	Attachments  []Attachments
}

func NewMailClient(username string, password string, from string,
	smtpHost string, imapHost string, token string) *Client {
	return &Client{
		Username:   username,
		Password:   password,
		From:       from,
		SmtpSecure: true,
		ImapSecure: true,
		SmtpServer: smtpHost,
		SmtpPort:   587, // 465 SSL/TLS | 587 STARTTLS
		Pop3Server: imapHost,
		Pop3Port:   995, // 993 Imap | 995 pop3
		Token:      token,
	}
}

func (c *Client) SendEmail(msg Message) error {

	m := gomail.NewMessage()

	m.SetHeader("From", c.From)
	m.SetHeader("To", msg.To...)

	if len(msg.CC) > 0 {
		m.SetHeader("Cc", msg.CC...)
	}
	if len(msg.BCC) > 0 {
		m.SetHeader("Bcc", msg.BCC...)
	}

	m.SetHeader("Subject", msg.Subject)

	m.SetBody(msg.BodyMimeType, msg.Body)

	d := gomail.NewDialer(c.SmtpServer, int(c.SmtpPort), c.Username, c.Password)

	if msg.Attachments != nil && len(msg.Attachments) > 0 {
		for _, a := range msg.Attachments {
			m.Attach(a.FileName)
		}
	}

	// d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	// d.SSL = c.SmtpSecure
	// d.StartTLSPolicy = gomail.MandatoryStartTLS

	return d.DialAndSend(m)
}

func (c *Client) GetMails() ([]Message, error) {
	var messages []Message = []Message{}
	var idsToDelete []int = []int{}

	opt := pop3.Opt{
		Host:       c.Pop3Server,
		Port:       int(c.Pop3Port),
		TLSEnabled: true,
	}

	pop := pop3.New(opt)

	conn, err := pop.NewConn()
	if err != nil {
		return nil, err
	}
	defer conn.Quit()

	if err := conn.Auth(c.Username, c.Password); err != nil {
		return nil, err
	}

	count, size, _ := conn.Stat()
	if size > MAX_DOWNLOAD_SIZE {
		return nil, errors.New("max download size.")
	}

	for id := 1; id <= count; id++ {
		m, _ := conn.Retr(id)

		if m.Header.Get("subject") == c.Token {
			log.Println("qr mail found")
			idsToDelete = append(idsToDelete, id)

			to := []string{m.Header.Get("from")}
			body, err := io.ReadAll(m.Body)

			if err != nil {
				return nil, err
			}
			msg := Message{
				To:      to,
				Subject: m.Header.Get("subject"),
			}

			newL := strings.Index(string(body), "\n")
			hash := string(body)[0:newL]
			fmt.Println(hash)
			splitedMimes := strings.Split(string(body), hash)

			for _, mime := range splitedMimes {
				cT := getKey(mime, "Content-Type")
				log.Println("Content Type of MIME: ", cT)
				if strings.Contains(cT, MIME_TYPE_PDF) {
					log.Printf("pdf found")
					enc := getKey(mime, "Content-Transfer-Encoding")
					content, _ := extractContentFromHeader(mime)
					if strings.Contains(enc, "base64") {
						log.Println("base64 encoded")
						dec, err := base64.StdEncoding.DecodeString(content)
						if err != nil {
							return nil, err
						}
						content = string(dec)
					}
					start := strings.Index(cT, `name=`)
					end := strings.Index(cT, `.pdf`)
					fileName := cT[start+6 : end+4]
					os.WriteFile(fileName, []byte(content), 0644)
					msg.Attachments = []Attachments{}
					msg.Attachments = append(msg.Attachments, Attachments{
						FileName: fileName,
						MimeTyoe: MIME_TYPE_PDF,
					})
				} else if strings.Contains(cT, MIME_TYPE_TEXT) {
					content, _ := extractContentFromHeader(mime)
					msg.BodyMimeType = MIME_TYPE_TEXT
					msg.Body = content
					log.Println(content)

				} else if strings.Contains(cT, MIME_TYPE_JSON) {
					content, _ := extractContentFromHeader(mime)
					msg.BodyMimeType = MIME_TYPE_JSON
					msg.Body = content
					log.Println(content)
				}
			}

			messages = append(messages, msg)
		}
	}

	for _, idToDel := range idsToDelete {
		conn.Dele(idToDel)
	}

	return messages, nil
}

func getKey(msg string, key string) string {
	start := strings.Index(msg, key+":")
	if start == -1 {
		return ""
	}

	stop := strings.Index(msg[start:], "\n")
	if stop == -1 {
		return ""
	}

	stop += start
	start += len(key + ": ")

	return msg[start:stop]
}

func extractContentFromHeader(part string) (string, string) {
	var (
		err             error
		line            string
		contentSel      bool = false
		header, content string
		lineNum         int = 0
	)

	reader := strings.NewReader(part)

	in := bufio.NewReader(reader)

	for err == nil {
		line, err = in.ReadString('\n')
		if (line == "\n" || line == "\r\n" || line == "\n\r") && lineNum > 0 {
			contentSel = true
		} else if !strings.Contains(line, "--") {
			if contentSel {
				content += line
			} else {
				header += line
			}
		}
		lineNum++
	}
	fmt.Println(err)
	return strings.Replace(content, line, "", 1), header
}

func (c *Client) GetBillingInformationsFromBody(body string) (specs.AccountDetails, specs.BillingDetails, error) {
	account := specs.AccountDetails{
		Name:        removeNewlines(getKey(body, "Name")),
		AddressType: qr.ADDRESS_TYPE_STRUCTURED,
		Address1:    removeNewlines(getKey(body, "Address1")),
		Address2:    removeNewlines(getKey(body, "Address2")),
		Zip:         removeNewlines(getKey(body, "Zip")),
		Location:    removeNewlines(getKey(body, "Location")),
		Country:     removeNewlines(getKey(body, "Country")),
	}

	amount, err := strconv.ParseFloat(getKey(body, "Amount"), 64)

	billingDetails := specs.BillingDetails{
		RefenreceType:  removeNewlines(getKey(body, "ReferenceType")),
		Referece:       removeNewlines(getKey(body, "Reference")),
		AdditionalInfo: removeNewlines(getKey(body, "AdditionalInformations")),
		Currency:       removeNewlines(getKey(body, "Currency")),
		Amount:         amount,
	}

	return account, billingDetails, err
}

func removeNewlines(in string) string {
	out := strings.ReplaceAll(in, "\r", "")
	out = strings.ReplaceAll(out, "\n", "")
	return out
}
