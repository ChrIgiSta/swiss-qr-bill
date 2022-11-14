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
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ChrIgiSta/swiss-qr-bill/bill"
	"github.com/ChrIgiSta/swiss-qr-bill/mail"
	"github.com/ChrIgiSta/swiss-qr-bill/qr"
	"github.com/ChrIgiSta/swiss-qr-bill/specs"
	"github.com/ChrIgiSta/swiss-qr-bill/utils"
)

const (
	QR_OUT                = `test_qr.png`
	PDF_OUT_NO_SUBMISSION = `test_pdf.pdf`
	PDF_OUT_SUBMISSION    = `test_pdf_subm.pdf`

	PDF_TEST_SUBMISSION = `graphics/pdf-bill-example.pdf`

	QR_SHOULD = `SPC
0200
1
CH0011112222333344446
S
Mr Testing Cowboy
Pinky Range 56
With Cows
3456
Behind the Mountains
Switzerland







674.45
EUR
S
Mr Dont Pay
Unknow Street 4

1998
BelowTheBridge
BeautyIland
NON

DONT USE
EPD
`
)

func TestToken(t *testing.T) {
	testToken := "d4bd63463818bcf87ebed218a1cbeae46e5bc3d97345790f65d2a2082443c5f52afc"
	sha256TestToken := "252d5ef815195e21576464c702bdc53bf1f44ca73a44c3c77f5fb454469efb43"

	genSha := utils.GetSha256([]byte(testToken))
	if string(genSha) != sha256TestToken {
		t.Error("sha256: ", string(genSha), "!=", string(sha256TestToken))
	}

	randToken := utils.CreateNewToken(32)
	randToken2 := utils.CreateNewToken(32)
	if len(randToken.Plain) != 32 {
		t.Error("token len")
	}
	if utils.GetSha256([]byte(randToken.Plain)) != randToken.Sha256 {
		t.Error("token hash")
	}
	if randToken.Plain == randToken2.Plain {
		t.Error("token not random (seed)")
	}
}

func TestUtils(t *testing.T) {
	var (
		IBAN_WRONG   string = "CH00 1111 2222 3333 4444"
		IBAN_WRONG2  string = "C100 1111 2222 3333 4444 4"
		IBAN_CORRECT string = "CH10 2345 6744 2356 3555 2"
		REF_WRONG    string = "56 67000 07803 17847 13400 09017"
		REF_WRONG2   string = "56 67000 07803 17847 13400 090174"
		REF_CORRECT  string = "21 00000 00003 13947 14300 09017"
	)
	err := utils.ValidateIban(IBAN_WRONG)
	if err == nil {
		t.Error("iban validate a wrong iban as correct")
	}
	err = utils.ValidateIban(IBAN_WRONG2)
	if err == nil {
		t.Error("iban validate a wrong iban as correct")
	}
	err = utils.ValidateIban(IBAN_CORRECT)
	if err != nil {
		t.Error("iban validate a correct iban as wrong")
	}

	err = utils.ValidateReference(qr.REFERENCE_TYPE_QR, REF_WRONG)
	if err == nil {
		t.Error("ref validate a wrong iban as correct")
	}
	err = utils.ValidateReference(qr.REFERENCE_TYPE_QR, REF_WRONG2)
	if err == nil {
		t.Error("ref validate a wrong iban as correct")
	}
	err = utils.ValidateReference(qr.REFERENCE_TYPE_QR, REF_CORRECT)
	if err != nil {
		t.Error("ref validate a wrong iban as correct")
	}
}

