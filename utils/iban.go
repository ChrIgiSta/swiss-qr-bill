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
	"strconv"
	"strings"
)

const (
	IBAN_LEN       = 21
	COUNTRY_ID_LEN = 2
)

func ValidateIban(iban string) error {
	iban = strings.ReplaceAll(iban, " ", "")

	if len(iban) != IBAN_LEN {
		return errors.New("iban len mismatch")
	}

	for i := 0; i < COUNTRY_ID_LEN; i++ {
		_, err := strconv.Atoi(string(iban[i]))
		if err == nil {
			return errors.New("country id is a number")
		}
	}
	for i := COUNTRY_ID_LEN; i < IBAN_LEN; i++ {
		_, err := strconv.Atoi(string(iban[i]))
		if err != nil {
			return errors.New("iban contains characters after country id")
		}
	}

	return nil
}
