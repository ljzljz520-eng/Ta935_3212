package catalog

import "testing"

func TestCategoryCatalogNormalizesKnownValues(t *testing.T) {
	items := Normalize([]string{" Design ", "design", "safety"})
	if len(items) != 2 || !IsKnownCategory(items[0]) || Label("safety") == "未分类" {
		t.Fatalf("unexpected catalog: %+v", items)
	}
}
