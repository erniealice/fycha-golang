package action

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"context"

	withholdingcertificatepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/withholding_certificate"
	"github.com/erniealice/fycha-golang/domain/treasury/withholding_certificate/form"
	"github.com/erniealice/pyeza-golang/view"
)

// NewEditAction handles GET (form pre-populated from existing record) and
// POST (submit) for editing a WithholdingCertificate.
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("withholding_certificate", "update") {
			return view.HTMXError("You do not have permission to edit withholding certificates")
		}

		id := viewCtx.Request.PathValue("id")
		labels := form.BuildLabels(deps.Labels)

		if viewCtx.Request.Method == http.MethodGet {
			if deps.ReadWithholdingCertificate == nil {
				return view.HTMXError("Read use case not available")
			}
			resp, err := deps.ReadWithholdingCertificate(ctx, &withholdingcertificatepb.ReadWithholdingCertificateRequest{
				Data: &withholdingcertificatepb.WithholdingCertificate{Id: id},
			})
			if err != nil {
				log.Printf("ReadWithholdingCertificate %s: %v", id, err)
				return view.HTMXError("Certificate not found")
			}
			data := resp.GetData()
			if len(data) == 0 {
				return view.HTMXError("Certificate not found")
			}
			wc := data[0]
			certifiedDisplay := fmt.Sprintf("%.2f", float64(wc.GetActualAmount())/100.0)
			return view.OK("withholding-certificate-drawer-form", &form.Data{
				FormAction:         deps.Routes.EditFor(id),
				IsEdit:             true,
				ID:                 id,
				RevenueID:          wc.GetRevenueId(),
				CertificateNumber:  wc.GetCertificateNumber(),
				PayorTin:           wc.GetBuyerTinSnapshot(),
				DateIssued:         wc.GetIssuedDate(),
				WhtAmountCertified: certifiedDisplay,
				StatusOptions:      form.BuildStatusOptions(wc.GetStatus().String()),
				Labels:             labels,
				CommonLabels:       deps.CommonLabels,
			})
		}

		// POST — update
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError("Invalid form data")
		}
		r := viewCtx.Request

		certifiedStr := r.FormValue("wht_amount_certified")
		var certifiedCents int64
		if certifiedStr != "" {
			f, err := strconv.ParseFloat(certifiedStr, 64)
			if err == nil {
				certifiedCents = int64(f * 100)
			}
		}

		statusVal := r.FormValue("status")
		var status withholdingcertificatepb.WithholdingCertificateStatus
		if sv, ok := withholdingcertificatepb.WithholdingCertificateStatus_value[statusVal]; ok {
			status = withholdingcertificatepb.WithholdingCertificateStatus(sv)
		}

		updated := &withholdingcertificatepb.WithholdingCertificate{
			Id:                id,
			CertificateNumber: r.FormValue("certificate_number"),
			BuyerTinSnapshot:  optionalString(r.FormValue("payor_tin")),
			IssuedDate:        optionalString(r.FormValue("date_issued")),
			ActualAmount:      certifiedCents,
			Status:            status,
		}

		if deps.UpdateWithholdingCertificate == nil {
			return view.HTMXError("Update use case not available")
		}
		if _, err := deps.UpdateWithholdingCertificate(ctx, &withholdingcertificatepb.UpdateWithholdingCertificateRequest{
			Data: updated,
		}); err != nil {
			log.Printf("UpdateWithholdingCertificate %s: %v", id, err)
			return view.HTMXError(err.Error())
		}

		return view.HTMXSuccess("withholding-certs-table")
	})
}
