package validator

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		checkErr func(err error) bool
	}{
		{
			name: "invalid struct: interface",
			args: args{
				v: new(any),
			},
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrNotStruct)
			},
		},
		{
			name: "invalid struct: map",
			args: args{
				v: map[string]string{},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrNotStruct)
			},
		},
		{
			name: "invalid struct: string",
			args: args{
				v: "some string",
			},
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrNotStruct)
			},
		},
		{
			name: "valid struct with no fields",
			args: args{
				v: struct{}{},
			},
			wantErr: false,
		},
		{
			name: "valid struct with untagged fields",
			args: args{
				v: struct {
					f1 string
					f2 string
				}{},
			},
			wantErr: false,
		},
		{
			name: "valid struct with unexported fields",
			args: args{
				v: struct {
					foo string `validate:"len:10"`
				}{},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				e := &ValidationErrors{}
				return errors.As(err, e) && e.Error() == ErrValidateForUnexportedFields.Error()
			},
		},
		{
			name: "invalid validator syntax",
			args: args{
				v: struct {
					Foo string `validate:"len:abcdef"`
				}{},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				e := &ValidationErrors{}
				return errors.As(err, e) && e.Error() == ErrInvalidValidatorSyntax.Error()
			},
		},
		{
			name: "valid struct with tagged fields",
			args: args{
				v: struct {
					Len       string `validate:"len:20"`
					LenZ      string `validate:"len:0"`
					InInt     int    `validate:"in:20,25,30"`
					InNeg     int    `validate:"in:-20,-25,-30"`
					InStr     string `validate:"in:foo,bar"`
					MinInt    int    `validate:"min:10"`
					MinIntNeg int    `validate:"min:-10"`
					MinStr    string `validate:"min:10"`
					MinStrNeg string `validate:"min:-1"`
					MaxInt    int    `validate:"max:20"`
					MaxIntNeg int    `validate:"max:-2"`
					MaxStr    string `validate:"max:20"`
				}{
					Len:       "abcdefghjklmopqrstvu",
					LenZ:      "",
					InInt:     25,
					InNeg:     -25,
					InStr:     "bar",
					MinInt:    15,
					MinIntNeg: -9,
					MinStr:    "abcdefghjkl",
					MinStrNeg: "abc",
					MaxInt:    16,
					MaxIntNeg: -3,
					MaxStr:    "abcdefghjklmopqrst",
				},
			},
			wantErr: false,
		},
		{
			name: "Int slices ok",
			args: args{
				v: struct {
					InInt     []int `validate:"in:20,25,30"`
					InNeg     []int `validate:"in:-20,-25,-30"`
					MinInt    []int `validate:"min:10"`
					MinIntNeg []int `validate:"min:-10"`
					MaxInt    []int `validate:"max:20"`
					MaxIntNeg []int `validate:"max:-2"`
				}{
					InInt:     []int{20, 25, 30},
					InNeg:     []int{-20, -25, -30},
					MinInt:    []int{10, 100, 1000},
					MinIntNeg: []int{-10, -1, 10},
					MaxInt:    []int{20, 10, 2},
					MaxIntNeg: []int{-10, -2},
				},
			},
			wantErr: false,
		},
		{
			name: "Str slices ok",
			args: args{
				v: struct {
					Len       []string `validate:"len:20"`
					LenZ      []string `validate:"len:0"`
					InStr     []string `validate:"in:foo,bar"`
					MinStr    []string `validate:"min:10"`
					MinStrNeg []string `validate:"min:-1"`
					MaxStr    []string `validate:"max:20"`
				}{
					Len:       []string{"abcdefghjklmopqrstvu"},
					LenZ:      []string{"", "", ""},
					InStr:     []string{"bar", "bar"},
					MinStr:    []string{"abcdefghjkl", "abcdefghjk", "abcdefghjkabcdefghjk"},
					MinStrNeg: []string{"abc", "", " "},
					MaxStr:    []string{"abcdefghjklmopqrst", "", " ", "abcdefghjklmopqrstmm"},
				},
			},
			wantErr: false,
		},
		{
			name: "wrong length",
			args: args{
				v: struct {
					Lower    string `validate:"len:24"`
					Higher   string `validate:"len:5"`
					Zero     string `validate:"len:3"`
					BadSpec  string `validate:"len:%12"`
					Negative string `validate:"len:-6"`
				}{
					Lower:    "abcdef",
					Higher:   "abcdef",
					Zero:     "",
					BadSpec:  "abc",
					Negative: "abcd",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 5)
				return true
			},
		},
		{
			name: "wrong len on slices",
			args: args{
				v: struct {
					Lower    []string `validate:"len:24"`
					Higher   []string `validate:"len:5"`
					Zero     []string `validate:"len:3"`
					BadSpec  []string `validate:"len:%12"`
					Negative []string `validate:"len:-6"`
				}{
					Lower:    []string{"abcdef", ""},                          // + 2
					Higher:   []string{"abcdef", "                         "}, // + 2
					Zero:     []string{""},                                    // + 1
					BadSpec:  []string{"abc", "sdv", "", "+1"},                // + 1
					Negative: []string{"abc", "sdv", "", "+1"},                // + 1
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 7)
				return true
			},
		},
		{
			name: "wrong in",
			args: args{
				v: struct {
					InA     string `validate:"in:ab,cd"`
					InB     string `validate:"in:aa,bb,cd,ee"`
					InC     int    `validate:"in:-1,-3,5,7"`
					InD     int    `validate:"in:5-"`
					InEmpty string `validate:"in:"`
				}{
					InA:     "ef",
					InB:     "ab",
					InC:     2,
					InD:     12,
					InEmpty: "",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 5)
				return true
			},
		},
		{
			name: "wrong in on slices",
			args: args{
				v: struct {
					InA     []string `validate:"in:ab,cd"`
					InB     []string `validate:"in:aa,bb,cd,ee"`
					InC     []int    `validate:"in:-1,-3,5,7"`
					InD     []int    `validate:"in:5-"`
					InEmpty []string `validate:"in:"`
				}{
					InA:     []string{"ef", ""},       // + 1
					InB:     []string{"ab", "", "ba"}, // + 2
					InC:     []int{2, -1, 4, 8},       // + 3
					InD:     []int{12, 15, 56},        // + 1
					InEmpty: []string{"", "asdafd"},   // + 1
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 8)
				return true
			},
		},
		{
			name: "wrong min",
			args: args{
				v: struct {
					MinA string `validate:"min:12"`
					MinB int    `validate:"min:-12"`
					MinC int    `validate:"min:5-"`
					MinD int    `validate:"min:"`
					MinE string `validate:"min:"`
				}{
					MinA: "ef",
					MinB: -22,
					MinC: 12,
					MinD: 11,
					MinE: "abc",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 5)
				return true
			},
		},
		{
			name: "wrong min on slices",
			args: args{
				v: struct {
					MinA []string `validate:"min:10"`
					MinB []int    `validate:"min:-10"`
					MinC []int    `validate:"min:5-"`
					MinD []int    `validate:"min:"`
					MinE []string `validate:"min:"`
				}{
					MinA: []string{"ef", "", "abcdeabcde"}, // + 2
					MinB: []int{-22, -11, -10, 10},         // + 2
					MinC: []int{12, 0, -12},                // + 1
					MinD: []int{11, 0, -11},                // + 1
					MinE: []string{"abc", "", ","},         // + 1
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 7)
				return true
			},
		},
		{
			name: "wrong max",
			args: args{
				v: struct {
					MaxA string `validate:"max:2"`
					MaxB string `validate:"max:-7"`
					MaxC int    `validate:"max:-12"`
					MaxD int    `validate:"max:5-"`
					MaxE int    `validate:"max:"`
					MaxF string `validate:"max:"`
				}{
					MaxA: "efgh",
					MaxB: "ab",
					MaxC: 22,
					MaxD: 12,
					MaxE: 11,
					MaxF: "abc",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 6)
				return true
			},
		},
		{
			name: "wrong max on slices",
			args: args{
				v: struct {
					MaxA []string `validate:"max:2"`
					MaxB []string `validate:"max:-7"`
					MaxC []int    `validate:"max:-12"`
					MaxD []int    `validate:"max:5-"`
					MaxE []int    `validate:"max:"`
					MaxF []string `validate:"max:"`
				}{
					MaxA: []string{"efgh", "             "}, // + 2
					MaxB: []string{"ab", "             "},   // + 2
					MaxC: []int{22, 20, -11},                // + 3
					MaxD: []int{12, 456, 0, -111},           // + 1
					MaxE: []int{12, 456, 0, -111},           // + 1
					MaxF: []string{"abc", "", ",", "dsav"},  // + 1
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 10)
				return true
			},
		},
		{
			name: "test multiply tags string: max + min + in",
			args: args{
				v: struct {
					MaxMinIn []string `validate:"max:5 min:1 in:a,aa,bbb,bbbb,bbabb"`
				}{
					MaxMinIn: []string{"a", "aa", "bbb", "bbabb"},
				},
			},
			wantErr: false,
		},
		{
			name: "test multiply tags string: len + in",
			args: args{
				v: struct {
					LenIn []string `validate:"len:2 in:aa,bb"`
				}{
					LenIn: []string{"aa", "bb"},
				},
			},
			wantErr: false,
		},
		{
			name: "test multiply tags string err: wrong",
			args: args{
				v: struct {
					StrA []string `validate:"len:5 in:aa,bb"`
				}{
					StrA: []string{"aa", "bbbbb"},
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 2)
				return true
			},
		},
		{
			name: "test multiply tags int: max + min + in",
			args: args{
				v: struct {
					MaxMinIn []int `validate:"max:50 min:1 in:3,25"`
				}{
					MaxMinIn: []int{25, 3},
				},
			},
			wantErr: false,
		},

		{
			name: "test multiply tags err int: wrong",
			args: args{
				v: struct {
					MaxMin []int `validate:"max:35 min:25"`
				}{
					MaxMin: []int{36, 12},
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 2)
				return true
			},
		},
		{
			name: "Invalid tags",
			args: args{
				v: struct {
					Int1 int    `validate:"len:2"`
					Str1 string `validate:"mini:2"`
				}{
					Int1: 36,
					Str1: "as",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 2)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.args.v)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, tt.checkErr(err), "test expect an error, but got wrong error type")
			} else {
				assert.NoError(t, err)
			}
		})
	}

}
