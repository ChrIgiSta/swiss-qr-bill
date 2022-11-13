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
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"time"
)

const CharSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
	"abcdefghijklmnopqrstuvwxyz" +
	"1234567890" +
	`+-_"*%&<>()=?!:;@#/\|{}[].,`

type Token struct {
	Plain  string
	Sha256 string
}

func CreateNewToken(length int) Token {
	rand.Seed(time.Now().UnixNano())

	t := make([]byte, length)

	for i := range t {
		t[i] = CharSet[rand.Intn(len(CharSet))]
	}

	token := Token{
		Plain:  string(t),
		Sha256: GetSha256([]byte(t)),
	}

	return token
}

func GetSha256(plain []byte) string {
	hasher := sha256.New()
	hasher.Write(plain)
	return hex.EncodeToString(hasher.Sum(nil))
}
