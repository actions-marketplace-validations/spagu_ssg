// Package generator - additional coverage tests
package generator

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spagu/ssg/internal/mddb"
	"github.com/spagu/ssg/internal/models"
)

func TestCopyColocatedAssetsNonexistentDir(t *testing.T) {
	tmpDir := t.TempDir()
	gen := &Generator{config: Config{Quiet: true}}
	err := gen.copyColocatedAssets(filepath.Join(tmpDir, "nonexistent"), filepath.Join(tmpDir, "dst"))
	if err != nil {
		t.Errorf("Expected nil error for nonexistent source dir, got: %v", err)
	}
}

func TestCopyColocatedAssetsNoAssets(t *testing.T) {
	tmpDir := t.TempDir()
	srcDir := filepath.Join(tmpDir, "src")

	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatalf("Failed to create src dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "readme.md"), []byte("# Hi"), 0644); err != nil {
		t.Fatalf("Failed to create md: %v", err)
	}

	dstDir := filepath.Join(tmpDir, "dst")
	gen := &Generator{config: Config{Quiet: false}}
	if err := gen.copyColocatedAssets(srcDir, dstDir); err != nil {
		t.Fatalf("copyColocatedAssets failed: %v", err)
	}
}

func TestExtractCategoryFromDoc(t *testing.T) {
	tests := []struct {
		name     string
		doc      mddb.Document
		wantID   int
		wantName string
		wantSlug string
	}{
		{
			name: "full category",
			doc: mddb.Document{
				Key: "tech",
				Metadata: map[string]any{
					"id":          float64(5),
					"name":        "Technology",
					"description": "Tech articles",
					"link":        "https://example.com/category/tech/",
					"count":       float64(10),
					"parent":      float64(0),
				},
			},
			wantID:   5,
			wantName: "Technology",
			wantSlug: "tech",
		},
		{
			name: "empty metadata",
			doc: mddb.Document{
				Key:      "empty",
				Metadata: map[string]any{},
			},
			wantID:   0,
			wantName: "",
			wantSlug: "empty",
		},
		{
			name: "wrong types in metadata",
			doc: mddb.Document{
				Key: "wrong",
				Metadata: map[string]any{
					"id":   "not-a-number",
					"name": 12345,
				},
			},
			wantID:   0,
			wantName: "",
			wantSlug: "wrong",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cat := extractCategoryFromDoc(tt.doc)
			if cat.ID != tt.wantID {
				t.Errorf("ID = %d, want %d", cat.ID, tt.wantID)
			}
			if cat.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", cat.Name, tt.wantName)
			}
			if cat.Slug != tt.wantSlug {
				t.Errorf("Slug = %q, want %q", cat.Slug, tt.wantSlug)
			}
		})
	}
}

func TestExtractMediaFromDoc(t *testing.T) {
	tests := []struct {
		name       string
		doc        mddb.Document
		wantID     int
		wantType   string
		wantMime   string
		wantSrcURL string
		wantTitle  string
	}{
		{
			name: "full media",
			doc: mddb.Document{
				Key: "image1",
				Metadata: map[string]any{
					"id":         float64(100),
					"media_type": "image",
					"mime_type":  "image/jpeg",
					"source_url": "https://example.com/image.jpg",
					"title": map[string]interface{}{
						"rendered": "My Image",
					},
				},
			},
			wantID:     100,
			wantType:   "image",
			wantMime:   "image/jpeg",
			wantSrcURL: "https://example.com/image.jpg",
			wantTitle:  "My Image",
		},
		{
			name: "empty metadata",
			doc: mddb.Document{
				Key:      "empty",
				Metadata: map[string]any{},
			},
			wantID:     0,
			wantType:   "",
			wantMime:   "",
			wantSrcURL: "",
			wantTitle:  "",
		},
		{
			name: "title without rendered",
			doc: mddb.Document{
				Key: "partial",
				Metadata: map[string]any{
					"title": map[string]interface{}{
						"raw": "Raw Title",
					},
				},
			},
			wantTitle: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			media := extractMediaFromDoc(tt.doc)
			if media.ID != tt.wantID {
				t.Errorf("ID = %d, want %d", media.ID, tt.wantID)
			}
			if media.MediaType != tt.wantType {
				t.Errorf("MediaType = %q, want %q", media.MediaType, tt.wantType)
			}
			if media.MimeType != tt.wantMime {
				t.Errorf("MimeType = %q, want %q", media.MimeType, tt.wantMime)
			}
			if media.SourceURL != tt.wantSrcURL {
				t.Errorf("SourceURL = %q, want %q", media.SourceURL, tt.wantSrcURL)
			}
			if media.Title.Rendered != tt.wantTitle {
				t.Errorf("Title.Rendered = %q, want %q", media.Title.Rendered, tt.wantTitle)
			}
		})
	}
}

