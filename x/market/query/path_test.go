package query

import (
	"fmt"
	"path"
	"strings"
	"testing"

	"github.com/ovrclk/akash/testutil"
	"github.com/ovrclk/akash/x/market/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var errWildcard = errors.New("wildcard error to assert the test should return an untyped error")

func mkPath(t *testing.T, stateStr string) string {
	return path.Join(
		testutil.AccAddress(t).String(),
		stateStr,
		testutil.AccAddress(t).String(),
	)
}

func TestBidPath(t *testing.T) {
	type testCase struct {
		desc     string
		path     string
		expState types.BidState
		expOk    bool
		expErr   error
	}

	tests := []testCase{
		{
			path:   mkPath(t, "open"),
			expOk:  true,
			expErr: nil,
		},
		{
			path:   mkPath(t, "closed"),
			expOk:  true,
			expErr: nil,
		},
		{
			path:   mkPath(t, "matched"),
			expOk:  true,
			expErr: nil,
		},
		{
			path:   mkPath(t, "neh"),
			expErr: ErrStateValue,
		},
		{
			path:   fmt.Sprintf("%s/%s", testutil.AccAddress(t).String(), "open"),
			expOk:  true,
			expErr: nil,
		},
		{
			path: fmt.Sprintf("%s/%s/%s",
				testutil.AccAddress(t).String(),
				"open",
				testutil.AccAddress(t).String()),
			expOk:  true,
			expErr: nil,
		},
		{
			desc: "invalid owner address",
			path: fmt.Sprintf("%s/%s/%s",
				"foo",
				"open",
				testutil.AccAddress(t).String()),
			expOk:  false,
			expErr: errWildcard,
		},
		{
			desc: "invalid provider address",
			path: fmt.Sprintf("%s/%s/%s",
				testutil.AccAddress(t).String(),
				"open",
				"foo"),
			expOk:  false,
			expErr: errWildcard,
		},
	}
	for _, test := range tests {
		tf := func(t *testing.T, test testCase) func(*testing.T) {
			return func(t *testing.T) {
				parts := strings.Split(test.path, "/")
				filters, ok, err := parseBidFiltersPath(parts)
				if test.expErr == errWildcard {
					assert.Error(t, err)
				} else {
					assert.Equal(t, err, test.expErr)
				}
				assert.Equal(t, ok, test.expOk)
				t.Logf("%#v", filters) // TODO: rm
			}
		}
		t.Run(test.desc, tf(t, test))
	}
}
