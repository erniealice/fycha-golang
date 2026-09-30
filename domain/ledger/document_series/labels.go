package document_series

import (
	"strings"

	enumspb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/enums"
	pb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/document_series"
)

// Labels holds every translatable string of the document series pages.
// Lyngua file: general/document_series.json, root key "document_series". JSON
// tags byte-match the lyngua key tree; every string field has a general-tier
// key (TestDefaultLabelsEqualGeneralLyngua) — the Go defaults are the English
// fallback.
type Labels struct {
	Page    PageLabels    `json:"page"`
	Tabs    TabLabels     `json:"tabs"`
	Columns ColumnLabels  `json:"columns"`
	Buttons ButtonLabels  `json:"buttons"`
	Form    FormLabels    `json:"form"`
	Enums   EnumLabels    `json:"enums"`
	Detail  DetailLabels  `json:"detail"`
	Confirm ConfirmLabels `json:"confirm"`
	Empty   EmptyLabels   `json:"empty"`
	Errors  ErrorLabels   `json:"errors"`
}

type PageLabels struct {
	Title        string `json:"title"`
	Subtitle     string `json:"subtitle"`
	TitleActive  string `json:"title_active"`
	TitleRetired string `json:"title_retired"`
}

type TabLabels struct {
	Active  string `json:"active"`
	Retired string `json:"retired"`
	Info    string `json:"info"`
}

type ColumnLabels struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	Issuer       string `json:"issuer"`
	DocumentKind string `json:"document_kind"`
	Prefix       string `json:"prefix"`
	NextNumber   string `json:"next_number"`
	FiscalReset  string `json:"fiscal_reset"`
	Status       string `json:"status"`
	Branch       string `json:"branch"`
	TaxID        string `json:"tax_id"`
}

type ButtonLabels struct {
	Add    string `json:"add"`
	Edit   string `json:"edit"`
	Retire string `json:"retire"`
}

type FormLabels struct {
	CodeLabel               string `json:"code_label"`
	CodePlaceholder         string `json:"code_placeholder"`
	NameLabel               string `json:"name_label"`
	NamePlaceholder         string `json:"name_placeholder"`
	IssuerNameLabel         string `json:"issuer_name_label"`
	IssuerTaxIDLabel        string `json:"issuer_tax_id_label"`
	DocumentKindLabel       string `json:"document_kind_label"`
	DocumentKindPlaceholder string `json:"document_kind_placeholder"`
	PrefixLabel             string `json:"prefix_label"`
	PrefixPlaceholder       string `json:"prefix_placeholder"`
	BranchCodeLabel         string `json:"branch_code_label"`
	FiscalResetLabel        string `json:"fiscal_reset_label"`
	NextNumberLabel         string `json:"next_number_label"`
	NumberPaddingLabel      string `json:"number_padding_label"`
}

type EnumLabels struct {
	StatusActive         string `json:"status_active"`
	StatusRetired        string `json:"status_retired"`
	FiscalResetNone      string `json:"fiscal_reset_none"`
	FiscalResetYearly    string `json:"fiscal_reset_yearly"`
	DocumentKindInvoice  string `json:"document_kind_invoice"`
	DocumentKindRecovery string `json:"document_kind_recovery_document"`
}

type DetailLabels struct {
	NextPreview  string `json:"next_preview"` // "Next document: {0}"
	LockedNotice string `json:"locked_notice"`
}

type ConfirmLabels struct {
	RetireTitle string `json:"retire_title"`
	RetireMsg   string `json:"retire_msg"`
}

type EmptyLabels struct {
	ActiveTitle   string `json:"active_title"`
	ActiveMessage string `json:"active_message"`
	RetiredTitle  string `json:"retired_title"`
}

// ErrorLabels: the first six fields are the document_series use-case error
// codes (`ErrorCode()`), keyed `document_series.errors.<code>` and shared with
// espyna's Translator; the rest are view-only messages.
type ErrorLabels struct {
	NotFound            string `json:"not_found"`
	Validation          string `json:"validation"`
	CodeTaken           string `json:"code_taken"`
	Retired             string `json:"retired"`
	TransactionRequired string `json:"transaction_required"`
	NumberingLocked     string `json:"numbering_locked"`

	FormInvalid string `json:"form_invalid"`
	Unavailable string `json:"unavailable"`
	Generic     string `json:"generic"`
}

// Format substitutes {0}, {1} placeholders (Lyngua positional convention).
func Format(tmpl string, args ...string) string {
	for i, a := range args {
		tmpl = strings.ReplaceAll(tmpl, "{"+string(rune('0'+i))+"}", a)
	}
	return tmpl
}

