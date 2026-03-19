package export

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"crona/kernel/internal/runtime"
	sharedtypes "crona/shared/types"
)

type templateMeta struct {
	BaseHash   string  `json:"baseHash"`
	LastSynced *string `json:"lastSynced,omitempty"`
}

type exportConfig struct {
	ReportsDir string `json:"reportsDir,omitempty"`
}

type templateStatus struct {
	userPath        string
	bundledPath     string
	resettable      bool
	engine          string
	exists          bool
	customized      bool
	updateAvailable bool
	baseHash        string
	defaultHash     string
	activeSource    string
	lastSyncedAt    *string
}

type pdfRenderer struct {
	available bool
	name      string
	path      string
	engine    string
}

func EnsureAssets(paths runtime.Paths) (sharedtypes.ExportAssetStatus, error) {
	exportDir := filepath.Join(paths.UserAssetsDir, "export")
	if err := os.MkdirAll(exportDir, 0o700); err != nil {
		return sharedtypes.ExportAssetStatus{}, err
	}

	reportsDir, reportsCustomized, err := resolveReportsDir(paths)
	if err != nil {
		return sharedtypes.ExportAssetStatus{}, err
	}
	if err := os.MkdirAll(reportsDir, 0o700); err != nil {
		return sharedtypes.ExportAssetStatus{}, err
	}

	templateAssets := make([]sharedtypes.ExportTemplateAsset, 0, len(assetDescriptors()))
	var dailyMarkdownStatus templateStatus
	var dailyPDFStatus templateStatus
	var dailyDocsPath string
	for _, descriptor := range assetDescriptors() {
		status, err := ensureTemplateAsset(paths, descriptor)
		if err != nil {
			return sharedtypes.ExportAssetStatus{}, err
		}
		templateAssets = append(templateAssets, sharedtypes.ExportTemplateAsset{
			ReportKind:      descriptor.reportKind,
			AssetKind:       descriptor.assetKind,
			Label:           descriptor.label,
			Name:            descriptor.name,
			Engine:          descriptor.engine,
			UserPath:        status.userPath,
			BundledPath:     status.bundledPath,
			Resettable:      descriptor.resettable,
			Exists:          status.exists,
			Customized:      status.customized,
			UpdateAvailable: status.updateAvailable,
			BaseHash:        status.baseHash,
			DefaultHash:     status.defaultHash,
			ActiveSource:    status.activeSource,
			LastSyncedAt:    status.lastSyncedAt,
		})
		if descriptor.reportKind == sharedtypes.ExportReportKindDaily && descriptor.assetKind == sharedtypes.ExportAssetKindTemplateMarkdown {
			dailyMarkdownStatus = status
		}
		if descriptor.reportKind == sharedtypes.ExportReportKindDaily && descriptor.assetKind == sharedtypes.ExportAssetKindTemplatePDF {
			dailyPDFStatus = status
		}
		if descriptor.reportKind == sharedtypes.ExportReportKindDaily && descriptor.assetKind == sharedtypes.ExportAssetKindVariableDocs {
			dailyDocsPath = status.userPath
		}
	}

	renderer := detectPDFRenderer()

	return sharedtypes.ExportAssetStatus{
		TemplatePath:           dailyMarkdownStatus.userPath,
		TemplateDocsPath:       dailyDocsPath,
		BundledTemplatePath:    dailyMarkdownStatus.bundledPath,
		PDFTemplatePath:        dailyPDFStatus.userPath,
		PDFBundledTemplatePath: dailyPDFStatus.bundledPath,
		ReportsDir:             reportsDir,
		DefaultReportsDir:      defaultReportsDir(paths),
		ReportsDirCustomized:   reportsCustomized,
		UserTemplateExists:     dailyMarkdownStatus.exists,
		UserTemplateCustomized: dailyMarkdownStatus.customized,
		DefaultUpdateAvailable: dailyMarkdownStatus.updateAvailable,
		PDFUserTemplateExists:  dailyPDFStatus.exists,
		PDFTemplateCustomized:  dailyPDFStatus.customized,
		PDFUpdateAvailable:     dailyPDFStatus.updateAvailable,
		TemplateBaseHash:       dailyMarkdownStatus.baseHash,
		CurrentDefaultHash:     dailyMarkdownStatus.defaultHash,
		PDFTemplateBaseHash:    dailyPDFStatus.baseHash,
		PDFCurrentDefaultHash:  dailyPDFStatus.defaultHash,
		TemplateName:           "daily/report.hbs",
		TemplateEngine:         dailyMarkdownStatus.engine,
		ActiveTemplateSource:   dailyMarkdownStatus.activeSource,
		PDFTemplateName:        "daily/report.pdf.hbs",
		PDFTemplateEngine:      dailyPDFStatus.engine,
		PDFTemplateSource:      dailyPDFStatus.activeSource,
		PDFRendererAvailable:   renderer.available,
		PDFRendererName:        renderer.name,
		PDFRendererPath:        renderer.path,
		LastSyncedAt:           dailyMarkdownStatus.lastSyncedAt,
		PDFLastSyncedAt:        dailyPDFStatus.lastSyncedAt,
		TemplateAssets:         templateAssets,
	}, nil
}

