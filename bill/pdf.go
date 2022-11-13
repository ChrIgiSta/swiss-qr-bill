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

package bill

import (
	"fmt"
	"log"

	"github.com/ChrIgiSta/swiss-qr-bill/qr"
	"github.com/ChrIgiSta/swiss-qr-bill/specs"

	"github.com/signintech/gopdf"
)

const (
	SIX_PT_MM    = 2.12
	SEVEN_PT_MM  = 2.47
	EIGHT_PT_MM  = 2.82
	NINE_PT_MM   = 3.18
	TEN_PT_MM    = 3.52
	ELEVEN_PT_MM = 3.88

	Y_SHIFT = 297 - 105 // position at bottom of A4

	A4_WIDE   = 210
	A4_HEIGHT = 297

	X_RECEIPT = 5
	X_PAYMENT = 62 + 5 + 51 + 1 // +1 security matchin
	FONT_REG  = "liberation-sans"
	FONT_BOLD = "liberation-sans-bold"

	QR_WIDE   = 46
	QR_HEIGHT = 46
)

func loadFonts(pdf *gopdf.GoPdf) error {
	// maybe downoad into container..?
	err := pdf.AddTTFFont(FONT_REG, "ttf/LiberationSans-Regular.ttf")
	if err != nil {
		log.Print(err.Error())
		return err
	}
	err = pdf.AddTTFFont(FONT_BOLD, "ttf/LiberationSans-Bold.ttf")
	if err != nil {
		log.Print(err.Error())
		return err
	}
	return nil
}

func drawBoarder(pdf *gopdf.GoPdf) {
	// outlines
	pdf.SetLineType("dotted")
	pdf.SetLineWidth(0.1)
	pdf.Line(0, Y_SHIFT, A4_WIDE, Y_SHIFT)
	// pdf.Line(0, Y_SHIFT+105, 210, Y_SHIFT+105) // bottom line not needed
	pdf.Line(62, Y_SHIFT, 62, Y_SHIFT+105)
}

func drawBillQr(pdf *gopdf.GoPdf, qrFile string) {
	// qr (add 2mm border to the 46 mm)
	pdf.Image(qrFile, 62+5-2, 5+7+5+Y_SHIFT-2, &gopdf.Rect{
		W: QR_WIDE + 2*2,
		H: QR_HEIGHT + 2*2,
	})
}

func drawScissors(pdf *gopdf.GoPdf) {
	var h float64 = 4

	pdf.Image("graphics/scissors.png", 20, Y_SHIFT-h/2, &gopdf.Rect{
		W: h * 1.6,
		H: h,
	})
}

func addExistingBillPdf(pdf *gopdf.GoPdf, pdfPath string) {
	pdf.SetXY(0, 0)
	tpl1 := pdf.ImportPage(pdfPath, 1, "/MediaBox")
	pdf.UseImportedTemplate(tpl1, 0, 0, A4_WIDE, A4_HEIGHT)
}

func drawTitle(pdf *gopdf.GoPdf, dictionary *specs.TranslationTable) error {
	// titles
	fontSize := 11
	err := pdf.SetFont(FONT_BOLD, "", fontSize)
	if err != nil {
		log.Print(err.Error())
		return err
	}

	pdf.SetXY(X_RECEIPT, Y_SHIFT+5+ELEVEN_PT_MM)
	pdf.Text(dictionary.Receipt)
	pdf.SetXY(62+5, Y_SHIFT+5+ELEVEN_PT_MM)
	pdf.Text(dictionary.PaymentPart)

	return err
}