// DefaultLabels returns the English defaults (identical to lyngua
// general/document_series.json).
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			Title: "Document series", Subtitle: "Numbering for the documents you issue",
			TitleActive: "Active series", TitleRetired: "Retired series",
		},
		Tabs: TabLabels{Active: "Active", Retired: "Retired", Info: "Info"},
		Columns: ColumnLabels{
			Code: "Code", Name: "Name", Issuer: "Issuer", DocumentKind: "Document", Prefix: "Prefix",
			NextNumber: "Next number", FiscalReset: "Restarts", Status: "Status", Branch: "Branch", TaxID: "Tax ID",
		},
		Buttons: ButtonLabels{Add: "Add series", Edit: "Edit", Retire: "Retire"},
		Form: FormLabels{
			CodeLabel: "Code", CodePlaceholder: "e.g. SOA-2026", NameLabel: "Name", NamePlaceholder: "e.g. Statements 2026",
			IssuerNameLabel: "Issuer name", IssuerTaxIDLabel: "Issuer tax ID", DocumentKindLabel: "Document",
			DocumentKindPlaceholder: "Select a document", PrefixLabel: "Prefix", PrefixPlaceholder: "e.g. SOA-",
			BranchCodeLabel: "Branch code", FiscalResetLabel: "Restart numbering", NextNumberLabel: "Next number",
			NumberPaddingLabel: "Number length",
		},
		Enums: EnumLabels{
			StatusActive: "Active", StatusRetired: "Retired", FiscalResetNone: "Never", FiscalResetYearly: "Every year",
			DocumentKindInvoice: "Invoice", DocumentKindRecovery: "Recovery document",
		},
		Detail: DetailLabels{NextPreview: "Next document: {0}", LockedNotice: "Numbers already issued cannot be changed."},
		Confirm: ConfirmLabels{
			RetireTitle: "Retire this series?",
			RetireMsg:   "No new documents can be issued from it. Existing documents are unaffected.",
		},
		Empty: EmptyLabels{ActiveTitle: "No active series", ActiveMessage: "Add a series before issuing documents.", RetiredTitle: "Nothing retired"},
		Errors: ErrorLabels{
			NotFound: "Series not found.", Validation: "The document series is not valid.",
			CodeTaken: "That code is already used.", Retired: "This document series is retired.",
			TransactionRequired: "This action needs a database transaction, which is not available.",
			NumberingLocked:     "The prefix and number padding cannot change once numbers have been issued.",
			FormInvalid:         "Check the form and try again.",
			Unavailable:         "Document series is not available.",
			Generic:             "Something went wrong. Please try again.",
		},
	}
}

// ErrorMessage maps a use-case error code (ErrorCode()) to its Lyngua message.
func (l Labels) ErrorMessage(code string) string {
	switch code {
	case ErrKindValidation:
		return l.Errors.Validation
	case ErrKindNotFound:
		return l.Errors.NotFound
	case ErrKindCodeTaken:
		return l.Errors.CodeTaken
	case ErrKindRetired:
		return l.Errors.Retired
	case ErrKindTransactionRequired:
		return l.Errors.TransactionRequired
	case ErrKindNumberingLocked:
		return l.Errors.NumberingLocked
	default:
		return l.Errors.Generic
	}
}

// DocumentKindLabel returns the label of a series document kind.
func (l Labels) DocumentKindLabel(k enumspb.ChargeDocumentKind) string {
	switch k {
	case enumspb.ChargeDocumentKind_CHARGE_DOCUMENT_KIND_INVOICE:
		return l.Enums.DocumentKindInvoice
	case enumspb.ChargeDocumentKind_CHARGE_DOCUMENT_KIND_RECOVERY_DOCUMENT:
		return l.Enums.DocumentKindRecovery
	default:
		return ""
	}
}

// FiscalResetLabel returns the label of a fiscal-reset mode.
func (l Labels) FiscalResetLabel(r pb.DocumentSeriesFiscalReset) string {
	if r == pb.DocumentSeriesFiscalReset_DOCUMENT_SERIES_FISCAL_RESET_YEARLY {
		return l.Enums.FiscalResetYearly
	}
	return l.Enums.FiscalResetNone
}

// NextDocumentNumber previews the number the next issuance receives
// (prefix + zero-padded next number).
func NextDocumentNumber(s *pb.DocumentSeries) string {
	pad := int(s.GetNumberPadding())
	num := []byte{}
	n := s.GetNextNumber()
	for n > 0 {
		num = append([]byte{byte('0' + n%10)}, num...)
		n /= 10
	}
	for len(num) < pad {
		num = append([]byte{'0'}, num...)
	}
	if len(num) == 0 {
		num = []byte("0")
	}
	return s.GetPrefix() + string(num)
}
