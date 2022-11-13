-- Copyright © 2022, Staufi Tech - Switzerland
-- All rights reserved.
--  THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
--  AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
--  IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
--  ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
--  LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
--  CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
--  SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
--  INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
--  CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
--  ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
--  POSSIBILITY OF SUCH DAMAGE.
 

CREATE TABLE IF NOT EXISTS issuer 
(
    id       BIGINT PRIMARY KEY NOT NULL UNIQUE AUTO_INCREMENT,
    fistname TEXT NOT NULL, 
    lastname TEXT, 
    address1 TEXT NOT NULL,
    address2 TEXT,
    zip      TEXT NOT NULL, 
    location TEXT NOT NULL, 
    country  ENUM ('CH', 'FL'),
    iban      VARCHAR(26) NOT NULL -- (21 + 5) with spaces
);

CREATE TABLE IF NOT EXISTS version
(
    version TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS customer
(
    cust_id  BIGINT PRIMARY KEY NOT NULL UNIQUE AUTO_INCREMENT,
    fistname TEXT NOT NULL, 
    lastname TEXT, 
    address1 TEXT NOT NULL,
    address2 TEXT,
    zip      TEXT NOT NULL, 
    location TEXT NOT NULL, 
    country  ENUM ('CH', 'FL')
);

CREATE TABLE IF NOT EXISTS bill
(
    id        BIGINT PRIMARY KEY NOT NULL UNIQUE AUTO_INCREMENT,
    cust_id   BIGINT REFERENCES customer(cust_id),
    ts        TIMESTAMP NOT NULL,
    ref_type  ENUM ('NON','QRR','SCOR'),
    reference TEXT,
    add_msg   TEXT,
    currency  ENUM ('CHF', 'EUR'),
    amount    DOUBLE
);

CREATE TABLE IF NOT EXISTS api 
(
    token  TEXT, -- if NULL, token is diabled
    enable BOOLEAN
);

CREATE TABLE IF NOT EXISTS mail 
(
    token         TEXT,
    enable        BOOLEAN NOT NULL DEFAULT false,
    issuer_id     BIGINT  REFERENCES issuer(id),
    username      TEXT NOT NULL,
    sender_name   TEXT,
    email         TEXT NOT NULL CHECK ( email regexp '^([a-zA-Z0-9_.-])+@(([a-zA-Z0-9-])+.)+([a-zA-Z0-9]{2,4})+\$' ),
    password      TEXT,
    smtp_secure   BOOLEAN NOT NULL DEFAULT true,
    smtp_host     TEXT,
    smtp_port     INT NOT NULL DEFAULT 587,
    pop3_secure   BOOLEAN NOT NULL DEFAULT true,
    pop3_host     TEXT,
    pop3_port     INT NOT NULL DEFAULT 993,
    use_whitelist BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS mail_whitelist
(
    email TEXT NOT NULL CHECK ( email regexp '^([a-zA-Z0-9_.-])+@(([a-zA-Z0-9-])+.)+([a-zA-Z0-9]{2,4})+\$' )
);

CREATE TABLE IF NOT EXISTS translation
(
    language_code        ENUM ('DE', 'FR', 'IT', 'EN') NOT NULL UNIQUE,
    payment_part         TEXT NOT NULL,
    account              TEXT NOT NULL,
    reference            TEXT NOT NULL,
    additional_infos     TEXT NOT NULL,
    further_infos        TEXT NOT NULL,
    currency             TEXT NOT NULL,
    amount               TEXT NOT NULL,
    receipt              TEXT NOT NULL,
    acceptance_point     TEXT NOT NULL,
    sep_before_pay       TEXT NOT NULL,
    payable_by           TEXT NOT NULL,
    payable_by_name_addr TEXT NOT NULL,
    in_favour             TEXT NOT NULL
);

--
-- data
--
INSERT INTO translation (language_code, payment_part, account, reference, additional_infos, further_infos, currency, amount, receipt, acceptance_point, sep_before_pay, payable_by, payable_by_name_addr, in_favour)
VALUES 
    ('EN', 'Payment part',      'Account / Payable to', 'Reference',   'Additional information',      'Further information',          'Currency', 'Amount',  'Receipt',        'Acceptance point',      'Separate before paying in',        'Payable by',    'Payable by (name/address)',    'In favour of'),
    ('DE', 'Zahlteil',          'Konto / Zahlbar an',   'Referenz',    'Zusätzliche Informationen',   'Weitere Informationen',        'Währung',  'Betrag',  'Empfangsschein', 'Annahmestelle',         'Vor der Einzahlung abzutrennen',   'Zahlbar durch', 'Zahlbar durch (Name/Adresse)', 'Zugunsten'),
    ('FR', 'Section paiement',  'Compte / Payable à',   'Référence',   'Informations additionnelles', 'Informations supplémentaires', 'Monnaie',  'Montant', 'Récépissé',      'Point de dépôt',        'A détacher avant le versement',    'Payable par',   'Payable par (nom/adresse)',    'En faveur de'),
    ('IT', 'Sezione pagamento', 'Conto / Pagabile a',   'Riferimento', 'Informazioni aggiuntive',     'Informazioni supplementari',   'Valuta',   'Importo', 'Ricevuta',       'Punto di accettazione', 'Da staccare prima del versamento', 'Pagabile da',   'Pagabile da (nome/indirizzo)', 'A favore di');
 
INSERT INTO version (version) VALUES ('v1.0.0');
 