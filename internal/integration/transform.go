package integration

import (
	"fmt"
	"strings"
	"time"
)

// Transform applies a transformation rule to a value
func Transform(rule, value, arg string) (string, error) {
	switch rule {
	case "direct":
		return value, nil
	case "lab_map":
		return labMap(value), nil
	case "date_parse":
		return parseDate(value, arg), nil
	case "json_extract":
		return jsonExtract(value, arg), nil
	case "upper":
		return strings.ToUpper(value), nil
	case "lower":
		return strings.ToLower(value), nil
	case "abnormal_flag":
		return mapAbnormal(value), nil
	case "gender_map":
		return mapGender(value), nil
	default:
		return value, nil
	}
}

func labMap(hospitalExamType string) string {
	examType := strings.ToLower(hospitalExamType)
	switch {
	case strings.Contains(examType, "lab") || strings.Contains(examType, "检验"):
		return "lab"
	case strings.Contains(examType, "us") || strings.Contains(examType, "ultrasound") || strings.Contains(examType, "超声") || strings.Contains(examType, "b超"):
		return "imaging"
	case strings.Contains(examType, "ct") || strings.Contains(examType, "mri") || strings.Contains(examType, "x-ray"):
		return "imaging"
	case strings.Contains(examType, "出院") || strings.Contains(examType, "discharge") || strings.Contains(examType, "小结"):
		return "discharge"
	default:
		return "lab"
	}
}

func parseDate(value, format string) string {
	if format == "" {
		format = "20060102"
	}
	// Try common hospital date formats
	formats := []string{format, "2006-01-02", "2006/01/02", "20060102", "2006-01-02 15:04:05", "02/Jan/2006"}
	for _, f := range formats {
		if t, err := time.Parse(f, value); err == nil {
			return t.Format("2006-01-02")
		}
	}
	// Return as-is if unparseable
	return value
}

func jsonExtract(value, path string) string {
	// Simplified JSON path extraction
	// path: e.g., "result.key1.key2"
	return value
}

func mapAbnormal(value string) string {
	v := strings.ToLower(value)
	// Common hospital abnormal indicators
	switch {
	case v == "h" || v == "high" || v == "↑" || strings.Contains(v, "异常") || strings.Contains(v, "阳性"):
		return "high"
	case v == "l" || v == "low" || v == "↓":
		return "low"
	case strings.Contains(v, "危急"):
		return "critical"
	default:
		return "normal"
	}
}

func mapGender(value string) string {
	v := strings.ToLower(value)
	switch {
	case v == "m" || v == "男" || v == "male" || v == "1":
		return "1"
	case v == "f" || v == "女" || v == "female" || v == "0":
		return "0"
	default:
		return ""
	}
}

// NormalizeReportType converts various hospital report type strings to standard types
func NormalizeReportType(raw string) string {
	return labMap(raw)
}

// MapSource maps hospital system codes to standard source strings
func MapSource(sourceType string) string {
	m := map[string]string{
		"his":  "his",
		"lis":  "lis",
		"pacs": "pacs",
		"emr":  "manual",
		"manual": "manual",
	}
	if v, ok := m[strings.ToLower(sourceType)]; ok {
		return v
	}
	return "manual"
}

// TryParseDate attempts to parse a date string in common hospital formats
func TryParseDate(s string) string {
	if s == "" {
		return ""
	}
	return parseDate(s, "")
}

// ExtractAbnormalFlag extracts abnormal indicator from single test result
func ExtractAbnormalFlag(result, refRange string) string {
	if result == "" || refRange == "" {
		return ""
	}
	_ = refRange // In production: parse ref range and compare
	return ""
}

// BuildSummaryJSON creates a JSON summary from parsed report fields
func BuildSummaryJSON(fields map[string]string) string {
	// In production: marshal to structured JSON
	return fmt.Sprintf("%v", fields)
}
