package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

const (
	coverageFile     = "coverage.out"
	packageTargetPct = 60.0
	fileTargetPct    = 60.0
	packageDepthV2   = 2
)

var packageScope = []string{"internal"}
var packageScopeExclude = []string{"*/api*", "*/infra*"}
var fileNameIgnoreScope = []string{"*query*", "*repository*"}

type covStat struct {
	total   int64
	covered int64
}

type covBlock struct {
	file          string
	pkg           string
	numStatements int64
	covered       bool
}

type pkgRow struct {
	name   string
	cover  float64
	status string
	result string
}

type fileRow struct {
	name   string
	cover  float64
	result string
}

var useColor bool
var orderBy string

func main() {
	flag.BoolVar(&useColor, "color", false, "enable colored output")
	flag.StringVar(&orderBy, "order", "name", "order by name or coverage")
	flag.Parse()
	if orderBy != "name" && orderBy != "coverage" {
		orderBy = "name"
	}

	os.Remove(coverageFile)

	testPkgs, err := packagesWithTests()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to list packages with tests: %v\n", err)
		os.Exit(1)
	}

	logOutput, exitCode := runGoTest()

	if _, err := os.Stat(coverageFile); err == nil {
		pkgStats, fileStats, err := parseCoverageProfile(coverageFile, testPkgs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse coverage profile: %v\n", err)
			fmt.Println()
			fmt.Println("== Error List ==")
			printErrorList(logOutput)
			os.Exit(1)
		}

		printPackageCoverageV2(pkgStats, testPkgs)
		fmt.Println()
		printFileCoverage(fileStats)
	} else {
		printPackageCoverageEmpty()
		fmt.Println()
		printFileCoverageEmpty()
	}

	fmt.Println()
	fmt.Println("== Error List ==")
	printErrorList(logOutput)

	if exitCode != 0 {
		os.Exit(exitCode)
	}
}

func packagesWithTests() (map[string]bool, error) {
	cmd := exec.Command("go", "list", "-f", "{{if or (gt (len .TestGoFiles) 0) (gt (len .XTestGoFiles) 0)}}{{.ImportPath}}{{end}}", "./...")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	result := make(map[string]bool)
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		result[line] = true
	}
	return result, scanner.Err()
}

func runGoTest() (string, int) {
	cmd := exec.Command("go", "test", "-coverpkg=./...", "-coverprofile="+coverageFile, "./...")
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	err := cmd.Run()
	if err == nil {
		return buf.String(), 0
	}

	var exitErr *exec.ExitError
	if ok := asExitError(err, &exitErr); ok {
		return buf.String(), exitErr.ExitCode()
	}

	fmt.Fprintf(os.Stderr, "failed to run go test: %v\n", err)
	return buf.String(), 1
}

func asExitError(err error, target **exec.ExitError) bool {
	e, ok := err.(*exec.ExitError)
	if ok {
		*target = e
		return true
	}
	return false
}

func parseCoverageProfile(path string, testPkgs map[string]bool) (map[string]covStat, map[string]covStat, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	pkgStats := make(map[string]covStat)
	fileStats := make(map[string]covStat)
	blocks := make(map[string]covBlock)

	scanner := bufio.NewScanner(bytes.NewReader(data))
	first := true
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if first {
			first = false
			if strings.HasPrefix(line, "mode:") {
				continue
			}
		}

		// format:
		// <file>:startLine.startCol,endLine.endCol numStatements count
		fields := strings.Fields(line)
		if len(fields) != 3 {
			continue
		}

		fileAndRange := fields[0]
		numStatementsStr := fields[1]
		countStr := fields[2]

		colonIdx := strings.Index(fileAndRange, ":")
		if colonIdx < 0 {
			continue
		}
		file := fileAndRange[:colonIdx]
		pkg := filepath.ToSlash(filepath.Dir(file))

		numStatements, err := strconv.ParseInt(numStatementsStr, 10, 64)
		if err != nil {
			continue
		}
		count, err := strconv.ParseInt(countStr, 10, 64)
		if err != nil {
			continue
		}

		key := fmt.Sprintf("%s %d", fileAndRange, numStatements)
		block := blocks[key]
		if block.file == "" {
			block = covBlock{
				file:          file,
				pkg:           pkg,
				numStatements: numStatements,
				covered:       count > 0,
			}
		} else if count > 0 {
			block.covered = true
		}
		blocks[key] = block
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	for _, block := range blocks {
		if isIgnoredFile(block.file) {
			continue
		}
		if isOutOfScopePackage(block.pkg) {
			continue
		}
		ps := pkgStats[block.pkg]
		ps.total += block.numStatements
		if block.covered {
			ps.covered += block.numStatements
		}
		pkgStats[block.pkg] = ps

		if testPkgs[block.pkg] {
			fs := fileStats[block.file]
			fs.total += block.numStatements
			if block.covered {
				fs.covered += block.numStatements
			}
			fileStats[block.file] = fs
		}
	}

	return pkgStats, fileStats, nil
}