func drawIssuer(pdf *gopdf.GoPdf, dictionary *specs.TranslationTable,
	issuer *specs.AccountDetails, billingDetails *specs.BillingDetails, fontSizeTitle int,
	fontSizeContent int, spacing float64, x float64) (float64, error) {

	err := pdf.SetFont(FONT_BOLD, "", fontSizeTitle)
	if err != nil {
		log.Print(err.Error())
		return -1, err
	}
	y := Y_SHIFT + 5 + 7 + SIX_PT_MM
	pdf.SetXY(x, y)
	pdf.Text(dictionary.Account)

	err = pdf.SetFont(FONT_REG, "", fontSizeContent)
	if err != nil {
		log.Print(err.Error())
		return y, err
	}
	y += spacing
	pdf.SetXY(x, y)
	pdf.Text(billingDetails.IBAN)
	y += spacing
	pdf.SetXY(x, y)
	pdf.Text(issuer.Name)
	y += spacing
	pdf.SetXY(x, y)
	pdf.Text(issuer.Address1)
	if issuer.Address2 != "" {
		y += spacing
		pdf.SetXY(x, y)
		pdf.Text(issuer.Address2)
	}
	y += spacing
	pdf.SetXY(x, y)
	pdf.Text(issuer.Zip + " " + issuer.Location)

	return y, nil
}

func drawReference(pdf *gopdf.GoPdf, dictionary *specs.TranslationTable,
	billingDetails *specs.BillingDetails, yStart float64, x float64, fontSizeTitle int,
	fontSizeContent int, spacing float64) (float64, error) {

	err := pdf.SetFont(FONT_BOLD, "", fontSizeTitle)
	if err != nil {
		log.Print(err.Error())
		return yStart, err
	}
	yStart += 2 * spacing
	pdf.SetXY(x, yStart)
	pdf.Text(dictionary.Reference)

	err = pdf.SetFont(FONT_REG, "", fontSizeContent)
	if err != nil {
		log.Print(err.Error())
		return yStart, err
	}
	yStart += spacing
	pdf.SetXY(x, yStart)
	pdf.Text(billingDetails.Referece)

	return yStart, nil
}

func drawPayableByReceipt(pdf *gopdf.GoPdf, dictionary *specs.TranslationTable,
	receipt *specs.AccountDetails, yStart float64, x float64, fontSizeTitle int,
	fontSizeContent int, spacing float64) (float64, error) {

	err := pdf.SetFont(FONT_BOLD, "", fontSizeTitle)
	if err != nil {
		log.Print(err.Error())
		return yStart, err
	}
	yStart += 2 * spacing
	pdf.SetXY(x, yStart)
	pdf.Text(dictionary.PayableBy)

	err = pdf.SetFont(FONT_REG, "", fontSizeContent)
	if err != nil {
		log.Print(err.Error())
		return yStart, err
	}
	yStart += spacing
	pdf.SetXY(x, yStart)
	pdf.Text(receipt.Name)
	yStart += spacing
	pdf.SetXY(x, yStart)
	pdf.Text(receipt.Address1)
	if receipt.Address2 != "" {
		yStart += spacing
		pdf.SetXY(x, yStart)
		pdf.Text(receipt.Address2)
	}
	yStart += spacing
	pdf.SetXY(x, yStart)
	pdf.Text(receipt.Zip + " " + receipt.Location)

	return yStart, nil
}

