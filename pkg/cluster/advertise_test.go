package cluster

import (
	"net"
	"testing"

	"github.com/SimonRichardson/alchemy/pkg/cluster/mocks"
	"github.com/golang/mock/gomock"

	"github.com/go-kit/kit/log"
)

func TestCalculateAdvertiseAddr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type expectation struct {
		required bool
		host     string
		addrs    []net.IPAddr
	}

	for _, testcase := range []struct {
		name          string
		bindAddr      string
		advertiseAddr string
		expect        expectation
		want          string
	}{
		{"Public bind no advertise",
			"1.2.3.4", "", expectation{}, "1.2.3.4",
		},
		{"Private bind no advertise",
			"10.1.2.3", "", expectation{}, "10.1.2.3",
		},
		{"Zeroes bind public advertise",
			"0.0.0.0", "2.3.4.5", expectation{}, "2.3.4.5",
		},
		{"Zeroes bind private advertise",
			"0.0.0.0", "172.16.1.9", expectation{}, "172.16.1.9",
		},
		{"Public bind private advertise",
			"188.177.166.155", "10.11.12.13", expectation{}, "10.11.12.13",
		},
		{"IPv6 bind no advertise",
			"::", "", expectation{}, "::",
		},
		{"IPv6 bind private advertise",
			"::", "172.16.1.1", expectation{}, "172.16.1.1",
		},
		{"Valid hostname as bind addr",
			"validhost.com", "", expectation{
				required: true,
				host:     "validhost.com",
				addrs:    []net.IPAddr{{IP: net.ParseIP("10.21.32.43")}},
			}, "10.21.32.43",
		},
		{"Valid hostname as advertise addr",
			"0.0.0.0", "validhost.com", expectation{
				required: true,
				host:     "validhost.com",
				addrs:    []net.IPAddr{{IP: net.ParseIP("10.21.32.43")}},
			}, "10.21.32.43",
		},
		{"Valid multi-hostname as bind addr",
			"multihost.com", "", expectation{
				required: true,
				host:     "multihost.com",
				addrs:    []net.IPAddr{{IP: net.ParseIP("10.1.0.1")}, {IP: net.ParseIP("10.1.0.2")}},
			}, "10.1.0.1",
		},
		{"Valid multi-hostname as advertise addr",
			"0.0.0.0", "multihost.com", expectation{
				required: true,
				host:     "multihost.com",
				addrs:    []net.IPAddr{{IP: net.ParseIP("10.1.0.1")}, {IP: net.ParseIP("10.1.0.2")}},
			}, "10.1.0.1",
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			r := mocks.NewMockResolver(ctrl)

			if expect := testcase.expect; expect.required {
				r.EXPECT().LookupIPAddr(gomock.Any(), expect.host).Return(expect.addrs, nil)
			}

			ip, err := CalculateAdvertiseIP(testcase.bindAddr, testcase.advertiseAddr, r, log.NewNopLogger())
			if err != nil {
				t.Fatal(err)
			}
			if want, have := testcase.want, ip.String(); want != have {
				t.Fatalf("want '%s', have '%s'", want, have)
			}
		})
	}
}