func printPackageCoverage(pkgStats map[string]covStat, testPkgs map[string]bool) {
	fmt.Println("== Package coverage ==")
	headers := []string{"Package", "Coverage", "Test Status", "Result"}

	pkgs := make([]string, 0, len(pkgStats))
	for pkg := range pkgStats {
		pkgs = append(pkgs, pkg)
	}
	sort.Strings(pkgs)

	rows := make([]pkgRow, 0, len(pkgs))
	for _, pkg := range pkgs {
		stat := pkgStats[pkg]
		pct := percent(stat)
		status := "no tests"
		if testPkgs[pkg] {
			status = "has tests"
		}

		judge := "PASS"
		if pct < packageTargetPct {
			judge = fmt.Sprintf("FAIL(<%.0f%%)", packageTargetPct)
		}

		rows = append(rows, pkgRow{
			name:   pkg,
			cover:  pct,
			status: status,
			result: judge,
		})
	}
	sortPackageRows(rows)

	display := make([][]string, 0, len(rows))
	for _, row := range rows {
		display = append(display, []string{
			colorizeByResult(normalizePackageForScope(row.name), row.result),
			fmt.Sprintf("%.2f%%", row.cover),
			row.status,
			colorizeResult(row.result),
		})
	}
	printTable(headers, display)
}

func printPackageCoverageV2(pkgStats map[string]covStat, testPkgs map[string]bool) {
	fmt.Println("== Package coverage (depth 3) ==")
	headers := []string{"Package", "Coverage", "Test Status", "Result"}

	groupedStats, groupedTestPkgs := aggregatePackageCoverage(pkgStats, testPkgs, packageDepthV2)

	pkgs := make([]string, 0, len(groupedStats))
	for pkg := range groupedStats {
		pkgs = append(pkgs, pkg)
	}
	sort.Strings(pkgs)

	rows := make([]pkgRow, 0, len(pkgs))
	for _, pkg := range pkgs {
		stat := groupedStats[pkg]
		pct := percent(stat)
		status := "no tests"
		if groupedTestPkgs[pkg] {
			status = "has tests"
		}

		judge := "PASS"
		if pct < packageTargetPct {
			judge = fmt.Sprintf("FAIL(<%.0f%%)", packageTargetPct)
		}

		rows = append(rows, pkgRow{
			name:   pkg,
			cover:  pct,
			status: status,
			result: judge,
		})
	}
	sortPackageRows(rows)

	display := make([][]string, 0, len(rows))
	for _, row := range rows {
		display = append(display, []string{
			colorizeByResult(normalizePackageForScope(row.name), row.result),
			fmt.Sprintf("%.2f%%", row.cover),
			row.status,
			colorizeResult(row.result),
		})
	}
	printTable(headers, display)
}

func printFileCoverage(fileStats map[string]covStat) {
	fmt.Println("== File coverage (only packages that have tests) ==")
	headers := []string{"File", "Coverage", "Result"}

	files := make([]string, 0, len(fileStats))
	for file := range fileStats {
		files = append(files, file)
	}
	sort.Strings(files)

	rows := make([]fileRow, 0, len(files))
	for _, file := range files {
		stat := fileStats[file]
		pct := percent(stat)

		judge := "PASS"
		if pct < fileTargetPct {
			judge = fmt.Sprintf("FAIL(<%.0f%%)", fileTargetPct)
		}

		rows = append(rows, fileRow{
			name:   normalizeFileForScope(file),
			cover:  pct,
			result: judge,
		})
	}
	sortFileRows(rows)

	display := make([][]string, 0, len(rows))
	for _, row := range rows {
		display = append(display, []string{
			colorizeByResult(row.name, row.result),
			fmt.Sprintf("%.2f%%", row.cover),
			colorizeResult(row.result),
		})
	}
	printTable(headers, display)
}

func percent(stat covStat) float64 {
	if stat.total == 0 {
		return 0
	}
	return float64(stat.covered) * 100 / float64(stat.total)
}

