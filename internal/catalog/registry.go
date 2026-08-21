package catalog

import (
	"sort"
	"strings"
)

type CategoryRule struct {
	Code        string
	Label       string
	Retention   int
	NeedsReview bool
	Description string
}

var rules = map[string]CategoryRule{
	"design":   {Code: "design", Label: "设计文件", Retention: 10, NeedsReview: true, Description: "图纸、设计说明与计算书"},
	"contract": {Code: "contract", Label: "合同文件", Retention: 15, NeedsReview: true, Description: "合同、变更与结算依据"},
	"safety":   {Code: "safety", Label: "安全文件", Retention: 20, NeedsReview: true, Description: "安全交底、检查和整改"},
	"quality":  {Code: "quality", Label: "质量文件", Retention: 20, NeedsReview: true, Description: "验收、检测与不合格项"},
}

func Rule(code string) (CategoryRule, bool) {
	rule, ok := rules[strings.ToLower(strings.TrimSpace(code))]
	return rule, ok
}

func AllRules() []CategoryRule {
	result := make([]CategoryRule, 0, len(rules))
	for _, rule := range rules {
		result = append(result, rule)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Code < result[j].Code })
	return result
}

func RetentionYears(code string) int {
	if rule, ok := Rule(code); ok {
		return rule.Retention
	}
	return 0
}

func RequiresReview(categories []string) bool {
	for _, category := range categories {
		if rule, ok := Rule(category); ok && rule.NeedsReview {
			return true
		}
	}
	return false
}

func Describe(categories []string) []string {
	result := []string{}
	seen := map[string]bool{}
	for _, category := range categories {
		key := strings.ToLower(strings.TrimSpace(category))
		if seen[key] {
			continue
		}
		if rule, ok := Rule(key); ok {
			result = append(result, rule.Label+": "+rule.Description)
			seen[key] = true
		}
	}
	sort.Strings(result)
	return result
}

func ValidateRetention(code string, years int) bool {
	minimum := RetentionYears(code)
	if minimum == 0 || years < minimum {
		return false
	}
	return true
}

func NormalizeWithDefaults(categories []string) []string {
	result := Normalize(categories)
	if len(result) == 0 {
		result = append(result, "design")
	}
	return result
}
