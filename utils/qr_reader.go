package utils

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/ChrIgiSta/swiss-qr-bill/qr"
	"github.com/ChrIgiSta/swiss-qr-bill/specs"
	"github.com/liyue201/goqr"
)

func ReadQr(qrPath string) (string, error) {
	out := ""

	imgBin, err := ioutil.ReadFile(qrPath)
	if err != nil {
		return "", err
	}

	imgDec, _, err := image.Decode(bytes.NewReader(imgBin))
	if err != nil {
		return "", err
	}

	qrCodes, err := goqr.Recognize(imgDec)
	if err != nil {
		return "", err
	}
	for _, qrCode := range qrCodes {
		out += fmt.Sprintln(string(qrCode.Payload))
	}

	return out, err
}

func EncodeQrText(qrTxt string) (string, *specs.AccountDetails, *specs.AccountDetails, *specs.BillingDetails, error) {
	var (
		issuer           specs.AccountDetails = specs.AccountDetails{}
		finalBeneficiary specs.AccountDetails = specs.AccountDetails{}
		receipt          specs.AccountDetails = specs.AccountDetails{}
		details          specs.BillingDetails = specs.BillingDetails{}
		iban             string               = ""
	)

	lineNum := 0
	scanner := bufio.NewScanner(strings.NewReader(qrTxt))

	for scanner.Scan() {
		line := scanner.Text()
		switch lineNum {
		case 0:
			if line != qr.QR_TYPE {
				return "", nil, nil, nil, errors.New("no swiss payment code")
			}
		case 1:
			if line != qr.VERSION {
				return "", nil, nil, nil, errors.New("unsupported version")
			}
		case 2:
			if line != strconv.Itoa(qr.UTF8) {
				return "", nil, nil, nil, errors.New("unsupported coding type (not utf8)")
			}
		case 3:
			err := ValidateIban(line)
			if err != nil {
				return "", nil, nil, nil, err
			}
			iban = line
		case 4:
			issuer.AddressType = line
		case 5:
			issuer.Name = line
		case 6:
			issuer.Address1 = line
		case 7:
			issuer.Address2 = line
		case 8:
			issuer.Zip = line
		case 9:
			issuer.Location = line
		case 10:
			issuer.Country = line
		case 11:
			finalBeneficiary.AddressType = line
		case 12:
			finalBeneficiary.Name = line
		case 13:
			finalBeneficiary.Address1 = line
		case 14:
			finalBeneficiary.Address2 = line
		case 15:
			finalBeneficiary.Zip = line
		case 16:
			finalBeneficiary.Location = line
		case 17:
			finalBeneficiary.Country = line
		case 18:
			amount, err := strconv.ParseFloat(line, 64)
			if err != nil {
				return "", nil, nil, nil, err
			}
			details.Amount = amount
		case 19:
			details.Currency = line
		case 20:
			receipt.AddressType = line
		case 21:
			receipt.Name = line
		case 22:
			receipt.Address1 = line
		case 23:
			receipt.Address2 = line
		case 24:
			receipt.Zip = line
		case 25:
			receipt.Location = line
		case 26:
			receipt.Country = line
		case 27:
			details.RefenreceType = line
		case 28:
			details.Referece = line
		case 29:
			details.AdditionalInfo = line
		case 30:
			if line != qr.TRAILER {
				return "", nil, nil, nil, errors.New("no trailer")
			}
		}

		lineNum++
	}
	return iban, &issuer, &receipt, &details, nil
}
