package itertools_test

import (
	"context"
	"encoding/csv"
	"strings"
	"testing"
	"time"

	"github.com/amjadjibon/itertools"
	"github.com/stretchr/testify/assert"
)

func TestFromCSV(t *testing.T) {
	csvData := `name,age,city
Alice,30,NYC
Bob,25,LA
Charlie,35,Chicago`

	reader := csv.NewReader(strings.NewReader(csvData))
	iter := itertools.FromCSV(reader)
	records := iter.Collect()

	assert.Equal(t, 4, len(records))
	assert.Equal(t, []string{"name", "age", "city"}, records[0])
	assert.Equal(t, []string{"Alice", "30", "NYC"}, records[1])
	assert.Equal(t, []string{"Bob", "25", "LA"}, records[2])
	assert.Equal(t, []string{"Charlie", "35", "Chicago"}, records[3])
}

func TestFromCSV_WithFilter(t *testing.T) {
	csvData := `name,age,city
Alice,30,NYC
Bob,25,LA
Charlie,35,Chicago
Dave,28,NYC
Eve,32,LA`

	reader := csv.NewReader(strings.NewReader(csvData))
	iter := itertools.FromCSV(reader)

	// Skip header and filter by age > 28
	records := iter.
		Drop(1). // Skip header
		Filter(func(row []string) bool {
			return len(row) >= 2 && row[1] > "28" // Simple string comparison
		}).
		Collect()

	assert.Equal(t, 3, len(records))
	assert.Equal(t, []string{"Alice", "30", "NYC"}, records[0])
	assert.Equal(t, []string{"Charlie", "35", "Chicago"}, records[1])
	assert.Equal(t, []string{"Eve", "32", "LA"}, records[2])
}

func TestFromCSV_EarlyTermination(t *testing.T) {
	csvData := strings.Repeat("a,b,c\n", 1000)
	reader := csv.NewReader(strings.NewReader(csvData))
	iter := itertools.FromCSV(reader)

	records := iter.Take(5).Collect()

	assert.Equal(t, 5, len(records))
}

func TestFromCSVWithContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	csvData := strings.Repeat("a,b,c\n", 1000000)
	reader := csv.NewReader(strings.NewReader(csvData))

	// Cancel immediately
	cancel()

	iter := itertools.FromCSVWithContext(ctx, reader)
	records := iter.Collect()

	// Should have stopped early due to cancelled context
	assert.Less(t, len(records), 1000000)
}

func TestFromCSVWithHeaders(t *testing.T) {
	csvData := `name,age,city
Alice,30,NYC
Bob,25,LA
Charlie,35,Chicago`

	reader := csv.NewReader(strings.NewReader(csvData))
	iter, headers, err := itertools.FromCSVWithHeaders(reader)

	assert.NoError(t, err)
	assert.Equal(t, []string{"name", "age", "city"}, headers)

	records := iter.Collect()
	assert.Equal(t, 3, len(records))
	assert.Equal(t, []string{"Alice", "30", "NYC"}, records[0].Fields)
	assert.Equal(t, 0, records[0].Index)
	assert.Equal(t, []string{"Bob", "25", "LA"}, records[1].Fields)
	assert.Equal(t, 1, records[1].Index)
}

func TestCSVRow_Get(t *testing.T) {
	row := itertools.CSVRow{
		Fields: []string{"Alice", "30", "NYC"},
		Index:  0,
	}

	assert.Equal(t, "Alice", row.Get(0))
	assert.Equal(t, "30", row.Get(1))
	assert.Equal(t, "NYC", row.Get(2))
	assert.Equal(t, "", row.Get(3))  // Out of bounds
	assert.Equal(t, "", row.Get(-1)) // Negative index
}

func TestCSVRow_GetByHeader(t *testing.T) {
	headers := []string{"name", "age", "city"}
	row := itertools.CSVRow{
		Fields: []string{"Alice", "30", "NYC"},
		Index:  0,
	}

	assert.Equal(t, "Alice", row.GetByHeader(headers, "name"))
	assert.Equal(t, "30", row.GetByHeader(headers, "age"))
	assert.Equal(t, "NYC", row.GetByHeader(headers, "city"))
	assert.Equal(t, "", row.GetByHeader(headers, "country")) // Not found
}