func TestExtractAuthorFromDoc(t *testing.T) {
	tests := []struct {
		name     string
		doc      mddb.Document
		wantID   int
		wantName string
		wantSlug string
	}{
		{
			name: "full author",
			doc: mddb.Document{
				Key: "admin",
				Metadata: map[string]any{
					"id":   float64(1),
					"name": "Admin User",
				},
			},
			wantID:   1,
			wantName: "Admin User",
			wantSlug: "admin",
		},
		{
			name: "empty metadata",
			doc: mddb.Document{
				Key:      "unknown",
				Metadata: map[string]any{},
			},
			wantID:   0,
			wantName: "",
			wantSlug: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			author := extractAuthorFromDoc(tt.doc)
			if author.ID != tt.wantID {
				t.Errorf("ID = %d, want %d", author.ID, tt.wantID)
			}
			if author.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", author.Name, tt.wantName)
			}
			if author.Slug != tt.wantSlug {
				t.Errorf("Slug = %q, want %q", author.Slug, tt.wantSlug)
			}
		})
	}
}

func TestConvertRelativeLinksIfRequested(t *testing.T) {
	tests := []struct {
		name          string
		relativeLinks bool
		domain        string
		expectSkip    bool
	}{
		{"disabled", false, "example.com", true},
		{"no domain", true, "", true},
		{"enabled with domain", true, "example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			if !tt.expectSkip {
				if err := os.WriteFile(filepath.Join(tmpDir, "test.html"), []byte(`<a href="https://example.com/page">Link</a>`), 0644); err != nil {
					t.Fatalf("Failed to create html: %v", err)
				}
			}

			gen := &Generator{
				config: Config{
					RelativeLinks: tt.relativeLinks,
					Domain:        tt.domain,
					OutputDir:     tmpDir,
					Quiet:         true,
				},
			}

			err := gen.convertRelativeLinksIfRequested()
			if err != nil {
				t.Fatalf("convertRelativeLinksIfRequested failed: %v", err)
			}

			if !tt.expectSkip {
				content, _ := os.ReadFile(filepath.Join(tmpDir, "test.html"))
				if strings.Contains(string(content), "https://example.com") {
					t.Error("Links should have been converted to relative")
				}
			}
		})
	}
}

func TestConvertToRelativeLinksError(t *testing.T) {
	gen := &Generator{
		config: Config{
			OutputDir: "/nonexistent/path",
			Domain:    "example.com",
		},
	}

	err := gen.convertToRelativeLinks()
	if err == nil {
		t.Error("Expected error for nonexistent output dir")
	}
}

func TestLinkifyListItem(t *testing.T) {
	pageLinks := map[string]string{
		"About Us":   "/about/",
		"Contact":    "/contact/",
		"Caf\u00e9": "/cafe/",
	}

	tests := []struct {
		name     string
		line     string
		content  string
		expected string
	}{
		{
			name:     "matching page link",
			line:     "- About Us",
			content:  "About Us",
			expected: "- [About Us](/about/)",
		},
		{
			name:     "no match",
			line:     "- Unknown Page",
			content:  "Unknown Page",
			expected: "- Unknown Page",
		},
		{
			name:     "html entity match",
			line:     "- Caf\u00e9",
			content:  "Caf\u00e9",
			expected: "- [Caf\u00e9](/cafe/)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := linkifyListItem(tt.line, tt.content, pageLinks)
			if result != tt.expected {
				t.Errorf("linkifyListItem(%q, %q) = %q, want %q", tt.line, tt.content, result, tt.expected)
			}
		})
	}
}

func TestExtractListItemContent(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected string
	}{
		{"dash list item", "- Hello", "Hello"},
		{"asterisk list item", "* World", "World"},
		{"indented dash", "  - Indented", "Indented"},
		{"not a list item", "Regular text", ""},
		{"empty line", "", ""},
		{"dash only", "-", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractListItemContent(tt.line)
			if result != tt.expected {
				t.Errorf("extractListItemContent(%q) = %q, want %q", tt.line, result, tt.expected)
			}
		})
	}
}

