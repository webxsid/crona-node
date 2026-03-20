package app

import (
	sharedtypes "crona/shared/types"
	helperpkg "crona/tui/internal/tui/app/helpers"
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
	dialogKind  string
}

func (m Model) configItems() []configItem {
	if m.exportAssets == nil {
		return nil
	}
	items := make([]configItem, 0, len(m.exportAssets.TemplateAssets)+3)
	for _, asset := range m.exportAssets.TemplateAssets {
		state := helperpkg.ExportAssetStateLabel(asset)
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
		detailMeta:  helperpkg.ReportsDirMeta(m.exportAssets),
		detailBody:  "Generated reports are written under\n" + m.exportAssets.ReportsDir + "\n\nDefault\n" + m.exportAssets.DefaultReportsDir + "\n\nPress c to change the directory.\nPress r to restore the default directory.",
		mutable:     true,
		actionHint:  "change dir",
		dialogKind:  "edit_export_reports_dir",
	})
	items = append(items, configItem{
		label:       "ICS export directory",
		value:       m.exportAssets.ICSDir,
		detailTitle: "ICS Export Directory",
		detailMeta:  helperpkg.ICSDirMeta(m.exportAssets),
		detailBody:  "Calendar exports are written under\n" + m.exportAssets.ICSDir + "\n\nDefault\n" + m.exportAssets.DefaultICSDir + "\n\nUse this directory for Shortcuts, Folder Actions, or other local automations.\nPress c to change the directory.\nPress r to restore the default directory.",
		mutable:     true,
		actionHint:  "change dir",
		dialogKind:  "edit_export_ics_dir",
	})
	items = append(items, configItem{
		label:       "PDF renderer",
		value:       helperpkg.PDFRendererStateLabel(m.exportAssets),
		detailTitle: "PDF Renderer",
		detailMeta:  "External renderer discovery",
		detailBody:  helperpkg.PDFRendererDetailBody(m.exportAssets),
	})
	return items
}
