package status

import (
	"bytes"
	"context"

	"github.com/jung-kurt/gofpdf"
	"github.com/yushafro/03.12.2025/pkg/logger"
)

const (
	pageMarginTop   = 40
	pageMarginLeft  = 10
	fontSizeTitle   = 16
	fontSizeLink    = 12
	linkMarginLeft  = 80
	linkMarginRight = 30
	lineMargin      = 8
)

func generatePDFReport(ctx context.Context, links Links) ([]byte, error) {
	log := logger.FromContext(ctx)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", fontSizeTitle)
	pdf.Cell(pageMarginLeft, pageMarginTop, "Status report")
	pdf.Ln(lineMargin)
	pdf.SetFont("Arial", "", fontSizeLink)

	if len(links) == 0 {
		pdf.Cell(linkMarginLeft, pageMarginTop+lineMargin, "No links found")
		pdf.Ln(lineMargin)
	} else {
		for link, status := range links {
			statusText := status

			pdf.Cell(linkMarginLeft, pageMarginTop+lineMargin, link)
			pdf.Cell(linkMarginRight, pageMarginTop+lineMargin, statusText)
			pdf.Ln(lineMargin)
		}
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	log.Info(ctx, "PDF report generated")

	return buf.Bytes(), nil
}
