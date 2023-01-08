package dao

import "testing"

func TestLoginAuth(t *testing.T) {
	var tests = []struct {
		userName, password string
		convey             bool
	}{
		{"1", "1", true},
		{"2", "2", true},
		{"-1", "-1", false},
		{"", "", false},
		{"1", "", false},
		{"", "1", false},
	}
	for _, test := range tests {
		if ok, err := LoginAuth(test.userName, test.password); err != nil || ok != test.convey {
			t.Errorf("LoginAuth failed. userName:%s, password:%s, convey:%t", test.userName, test.password, test.convey)
		}
	}
}
