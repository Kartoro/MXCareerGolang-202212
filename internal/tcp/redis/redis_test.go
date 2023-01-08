package redis

import "testing"

func TestSetToken(t *testing.T) {
	var tests = []struct {
		userName string
		token    string
		exp      int64
	}{
		{"1", "token", 5},
	}
	for _, test := range tests {
		if err := SetToken(test.userName, test.token, test.exp); err != nil {
			t.Errorf("SetToken not pass. userName:%s, token:%s, exp:%d, err:%q", test.userName, test.token, test.exp, err)
		}
	}
}

func TestCheckToken(t *testing.T) {
	var tests = []struct {
		userName string
		token    string
		convey   bool
	}{
		{"1", "token", true},
		{"2", "token", false},
	}
	_ = SetToken("1", "token", 5)
	for _, test := range tests {
		if ok, err := CheckToken(test.userName, test.token); ok != test.convey {
			t.Errorf("CheckToken no pass. userName:%s, token:%s, convey:%t, err:%q", test.userName, test.token, test.convey, err)
		}
	}
}