func TestFromCSVWithHeaders_Filter(t *testing.T) {
	csvData := `name,age,city
Alice,30,NYC
Bob,25,LA
Charlie,35,Chicago
Dave,28,NYC
Eve,32,LA`

	reader := csv.NewReader(strings.NewReader(csvData))
	iter, headers, err := itertools.FromCSVWithHeaders(reader)

	assert.NoError(t, err)

	// Filter people from NYC
	nycPeople := iter.Filter(func(row itertools.CSVRow) bool {
		return row.GetByHeader(headers, "city") == "NYC"
	}).Collect()

	assert.Equal(t, 2, len(nycPeople))
	assert.Equal(t, "Alice", nycPeople[0].GetByHeader(headers, "name"))
	assert.Equal(t, "Dave", nycPeople[1].GetByHeader(headers, "name"))
}

func TestFromCSVWithHeaders_Map(t *testing.T) {
	csvData := `name,age
Alice,30
Bob,25
Charlie,35`

	reader := csv.NewReader(strings.NewReader(csvData))
	iter, headers, err := itertools.FromCSVWithHeaders(reader)

	assert.NoError(t, err)

	// Transform to uppercase names in CSVRow
	transformed := iter.Map(func(row itertools.CSVRow) itertools.CSVRow {
		name := row.GetByHeader(headers, "name")
		row.Fields[0] = strings.ToUpper(name)
		return row
	}).Collect()

	assert.Equal(t, 3, len(transformed))
	assert.Equal(t, "ALICE", transformed[0].GetByHeader(headers, "name"))
	assert.Equal(t, "BOB", transformed[1].GetByHeader(headers, "name"))
	assert.Equal(t, "CHARLIE", transformed[2].GetByHeader(headers, "name"))
}

func TestFromCSVWithHeaders_Complex(t *testing.T) {
	csvData := `product,price,quantity,category
Laptop,1200,5,Electronics
Mouse,25,50,Electronics
Desk,300,10,Furniture
Chair,150,20,Furniture
Keyboard,75,30,Electronics`

	reader := csv.NewReader(strings.NewReader(csvData))
	iter, headers, err := itertools.FromCSVWithHeaders(reader)

	assert.NoError(t, err)

	// Find all electronics
	electronics := iter.Filter(func(row itertools.CSVRow) bool {
		category := row.GetByHeader(headers, "category")
		return category == "Electronics"
	}).Collect()

	assert.Equal(t, 3, len(electronics))
	assert.Equal(t, "Laptop", electronics[0].GetByHeader(headers, "product"))
	assert.Equal(t, "Mouse", electronics[1].GetByHeader(headers, "product"))
	assert.Equal(t, "Keyboard", electronics[2].GetByHeader(headers, "product"))
}

func TestFromCSVWithHeadersContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	csvData := strings.Repeat("name,age,city\n", 1000000)
	reader := csv.NewReader(strings.NewReader(csvData))

	iter, headers, err := itertools.FromCSVWithHeadersContext(ctx, reader)

	assert.NoError(t, err)
	assert.Equal(t, []string{"name", "age", "city"}, headers)

	// This should timeout before reading all rows
	records := iter.Collect()
	assert.Less(t, len(records), 1000000)
}

func TestFromCSV_MalformedRows(t *testing.T) {
	csvData := `name,age,city
Alice,30,NYC
Bob,25
Charlie,35,Chicago,Extra
Dave,28,LA`

	reader := csv.NewReader(strings.NewReader(csvData))
	reader.FieldsPerRecord = -1 // Allow variable number of fields

	iter := itertools.FromCSV(reader)
	records := iter.Collect()

	// Should collect all rows including malformed ones
	assert.Equal(t, 5, len(records))
}

func TestFromCSV_LargeFile(t *testing.T) {
	// Simulate a large CSV file
	var sb strings.Builder
	sb.WriteString("id,name,value\n")
	for i := 0; i < 10000; i++ {
		sb.WriteString(strings.Join([]string{
			string(rune(i)),
			"item" + string(rune(i)),
			string(rune(i * 10)),
		}, ","))
		sb.WriteString("\n")
	}

	reader := csv.NewReader(strings.NewReader(sb.String()))
	iter := itertools.FromCSV(reader)

	// Process with chaining - skip header, take every 100th row
	samples := iter.
		Drop(1).
		StepBy(100).
		Take(10).
		Collect()

	assert.Equal(t, 10, len(samples))
}
