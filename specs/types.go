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

package specs

type AccountDetails struct {
	AddressType string `json:"address_type"`
	Name        string `json:"name"` // lastname + fistname
	Address1    string `json:"address1"`
	Address2    string `json:"address2"`
	Zip         string `json:"zip"`
	Location    string `json:"location"`
	Country     string `json:"country"`
}

type BillingDetails struct {
	// billing details
	IBAN           string  `json:"iban"`
	RefenreceType  string  `json:"reference_type"`
	Referece       string  `json:"reference"`
	AdditionalInfo string  `json:"additional_info"`
	Currency       string  `json:"currency"`
	Amount         float64 `json:"amount"`
}

type TranslationTable struct {
	PaymentPart       string `json:"payment_part"`
	Account           string `json:"account"`
	Reference         string `json:"reference"`
	AdditionalInfos   string `json:"additional_infos"`
	FurtherInfos      string `json:"further_infos"`
	Currency          string `json:"currency"`
	Amount            string `json:"amount"`
	Receipt           string `json:"receipt"`
	AcceptancePoint   string `json:"acceptance_point"`
	SepBeforePay      string `json:"separate_before_pay"`
	PayableBy         string `json:"payable_by"`
	PayableByNameAddr string `json:"payable_by_with_name_address"`
	InFavour          string `json:"in_favour"`
}

type MailConfig struct {
	Enable       bool   `json:"enable"`
	IssuerId     int    `json:"issuer_id"`
	Username     string `json:"username"`
	SenderName   string `json:"sender_name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	SmtpSecure   bool   `json:"smtp_secure"`
	SmtpHost     string `json:"smtp_host"`
	SmtpPort     int    `json:"smtp_port"`
	Pop3Secure   bool   `json:"pop3_secure"`
	Pop3Host     string `json:"pop3_host"`
	Pop3Port     int    `json:"pop3_port"`
	Token        string `json:"token"`
	UseWhitelist bool   `json:"use_whitelist"`
}