func printErrorList(logOutput string) {
	lines := strings.Split(strings.ReplaceAll(logOutput, "\r\n", "\n"), "\n")
	var out []string

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		if strings.HasPrefix(line, "# ") {
			out = append(out, line)
			for j := i + 1; j < len(lines); j++ {
				next := lines[j]
				if next == "" {
					i = j
					break
				}
				if strings.HasPrefix(next, "# ") || strings.HasPrefix(next, "FAIL\t") || strings.HasPrefix(next, "FAIL ") {
					i = j - 1
					break
				}
				if looksLikeCoverageLine(next) || strings.HasPrefix(next, "ok  \t") || strings.HasPrefix(next, "ok\t") {
					i = j - 1
					break
				}
				out = append(out, next)
				if j == len(lines)-1 {
					i = j
				}
			}
			continue
		}

		if strings.HasPrefix(line, "FAIL\t") || strings.HasPrefix(line, "FAIL ") {
			out = append(out, line)
		}
	}

	if len(out) == 0 {
		printTable([]string{"No.", "Error Message"}, [][]string{
			{"1", "none"},
		})
		return
	}

	rows := make([][]string, 0, len(out))
	for i, line := range out {
		rows = append(rows, []string{strconv.Itoa(i + 1), line})
	}
	printTable([]string{"No.", "Error Message"}, rows)
}

func looksLikeCoverageLine(s string) bool {
	return strings.Contains(s, "coverage:") && !strings.HasPrefix(strings.TrimSpace(s), "#")
}

func printPackageCoverageEmpty() {
	fmt.Println("== Package coverage (depth 3) ==")
	printTable(
		[]string{"Package", "Coverage", "Test Status", "Result"},
		[][]string{{"-", "(coverage.out not generated)", "-", "-"}},
	)
}

func printFileCoverageEmpty() {
	fmt.Println("== File coverage (only packages that have tests) ==")
	printTable(
		[]string{"File", "Coverage", "Result"},
		[][]string{{"-", "(coverage.out not generated)", "-"}},
	)
}

func printTable(headers []string, rows [][]string) {
	fmt.Println("[//]: # (@formatter:off)")
	widths := tableWidths(headers, rows)
	printTableLine(headers, widths)
	printTableSeparator(widths)
	for _, row := range rows {
		printTableLine(row, widths)
	}
	fmt.Println("[//]: # (@formatter:on)")
}

func printTableLine(cells []string, widths []int) {
	var b strings.Builder
	b.WriteString("|")
	for i := range widths {
		cell := ""
		if i < len(cells) {
			cell = cells[i]
		}
		padding := widths[i] - displayWidth(cell)
		if padding < 0 {
			padding = 0
		}
		b.WriteString(" ")
		b.WriteString(cell)
		if padding > 0 {
			b.WriteString(strings.Repeat(" ", padding))
		}
		b.WriteString(" |")
	}
	fmt.Println(b.String())
}

func printTableSeparator(widths []int) {
	var b strings.Builder
	b.WriteString("|")
	for _, w := range widths {
		b.WriteString(" ")
		b.WriteString(strings.Repeat("-", w))
		b.WriteString(" |")
	}
	fmt.Println(b.String())
}

func tableWidths(headers []string, rows [][]string) []int {
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = displayWidth(h)
	}
	for _, row := range rows {
		for i := range headers {
			cell := ""
			if i < len(row) {
				cell = row[i]
			}
			if w := displayWidth(cell); w > widths[i] {
				widths[i] = w
			}
		}
	}
	return widths
}

func displayWidth(s string) int {
	return len([]rune(stripAnsi(s)))
}

func aggregatePackageCoverage(pkgStats map[string]covStat, testPkgs map[string]bool, depth int) (map[string]covStat, map[string]bool) {
	groupedStats := make(map[string]covStat)
	groupedTestPkgs := make(map[string]bool)

	for pkg, stat := range pkgStats {
		if isOutOfScopePackage(pkg) {
			continue
		}
		normalized := normalizePackageForScope(pkg)
		if normalized == "" {
			continue
		}
		key := packageDepthKey(normalized, depth)
		merged := groupedStats[key]
		merged.total += stat.total
		merged.covered += stat.covered
		groupedStats[key] = merged
	}

	for pkg, hasTests := range testPkgs {
		if !hasTests {
			continue
		}
		if isOutOfScopePackage(pkg) {
			continue
		}
		normalized := normalizePackageForScope(pkg)
		if normalized == "" {
			continue
		}
		key := packageDepthKey(normalized, depth)
		groupedTestPkgs[key] = true
	}

	return groupedStats, groupedTestPkgs
}

