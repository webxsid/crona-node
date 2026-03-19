package app

import (
	sharedtypes "crona/shared/types"
	"crona/tui/internal/api"
)

type configItem struct {
	label       string
	value       string
	path        string
	detailTitle string
	detailMeta  string
	detailBody  string
	editable    bool
	mutable     bool
	actionHint  string
	reportKind  sharedtypes.ExportReportKind
	assetKind   sharedtypes.ExportAssetKind
	resettable  bool
}

func (m Model) configItems() []configItem {
	if m.exportAssets == nil {
		return nil
	}
	items := make([]configItem, 0, len(m.exportAssets.TemplateAssets)+2)
	for _, asset := range m.exportAssets.TemplateAssets {
		state := exportAssetStateLabel(asset)
		detailBody := "Path\n" + asset.UserPath + "\n\nBundled\n" + asset.BundledPath
		detailBody += "\n\nPress e to open in $EDITOR."
		if asset.Resettable {
			detailBody += "\nPress r to replace it with the bundled default."
		}
		items = append(items, configItem{
			label:       asset.Label,
			value:       state,
			path:        asset.UserPath,
			detailTitle: asset.Label,
			detailMeta:  "Engine " + asset.Engine + "   Source " + asset.ActiveSource + "   State " + state,
			detailBody:  detailBody,
			editable:    true,
			reportKind:  asset.ReportKind,
			assetKind:   asset.AssetKind,
			resettable:  asset.Resettable && (asset.Customized || asset.UpdateAvailable),
		})
	}
	items = append(items, configItem{
		label:       "Reports directory",
		value:       m.exportAssets.ReportsDir,
		detailTitle: "Report Output Directory",
		detailMeta:  reportsDirMeta(m.exportAssets),
		detailBody:  "Generated reports are written under\n" + m.exportAssets.ReportsDir + "\n\nDefault\n" + m.exportAssets.DefaultReportsDir + "\n\nPress c to change the directory.\nPress r to restore the default directory.",
		mutable:     true,
		actionHint:  "change dir",
	})
	items = append(items, configItem{
		label:       "PDF renderer",
		value:       pdfRendererStateLabel(m.exportAssets),
		detailTitle: "PDF Renderer",
		detailMeta:  "External renderer discovery",
		detailBody:  pdfRendererDetailBody(m.exportAssets),
	})
	return items
}

func reportsDirMeta(status *api.ExportAssetStatus) string {
	if status == nil {
		return ""
	}
	if status.ReportsDirCustomized {
		return "Mode file export   Source custom"
	}
	return "Mode file export   Source default"
}

func exportAssetStateLabel(asset sharedtypes.ExportTemplateAsset) string {
	if asset.Resettable {
		switch {
		case asset.Customized:
			return "customized"
		case asset.UpdateAvailable:
			return "new default available"
		default:
			return "default"
		}
	}
	return truncate(asset.UserPath, 28)
}

func pdfRendererStateLabel(status *api.ExportAssetStatus) string {
	if status == nil {
		return ""
	}
	if status.PDFRendererAvailable {
		return status.PDFRendererName
	}
	return "unavailable"
}

func pdfRendererDetailBody(status *api.ExportAssetStatus) string {
	if status == nil {
		return ""
	}
	if !status.PDFRendererAvailable {
		return "No supported PDF renderer detected.\n\nInstall pandoc with a supported PDF engine and press R in Config to rescan."
	}
	return "Renderer\n" + status.PDFRendererName + "\n\nPath\n" + status.PDFRendererPath + "\n\nPress R in Config to rescan available PDF tools."
}
