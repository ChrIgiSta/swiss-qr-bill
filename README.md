# Swiss QR-Bill Generator
## Main Goal of the App
 - Provides an API to geneate Swiss QR Bills
 - Provides an Mail Client witch sends back an pdf with QR Bill attached
 - Simply generate QR bills directly (see cmd/example)

## how those it seens
<object data="https://github.com/ChrIgiSta/swiss-qr-bill/blob/main/out/example-bill-from-pdf.pdf" type="application/pdf" width="700px" height="700px">
    <embed src="https://github.com/ChrIgiSta/swiss-qr-bill/blob/main/out/example-bill-from-pdf.pdf">
        <p>Look up: <a href="https://github.com/ChrIgiSta/swiss-qr-bill/blob/main/out/example-bill-from-pdf.pdf">Download PDF</a>.</p>
    </embed>
</object>


This application provides a RESTfull API to generate your own Swiss QR-Bill for free.

So you can automate your billing system and append your QR-Bill to your Invoice to send it to our costumers.

## Contibution
 - are very welcome -> make a PR

## Specs
`https://www.paymentstandards.ch/de/shared/communication-grid.html?utm_campaign=vanity%20url&utm_medium=redirect&utm_source=www.paymentstandards.ch&utm_medium=redirect&utm_source=www.paymentstandards.ch/kommunikationsmatrix`

You can find the specification under the following link:
`https://www.paymentstandards.ch/dam/downloads/ig-qr-bill-de.pdf`
You can find the graphics under the following link:
- Swiss Cross (7mm): `https://www.paymentstandards.ch/de/shared/communication-grid/swiss-qr-code.html`
- Markers and scissors symbols: `https://www.paymentstandards.ch/de/shared/communication-grid/eckmarken.html`


ToDo: PDF Signieren
      Multilingual (DB)
      Referenz

Not Supported until Now:
      Unstructured Addresses
      SCOR Reference Validation