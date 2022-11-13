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

package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/ChrIgiSta/swiss-qr-bill/sql"
)

const (
	TOKEN_KEY = "X-API-Key"
)

type BillInformation struct {
	IssuerId     int     `json:"issuer_id"`
	Name         string  `json:"name"`
	FirstName    string  `json:"firstname"`
	Street       string  `json:"street"`
	StreetNumber string  `json:"streetNumber"`
	Postal       string  `json:"postal"`
	City         string  `json:"city"`
	Amount       float32 `json:"amount"`
	Message      string  `json:"message"`
}

type Api struct {
	apiPath string
	port    int
	db      sql.Db
}

func NewApi(apiPath string, port int, db *sql.Db) *Api {
	return &Api{
		apiPath: apiPath,
		port:    port,
		db:      *db,
	}
}

func (api *Api) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	http.HandleFunc(api.apiPath+"/bill", api.GetBill)

	log.Println("start server on port", api.port)
	listen := fmt.Sprintf(":%d", api.port)
	err := http.ListenAndServe(listen, nil)
	if err != nil {
		log.Println("http server shutdown todue an error. ", err)
	}
}

func (api *Api) GetBill(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")

	token, err := api.db.GetApiToken()
	if err != nil {
		log.Println("unable to validate token", err.Error())
		http.Error(w, "unable to validate token", http.StatusInternalServerError)
		return
	}
	if auth != TOKEN_KEY+" "+token {
		log.Println("unauthorized client wanted to gen bill. ", r.Host)
		http.Error(w, "not authorized", http.StatusUnauthorized)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("error while reading body. ", err)
		http.Error(w, "cannot read body", http.StatusNoContent)
		return
	}
	defer r.Body.Close()

	billInfo := BillInformation{}
	err = json.Unmarshal(body, &billInfo)
	if err != nil {
		log.Println("error while encoding json", err)
		http.Error(w, "cannot unmarshal json", http.StatusNotAcceptable)
		return
	}

	log.Println("generate new bill")

	log.Fatal("not implemented: -> ToDo: get issuer, generate qr and send it back")
}
