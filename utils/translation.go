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

import "github.com/ChrIgiSta/swiss-qr-bill/specs"

func GetEnglishTranslationTable() *specs.TranslationTable {
	return &specs.TranslationTable{
		PaymentPart:       "Payment part",
		Account:           "Account / Payable to",
		Reference:         "Reference",
		AdditionalInfos:   "Additional information",
		FurtherInfos:      "Further information",
		Currency:          "Currency",
		Amount:            "Amount",
		Receipt:           "Receipt",
		AcceptancePoint:   "Acceptance point",
		SepBeforePay:      "Separate before paying in",
		PayableBy:         "Payable by",
		PayableByNameAddr: "Payable by (name/address)",
		InFavour:          "In favour of",
	}
}