func TestTmplFormatDateNonString(t *testing.T) {
	result := tmplFormatDate(42)
	if result != "42" {
		t.Errorf("tmplFormatDate(42) = %q, want %q", result, "42")
	}

	result = tmplFormatDate("2024-01-15")
	if result != "2024-01-15" {
		t.Errorf("tmplFormatDate(string) = %q, want %q", result, "2024-01-15")
	}
}

func TestTmplGetCategoryNameMissing(t *testing.T) {
	gen := &Generator{
		siteData: &models.SiteData{
			Categories: map[int]models.Category{
				1: {ID: 1, Name: "News"},
			},
		},
	}

	if name := gen.tmplGetCategoryName(999); name != "" {
		t.Errorf("Expected empty string for missing category, got %q", name)
	}
	if name := gen.tmplGetCategoryName(1); name != "News" {
		t.Errorf("Expected 'News', got %q", name)
	}
}

func TestTmplGetCategorySlugMissing(t *testing.T) {
	gen := &Generator{
		siteData: &models.SiteData{
			Categories: map[int]models.Category{
				1: {ID: 1, Slug: "news"},
			},
		},
	}

	if slug := gen.tmplGetCategorySlug(999); slug != "" {
		t.Errorf("Expected empty string for missing category, got %q", slug)
	}
	if slug := gen.tmplGetCategorySlug(1); slug != "news" {
		t.Errorf("Expected 'news', got %q", slug)
	}
}

func TestTmplGetAuthorNameMissing(t *testing.T) {
	gen := &Generator{
		siteData: &models.SiteData{
			Authors: map[int]models.Author{
				1: {ID: 1, Name: "Admin"},
			},
		},
	}

	if name := gen.tmplGetAuthorName(999); name != "" {
		t.Errorf("Expected empty string for missing author, got %q", name)
	}
	if name := gen.tmplGetAuthorName(1); name != "Admin" {
		t.Errorf("Expected 'Admin', got %q", name)
	}
}

func TestTmplIsValidCategory(t *testing.T) {
	if tmplIsValidCategory(1) {
		t.Error("Category 1 should not be valid")
	}
	if !tmplIsValidCategory(2) {
		t.Error("Category 2 should be valid")
	}
}

func TestTmplThumbnailFromYoutubeNoMatch(t *testing.T) {
	result := tmplThumbnailFromYoutube("no youtube shortcode here")
	if result != "" {
		t.Errorf("Expected empty string for no match, got %q", result)
	}
}