func packageDepthKey(pkg string, depth int) string {
	if depth <= 0 {
		return pkg
	}
	normalized := strings.TrimPrefix(filepath.ToSlash(pkg), "./")
	parts := strings.Split(normalized, "/")
	if len(parts) <= depth {
		return normalized
	}
	return strings.Join(parts[:depth], "/")
}

func isOutOfScopePackage(pkg string) bool {
	normalized := normalizePackageForScope(pkg)
	if normalized == "" {
		return true
	}
	if !isAllowedTopLevel(normalized) {
		return true
	}
	if isExcludedPackage(normalized) {
		return true
	}
	return false
}

func isAllowedTopLevel(pkg string) bool {
	for _, entry := range packageScope {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		if pkg == entry || strings.HasPrefix(pkg, entry+"/") {
			return true
		}
	}
	return false
}

func isExcludedPackage(pkg string) bool {
	for _, entry := range packageScopeExclude {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		if strings.Contains(entry, "*") {
			if wildcardMatch(entry, pkg) {
				return true
			}
			continue
		}
		if strings.Contains(entry, "/") {
			if pkg == entry || strings.HasPrefix(pkg, entry+"/") {
				return true
			}
			continue
		}
		if hasPackageSegment(pkg, entry) {
			return true
		}
	}
	return false
}

func hasPackageSegment(pkg string, segment string) bool {
	for _, part := range strings.Split(pkg, "/") {
		if part == segment {
			return true
		}
	}
	return false
}

func normalizePackageForScope(pkg string) string {
	normalized := strings.TrimPrefix(filepath.ToSlash(pkg), "./")
	if strings.HasPrefix(normalized, "internal/") || strings.HasPrefix(normalized, "pkg/") {
		return normalized
	}
	if idx := strings.Index(normalized, "/internal/"); idx >= 0 {
		return normalized[idx+1:]
	}
	if idx := strings.Index(normalized, "/pkg/"); idx >= 0 {
		return normalized[idx+1:]
	}
	return ""
}

func normalizeFileForScope(file string) string {
	normalizedFile := filepath.ToSlash(file)
	dir := filepath.Dir(normalizedFile)
	normalizedDir := normalizePackageForScope(dir)
	if normalizedDir == "" {
		return normalizedFile
	}
	return filepath.ToSlash(filepath.Join(normalizedDir, filepath.Base(normalizedFile)))
}

func isIgnoredFile(file string) bool {
	base := filepath.Base(file)
	for _, entry := range fileNameIgnoreScope {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		if wildcardMatch(entry, base) {
			return true
		}
	}
	return false
}

func wildcardMatch(pattern string, value string) bool {
	escaped := regexp.QuoteMeta(pattern)
	escaped = strings.ReplaceAll(escaped, "\\*", ".*")
	escaped = "^" + escaped + "$"
	matched, err := regexp.MatchString(escaped, value)
	if err != nil {
		return false
	}
	return matched
}

func colorizeResult(result string) string {
	if !useColor {
		return result
	}
	if strings.HasPrefix(result, "PASS") {
		return "\x1b[32m" + result + "\x1b[0m"
	}
	if strings.HasPrefix(result, "FAIL") {
		return "\x1b[31m" + result + "\x1b[0m"
	}
	return result
}

func colorizeByResult(value string, result string) string {
	if !useColor {
		return value
	}
	if strings.HasPrefix(result, "PASS") {
		return "\x1b[32m" + value + "\x1b[0m"
	}
	if strings.HasPrefix(result, "FAIL") {
		return "\x1b[31m" + value + "\x1b[0m"
	}
	return value
}

func stripAnsi(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	inEscape := false
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if inEscape {
			if ch == 'm' {
				inEscape = false
			}
			continue
		}
		if ch == 0x1b && i+1 < len(s) && s[i+1] == '[' {
			inEscape = true
			i++
			continue
		}
		b.WriteByte(ch)
	}
	return b.String()
}

func sortPackageRows(rows []pkgRow) {
	if orderBy == "coverage" {
		sort.Slice(rows, func(i, j int) bool {
			if rows[i].cover == rows[j].cover {
				return rows[i].name < rows[j].name
			}
			return rows[i].cover > rows[j].cover
		})
		return
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].name < rows[j].name
	})
}

func sortFileRows(rows []fileRow) {
	if orderBy == "coverage" {
		sort.Slice(rows, func(i, j int) bool {
			if rows[i].cover == rows[j].cover {
				return rows[i].name < rows[j].name
			}
			return rows[i].cover > rows[j].cover
		})
		return
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].name < rows[j].name
	})
}
func init() {
	_ = runtime.GOOS
}
