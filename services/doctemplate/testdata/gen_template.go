//go:build ignore

// This program generates test DOCX templates for doctemplate tests.
// Run: go run testdata/gen_template.go
package main

import (
	"archive/zip"
	"fmt"
	"os"
)

const invoiceDocumentXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:wpc="http://schemas.microsoft.com/office/word/2010/wordprocessingCanvas" xmlns:mc="http://schemas.openxmlformats.org/markup-compatibility/2006" xmlns:o="urn:schemas-microsoft-com:office:office" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:m="http://schemas.openxmlformats.org/officeDocument/2006/math" xmlns:v="urn:schemas-microsoft-com:vml" xmlns:wp14="http://schemas.microsoft.com/office/word/2010/wordprocessingDrawing" xmlns:wp="http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing" xmlns:w10="urn:schemas-microsoft-com:office:word" xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main" xmlns:w14="http://schemas.microsoft.com/office/word/2010/wordml" xmlns:w15="http://schemas.microsoft.com/office/word/2012/wordml" xmlns:wpg="http://schemas.microsoft.com/office/word/2010/wordprocessingGroup" xmlns:wpi="http://schemas.microsoft.com/office/word/2010/wordprocessingInk" xmlns:wne="http://schemas.microsoft.com/office/word/2006/wordml" xmlns:wps="http://schemas.microsoft.com/office/word/2010/wordprocessingShape" mc:Ignorable="w14 w15 wpc wpg wpi wps"><w:body><w:p w:rsidR="00000001" w:rsidRDefault="00000001"><w:pPr><w:pStyle w:val="Title"/></w:pPr><w:r><w:t>INVOICE</w:t></w:r></w:p><w:p w:rsidR="00000001" w:rsidRDefault="00000001"><w:r><w:t>Bill To: </w:t></w:r><w:r><w:rPr><w:b/></w:rPr><w:t>{{</w:t></w:r><w:r><w:rPr><w:b/></w:rPr><w:t>client.name</w:t></w:r><w:r><w:rPr><w:b/></w:rPr><w:t>}}</w:t></w:r></w:p><w:p w:rsidR="00000001" w:rsidRDefault="00000001"><w:r><w:t>Address: {{</w:t></w:r><w:r><w:t>client.address</w:t></w:r><w:r><w:t>}}</w:t></w:r></w:p><w:p w:rsidR="00000001" w:rsidRDefault="00000001"><w:r><w:t>Date: {{date}}</w:t></w:r></w:p><w:p w:rsidR="00000001" w:rsidRDefault="00000001"/><w:tbl><w:tblPr><w:tblStyle w:val="TableGrid"/><w:tblW w:w="9000" w:type="dxa"/><w:tblLook w:val="04A0" w:firstRow="1" w:lastRow="0" w:firstColumn="1" w:lastColumn="0" w:noHBand="0" w:noVBand="1"/></w:tblPr><w:tblGrid><w:gridCol w:w="6000"/><w:gridCol w:w="3000"/></w:tblGrid><w:tr w:rsidR="00000001"><w:trPr><w:tblHeader/></w:trPr><w:tc><w:tcPr><w:tcW w:w="6000" w:type="dxa"/><w:shd w:val="clear" w:color="auto" w:fill="4472C4"/></w:tcPr><w:p><w:pPr><w:rPr><w:b/><w:color w:val="FFFFFF"/></w:rPr></w:pPr><w:r><w:rPr><w:b/><w:color w:val="FFFFFF"/></w:rPr><w:t>Description</w:t></w:r></w:p></w:tc><w:tc><w:tcPr><w:tcW w:w="3000" w:type="dxa"/><w:shd w:val="clear" w:color="auto" w:fill="4472C4"/></w:tcPr><w:p><w:pPr><w:jc w:val="right"/><w:rPr><w:b/><w:color w:val="FFFFFF"/></w:rPr></w:pPr><w:r><w:rPr><w:b/><w:color w:val="FFFFFF"/></w:rPr><w:t>Amount</w:t></w:r></w:p></w:tc></w:tr><w:tr w:rsidR="00000001"><w:tc><w:tcPr><w:tcW w:w="6000" w:type="dxa"/></w:tcPr><w:p><w:r><w:t>{{#items}}</w:t></w:r></w:p></w:tc><w:tc><w:tcPr><w:tcW w:w="3000" w:type="dxa"/></w:tcPr><w:p><w:r><w:t></w:t></w:r></w:p></w:tc></w:tr><w:tr w:rsidR="00000001"><w:tc><w:tcPr><w:tcW w:w="6000" w:type="dxa"/></w:tcPr><w:p><w:r><w:t>{{description}}</w:t></w:r></w:p></w:tc><w:tc><w:tcPr><w:tcW w:w="3000" w:type="dxa"/></w:tcPr><w:p><w:pPr><w:jc w:val="right"/></w:pPr><w:r><w:t>{{amount}}</w:t></w:r></w:p></w:tc></w:tr><w:tr w:rsidR="00000001"><w:tc><w:tcPr><w:tcW w:w="6000" w:type="dxa"/></w:tcPr><w:p><w:r><w:t>{{/items}}</w:t></w:r></w:p></w:tc><w:tc><w:tcPr><w:tcW w:w="3000" w:type="dxa"/></w:tcPr><w:p><w:r><w:t></w:t></w:r></w:p></w:tc></w:tr><w:tr w:rsidR="00000001"><w:tc><w:tcPr><w:tcW w:w="6000" w:type="dxa"/><w:shd w:val="clear" w:color="auto" w:fill="D9E2F3"/></w:tcPr><w:p><w:pPr><w:jc w:val="right"/><w:rPr><w:b/></w:rPr></w:pPr><w:r><w:rPr><w:b/></w:rPr><w:t>Total</w:t></w:r></w:p></w:tc><w:tc><w:tcPr><w:tcW w:w="3000" w:type="dxa"/><w:shd w:val="clear" w:color="auto" w:fill="D9E2F3"/></w:tcPr><w:p><w:pPr><w:jc w:val="right"/><w:rPr><w:b/></w:rPr></w:pPr><w:r><w:rPr><w:b/></w:rPr><w:t>{{total}}</w:t></w:r></w:p></w:tc></w:tr></w:tbl><w:p w:rsidR="00000001" w:rsidRDefault="00000001"/><w:p w:rsidR="00000001" w:rsidRDefault="00000001"><w:r><w:t>Notes: {{notes}}</w:t></w:r></w:p><w:p w:rsidR="00000001" w:rsidRDefault="00000001"><w:r><w:t>Thank you for your business!</w:t></w:r></w:p><w:sectPr w:rsidR="00000001"><w:pgSz w:w="12240" w:h="15840"/><w:pgMar w:top="1440" w:right="1440" w:bottom="1440" w:left="1440" w:header="720" w:footer="720" w:gutter="0"/><w:cols w:space="720"/><w:docGrid w:linePitch="360"/></w:sectPr></w:body></w:document>`

const contentTypesXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/><Default Extension="xml" ContentType="application/xml"/><Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/></Types>`

const relsXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/></Relationships>`

const documentRelsXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"></Relationships>`

func main() {
	if err := createDocx("testdata/invoice-template.docx"); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Created testdata/invoice-template.docx")
}

func createDocx(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := zip.NewWriter(f)

	files := map[string]string{
		"[Content_Types].xml":          contentTypesXML,
		"_rels/.rels":                  relsXML,
		"word/document.xml":            invoiceDocumentXML,
		"word/_rels/document.xml.rels": documentRelsXML,
	}

	for name, content := range files {
		writer, err := w.Create(name)
		if err != nil {
			return err
		}
		if _, err := writer.Write([]byte(content)); err != nil {
			return err
		}
	}

	return w.Close()
}
