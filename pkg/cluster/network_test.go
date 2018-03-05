package cluster

import "testing"

func TestParseAddr(t *testing.T) {
	for _, testcase := range []struct {
		addr        string
		defaultPort int
		network     string
		address     string
		host        string
		port        int
	}{
		{"foo", 123, "tcp", "foo:123", "foo", 123},
		{"foo:80", 123, "tcp", "foo:80", "foo", 80},
		{"udp://foo", 123, "udp", "foo:123", "foo", 123},
		{"udp://foo:8080", 123, "udp", "foo:8080", "foo", 8080},
		{"tcp+dnssrv://testing:7650", 7650, "tcp+dnssrv", "testing:7650", "testing", 7650},
	} {
		network, address, host, port, err := ParseAddr(testcase.addr, testcase.defaultPort)
		if err != nil {
			t.Errorf("(%q, %d): %v", testcase.addr, testcase.defaultPort, err)
			continue
		}
		var (
			matchNetwork = network == testcase.network
			matchAddress = address == testcase.address
		)
		if !matchNetwork || !matchAddress {
			t.Errorf("(%q, %d): want [%s %s %s %d], have [%s %s %s %d]",
				testcase.addr, testcase.defaultPort,
				testcase.network, testcase.address,
				testcase.host, testcase.port,
				network, address,
				host, port,
			)
			continue
		}
	}
}

func TestHasNonlocal(t *testing.T) {
	for _, testcase := range []struct {
		name  string
		input []string
		want  bool
	}{
		{
			"empty",
			[]string{},
			false,
		},
		{
			"127",
			[]string{"127.0.0.9"},
			false,
		},
		{
			"127 with port",
			[]string{"127.0.0.1:1234"},
			false,
		},
		{
			"nonlocal IP",
			[]string{"1.2.3.4"},
			true,
		},
		{
			"nonlocal IP with port",
			[]string{"1.2.3.4:5678"},
			true,
		},
		{
			"nonlocal host",
			[]string{"foo.corp"},
			true,
		},
		{
			"nonlocal host with port",
			[]string{"foo.corp:7659"},
			true,
		},
		{
			"localhost",
			[]string{"localhost"},
			false,
		},
		{
			"localhost with port",
			[]string{"localhost:1234"},
			false,
		},
		{
			"multiple IP",
			[]string{"127.0.0.1", "1.2.3.4"},
			true,
		},
		{
			"multiple hostname",
			[]string{"localhost", "otherhost"},
			true,
		},
		{
			"multiple local",
			[]string{"localhost", "127.0.0.1", "127.128.129.130:4321", "localhost:10001", "localhost:10002"},
			false,
		},
		{
			"multiple mixed",
			[]string{"localhost", "127.0.0.1", "129.128.129.130:4321", "localhost:10001", "localhost:10002"},
			true,
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			if want, have := testcase.want, HasNonlocal(testcase.input); want != have {
				t.Errorf("want %v, have %v", want, have)
			}
		})
	}
}

func TestIsUnRoutable(t *testing.T) {
	for _, testcase := range []struct {
		input string
		want  bool
	}{
		{"0.0.0.0", true},
		{"127.0.0.1", true},
		{"127.128.129.130", true},
		{"localhost", true},
		{"foo", false},
		{"::", true},
	} {
		t.Run(testcase.input, func(t *testing.T) {
			if want, have := testcase.want, IsUnRoutable(testcase.input); want != have {
				t.Errorf("want %v, have %v", want, have)
			}
		})
	}
}
