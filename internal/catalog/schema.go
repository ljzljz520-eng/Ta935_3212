package catalog

import (
	"fmt"
	"strings"
)

type SchemaField struct {
	Code       string
	Label      string
	DataType   string
	Searchable bool
}

var schemaFields = []SchemaField{
	{Code: "project_code", Label: "工程编号", DataType: "text", Searchable: true},
	{Code: "project_name", Label: "工程名称", DataType: "text", Searchable: true},
	{Code: "discipline", Label: "专业类别", DataType: "enum", Searchable: true},
	{Code: "contract_ref", Label: "合同引用", DataType: "text", Searchable: true},
	{Code: "revision_note", Label: "修订说明", DataType: "text", Searchable: true},
	{Code: "approval_scope", Label: "批准范围", DataType: "enum", Searchable: false},
	{Code: "retention_class", Label: "保管类别", DataType: "enum", Searchable: true},
	{Code: "safety_level", Label: "安全等级", DataType: "enum", Searchable: true},
	{Code: "quality_gate", Label: "质量关口", DataType: "enum", Searchable: true},
	{Code: "site_location", Label: "现场位置", DataType: "text", Searchable: true},
	{Code: "responsible_unit", Label: "责任单位", DataType: "text", Searchable: true},
	{Code: "review_deadline", Label: "审核期限", DataType: "date", Searchable: false},
	{Code: "archive_box", Label: "归档盒号", DataType: "text", Searchable: true},
	{Code: "archive_checksum", Label: "归档校验", DataType: "digest", Searchable: false},
}

func SchemaFields() []SchemaField { return append([]SchemaField(nil), schemaFields...) }

func SearchableFields() []SchemaField {
	result := []SchemaField{}
	for _, field := range schemaFields {
		if field.Searchable {
			result = append(result, field)
		}
	}
	return result
}

func FindField(code string) (SchemaField, bool) {
	for _, field := range schemaFields {
		if field.Code == strings.ToLower(strings.TrimSpace(code)) {
			return field, true
		}
	}
	return SchemaField{}, false
}

func ValidateValue(code, value string) error {
	field, ok := FindField(code)
	if !ok {
		return fmt.Errorf("unknown schema field: %s", code)
	}
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("value for %s is empty", field.Label)
	}
	switch field.DataType {
	case "digest":
		if len(strings.TrimSpace(value)) < 6 {
			return fmt.Errorf("value for %s is too short", field.Label)
		}
	case "date":
		if len(strings.TrimSpace(value)) != 10 {
			return fmt.Errorf("value for %s must be yyyy-mm-dd", field.Label)
		}
	}
	return nil
}

func SchemaCodes() []string {
	result := make([]string, 0, len(schemaFields))
	for _, field := range schemaFields {
		result = append(result, field.Code)
	}
	return result
}

func HasCode(code string) bool { _, ok := FindField(code); return ok }

func CategoryRetentionLabel(code string) string {
	years := RetentionYears(code)
	if years == 0 {
		return "未定义"
	}
	return fmt.Sprintf("保管%d年", years)
}
