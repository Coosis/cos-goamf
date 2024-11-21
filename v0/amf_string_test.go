package v0

import (
	"testing"
)

func genStringTestcases() ([]string, [][]byte) {
	strings := []string{
		"hello",
		"world",
		"foo",
		"bar",
	}

	bytes := [][]byte{
		{AMF_STRING, 0, 5, 'h', 'e', 'l', 'l', 'o'},
		{AMF_STRING, 0, 5, 'w', 'o', 'r', 'l', 'd'},
		{AMF_STRING, 0, 3, 'f', 'o', 'o'},
		{AMF_STRING, 0, 3, 'b', 'a', 'r'},
	}

	return strings, bytes
}

func genLongStringTestcases() ([]string, [][]byte) {
	strings := []string{
		"hello",
		"world",
		"foo",
		"bar",
	}

	bytes := [][]byte{
		{AMF_LONGSTRING, 0, 0, 0, 5, 'h', 'e', 'l', 'l', 'o'},
		{AMF_LONGSTRING, 0, 0, 0, 5, 'w', 'o', 'r', 'l', 'd'},
		{AMF_LONGSTRING, 0, 0, 0, 3, 'f', 'o', 'o'},
		{AMF_LONGSTRING, 0, 0, 0, 3, 'b', 'a', 'r'},
	}

	return strings, bytes
}

func TestAmfStringEncode(t *testing.T) {
	testcases, expected := genStringTestcases()
	for i, value := range testcases {
		result := AmfStringEncode(value)
		if !bytesEqual(result, expected[i]) {
			t.Errorf("AmfStringEncode failed: expected %v, got %v", expected[i], result)
		}
	}
}

func TestAmfStringDecode(t *testing.T) {
	testcases, expected := genStringTestcases()
	for i, bytes := range expected {
		result, _, err := AmfStringDecode(bytes)
		if err != nil {
			t.Errorf("AmfStringDecode failed: %v", err)
		}
		if result != testcases[i] {
			t.Errorf("AmfStringDecode failed: expected %v, got %v", testcases[i], result)
		}
	}
}


func TestAmfLongStringEncode(t *testing.T) {
	testcases, expected := genLongStringTestcases()
	for i, value := range testcases {
		result := AmfLongStringEncode(value)
		if !bytesEqual(result, expected[i]) {
			t.Errorf("AmfLongStringEncode failed: expected %v, got %v", expected[i], result)
		}
	}
}

func TestAmfLongStringDecode(t *testing.T) {
	testcases, expected := genLongStringTestcases()
	for i, bytes := range expected {
		result, _, err := AmfLongStringDecode(bytes)
		if err != nil {
			t.Errorf("AmfLongStringDecode failed: %v", err)
		}
		if result != testcases[i] {
			t.Errorf("AmfLongStringDecode failed: expected %v, got %v", testcases[i], result)
		}
	}
}
