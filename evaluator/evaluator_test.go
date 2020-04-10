package evaluator

import (
	"testing"

	"github.com/ziyoung/lox-go/parser"
	"github.com/ziyoung/lox-go/valuer"
)

func TestEvalNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"0", 0},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"0.01", 0.01},

		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},

		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"50 / (2 * 2) + 10", 22.5},
	}

	for i, test := range tests {
		v, err := testEvalExpr(test.input)
		if err != nil {
			t.Fatalf("test[%d] failed. error: %s", i, err.Error())
		}
		if !testNumberValuer(t, v, test.expected) {
			t.Fatalf("test[%d] failed. input is %s", i, test.input)
		}
	}
}

func TestEvalBoolean(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"!true", false},
		{"!false", true},
		{"!!true", true},

		{"true == true", true},
		{"false == false", true},
		{"false == true", false},
		{"1 == 1", true},
		{"1 != 1", false},

		{"1 >= 1", true},
		{"1 < 1", false},
		{`"" == ""`, true},
		{`"" == " "`, false},
		{`"" != " "`, true},

		{"nil == nil", true},
		{"nil != nil", false},
		{"nil == false", true},
		{"nil == true", false},
		{"0 == true", false},

		{"1 == true", true},
		{`"" == true`, false},
		{`"x" == true`, true},
	}

	for i, test := range tests {
		v, err := testEvalExpr(test.input)
		if err != nil {
			t.Fatalf("test[%d] failed. error: %s", i, err.Error())
		}
		testBooleanValuer(t, v, test.expected)
	}
}

func TestEvalString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`""`, ""},
		{`"x"`, "x"},
		{`"" + ""`, ""},
		{`"" + "x"`, "x"},
		{`"" + 1`, "1"},
		{`123 + "x"`, "123x"},
	}

	for i, test := range tests {
		v, err := testEvalExpr(test.input)
		if err != nil {
			t.Fatalf("test[%d] failed. error: %s", i, err.Error())
		}
		if !testStringValuer(t, v, test.expected) {
			t.Fatalf("test[%d] failed. input is %s", i, test.input)
		}
	}
}

func testEvalExpr(input string) (v valuer.Valuer, err error) {
	expr, err := parser.ParseExpr(input)
	if err != nil {
		panic(err)
	}
	defer func() {
		if r := recover(); r != nil {
			if runErr, ok := r.(runtimeError); ok {
				err = &runErr
				expr = nil
			} else {
				panic(r)
			}
		}
	}()
	return Eval(expr), nil
}

func testNumberValuer(t *testing.T, val valuer.Valuer, expected float64) bool {
	v, ok := val.(*valuer.Number)
	if !ok {
		t.Errorf("expected is Number. got %T (%+[1]v)", val)
		return false
	}
	if v.Value != expected {
		t.Errorf("expected value is %f. got %f", expected, v.Value)
		return false
	}
	return true
}

func testBooleanValuer(t *testing.T, val valuer.Valuer, expected bool) bool {
	v, ok := val.(*valuer.Boolean)
	if !ok {
		t.Errorf("expected is Boolean. got %T (%[1]v)", val)
		return false
	}
	if v.Value != expected {
		t.Errorf("expected value is %t. got %t", expected, v.Value)
		return false
	}
	return true
}

func testStringValuer(t *testing.T, val valuer.Valuer, expected string) bool {
	s, ok := val.(*valuer.String)
	if !ok {
		t.Errorf("expected is String. got %T (%[1]v)", val)
		return false
	}
	if s.Value != expected {
		t.Errorf("expected value is %s. got %s", expected, s.Value)
		return false
	}
	return true
}