func TestQr(t *testing.T) {
	issuer := specs.AccountDetails{
		AddressType: qr.ADDRESS_TYPE_STRUCTURED,
		Name:        "Mr Testing Cowboy",
		Address1:    "Pinky Range 56",
		Address2:    "With Cows",
		Zip:         "3456",
		Location:    "Behind the Mountains",
		Country:     "Switzerland",
	}
	receipt := specs.AccountDetails{
		AddressType: qr.ADDRESS_TYPE_STRUCTURED,
		Name:        "Mr Dont Pay",
		Address1:    "Unknow Street 4",
		Zip:         "1998",
		Location:    "BelowTheBridge",
		Country:     "BeautyIland",
	}
	billingDetails := specs.BillingDetails{
		IBAN:           "CH00 1111 2222 3333 4444 6",
		RefenreceType:  qr.REFERENCE_TYPE_NO_REF,
		AdditionalInfo: "DONT USE",
		Currency:       qr.CURRENCY_EURO,
		Amount:         674.45,
	}
	billQr := qr.NewSwissBillQr(&issuer)
	billQr.GetSwissPaymentQR(&receipt, &billingDetails, QR_OUT)

	qrTxtIs, err := utils.ReadQr(QR_OUT)
	if err != nil {
		t.Log(err)
		t.Error("reading qr")
	}
	if strings.Compare(qrTxtIs, QR_SHOULD) != 0 {
		t.Log(qrTxtIs)
		t.Log("\r\n")
		t.Log(QR_SHOULD)
		t.Log("\r\n")
		t.Error("qr not as excepted")
	}
}

func TestPDF(t *testing.T) {
	txt, _ := utils.ReadQr(QR_OUT)
	iban, issuer, receipt, detail, err := utils.EncodeQrText(txt)
	if err != nil {
		t.Error("cannot encode qr", err)
	}

	t.Log(iban)
	err = bill.CreatePDF(issuer, receipt, detail, QR_OUT, PDF_OUT_NO_SUBMISSION,
		utils.GetEnglishTranslationTable(), nil)
	if err != nil {
		t.Error("cannot create billing pdf")
	}
	err = bill.CreatePDF(issuer, receipt, detail, QR_OUT, PDF_OUT_SUBMISSION,
		utils.GetEnglishTranslationTable(), PDF_TEST_SUBMISSION)
	if err != nil {
		t.Error("cannot create billing pdf")
	}
}

func TestClean(t *testing.T) {
	err := os.Remove(QR_OUT)
	if err != nil {
		t.Error("no qr to delete")
	}
	err = os.Remove(PDF_OUT_NO_SUBMISSION)
	if err != nil {
		t.Error("no pdf without submission to delete")
	}
	err = os.Remove(PDF_OUT_SUBMISSION)
	if err != nil {
		t.Error("no pdf with submission to delete")
	}
}

func TestMail(t *testing.T) {
	token := "makeMeQrBill"

	mailClient := mail.NewMailClient(os.Getenv("MAIL_USER"), os.Getenv("MAIL_PASSWORD"),
		os.Getenv("MAIL_SENDER"), os.Getenv("SMTP_HOST"), os.Getenv("POP_HOST"), token)

	msg := mail.Message{
		Subject: token,
		Body: `
		Name: Mister Receipt
		Address1: Hello World 4
		Location: NoLoc
		Zip: 4523
		Country: Schweiz
		Currency: CHF
		Amount: 3445.34
		ReferenceType: NON
		AdditionalInformations: Please pay fast`,
		Attachments:  []mail.Attachments{},
		BodyMimeType: mail.MIME_TYPE_TEXT,
	}
	msg.To = []string{os.Getenv("MAIL_USER")}

	msg.Attachments = append(msg.Attachments, mail.Attachments{
		FileName: "graphics/pdf-bill-example.pdf",
		MimeTyoe: mail.MIME_TYPE_PDF,
	})

	err := mailClient.SendEmail(msg)
	if err != nil {
		t.Error("send mail", err)
	}

	// critical
	time.Sleep(5 * time.Second)

	mails, err := mailClient.GetMails()
	if err != nil {
		t.Error("get mails", err)
	}
	for _, mail := range mails {
		receipt, bDetails, err := mailClient.GetBillingInformationsFromBody(mail.Body)
		// ToDo: Check
		fmt.Println("RECEIPT: *************", receipt)
		fmt.Println("Details", bDetails)
		for _, att := range mail.Attachments {
			if att.FileName != "pdf-bill-example.pdf" {
				t.Error("pdf filename missmatch")
			}
			err = os.Remove(att.FileName)
			if err != nil {
				t.Error("removing pdf", err)
			}
		}
	}
}
