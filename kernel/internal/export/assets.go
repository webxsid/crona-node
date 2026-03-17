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

	markdownStatus, err := ensureTemplateAsset(paths, sharedtypes.ExportFormatMarkdown)
	if err != nil {
		return sharedtypes.ExportAssetStatus{}, err
	}
	pdfStatus, err := ensureTemplateAsset(paths, sharedtypes.ExportFormatPDF)
	if err != nil {
		return sharedtypes.ExportAssetStatus{}, err
	}

	defaultDocs, _, err := defaultDocsSource(paths)
	if err != nil {
		return sharedtypes.ExportAssetStatus{}, err
	}
	if err := ensureDocsFile(docsTemplatePath(paths), defaultDocs); err != nil {
		return sharedtypes.ExportAssetStatus{}, err
	}

	renderer := detectPDFRenderer()

	return sharedtypes.ExportAssetStatus{
		TemplatePath:           markdownStatus.userPath,
		TemplateDocsPath:       docsTemplatePath(paths),
		BundledTemplatePath:    markdownStatus.bundledPath,
		PDFTemplatePath:        pdfStatus.userPath,
		PDFBundledTemplatePath: pdfStatus.bundledPath,
		ReportsDir:             reportsDir,
		DefaultReportsDir:      defaultReportsDir(paths),
		ReportsDirCustomized:   reportsCustomized,
		UserTemplateExists:     markdownStatus.exists,
		UserTemplateCustomized: markdownStatus.customized,
		DefaultUpdateAvailable: markdownStatus.updateAvailable,
		PDFUserTemplateExists:  pdfStatus.exists,
		PDFTemplateCustomized:  pdfStatus.customized,
		PDFUpdateAvailable:     pdfStatus.updateAvailable,
		TemplateBaseHash:       markdownStatus.baseHash,
		CurrentDefaultHash:     markdownStatus.defaultHash,
		PDFTemplateBaseHash:    pdfStatus.baseHash,
		PDFCurrentDefaultHash:  pdfStatus.defaultHash,
		TemplateName:           "daily-report.hbs",
		TemplateEngine:         "hbs",
		ActiveTemplateSource:   markdownStatus.activeSource,
		PDFTemplateName:        "daily-report.pdf.hbs",
		PDFTemplateEngine:      "hbs",
		PDFTemplateSource:      pdfStatus.activeSource,
		PDFRendererAvailable:   renderer.available,
		PDFRendererName:        renderer.name,
		PDFRendererPath:        renderer.path,
		LastSyncedAt:           markdownStatus.lastSyncedAt,
		PDFLastSyncedAt:        pdfStatus.lastSyncedAt,
	}, nil
}

func ResetTemplate(paths runtime.Paths, format sharedtypes.ExportFormat) (sharedtypes.ExportAssetStatus, error) {
	body, _, err := defaultTemplateSource(paths, normalizeExportFormat(format))
	if err != nil {
		return sharedtypes.ExportAssetStatus{}, err
	}
	now := time.Now().UTC().Format(time.RFC3339)
	if err := os.WriteFile(userTemplatePath(paths, format), body, runtime.FilePerm()); err != nil {
		return sharedtypes.ExportAssetStatus{}, err
	}
	if err := writeTemplateMeta(templateMetaPath(paths, format), templateMeta{
		BaseHash:   hashBytes(body),
		LastSynced: &now,
	}); err != nil {
		return sharedtypes.ExportAssetStatus{}, err
	}
	return EnsureAssets(paths)
}

func LoadActiveTemplate(paths runtime.Paths, format sharedtypes.ExportFormat) ([]byte, sharedtypes.ExportAssetStatus, error) {
	status, err := EnsureAssets(paths)
	if err != nil {
		return nil, sharedtypes.ExportAssetStatus{}, err
	}
	body, err := os.ReadFile(activeTemplatePathForStatus(status, normalizeExportFormat(format)))
	if err != nil {
		return nil, sharedtypes.ExportAssetStatus{}, err
	}
	return body, status, nil
}

