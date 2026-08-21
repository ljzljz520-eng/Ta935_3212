package catalog

import "strings"

var classifications = map[string]string{
	"design":   "设计文件",
	"contract": "合同文件",
	"safety":   "安全文件",
	"quality":  "质量文件",
}

func IsKnownCategory(category string) bool {
	_, ok := classifications[strings.ToLower(strings.TrimSpace(category))]
	return ok
}

func Label(category string) string {
	if label, ok := classifications[strings.ToLower(strings.TrimSpace(category))]; ok {
		return label
	}
	return "未分类"
}

func Normalize(categories []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(categories))
	for _, category := range categories {
		value := strings.ToLower(strings.TrimSpace(category))
		if value != "" && !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

func AllowedCategories() []string {
	return []string{"design", "contract", "safety", "quality"}
}