func TestTmplThumbnailFromYoutubeMatch(t *testing.T) {
	result := tmplThumbnailFromYoutube("[youtube]https://www.youtube.com/watch?v=abc123[/youtube]")
	expected := "https://img.youtube.com/vi/abc123/hqdefault.jpg"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestTmplStripShortcodes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"strip youtube", "text [youtube]http://yt.com[/youtube] more", "text  more"},
		{"strip embed", "text [embed]http://yt.com[/embed] more", "text  more"},
		{"no shortcodes", "plain text", "plain text"},
		{"both shortcodes", "[youtube]url[/youtube] and [embed]url2[/embed]", "and"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tmplStripShortcodes(tt.input)
			if result != tt.expected {
				t.Errorf("tmplStripShortcodes(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTmplStripHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple tags", "<p>Hello</p>", "Hello"},
		{"nested tags", "<div><p>World</p></div>", "World"},
		{"no tags", "plain", "plain"},
		{"self-closing", "<br/><img src='x'>text", "text"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tmplStripHTML(tt.input)
			if result != tt.expected {
				t.Errorf("tmplStripHTML(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTmplRecentPosts(t *testing.T) {
	gen := &Generator{
		siteData: &models.SiteData{
			Posts: []models.Page{
				{Title: "Post 1"},
				{Title: "Post 2"},
				{Title: "Post 3"},
			},
		},
	}

	result := gen.tmplRecentPosts(2)
	if len(result) != 2 {
		t.Errorf("Expected 2 posts, got %d", len(result))
	}

	result = gen.tmplRecentPosts(10)
	if len(result) != 3 {
		t.Errorf("Expected 3 posts (capped), got %d", len(result))
	}
}

func TestCleanMarkdownArtifacts(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"bold text", "**hello**", "<strong>hello</strong>"},
		{"standalone stars", "  **  \ntext", "\ntext"},
		{"no artifacts", "plain text", "plain text"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanMarkdownArtifacts(tt.input)
			if result != tt.expected {
				t.Errorf("cleanMarkdownArtifacts(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestAutolinkListItems(t *testing.T) {
	pageLinks := map[string]string{
		"About": "/about/",
	}

	input := "- About\n- Unknown\nRegular text"
	result := autolinkListItems(input, pageLinks)

	if !strings.Contains(result, "[About](/about/)") {
		t.Errorf("Expected autolinked About, got: %s", result)
	}
	if !strings.Contains(result, "- Unknown") {
		t.Errorf("Expected unchanged Unknown, got: %s", result)
	}
}

func TestConvertMarkdownToHTML(t *testing.T) {
	result := convertMarkdownToHTML("# Hello")
	if !strings.Contains(result, "<h1>Hello</h1>") {
		t.Errorf("Expected h1, got: %s", result)
	}
}

func TestGenerateWithRelativeLinks(t *testing.T) {
	tmpDir := t.TempDir()

	contentDir := filepath.Join(tmpDir, "content")
	sourceDir := filepath.Join(contentDir, "test-source")
	pagesDir := filepath.Join(sourceDir, "pages")
	postsDir := filepath.Join(sourceDir, "posts")
	templateDir := filepath.Join(tmpDir, "templates", "simple")
	outputDir := filepath.Join(tmpDir, "output")

	for _, dir := range []string{pagesDir, postsDir, templateDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create dir: %v", err)
		}
	}

	metadata := `{"categories":[],"media":[],"users":[]}`
	if err := os.WriteFile(filepath.Join(sourceDir, "metadata.json"), []byte(metadata), 0644); err != nil {
		t.Fatalf("Failed to create metadata: %v", err)
	}

	templates := map[string]string{
		"base.html":     `<!DOCTYPE html><html><body>base</body></html>`,
		"index.html":    `<!DOCTYPE html><html><body><a href="https://example.com/about">About</a></body></html>`,
		"page.html":     `<!DOCTYPE html><html><body>Page</body></html>`,
		"post.html":     `<!DOCTYPE html><html><body>Post</body></html>`,
		"category.html": `<!DOCTYPE html><html><body>Cat</body></html>`,
	}
	for name, content := range templates {
		if err := os.WriteFile(filepath.Join(templateDir, name), []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create template: %v", err)
		}
	}

	cfg := Config{
		Source:        "test-source",
		Template:      "simple",
		Domain:        "https://example.com",
		ContentDir:    contentDir,
		TemplatesDir:  filepath.Join(tmpDir, "templates"),
		OutputDir:     outputDir,
		RelativeLinks: true,
		Quiet:         true,
	}

	gen, _ := New(cfg)
	if err := gen.Generate(); err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	indexContent, _ := os.ReadFile(filepath.Join(outputDir, "index.html"))
	if strings.Contains(string(indexContent), "https://example.com") {
		t.Error("Absolute links should be converted to relative")
	}
}

func TestPrettifyIfRequestedSkipsWhenMinify(t *testing.T) {
	gen := &Generator{
		config: Config{
			PrettyHTML: true,
			MinifyHTML: true,
		},
	}
	err := gen.prettifyIfRequested()
	if err != nil {
		t.Errorf("Expected nil when MinifyHTML overrides PrettyHTML, got: %v", err)
	}
}

func TestMinifyIfRequestedSkipsWhenDisabled(t *testing.T) {
	gen := &Generator{
		config: Config{
			MinifyHTML: false,
			MinifyCSS:  false,
			MinifyJS:   false,
		},
	}
	err := gen.minifyIfRequested()
	if err != nil {
		t.Errorf("Expected nil when all minify disabled, got: %v", err)
	}
}

func TestLogContentStats(t *testing.T) {
	gen := &Generator{
		siteData: &models.SiteData{
			Pages:      []models.Page{{Title: "A"}, {Title: "B"}},
			Posts:      []models.Page{{Title: "P"}},
			Categories: map[int]models.Category{1: {ID: 1}},
			Media:      map[int]models.MediaItem{1: {ID: 1}},
		},
	}
	gen.logContentStats()
}

func TestBuildPageLinks(t *testing.T) {
	gen := &Generator{
		siteData: &models.SiteData{
			Pages: []models.Page{
				{Title: "About Us", Slug: "about", Link: "https://example.com/about/"},
			},
			Posts: []models.Page{
				{Title: "Hello World", Slug: "hello", Link: "https://example.com/hello/"},
			},
		},
	}

	links := gen.buildPageLinks()
	if _, ok := links["About Us"]; !ok {
		t.Error("Expected 'About Us' in page links")
	}
	if _, ok := links["Hello World"]; !ok {
		t.Error("Expected 'Hello World' in page links")
	}
}

func TestGenerateCategoriesMultiplePosts(t *testing.T) {
	tmpDir := t.TempDir()

	tmpl := template.Must(template.New("category.html").Parse("<html>{{.Category.Name}} - {{len .Posts}} posts</html>"))

	gen := &Generator{
		config: Config{
			OutputDir: tmpDir,
			Domain:    "example.com",
		},
		siteData: &models.SiteData{
			Domain: "example.com",
			Posts: []models.Page{
				{Title: "Post 1", Slug: "post1", Categories: []int{2, 3}},
				{Title: "Post 2", Slug: "post2", Categories: []int{2}},
				{Title: "Post 3", Slug: "post3", Categories: []int{3}},
			},
			Categories: map[int]models.Category{
				2: {ID: 2, Name: "News", Slug: "news"},
				3: {ID: 3, Name: "Tech", Slug: "tech"},
			},
		},
		tmpl: tmpl,
	}

	if err := gen.generateCategories(); err != nil {
		t.Fatalf("generateCategories failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(tmpDir, "category", "news", "index.html")); err != nil {
		t.Error("News category page not generated")
	}
	if _, err := os.Stat(filepath.Join(tmpDir, "category", "tech", "index.html")); err != nil {
		t.Error("Tech category page not generated")
	}
}

func TestRunStep(t *testing.T) {
	gen := &Generator{config: Config{Quiet: true}}

	errFn := func() error { return fmt.Errorf("step failed") }
	err := gen.runStep("test", errFn, "test context")
	if err == nil {
		t.Error("Expected error from runStep")
	}
	if !strings.Contains(err.Error(), "test context") {
		t.Errorf("Expected error to contain 'test context', got: %v", err)
	}

	okFn := func() error { return nil }
	err = gen.runStep("test", okFn, "test context")
	if err != nil {
		t.Errorf("Expected nil error from runStep, got: %v", err)
	}
}

func TestNewWithShortcodes(t *testing.T) {
	cfg := Config{
		Shortcodes: []Shortcode{
			{Name: "promo", Text: "Promo text"},
			{Name: "banner", Text: "Banner text"},
		},
	}

	gen, err := New(cfg)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	if len(gen.shortcodeMap) != 2 {
		t.Errorf("Expected 2 shortcodes in map, got %d", len(gen.shortcodeMap))
	}
	if _, ok := gen.shortcodeMap["promo"]; !ok {
		t.Error("Expected 'promo' in shortcode map")
	}
}

func TestTmplDecodeHTML(t *testing.T) {
	result := tmplDecodeHTML("&amp; &lt; &gt; &#8211;")
	if !strings.Contains(result, "&") || !strings.Contains(result, "<") || !strings.Contains(result, ">") {
		t.Errorf("tmplDecodeHTML failed: %s", result)
	}
}

func TestTmplFormatDatePLCov(t *testing.T) {
	tests := []struct {
		name     string
		date     time.Time
		expected string
	}{
		{"january", parseDateCov("2024-01-15"), "15 stycznia 2024"},
		{"june", parseDateCov("2024-06-01"), "1 czerwca 2024"},
		{"december", parseDateCov("2024-12-31"), "31 grudnia 2024"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tmplFormatDatePL(tt.date)
			if result != tt.expected {
				t.Errorf("tmplFormatDatePL = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestTmplGetURL(t *testing.T) {
	page := models.Page{
		Slug: "test",
		Link: "https://example.com/test/",
	}
	result := tmplGetURL(page)
	if result == "" {
		t.Error("Expected non-empty URL")
	}
}

func TestTmplGetCanonical(t *testing.T) {
	page := models.Page{
		Slug: "test",
		Link: "https://example.com/test/",
	}
	result := tmplGetCanonical(page, "example.com")
	if result == "" {
		t.Error("Expected non-empty canonical URL")
	}
}

func TestTmplHasValidCategories(t *testing.T) {
	pageWithCats := models.Page{Categories: []int{2, 3}}
	if !tmplHasValidCategories(pageWithCats) {
		t.Error("Expected true for page with valid categories")
	}

	pageWithOnlyDefault := models.Page{Categories: []int{1}}
	result := tmplHasValidCategories(pageWithOnlyDefault)
	_ = result
}

func TestFixMediaPathsThumbnailInSrcset(t *testing.T) {
	media := make(map[int]models.MediaItem)

	input := `<img srcset="/media/123_photo-300x200.jpg 300w, /media/123_photo-600x400.jpg 600w">`
	result := fixMediaPaths(input, media)

	if strings.Contains(result, "-300x200") || strings.Contains(result, "-600x400") {
		t.Errorf("Thumbnail suffixes should be removed from srcset, got: %s", result)
	}
}

func TestLoadMarkdownDirSetsSourceDirCov(t *testing.T) {
	tmpDir := t.TempDir()

	mdContent := `---
title: Test
status: publish
slug: test
link: https://example.com/test/
---
Content`
	if err := os.WriteFile(filepath.Join(tmpDir, "test.md"), []byte(mdContent), 0644); err != nil {
		t.Fatalf("Failed to create md: %v", err)
	}

	gen := &Generator{}
	pages, err := gen.loadMarkdownDir(tmpDir)
	if err != nil {
		t.Fatalf("loadMarkdownDir failed: %v", err)
	}

	if len(pages) != 1 {
		t.Fatalf("Expected 1 page, got %d", len(pages))
	}
	if pages[0].SourceDir != tmpDir {
		t.Errorf("Expected SourceDir=%q, got %q", tmpDir, pages[0].SourceDir)
	}
}

func TestLoadMarkdownDirWithImageFiles(t *testing.T) {
	tmpDir := t.TempDir()

	mdContent := `---
title: Page With Images
status: publish
slug: with-images
link: https://example.com/with-images/
---
Content with images`
	if err := os.WriteFile(filepath.Join(tmpDir, "page.md"), []byte(mdContent), 0644); err != nil {
		t.Fatalf("Failed to create md: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "hero.jpg"), []byte("jpg"), 0644); err != nil {
		t.Fatalf("Failed to create jpg: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "diagram.png"), []byte("png"), 0644); err != nil {
		t.Fatalf("Failed to create png: %v", err)
	}

	gen := &Generator{}
	pages, err := gen.loadMarkdownDir(tmpDir)
	if err != nil {
		t.Fatalf("loadMarkdownDir failed: %v", err)
	}

	if len(pages) != 1 {
		t.Errorf("Expected 1 page (images skipped), got %d", len(pages))
	}
}

func TestMinifyCSSFileEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "media query",
			input:    "@media (max-width: 768px) {\n  body { color: red; }\n}",
			expected: "@media (max-width:768px){body{color:red;}}",
		},
		{
			name:     "empty file",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			cssPath := filepath.Join(tmpDir, "test.css")

			if err := os.WriteFile(cssPath, []byte(tt.input), 0644); err != nil {
				t.Fatalf("Failed to write: %v", err)
			}
			if err := minifyCSSFile(cssPath); err != nil {
				t.Fatalf("minifyCSSFile failed: %v", err)
			}
			result, _ := os.ReadFile(cssPath)
			if string(result) != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, string(result))
			}
		})
	}
}

func TestMinifyJSFileEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty file",
			input:    "",
			expected: "",
		},
		{
			name:     "multi-line comment then code",
			input:    "/* header\n * comment\n */\nvar x = 1;\nvar y = 2;",
			expected: "var x = 1;\nvar y = 2;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			jsPath := filepath.Join(tmpDir, "test.js")

			if err := os.WriteFile(jsPath, []byte(tt.input), 0644); err != nil {
				t.Fatalf("Failed to write: %v", err)
			}
			if err := minifyJSFile(jsPath); err != nil {
				t.Fatalf("minifyJSFile failed: %v", err)
			}
			result, _ := os.ReadFile(jsPath)
			if string(result) != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, string(result))
			}
		})
	}
}

func TestMinifyHTMLFileEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty file",
			input:    "",
			expected: "",
		},
		{
			name:     "multiple comments",
			input:    "<!-- a --><div><!-- b --></div>",
			expected: "<div></div>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			htmlPath := filepath.Join(tmpDir, "test.html")

			if err := os.WriteFile(htmlPath, []byte(tt.input), 0644); err != nil {
				t.Fatalf("Failed to write: %v", err)
			}
			if err := minifyHTMLFile(htmlPath); err != nil {
				t.Fatalf("minifyHTMLFile failed: %v", err)
			}
			result, _ := os.ReadFile(htmlPath)
			if string(result) != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, string(result))
			}
		})
	}
}

func parseDateCov(s string) (t time.Time) {
	t, _ = time.Parse("2006-01-02", s)
	return
}
