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
	"log"

	"github.com/ChrIgiSta/swiss-qr-bill/bill"
	"github.com/ChrIgiSta/swiss-qr-bill/qr"
	"github.com/ChrIgiSta/swiss-qr-bill/specs"
	"github.com/ChrIgiSta/swiss-qr-bill/utils"
)

const (
	IBAN            = "CH00 0000 0000 0000 0000 0"
	BILLING_MESSAGE = "01-0032-12"

	OUT_QR                = "out/example-qr.png"
	OUT_PDF               = "out/example-bill.pdf"
	OUT_PDF_FROM_BILL_PDF = "out/example-bill-from-pdf.pdf"
)

func main() {

	tr := utils.GetEnglishTranslationTable()

	issuer := specs.AccountDetails{
		AddressType: qr.ADDRESS_TYPE_STRUCTURED,
		Name:        "MyCompany",
		Address1:    "Römerstrasse 45a",
		Address2:    "Postfach",
		Zip:         "5432",
		Location:    "Luftighausen",
		Country:     qr.COUNTRY_SWITZERLAND,
	}

	receipt := specs.AccountDetails{
		AddressType: qr.ADDRESS_TYPE_STRUCTURED,
		Name:        "Hans Mustermann",
		Address1:    "Trämilweg 45",
		Zip:         "1234",
		Location:    "Pfupfighofen",
		Country:     qr.COUNTRY_SWITZERLAND,
	}

	billingDetails := specs.BillingDetails{
		AdditionalInfo: BILLING_MESSAGE,
		IBAN:           "CH54 0900 0000 1585 0580 7",
		RefenreceType:  qr.REFERENCE_TYPE_QR,
		Referece:       "21 00000 00003 13947 14300 09017",
		Currency:       qr.CURRENCY_SWISS_FRANCS,
		Amount:         660.80,
	}

	err := utils.ValidateReference(billingDetails.RefenreceType, billingDetails.Referece)
	if err != nil {
		log.Fatal("invalide reference: ", err)
	}

	err = utils.ValidateIban(billingDetails.IBAN)
	if err != nil {
		log.Fatal("invalide iban: ", err)
	}

	sQr := qr.NewSwissBillQr(&issuer)

	sQr.GetSwissPaymentQR(&receipt, &billingDetails, OUT_QR)

	bill.CreatePDF(&issuer, &receipt, &billingDetails, OUT_QR, OUT_PDF, tr, nil)
	bill.CreatePDF(&issuer, &receipt, &billingDetails, OUT_QR, OUT_PDF_FROM_BILL_PDF, tr, "graphics/pdf-bill-example.pdf")
}