func ResetTemplate(paths runtime.Paths, reportKind sharedtypes.ExportReportKind, assetKind sharedtypes.ExportAssetKind) (sharedtypes.ExportAssetStatus, error) {
	descriptor, ok := findAssetDescriptor(reportKind, assetKind)
	if !ok {
		return sharedtypes.ExportAssetStatus{}, errors.New("export asset not found")
	}
	if !descriptor.resettable {
		return sharedtypes.ExportAssetStatus{}, errors.New("export asset is not resettable")
	}
	body, _, err := defaultAssetSource(paths, descriptor)
	if err != nil {
		return sharedtypes.ExportAssetStatus{}, err
	}
	now := time.Now().UTC().Format(time.RFC3339)
	if err := os.MkdirAll(filepath.Dir(userTemplatePath(paths, descriptor)), 0o700); err != nil {
		return sharedtypes.ExportAssetStatus{}, err
	}
	if err := os.WriteFile(userTemplatePath(paths, descriptor), body, runtime.FilePerm()); err != nil {
		return sharedtypes.ExportAssetStatus{}, err
	}
	if err := writeTemplateMeta(templateMetaPath(paths, descriptor), templateMeta{
		BaseHash:   hashBytes(body),
		LastSynced: &now,
	}); err != nil {
		return sharedtypes.ExportAssetStatus{}, err
	}
	return EnsureAssets(paths)
}

func LoadActiveTemplate(paths runtime.Paths, format sharedtypes.ExportFormat) ([]byte, sharedtypes.ExportAssetStatus, error) {
	return LoadActiveReportTemplate(paths, sharedtypes.ExportReportKindDaily, format)
}

func LoadActiveReportTemplate(paths runtime.Paths, reportKind sharedtypes.ExportReportKind, format sharedtypes.ExportFormat) ([]byte, sharedtypes.ExportAssetStatus, error) {
	status, err := EnsureAssets(paths)
	if err != nil {
		return nil, sharedtypes.ExportAssetStatus{}, err
	}
	body, err := os.ReadFile(activeTemplatePathForStatus(status, reportKind, normalizeExportFormat(format)))
	if err != nil {
		return nil, sharedtypes.ExportAssetStatus{}, err
	}
	return body, status, nil
}

func WriteDailyReport(paths runtime.Paths, date string, format sharedtypes.ExportFormat, body []byte) (string, error) {
	return WriteReport(paths, reportWriteSpec{
		Kind:     sharedtypes.ExportReportKindDaily,
		Label:    "Daily Report",
		Date:     date,
		Format:   format,
		BaseName: "daily-" + date,
	}, body)
}

func WriteReport(paths runtime.Paths, spec reportWriteSpec, body []byte) (string, error) {
	reportsDir, _, err := resolveReportsDir(paths)
	if err != nil {
		return "", err
	}
	target := filepath.Join(reportsDir, spec.BaseName+reportExt(spec.Format))
	if err := os.MkdirAll(filepath.Dir(target), 0o700); err != nil {
		return "", err
	}
	if err := os.WriteFile(target, body, runtime.FilePerm()); err != nil {
		return "", err
	}
	if err := writeReportMetadata(metadataPathForReport(target), reportFileMetadata{
		Kind:       spec.Kind,
		Label:      spec.Label,
		ScopeLabel: spec.ScopeLabel,
		Date:       spec.Date,
		StartDate:  spec.StartDate,
		EndDate:    spec.EndDate,
		Format:     spec.Format,
	}); err != nil {
		return "", err
	}
	return target, nil
}

func SetReportsDir(paths runtime.Paths, raw string) (sharedtypes.ExportAssetStatus, error) {
	config, _ := readExportConfig(exportConfigPath(paths))
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		config.ReportsDir = ""
	} else {
		resolved, err := normalizeReportsDir(paths, trimmed)
		if err != nil {
			return sharedtypes.ExportAssetStatus{}, err
		}
		config.ReportsDir = resolved
	}
	if err := writeExportConfig(exportConfigPath(paths), config); err != nil {
		return sharedtypes.ExportAssetStatus{}, err
	}
	return EnsureAssets(paths)
}

