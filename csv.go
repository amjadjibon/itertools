package itertools

import (
	"context"
	"encoding/csv"
	"io"
)

// FromCSV creates a lazy Iterator that reads records from a CSV reader.
// Each element is a []string representing one CSV row.
// This is useful for processing large CSV files without loading them entirely into memory.
//
// Example:
//
//	file, _ := os.Open("large_data.csv")
//	defer file.Close()
//	iter := itertools.FromCSV(csv.NewReader(file))
//	records := iter.Filter(func(row []string) bool {
//	    return len(row) > 0 && row[0] != ""
//	}).Take(100).Collect()
func FromCSV(r *csv.Reader) *Iterator[[]string] {
	return &Iterator[[]string]{
		seq: func(yield func([]string) bool) {
			for {
				record, err := r.Read()
				if err == io.EOF {
					return
				}
				if err != nil {
					// Skip malformed rows
					continue
				}
				if !yield(record) {
					return
				}
			}
		},
	}
}

// FromCSVWithContext creates a lazy Iterator from a CSV reader with context support.
// The iterator will stop when either the CSV is exhausted or the context is cancelled.
//
// Example:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	defer cancel()
//	file, _ := os.Open("huge_data.csv")
//	defer file.Close()
//	iter := itertools.FromCSVWithContext(ctx, csv.NewReader(file))
//	records := iter.Collect()
func FromCSVWithContext(ctx context.Context, r *csv.Reader) *Iterator[[]string] {
	return &Iterator[[]string]{
		seq: func(yield func([]string) bool) {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					record, err := r.Read()
					if err == io.EOF {
						return
					}
					if err != nil {
						// Skip malformed rows
						continue
					}
					if !yield(record) {
						return
					}
				}
			}
		},
	}
}

// CSVRow represents a parsed CSV row with helper methods
type CSVRow struct {
	Fields []string
	Index  int
}

// Get returns the field at the given index, or empty string if out of bounds
func (r CSVRow) Get(index int) string {
	if index >= 0 && index < len(r.Fields) {
		return r.Fields[index]
	}
	return ""
}

// GetByHeader returns the field with the given header name
// headers should be the first row of the CSV
func (r CSVRow) GetByHeader(headers []string, name string) string {
	for i, h := range headers {
		if h == name {
			return r.Get(i)
		}
	}
	return ""
}

// FromCSVWithHeaders creates a lazy Iterator that reads CSV records with header support.
// The first row is treated as headers and subsequent rows are wrapped in CSVRow for easier access.
// Returns the iterator and the header row.
//
// Example:
//
//	file, _ := os.Open("data.csv")
//	defer file.Close()
//	iter, headers := itertools.FromCSVWithHeaders(csv.NewReader(file))
//	records := iter.Filter(func(row CSVRow) bool {
//	    age := row.GetByHeader(headers, "age")
//	    return age != "" && age > "30"
//	}).Collect()
func FromCSVWithHeaders(r *csv.Reader) (*Iterator[CSVRow], []string, error) {
	// Read header row
	headers, err := r.Read()
	if err != nil {
		return nil, nil, err
	}

	index := 0
	iter := &Iterator[CSVRow]{
		seq: func(yield func(CSVRow) bool) {
			for {
				record, err := r.Read()
				if err == io.EOF {
					return
				}
				if err != nil {
					continue
				}
				if !yield(CSVRow{Fields: record, Index: index}) {
					return
				}
				index++
			}
		},
	}

	return iter, headers, nil
}

// FromCSVWithHeadersContext creates a lazy Iterator from CSV with headers and context support.
//
// Example:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
//	defer cancel()
//	file, _ := os.Open("large.csv")
//	defer file.Close()
//	iter, headers, _ := itertools.FromCSVWithHeadersContext(ctx, csv.NewReader(file))
func FromCSVWithHeadersContext(ctx context.Context, r *csv.Reader) (*Iterator[CSVRow], []string, error) {
	// Read header row
	headers, err := r.Read()
	if err != nil {
		return nil, nil, err
	}

	index := 0
	iter := &Iterator[CSVRow]{
		seq: func(yield func(CSVRow) bool) {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					record, err := r.Read()
					if err == io.EOF {
						return
					}
					if err != nil {
						continue
					}
					if !yield(CSVRow{Fields: record, Index: index}) {
						return
					}
					index++
				}
			}
		},
	}

	return iter, headers, nil
}
