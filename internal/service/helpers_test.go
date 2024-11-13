package service

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteCsvRecord(t *testing.T) {
	tests := []struct {
		name       string
		record     []string
		wantOutput string
		wantErr    bool
	}{
		{
			name:       "success - single record",
			record:     []string{"Alice", "Bob", "Charlie"},
			wantOutput: "Alice,Bob,Charlie\n",
			wantErr:    false,
		},
		{
			name:       "empty record",
			record:     []string{},
			wantOutput: "\n",
			wantErr:    false,
		},
		{
			name:       "special characters",
			record:     []string{"a,b", "c\"d", "e\nf"},
			wantOutput: "\"a,b\",\"c\"\"d\",\"e\nf\"\n",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := writeCsvRecord(&buf, tt.record)

			if (err != nil) != tt.wantErr {
				t.Errorf("writeCsvRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOutput := buf.String(); gotOutput != tt.wantOutput {
				t.Errorf("writeCsvRecord() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestParseCsv(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    [][]string
		wantErr bool
	}{
		{
			name:    "success - single record",
			input:   "Alice,Bob,Charlie\n",
			want:    [][]string{{"Alice", "Bob", "Charlie"}},
			wantErr: false,
		},
		{
			name:    "multiple records",
			input:   "Alice,Bob\nJohn,Doe\n",
			want:    [][]string{{"Alice", "Bob"}, {"John", "Doe"}},
			wantErr: false,
		},
		{
			name:    "empty record",
			input:   "\n",
			want:    nil,
			wantErr: false,
		},
		{
			name:    "invalid format",
			input:   "\"unclosed quote\n",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBufferString(tt.input)
			got, err := parseCsv(buf)

			if (err != nil) != tt.wantErr {
				t.Errorf("parseCsv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
