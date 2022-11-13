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

package sql

import (
	"database/sql"
	"strings"

	"github.com/ChrIgiSta/swiss-qr-bill/qr"
	"github.com/ChrIgiSta/swiss-qr-bill/specs"

	_ "github.com/go-sql-driver/mysql"
)

const (
	DRIVER = "mysql"
)

type Db struct {
	ConnectionString string
	dbCon            *sql.DB
}

func NewMariaDb(connectionString string) *Db {
	return &Db{
		ConnectionString: connectionString,
	}
}

func (db *Db) Connect() error {
	var err error

	db.dbCon, err = sql.Open(DRIVER, db.ConnectionString)
	return err
}

func (db *Db) RegisterVersion(version string) error {
	var err error

	_, err = db.dbCon.Exec("INSERT INTO versioning (version) VALUES (?)", version)
	return err
}

func (db *Db) CheckVersion() (string, error) {
	var (
		version string = ""
	)

	rows, err := db.dbCon.Query("SELECT version FROM version")
	if err != nil {
		return "", err
	}

	defer rows.Close()

	rows.Next()
	rows.Scan(&version)

	return version, err
}

func (db *Db) GetIssuer(id int) (string, specs.AccountDetails, error) {
	var (
		iban                string               = ""
		firstname, lastname string               = "", ""
		issuer              specs.AccountDetails = specs.AccountDetails{}
	)

	rows, err := db.dbCon.Query("SELECT iban, fistname, lastname, address1, "+
		"address2, zip, location, country FROM issuer WHERE id = ?", id)
	if err != nil {
		return "", issuer, err
	}
	defer rows.Close()

	rows.Next()
	err = rows.Scan(&iban, &firstname, &lastname, &issuer.Address1,
		&issuer.Address2, &issuer.Zip, &issuer.Location, &issuer.Country)
	issuer.AddressType = qr.ADDRESS_TYPE_STRUCTURED
	issuer.Name = lastname + " " + firstname

	return iban, issuer, err
}

func (db *Db) InsertIssuer(iban string, details specs.AccountDetails) (*specs.AccountDetails, int, error) {
	id := -1
	nameSplit := strings.Split(details.Name, " ")
	row, err := db.dbCon.Query("INSERT INTO issuer (iban, fistname, lastname, address1, "+
		"address2, zip, location, country) VALUES (?, ?, ?, ?, ?, ?, ?, ?) RETURNING id",
		iban, nameSplit[1], nameSplit[0], details.Address1, details.Address2,
		details.Zip, details.Location, details.Country)

	if err != nil {
		return nil, id, err
	}

	row.Next()
	row.Scan(&id)
	return &details, id, err
}

func (db *Db) InsertMailConfig(cnf *specs.MailConfig) error {
	_, err := db.dbCon.Exec("INSERT INTO mail "+
		"(token, enable, issuer_id, username, email, sender_name, password, smtp_secure, smtp_host, smtp_port, pop3_secure, pop3_host, pop3_port, use_whitelist) "+
		"VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)",
		cnf.Token, cnf.Enable, cnf.IssuerId, cnf.Username, cnf.Email, cnf.SenderName, cnf.Password, cnf.SmtpSecure, cnf.SmtpHost, cnf.SmtpPort, cnf.Pop3Secure,
		cnf.Pop3Host, cnf.Pop3Port, cnf.UseWhitelist)
	return err
}

func (db *Db) GetTranslationTable(languageCode string) (specs.TranslationTable, error) {
	tr := specs.TranslationTable{}

	row, err := db.dbCon.Query("SELECT paymentPart, account, reference, additionalInfos, "+
		"furtherInfos, currency, amount, receipt, acceptancePoint, sepBeforePay, payableBy, "+
		"payableByNameAddr, inFavour FROM translation WHERE languageCode = ?", languageCode)

	if err != nil {
		return tr, err
	}

	defer row.Close()
	row.Next()
	err = row.Scan(&tr.PaymentPart, &tr.Account, &tr.Reference, &tr.AdditionalInfos,
		&tr.FurtherInfos, &tr.Currency, &tr.Amount, &tr.Receipt, &tr.AcceptancePoint,
		&tr.SepBeforePay, &tr.PayableBy, &tr.PayableByNameAddr, &tr.InFavour)

	return tr, err
}

func (db *Db) GetApiToken() (string, error) {
	var token string = ""
	row := db.dbCon.QueryRow("SELECT token FROM api")
	err := row.Scan(&token)

	return token, err
}

func (db *Db) GetMailConfigurations() ([]*specs.MailConfig, error) {
	mCnfs := []*specs.MailConfig{}

	row, err := db.dbCon.Query("SELECT enable, issuer_id, username,	email, password " +
		",smtp_secure ,smtp_host ,smtp_port ,pop3_secure ,pop3_host ,pop3_port ,token ,use_whitelist " +
		"FROM mail WHERE enable = true")

	if err != nil {
		return nil, err
	}

	defer row.Close()
	for row.Next() {
		mCnf := specs.MailConfig{}
		err = row.Scan(&mCnf.Enable, &mCnf.IssuerId, &mCnf.Username, &mCnf.Email, &mCnf.Password,
			&mCnf.SmtpSecure, &mCnf.SmtpHost, &mCnf.SmtpPort, &mCnf.Pop3Secure, &mCnf.Pop3Host, &mCnf.Pop3Port,
			&mCnf.Token, &mCnf.UseWhitelist)
		if err != nil {
			return nil, err
		}
		mCnfs = append(mCnfs, &mCnf)
	}

	return mCnfs, err
}
