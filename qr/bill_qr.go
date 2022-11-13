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

package qr

import (
	"bufio"
	"fmt"
	"image"
	"log"
	"os"
	"strings"

	"github.com/ChrIgiSta/swiss-qr-bill/specs"
	"github.com/divan/qrlogo"
)

const (
	SWISS_CROSS_FILE = "graphics/CH-Kreuz_7mm.png"

	QR_TYPE     = "SPC"  // Swiss Payment Code
	VERSION     = "0200" // Version 2.0.0
	UTF8        = 1
	CODING_TYPE = UTF8

	ADDRESS_TYPE_STRUCTURED = "S"
	ADDRESS_TYPE_COMBINED   = "K"

	// CountryCode ISO 3166-1
	COUNTRY_SWITZERLAND  = "CH"
	COUNTRY_LICHTENSTEIN = "LI"

	CURRENCY_SWISS_FRANCS = "CHF"
	CURRENCY_EURO         = "EUR"

	REFERENCE_TYPE_NO_REF   = "NON"
	REFERENCE_TYPE_CREDITOR = "SCOR"
	REFERENCE_TYPE_QR       = "QRR"

	TRAILER = "EPD"
)

type SwissBillQr struct {
	QRType     string
	Version    string
	CodingType int
	// issuer parts
	Issuer           *specs.AccountDetails
	FinalBeneficiary *specs.AccountDetails
}

func NewSwissBillQr(issuer *specs.AccountDetails) *SwissBillQr {
	return &SwissBillQr{
		QRType:           QR_TYPE,
		Version:          VERSION,
		CodingType:       CODING_TYPE,
		Issuer:           issuer,
		FinalBeneficiary: &specs.AccountDetails{},
	}
}

func (s *SwissBillQr) GetSwissPaymentQR(receipt *specs.AccountDetails,
	billingDetails *specs.BillingDetails, outFile string) {

	qrTxt := s.getSwissPaymentText(receipt, billingDetails)
	s.qrWithLogo(qrTxt, outFile)
}

func (s *SwissBillQr) getSwissPaymentText(receipt *specs.AccountDetails,
	billingDetails *specs.BillingDetails) string {

	swissPaymentTxt := fmt.Sprintf(`%s
%s
%d
%s
%s
%s
%s
%s
%s
%s
%s
%s
%s
%s
%s
%s
%s
%s
%.2f
%s
%s
%s
%s
%s
%s
%s
%s
%s
%s
%s
%s`,
		s.QRType, s.Version, s.CodingType,
		strings.ReplaceAll(billingDetails.IBAN, " ", ""),

		s.Issuer.AddressType, s.Issuer.Name, s.Issuer.Address1, s.Issuer.Address2,
		s.Issuer.Zip, s.Issuer.Location, s.Issuer.Country, // ZE

		s.FinalBeneficiary.AddressType, s.FinalBeneficiary.Name,
		s.FinalBeneficiary.Address1, s.FinalBeneficiary.Address2,
		s.FinalBeneficiary.Zip, s.FinalBeneficiary.Location, s.FinalBeneficiary.Country, // EZE

		billingDetails.Amount, billingDetails.Currency,

		receipt.AddressType, receipt.Name, receipt.Address1, receipt.Address2,
		receipt.Zip, receipt.Location, receipt.Country, //EZP

		billingDetails.RefenreceType, strings.ReplaceAll(billingDetails.Referece, " ", ""),

		billingDetails.AdditionalInfo,

		TRAILER,
	)
	return swissPaymentTxt
}

func (s *SwissBillQr) qrWithLogo(txt string, outFile string) {
	inFile, err := os.Open(SWISS_CROSS_FILE)
	if err != nil {
		log.Println("cannot read logo file.", err)
		return
	}
	defer inFile.Close()
	logo, _, err := image.Decode(inFile)
	if err != nil {
		log.Println("cannot decode in file (logo) as image.", err)
	}
	buf, err := qrlogo.Encode(txt, logo, 1024)
	if err != nil {
		log.Println("cannot encode qr with logo.", err)
		return
	}

	file, err := os.Create(outFile)
	if err != nil {
		log.Println("cannot create output file.", err)
		return
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	_, err = writer.Write(buf.Bytes())
	if err != nil {
		log.Println("cannot write to file.", err)
		return
	}
	writer.Flush()
}
