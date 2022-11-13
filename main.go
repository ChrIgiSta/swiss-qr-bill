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

package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/ChrIgiSta/swiss-qr-bill/api"
	"github.com/ChrIgiSta/swiss-qr-bill/mail"
	"github.com/ChrIgiSta/swiss-qr-bill/qr"
	"github.com/ChrIgiSta/swiss-qr-bill/specs"
	"github.com/ChrIgiSta/swiss-qr-bill/sql"
	"github.com/ChrIgiSta/swiss-qr-bill/utils"
)

const (
	API_LISTEN_PORT = 3000
)

func main() {
	var wg sync.WaitGroup = sync.WaitGroup{}
	defer wg.Wait()
	defer log.Println("qr app stopped")

	fmt.Println("start qr app")

	connectionString := os.Getenv("SQL_USER") + ":" + os.Getenv("SQL_PASSWORD") + "@" +
		"tcp(" + os.Getenv("SQL_HOST") + ":" + os.Getenv("SQL_PORT") + ")/" + os.Getenv("SQL_DATABASE")
	log.Println("use db: " + connectionString)

	db := sql.NewMariaDb(connectionString)
	err := db.Connect()

	if err != nil {
		log.Fatal("cannot connect to db. ", err)
	}

	ver, err := db.CheckVersion()
	if err != nil && !strings.ContainsAny(err.Error(), "Error 1146:") {
		log.Fatal("cannot get version. ", err)
	}

	// init from env
	primAcc, iban := getPrimaryIssuerFromEnv()
	issuerId := -1
	if primAcc != nil {
		_, issuerId, err = db.InsertIssuer(iban, *primAcc)
		if err != nil {
			log.Println("couldn't add primary issuer from env. ", err.Error())
		} else {
			log.Println("added primary issuer with id ", issuerId)
		}
	}
	mailCnf := getPrimaryMailSettings(issuerId)
	if mailCnf != nil {
		log.Println("insert primary mail config")
		err = db.InsertMailConfig(mailCnf)
		if err != nil {
			log.Println("couldn't add primary email config from env.", err.Error())
		}
	}

	log.Println("starting @ version ", ver)

	// run api server

	log.Println("start api server")
	qrApi := api.NewApi("v1", 3000, db)
	wg.Add(1)
	go qrApi.Run(&wg)

	// run mail cient

	mailCnfs, err := db.GetMailConfigurations()
	if err != nil {
		log.Fatal("cannot get mail configs")
	}

	for _, mailCnf := range mailCnfs {
		if mailCnf != nil {
			if mailCnf.Enable {
				wg.Add(1)
				go mail.ServeMails(*mailCnf, db, &wg, 10)
			}
		}
	}
}

func getPrimaryIssuerFromEnv() (*specs.AccountDetails, string) {
	primatyAccount := specs.AccountDetails{}

	primatyAccount.Name = os.Getenv("FISTNAME") + " " + os.Getenv("LASTNAME")
	primatyAccount.Address1 = os.Getenv("ADDRESS_ROW1")
	primatyAccount.Address2 = os.Getenv("ADDRESS_ROW2")
	primatyAccount.AddressType = qr.ADDRESS_TYPE_STRUCTURED
	primatyAccount.Zip = os.Getenv("ZIP")
	primatyAccount.Location = os.Getenv("LOCATION")
	primatyAccount.Country = os.Getenv("COUNTRY")

	// check, if ok
	if primatyAccount.Name != "" && primatyAccount.Address1 != "" && primatyAccount.Zip != "" &&
		primatyAccount.Location != "" && primatyAccount.Country != "" {
		iban := os.Getenv("IBAN")
		if primatyAccount.Country == qr.COUNTRY_LICHTENSTEIN || primatyAccount.Country == qr.COUNTRY_SWITZERLAND {
			if err := utils.ValidateIban(iban); err == nil {
				return &primatyAccount, iban
			} else {
				log.Println("no valid iban in initial issuer", err.Error())
			}
		}
		log.Println("no valid country set")
	}
	log.Println("no primary account set")
	return nil, ""
}

func getPrimaryMailSettings(issuerId int) *specs.MailConfig {
	mailCnf := specs.MailConfig{}

	mailCnf.Enable = true
	mailCnf.Email = os.Getenv("MAIL_SENDER_ADDRESS")
	mailCnf.IssuerId = issuerId
	mailCnf.Password = os.Getenv("MAIL_PASSWORD")
	mailCnf.Pop3Host = os.Getenv("MAIL_POP3_HOST")
	mailCnf.Pop3Port = 995
	mailCnf.Pop3Secure = true
	mailCnf.SmtpHost = os.Getenv("MAIL_SMTP_HOST")
	mailCnf.SmtpPort = 587
	mailCnf.SmtpSecure = true
	mailCnf.Token = os.Getenv("MAIL_TOKEN")
	mailCnf.UseWhitelist = false

	if mailCnf.Email != "" && mailCnf.Pop3Host != "" && mailCnf.SmtpHost != "" {
		return &mailCnf
	}
	return nil
}