func ListReports(paths runtime.Paths) ([]sharedtypes.ExportReportFile, error) {
	reportsDir, _, err := resolveReportsDir(paths)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(reportsDir, 0o700); err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(reportsDir)
	if err != nil {
		return nil, err
	}
	out := make([]sharedtypes.ExportReportFile, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(strings.ToLower(name), ".meta.json") {
			continue
		}
		ext := strings.ToLower(filepath.Ext(name))
		if ext != ".md" && ext != ".pdf" && ext != ".csv" {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		path := filepath.Join(reportsDir, name)
		metadata, _ := readReportMetadata(metadataPathForReport(path))
		date := strings.TrimSuffix(name, ext)
		kind := sharedtypes.ExportReportKindDaily
		scopeLabel := ""
		startDate := ""
		endDate := ""
		dateLabel := date
		format := exportFormatFromExt(ext)
		if metadata != nil {
			kind = metadata.Kind
			scopeLabel = metadata.ScopeLabel
			if metadata.Date != "" {
				date = metadata.Date
			}
			startDate = metadata.StartDate
			endDate = metadata.EndDate
			dateLabel = reportDisplayDateLabel(date, startDate, endDate)
			format = string(metadata.Format)
		}
		out = append(out, sharedtypes.ExportReportFile{
			Name:       name,
			Path:       path,
			Kind:       kind,
			ScopeLabel: scopeLabel,
			Date:       date,
			StartDate:  startDate,
			EndDate:    endDate,
			DateLabel:  dateLabel,
			Format:     format,
			SizeBytes:  info.Size(),
			ModifiedAt: info.ModTime().UTC().Format(time.RFC3339),
		})
	}
	sort.Slice(out, func(i, j int) bool {
		left := reportDisplayDateLabel(out[i].Date, out[i].StartDate, out[i].EndDate)
		right := reportDisplayDateLabel(out[j].Date, out[j].StartDate, out[j].EndDate)
		if left == right {
			if out[i].Format == out[j].Format {
				return out[i].Name > out[j].Name
			}
			return out[i].Format < out[j].Format
		}
		return left > right
	})
	return out, nil
}

func DeleteReport(paths runtime.Paths, rawPath string) error {
	reportPath := strings.TrimSpace(rawPath)
	if reportPath == "" {
		return errors.New("report path is required")
	}
	reportsDir, _, err := resolveReportsDir(paths)
	if err != nil {
		return err
	}
	cleanReportsDir := filepath.Clean(reportsDir)
	if !filepath.IsAbs(reportPath) {
		reportPath = filepath.Join(cleanReportsDir, reportPath)
	}
	reportPath = filepath.Clean(reportPath)
	rel, err := filepath.Rel(cleanReportsDir, reportPath)
	if err != nil {
		return err
	}
	if rel == "." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) || rel == ".." {
		return errors.New("report path is outside reports directory")
	}
	info, err := os.Stat(reportPath)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return errors.New("report path must be a file")
	}
	if err := os.Remove(reportPath); err != nil {
		return err
	}
	metaPath := metadataPathForReport(reportPath)
	if err := os.Remove(metaPath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

func RenderPDF(paths runtime.Paths, date string, markdown string) (string, string, error) {
	return RenderPDFReport(paths, reportWriteSpec{
		Kind:     sharedtypes.ExportReportKindDaily,
		Label:    "Daily Report",
		Date:     date,
		Format:   sharedtypes.ExportFormatPDF,
		BaseName: "daily-" + date,
	}, markdown)
}

func RenderPDFReport(paths runtime.Paths, spec reportWriteSpec, markdown string) (string, string, error) {
	renderer := detectPDFRenderer()
	if !renderer.available {
		return "", "", errors.New("no supported PDF renderer found; install pandoc with a PDF engine such as tectonic, weasyprint, wkhtmltopdf, xelatex, or pdflatex")
	}
	tmpDir, err := os.MkdirTemp("", "crona-export-pdf-*")
	if err != nil {
		return "", "", err
	}
	defer os.RemoveAll(tmpDir)

	baseName := spec.BaseName
	if strings.TrimSpace(baseName) == "" {
		baseName = "report"
	}
	inputPath := filepath.Join(tmpDir, baseName+".md")
	if err := os.WriteFile(inputPath, []byte(markdown), 0o600); err != nil {
		return "", "", err
	}
	outputPath := filepath.Join(tmpDir, baseName+".pdf")

	cmd := exec.Command(renderer.path, inputPath, "-o", outputPath, "--pdf-engine="+renderer.engine)
	if output, err := cmd.CombinedOutput(); err != nil {
		msg := strings.TrimSpace(string(output))
		if msg == "" {
			msg = err.Error()
		}
		return "", renderer.name, errors.New("pdf export failed: " + msg)
	}

	body, err := os.ReadFile(outputPath)
	if err != nil {
		return "", renderer.name, err
	}
	finalPath, err := WriteReport(paths, reportWriteSpec{
		Kind:       spec.Kind,
		Label:      spec.Label,
		ScopeLabel: spec.ScopeLabel,
		Date:       spec.Date,
		StartDate:  spec.StartDate,
		EndDate:    spec.EndDate,
		Format:     sharedtypes.ExportFormatPDF,
		BaseName:   spec.BaseName,
	}, body)
	if err != nil {
		return "", renderer.name, err
	}
	return finalPath, renderer.name, nil
}

func exportConfigPath(paths runtime.Paths) string {
	return filepath.Join(paths.UserAssetsDir, "export", "config.json")
}

func activeTemplateSource(customized bool) string {
	if customized {
		return "user"
	}
	return "default"
}

func ensureTemplateAsset(paths runtime.Paths, descriptor assetDescriptor) (templateStatus, error) {
	body, bundledPath, err := defaultAssetSource(paths, descriptor)
	if err != nil {
		return templateStatus{}, err
	}
	defaultHash := hashBytes(body)
	meta, _ := readTemplateMeta(templateMetaPath(paths, descriptor))
	userPath := userTemplatePath(paths, descriptor)
	userContent, err := os.ReadFile(userPath)
	switch {
	case errors.Is(err, os.ErrNotExist):
		legacyContent, legacyMeta, foundLegacy, legacyErr := readLegacyUserAsset(paths, descriptor)
		if legacyErr != nil {
			return templateStatus{}, legacyErr
		}
		if foundLegacy {
			if err := os.MkdirAll(filepath.Dir(userPath), 0o700); err != nil {
				return templateStatus{}, err
			}
			if err := os.WriteFile(userPath, legacyContent, runtime.FilePerm()); err != nil {
				return templateStatus{}, err
			}
			if descriptor.resettable && legacyMeta.BaseHash != "" {
				if err := writeTemplateMeta(templateMetaPath(paths, descriptor), legacyMeta); err != nil {
					return templateStatus{}, err
				}
				meta = legacyMeta
			}
			userContent = legacyContent
			break
		}
		if err := os.MkdirAll(filepath.Dir(userPath), 0o700); err != nil {
			return templateStatus{}, err
		}
		if err := os.WriteFile(userPath, body, runtime.FilePerm()); err != nil {
			return templateStatus{}, err
		}
		if descriptor.resettable {
			now := time.Now().UTC().Format(time.RFC3339)
			meta = templateMeta{BaseHash: defaultHash, LastSynced: &now}
			if err := writeTemplateMeta(templateMetaPath(paths, descriptor), meta); err != nil {
				return templateStatus{}, err
			}
		}
		userContent = body
	case err != nil:
		return templateStatus{}, err
	}

	userHash := hashBytes(userContent)
	customized := userHash != defaultHash
	if descriptor.resettable && meta.BaseHash != "" {
		customized = userHash != meta.BaseHash
	}
	updateAvailable := descriptor.resettable && meta.BaseHash != "" && meta.BaseHash != defaultHash
	if !customized && updateAvailable {
		now := time.Now().UTC().Format(time.RFC3339)
		if err := os.WriteFile(userPath, body, runtime.FilePerm()); err != nil {
			return templateStatus{}, err
		}
		meta = templateMeta{BaseHash: defaultHash, LastSynced: &now}
		if err := writeTemplateMeta(templateMetaPath(paths, descriptor), meta); err != nil {
			return templateStatus{}, err
		}
		customized = false
		updateAvailable = false
	}
	if descriptor.resettable && meta.BaseHash == "" {
		now := time.Now().UTC().Format(time.RFC3339)
		meta = templateMeta{BaseHash: defaultHash, LastSynced: &now}
		_ = writeTemplateMeta(templateMetaPath(paths, descriptor), meta)
	}

	return templateStatus{
		userPath:        userPath,
		bundledPath:     bundledPath,
		resettable:      descriptor.resettable,
		engine:          descriptor.engine,
		exists:          true,
		customized:      customized,
		updateAvailable: updateAvailable,
		baseHash:        meta.BaseHash,
		defaultHash:     defaultHash,
		activeSource:    activeTemplateSource(customized),
		lastSyncedAt:    meta.LastSynced,
	}, nil
}

func defaultReportsDir(paths runtime.Paths) string {
	return paths.ReportsDir
}

func legacyDailyReportsDir(paths runtime.Paths) string {
	return filepath.Join(paths.ReportsDir, "daily")
}

func resolveReportsDir(paths runtime.Paths) (string, bool, error) {
	config, err := readExportConfig(exportConfigPath(paths))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return "", false, err
	}
	if strings.TrimSpace(config.ReportsDir) == "" {
		return defaultReportsDir(paths), false, nil
	}
	resolved, err := normalizeReportsDir(paths, config.ReportsDir)
	if err != nil {
		return "", false, err
	}
	if resolved == legacyDailyReportsDir(paths) {
		return defaultReportsDir(paths), false, nil
	}
	return resolved, true, nil
}

