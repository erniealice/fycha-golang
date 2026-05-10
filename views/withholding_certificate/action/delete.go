package action

import (
	"context"
	"log"

	fycha "github.com/erniealice/fycha-golang"
	withholdingcertificatepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/withholding_certificate"
	"github.com/erniealice/pyeza-golang/view"
)

// NewDeleteAction handles POST for deleting a WithholdingCertificate.
// The certificate ID is read from the "id" query parameter (table-actions.js convention).
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("withholding_certificate", "delete") {
			return fycha.HTMXError("You do not have permission to delete withholding certificates")
		}

		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return fycha.HTMXError("Certificate ID is required")
		}

		if deps.DeleteWithholdingCertificate == nil {
			return fycha.HTMXError("Delete use case not available")
		}
		if _, err := deps.DeleteWithholdingCertificate(ctx, &withholdingcertificatepb.DeleteWithholdingCertificateRequest{
			Data: &withholdingcertificatepb.WithholdingCertificate{Id: id},
		}); err != nil {
			log.Printf("DeleteWithholdingCertificate %s: %v", id, err)
			return fycha.HTMXError(err.Error())
		}

		return fycha.HTMXSuccess("withholding-certs-table")
	})
}