func CreatePDF(issuer *specs.AccountDetails, receipt *specs.AccountDetails,
	billingDetails *specs.BillingDetails, qrFile string, output string,
	dictionary *specs.TranslationTable, existingPdf interface{}) error {

	var (
		err      error
		fontSize float64
		pdf      gopdf.GoPdf
	)

	pdf = gopdf.GoPdf{}
	pdf.Start(gopdf.Config{
		PageSize: gopdf.Rect{W: A4_WIDE, H: A4_HEIGHT},
		Unit:     gopdf.UnitMM,
	})

	pdf.AddPage()

	err = loadFonts(&pdf)
	if err != nil {
		return err
	}

	if existingPdf != nil {
		addExistingBillPdf(&pdf, existingPdf.(string))
	}

	drawBoarder(&pdf)

	drawBillQr(&pdf, qrFile)

	drawScissors(&pdf)

	err = drawTitle(&pdf, dictionary)
	if err != nil {
		return err
	}

	y, err := drawIssuer(&pdf, dictionary, issuer, billingDetails, 6, 8, NINE_PT_MM, X_RECEIPT)
	if err != nil {
		return err
	}

	if billingDetails.RefenreceType != qr.REFERENCE_TYPE_NO_REF {
		y, err = drawReference(&pdf, dictionary, billingDetails, y, X_RECEIPT, 6, 8, NINE_PT_MM)
		if err != nil {
			return err
		}
	}

	_, err = drawPayableByReceipt(&pdf, dictionary, receipt, y, X_RECEIPT, 6, 8, NINE_PT_MM)
	if err != nil {
		return err
	}

	fontSize = 6
	err = pdf.SetFont(FONT_BOLD, "", fontSize)
	if err != nil {
		log.Print(err.Error())
		return err
	}
	y = Y_SHIFT + 5 + 7 + 56 + SIX_PT_MM
	pdf.SetXY(X_RECEIPT, y)
	pdf.Text(dictionary.Currency)
	pdf.SetXY(X_RECEIPT+15, y)
	pdf.Text(dictionary.Amount)
	fontSize = 8
	err = pdf.SetFont(FONT_REG, "", fontSize)
	if err != nil {
		log.Print(err.Error())
		return err
	}
	y += ELEVEN_PT_MM
	pdf.SetXY(X_RECEIPT, y)
	pdf.Text(billingDetails.Currency)
	pdf.SetXY(X_RECEIPT+15, y)
	pdf.Text(fmt.Sprintf("%.2f", billingDetails.Amount))

	fontSize = 6
	err = pdf.SetFont(FONT_BOLD, "", fontSize)
	if err != nil {
		log.Print(err.Error())
		return err
	}

	len, _ := pdf.MeasureTextWidth(dictionary.AcceptancePoint)
	pdf.SetXY(62-5-len, Y_SHIFT+5+7+56+14+ELEVEN_PT_MM)
	pdf.Text(dictionary.AcceptancePoint)

	y, err = drawIssuer(&pdf, dictionary, issuer, billingDetails, 8, 10, ELEVEN_PT_MM, X_PAYMENT)
	if err != nil {
		return err
	}

	if billingDetails.RefenreceType != qr.REFERENCE_TYPE_NO_REF {
		y, err = drawReference(&pdf, dictionary, billingDetails, y, X_PAYMENT, 8, 10, ELEVEN_PT_MM)
		if err != nil {
			return err
		}
	}

	if billingDetails.AdditionalInfo != "" {
		fontSize = 8
		err = pdf.SetFont(FONT_BOLD, "", fontSize)
		if err != nil {
			log.Print(err.Error())
			return err
		}
		y += 2 * ELEVEN_PT_MM
		pdf.SetXY(X_PAYMENT, y)
		pdf.Text(dictionary.AdditionalInfos)
		fontSize = 10
		err = pdf.SetFont(FONT_REG, "", fontSize)
		if err != nil {
			log.Print(err.Error())
			return err
		}
		y += ELEVEN_PT_MM
		pdf.SetXY(X_PAYMENT, y)
		pdf.Text(billingDetails.AdditionalInfo)
	}

	y, err = drawPayableByReceipt(&pdf, dictionary, receipt, y, X_PAYMENT, 8, 10, ELEVEN_PT_MM)
	if err != nil {
		return err
	}

	fontSize = 8
	err = pdf.SetFont(FONT_BOLD, "", fontSize)
	if err != nil {
		log.Print(err.Error())
		return err
	}
	y = Y_SHIFT + 5 + 7 + 56 + ELEVEN_PT_MM
	pdf.SetXY(62+5, y)
	pdf.Text(dictionary.Currency)
	pdf.SetXY(62+5+15, y)
	pdf.Text(dictionary.Amount)
	fontSize = 10
	err = pdf.SetFont(FONT_REG, "", fontSize)
	if err != nil {
		log.Print(err.Error())
		return err
	}
	y += ELEVEN_PT_MM
	pdf.SetXY(62+5, y)
	pdf.Text(billingDetails.Currency)
	pdf.SetXY(62+5+15, y)
	pdf.Text(fmt.Sprintf("%.2f", billingDetails.Amount))

	pdf.WritePdf(output)

	return err
}
