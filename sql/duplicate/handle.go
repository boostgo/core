package duplicate

import (
	"regexp"
	"strings"

	"github.com/lib/pq"
)

type Error struct {
	Field      string
	Value      string
	Constraint string
}

func Handle(err *pq.Error) *Error {
	// Extract field name from constraint name
	// For constraint "partners_name_key", extract "name"
	field := extractFieldFromConstraint(err.Constraint)

	// Extract the duplicate value from the detail message
	value := extractValueFromDetail(err.Detail)

	return &Error{
		Field:      field,
		Value:      value,
		Constraint: err.Constraint,
	}
}

func extractFieldFromConstraint(constraint string) string {
	// Handle different constraint naming patterns
	if strings.Contains(constraint, "_key") {
		// Pattern: "table_field_key" -> extract "field"
		parts := strings.Split(constraint, "_")
		if len(parts) >= 2 {
			return parts[len(parts)-2] // Get the second-to-last part
		}
	}

	if strings.Contains(constraint, "_unique") {
		// Pattern: "table_field_unique" -> extract "field"
		parts := strings.Split(constraint, "_")
		if len(parts) >= 2 {
			return parts[len(parts)-2]
		}
	}

	return constraint
}

func extractValueFromDetail(detail string) string {
	// PostgreSQL detail format: "Key (field)=(value) already exists."
	re := regexp.MustCompile(`\(([^)]+)\)=\(([^)]+)\)`)
	matches := re.FindStringSubmatch(detail)
	if len(matches) >= 3 {
		return matches[2] // Return the value part
	}
	return ""
}
