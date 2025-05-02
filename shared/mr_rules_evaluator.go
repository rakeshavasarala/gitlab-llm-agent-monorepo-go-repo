package shared

import (
	"strings"
)

type MR struct {
	Title         string
	Author        string
	Files         []string
	Description   string
	Diff          string
	CoverageDelta float64
}

func EvaluateRules(ruleset *Ruleset, mr MR) []RuleResult {
	results := []RuleResult{}

	for _, rule := range ruleset.Rules {
		if !rule.Enabled {
			continue
		}

		switch rule.ID {
		case "major-version-bump-check":
			applicable := strings.Contains(mr.Description, "bump") || strings.Contains(mr.Diff, "=>")
			if applicable && strings.Contains(mr.Diff, "v1") && strings.Contains(mr.Diff, "v2") {
				results = append(results, RuleResult{
					RuleID:     rule.ID,
					Applicable: true,
					Passed:     false,
					Message:    "Major version bump detected. Manual review required.",
				})
			} else if applicable {
				results = append(results, RuleResult{
					RuleID:     rule.ID,
					Applicable: true,
					Passed:     true,
					Message:    "Version bump detected but not major.",
				})
			} else {
				results = append(results, RuleResult{RuleID: rule.ID, Applicable: false})
			}

		case "helm-chart-version-bump":
			applicable := containsHelmChartFiles(mr.Files)
			if applicable && !strings.Contains(mr.Diff, "version:") {
				results = append(results, RuleResult{
					RuleID:     rule.ID,
					Applicable: true,
					Passed:     false,
					Message:    "Helm chart updated without version bump in Chart.yaml.",
				})
			} else if applicable {
				results = append(results, RuleResult{
					RuleID:     rule.ID,
					Applicable: true,
					Passed:     true,
					Message:    "Chart.yaml version bump detected.",
				})
			} else {
				results = append(results, RuleResult{RuleID: rule.ID, Applicable: false})
			}

		// Add more rules here...
		}
	}

	return results
}

func containsHelmChartFiles(files []string) bool {
	for _, f := range files {
		if strings.Contains(f, "Chart.yaml") || strings.Contains(f, "templates/") {
			return true
		}
	}
	return false
}