func WriteDailyReport(paths runtime.Paths, date string, format sharedtypes.ExportFormat, body []byte) (string, error) {
	reportsDir, _, err := resolveReportsDir(paths)
	if err != nil {
		return "", err
	}
	target := filepath.Join(reportsDir, date+reportExt(format))
	if err := os.MkdirAll(filepath.Dir(target), 0o700); err != nil {
		return "", err
	}
	if err := os.WriteFile(target, body, runtime.FilePerm()); err != nil {
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
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext != ".md" && ext != ".pdf" {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		name := entry.Name()
		out = append(out, sharedtypes.ExportReportFile{
			Name:       name,
			Path:       filepath.Join(reportsDir, name),
			Date:       strings.TrimSuffix(name, ext),
			Format:     exportFormatFromExt(ext),
			SizeBytes:  info.Size(),
			ModifiedAt: info.ModTime().UTC().Format(time.RFC3339),
		})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Date == out[j].Date {
			if out[i].Format == out[j].Format {
				return out[i].Name > out[j].Name
			}
			return out[i].Format < out[j].Format
		}
		return out[i].Date > out[j].Date
	})
	return out, nil
}

func RenderPDF(paths runtime.Paths, date string, markdown string) (string, string, error) {
	renderer := detectPDFRenderer()
	if !renderer.available {
		return "", "", errors.New("no supported PDF renderer found; install pandoc with a PDF engine such as tectonic, weasyprint, wkhtmltopdf, xelatex, or pdflatex")
	}
	tmpDir, err := os.MkdirTemp("", "crona-export-pdf-*")
	if err != nil {
		return "", "", err
	}
	defer os.RemoveAll(tmpDir)

	inputPath := filepath.Join(tmpDir, date+".md")
	if err := os.WriteFile(inputPath, []byte(markdown), 0o600); err != nil {
		return "", "", err
	}
	outputPath := filepath.Join(tmpDir, date+".pdf")

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
	finalPath, err := WriteDailyReport(paths, date, sharedtypes.ExportFormatPDF, body)
	if err != nil {
		return "", renderer.name, err
	}
	return finalPath, renderer.name, nil
}

func userTemplatePath(paths runtime.Paths, format sharedtypes.ExportFormat) string {
	switch normalizeExportFormat(format) {
	case sharedtypes.ExportFormatPDF:
		return filepath.Join(paths.UserAssetsDir, "export", "daily-report.pdf.user.hbs")
	default:
		return filepath.Join(paths.UserAssetsDir, "export", "daily-report.user.hbs")
	}
}

func templateMetaPath(paths runtime.Paths, format sharedtypes.ExportFormat) string {
	switch normalizeExportFormat(format) {
	case sharedtypes.ExportFormatPDF:
		return filepath.Join(paths.UserAssetsDir, "export", "daily-report.pdf.user.meta.json")
	default:
		return filepath.Join(paths.UserAssetsDir, "export", "daily-report.user.meta.json")
	}
}