func normalizeReportsDir(paths runtime.Paths, raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return defaultReportsDir(paths), nil
	}
	if trimmed == "~" || strings.HasPrefix(trimmed, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		if trimmed == "~" {
			trimmed = home
		} else {
			trimmed = filepath.Join(home, strings.TrimPrefix(trimmed, "~/"))
		}
	}
	if !filepath.IsAbs(trimmed) {
		trimmed = filepath.Join(paths.BaseDir, trimmed)
	}
	cleaned := filepath.Clean(trimmed)
	if cleaned == legacyDailyReportsDir(paths) {
		return defaultReportsDir(paths), nil
	}
	return cleaned, nil
}

func readExportConfig(path string) (exportConfig, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return exportConfig{}, err
	}
	var config exportConfig
	if err := json.Unmarshal(body, &config); err != nil {
		return exportConfig{}, err
	}
	return config, nil
}

func writeExportConfig(path string, config exportConfig) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	body, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, body, runtime.FilePerm())
}

func readTemplateMeta(path string) (templateMeta, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return templateMeta{}, err
	}
	var meta templateMeta
	if err := json.Unmarshal(body, &meta); err != nil {
		return templateMeta{}, err
	}
	return meta, nil
}

func writeTemplateMeta(path string, meta templateMeta) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	body, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, body, runtime.FilePerm())
}

