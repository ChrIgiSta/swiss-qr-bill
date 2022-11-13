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

package utils

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/ChrIgiSta/swiss-qr-bill/qr"
)

var QR_MATRIX = [10][10]int{
	{0, 9, 4, 6, 8, 2, 7, 1, 3, 5},
	{9, 4, 6, 8, 2, 7, 1, 3, 5, 0},
	{4, 6, 8, 2, 7, 1, 3, 5, 0, 9},
	{6, 8, 2, 7, 1, 3, 5, 0, 9, 4},
	{8, 2, 7, 1, 3, 5, 0, 9, 4, 6},
	{2, 7, 1, 3, 5, 0, 9, 4, 6, 8},
	{7, 1, 3, 5, 0, 9, 4, 6, 8, 2},
	{1, 3, 5, 0, 9, 4, 6, 8, 2, 7},
	{3, 5, 0, 9, 4, 6, 8, 2, 7, 1},
	{5, 0, 9, 4, 6, 8, 2, 7, 1, 3},
}

func ValidateReference(referenceType string, reference string) error {
	reference = strings.ReplaceAll(reference, " ", "")
	switch referenceType {

	case qr.REFERENCE_TYPE_NO_REF:
		if reference == "" {
			return nil
		}
		return errors.New("reference set for non ref type")

	case qr.REFERENCE_TYPE_QR:
		// 26 characters, 1 checksum mod 10 recursive
		if len(reference) != 27 {
			return errors.New("QR ref should be a 27 dig long string")
		}

		pShould, err := GetQrReferenceCheckNum(reference[:26])

		if err != nil {
			return err
		}

		pIs := string(reference[26])

		if pShould != pIs {
			errorMsg := fmt.Sprintf("checksum of qr ref isn't valid (is: %s, should: %s)", pIs, pShould)
			return errors.New(errorMsg)
		}

	case qr.REFERENCE_TYPE_CREDITOR:
		// ISO-11649, mod 97-10
		return errors.New("ref credo not implemented")

	default:
		return errors.New("unknown ref type")
	}

	return nil
}

func GetQrReferenceCheckNum(qrReference string) (string, error) {
	z := 0

	if len(qrReference) != 26 {
		fmt.Println(len(qrReference))
		return "", errors.New("QR ref should be a 26 dig long string (without checknum)")
	}
	for i := 0; i < 26; i++ {
		num, err := strconv.Atoi(string(qrReference[i]))
		if err != nil {
			return "", errors.New("QR ref contains characters.")
		}
		z = QR_MATRIX[z][num]
	}

	chkNum := math.Mod(float64(10-z), 10)

	return strconv.Itoa(int(chkNum)), nil
}