func docsTemplatePath(paths runtime.Paths) string {
	return filepath.Join(paths.UserAssetsDir, "export", "daily-report.variables.md")
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

func defaultTemplateSource(paths runtime.Paths, format sharedtypes.ExportFormat) ([]byte, string, error) {
	switch normalizeExportFormat(format) {
	case sharedtypes.ExportFormatPDF:
		return resolveAssetSource(paths, "daily-report.pdf.default.hbs", fallbackDailyReportPDFTemplate)
	default:
		return resolveAssetSource(paths, "daily-report.default.hbs", fallbackDailyReportTemplate)
	}
}

func defaultDocsSource(paths runtime.Paths) ([]byte, string, error) {
	return resolveAssetSource(paths, "daily-report.variables.md", fallbackDailyReportVariables)
}

func resolveAssetSource(paths runtime.Paths, name string, fallback string) ([]byte, string, error) {
	candidates := []string{
		filepath.Join(paths.BundledAssetsDir, "export", name),
		filepath.Join("assets", "export", name),
		filepath.Join("..", "assets", "export", name),
		filepath.Join("..", "..", "assets", "export", name),
	}
	for _, candidate := range candidates {
		body, err := os.ReadFile(candidate)
		if err == nil {
			return body, candidate, nil
		}
	}
	return []byte(fallback), filepath.Join(paths.BundledAssetsDir, "export", name), nil
}

func ensureDocsFile(path string, content []byte) error {
	existing, err := os.ReadFile(path)
	if err == nil && strings.TrimSpace(string(existing)) == strings.TrimSpace(string(content)) {
		return nil
	}
	if errors.Is(err, os.ErrNotExist) || err == nil {
		return os.WriteFile(path, content, runtime.FilePerm())
	}
	return err
}

func ensureTemplateAsset(paths runtime.Paths, format sharedtypes.ExportFormat) (templateStatus, error) {
	body, bundledPath, err := defaultTemplateSource(paths, format)
	if err != nil {
		return templateStatus{}, err
	}
	defaultHash := hashBytes(body)
	meta, _ := readTemplateMeta(templateMetaPath(paths, format))
	userPath := userTemplatePath(paths, format)
	userContent, err := os.ReadFile(userPath)
	switch {
	case errors.Is(err, os.ErrNotExist):
		now := time.Now().UTC().Format(time.RFC3339)
		if err := os.WriteFile(userPath, body, runtime.FilePerm()); err != nil {
			return templateStatus{}, err
		}
		meta = templateMeta{BaseHash: defaultHash, LastSynced: &now}
		if err := writeTemplateMeta(templateMetaPath(paths, format), meta); err != nil {
			return templateStatus{}, err
		}
		userContent = body
	case err != nil:
		return templateStatus{}, err
	}

	userHash := hashBytes(userContent)
	customized := userHash != defaultHash
	if meta.BaseHash != "" {
		customized = userHash != meta.BaseHash
	}
	updateAvailable := meta.BaseHash != "" && meta.BaseHash != defaultHash
	if !customized && updateAvailable {
		now := time.Now().UTC().Format(time.RFC3339)
		if err := os.WriteFile(userPath, body, runtime.FilePerm()); err != nil {
			return templateStatus{}, err
		}
		meta = templateMeta{BaseHash: defaultHash, LastSynced: &now}
		if err := writeTemplateMeta(templateMetaPath(paths, format), meta); err != nil {
			return templateStatus{}, err
		}
		customized = false
		updateAvailable = false
	}
	if meta.BaseHash == "" {
		now := time.Now().UTC().Format(time.RFC3339)
		meta = templateMeta{BaseHash: defaultHash, LastSynced: &now}
		_ = writeTemplateMeta(templateMetaPath(paths, format), meta)
	}

	return templateStatus{
		userPath:        userPath,
		bundledPath:     bundledPath,
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
	return filepath.Clean(trimmed), nil
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
	if format == sharedtypes.ExportFormatPDF {
		return format
	}
	return sharedtypes.ExportFormatMarkdown
}

func activeTemplatePathForStatus(status sharedtypes.ExportAssetStatus, format sharedtypes.ExportFormat) string {
	if normalizeExportFormat(format) == sharedtypes.ExportFormatPDF {
		return status.PDFTemplatePath
	}
	return status.TemplatePath
}

func reportExt(format sharedtypes.ExportFormat) string {
	if normalizeExportFormat(format) == sharedtypes.ExportFormatPDF {
		return ".pdf"
	}
	return ".md"
}

func exportFormatFromExt(ext string) string {
	switch strings.ToLower(ext) {
	case ".pdf":
		return string(sharedtypes.ExportFormatPDF)
	default:
		return string(sharedtypes.ExportFormatMarkdown)
	}
}

func hashBytes(body []byte) string {
	sum := sha256.Sum256(body)
	return hex.EncodeToString(sum[:])
}