func detectPDFRenderer() pdfRenderer {
	pandocPath, err := exec.LookPath("pandoc")
	if err != nil {
		return pdfRenderer{}
	}
	engineCandidates := []string{"tectonic", "weasyprint", "wkhtmltopdf", "xelatex", "pdflatex"}
	for _, engine := range engineCandidates {
		if _, err := exec.LookPath(engine); err == nil {
			return pdfRenderer{
				available: true,
				name:      "pandoc (" + engine + ")",
				path:      pandocPath,
				engine:    engine,
			}
		}
	}
	return pdfRenderer{}
}

func normalizeExportFormat(format sharedtypes.ExportFormat) sharedtypes.ExportFormat {
	if format == sharedtypes.ExportFormatPDF || format == sharedtypes.ExportFormatCSV {
		return format
	}
	return sharedtypes.ExportFormatMarkdown
}

func activeTemplatePathForStatus(status sharedtypes.ExportAssetStatus, reportKind sharedtypes.ExportReportKind, format sharedtypes.ExportFormat) string {
	assetKind := sharedtypes.ExportAssetKindTemplateMarkdown
	if normalizeExportFormat(format) == sharedtypes.ExportFormatPDF {
		assetKind = sharedtypes.ExportAssetKindTemplatePDF
	}
	for _, asset := range status.TemplateAssets {
		if asset.ReportKind == reportKind && asset.AssetKind == assetKind {
			return asset.UserPath
		}
	}
	if reportKind == sharedtypes.ExportReportKindDaily && assetKind == sharedtypes.ExportAssetKindTemplatePDF {
		return status.PDFTemplatePath
	}
	return status.TemplatePath
}

func reportExt(format sharedtypes.ExportFormat) string {
	switch normalizeExportFormat(format) {
	case sharedtypes.ExportFormatPDF:
		return ".pdf"
	case sharedtypes.ExportFormatCSV:
		return ".csv"
	default:
		return ".md"
	}
}

func exportFormatFromExt(ext string) string {
	switch strings.ToLower(ext) {
	case ".pdf":
		return string(sharedtypes.ExportFormatPDF)
	case ".csv":
		return string(sharedtypes.ExportFormatCSV)
	default:
		return string(sharedtypes.ExportFormatMarkdown)
	}
}

func hashBytes(body []byte) string {
	sum := sha256.Sum256(body)
	return hex.EncodeToString(sum[:])
}
