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
	"log"
	"sync"
	"time"

	"github.com/ChrIgiSta/swiss-qr-bill/qr"
	"github.com/ChrIgiSta/swiss-qr-bill/specs"
	"github.com/ChrIgiSta/swiss-qr-bill/sql"
)

func ServeMails(mailConfig specs.MailConfig, db *sql.Db, wg *sync.WaitGroup, interval int) {
	log.Println("ToDo: not all configurations from email config set")
	client := NewMailClient(mailConfig.Username, mailConfig.Password, mailConfig.Email, mailConfig.SmtpHost, mailConfig.Pop3Host, mailConfig.Token)

	log.Println("mailer started")

	for true {
		mails, err := client.GetMails()
		if err != nil {
			log.Println("error while reciving email", err)
			goto pass
		}
		for _, mail := range mails {
			// gen qr and bill
			iban, issuer, err := db.GetIssuer(mailConfig.IssuerId)
			if err != nil {
				log.Println("error bet issuer from db", err)
			}
			receipt, billingDetails, err := client.GetBillingInformationsFromBody(mail.Body)
			q := qr.NewSwissBillQr(&issuer)
			billingDetails.IBAN = iban
			q.GetSwissPaymentQR(&receipt, &billingDetails, "mail_generated_qr.png")
			if len(mail.Attachments) > 0 {
				// get pdf name
			}
			// send back
			// delete pdf
			log.Println(mail)
		}

	pass:
		time.Sleep(time.Duration(interval) * time.Second)
	}

	log.Println("mailer exited")
}
