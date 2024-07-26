package args

import (
	"fmt"
	"reflect"
	"testing"
)

type Options struct {
	Name        string  `cli:"name"`
	Age         int     `cli:"age"`
	Verbose     bool    `cli:"verbose"`
	Rate        float64 `cli:"rate"`
	Description string  `cli:"description,default=This is a test"`
	Debug       bool    `cli:"debug,short=d"`
}

func TestParse(t *testing.T) {
	tests := []struct {
		args          []string
		expected      Options
		expectedError error
	}{
		{
			args: []string{"--name=John", "--age=30", "--verbose", "--rate=1.5", "--description=Test", "-d"},
			expected: Options{
				Name:        "John",
				Age:         30,
				Verbose:     true,
				Rate:        1.5,
				Description: "Test",
				Debug:       true,
			},
			expectedError: nil,
		},
		{
			args: []string{"-name=John", "-age=30", "-verbose", "-rate=1.5", "-description=Test", "-d"},
			expected: Options{
				Name:        "John",
				Age:         30,
				Verbose:     true,
				Rate:        1.5,
				Description: "Test",
				Debug:       true,
			},
			expectedError: nil,
		},
		{
			args: []string{"--name", "John", "--age", "30", "--verbose", "--rate", "1.5", "--description", "Test", "-d"},
			expected: Options{
				Name:        "John",
				Age:         30,
				Verbose:     true,
				Rate:        1.5,
				Description: "Test",
				Debug:       true,
			},
			expectedError: nil,
		},
		{
			args: []string{"-name", "John", "-age", "30", "-verbose", "-rate", "1.5", "-description", "Test", "-d"},
			expected: Options{
				Name:        "John",
				Age:         30,
				Verbose:     true,
				Rate:        1.5,
				Description: "Test",
				Debug:       true,
			},
			expectedError: nil,
		},
		{
			args: []string{"--name=Jane", "--age=25"},
			expected: Options{
				Name:        "Jane",
				Age:         25,
				Verbose:     false,
				Rate:        0,
				Description: "This is a test",
				Debug:       false,
			},
			expectedError: nil,
		},
		{
			args:          []string{"--unknown=unknown"},
			expected:      Options{},
			expectedError: fmt.Errorf("unknown option: %q", "unknown"),
		},
		{
			args:          []string{"-unknown", "unknown"},
			expected:      Options{},
			expectedError: fmt.Errorf("unknown option: %q", "unknown"),
		},
		{
			args:          []string{"--name"},
			expected:      Options{},
			expectedError: fmt.Errorf("value for the option %q is not set", "name"),
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			var opts Options
			err := Parse(tt.args, &opts)

			if !reflect.DeepEqual(err, tt.expectedError) {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}

			if !reflect.DeepEqual(opts, tt.expected) {
				t.Errorf("expected options %v, got %v", tt.expected, opts)
			}
		})
	}
}
