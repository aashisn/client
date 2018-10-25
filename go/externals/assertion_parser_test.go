// Copyright 2015 Keybase, Inc. All rights reserved. Use of
// this source code is governed by the included BSD license.

package externals

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalization(t *testing.T) {
	tc := setupTest(t, "Normalization", 1)
	defer tc.Cleanup()
	inp := "Web://A.AA || HttP://B.bb && dnS://C.cc || MaxFactor@reddit || zQueal@keyBASE || XanxA@hackernews || foO@TWITTER || 0123456789ABCDEF0123456789abcd19@uid || josh@gubble.SoCiAl"
	outp := "a.aa@web,b.bb@http+c.cc@dns,maxfactor@reddit,zqueal,XanxA@hackernews,foo@twitter,0123456789abcdef0123456789abcd19@uid,josh@gubble.social"
	expr, err := AssertionParse(tc.G, inp)
	require.NoError(t, err)
	require.Equal(t, expr.String(), outp)
}

type Pair struct {
	k, v string
}

func TestParserFail1(t *testing.T) {
	tc := setupTest(t, "ParserFail1", 1)
	defer tc.Cleanup()
	bads := []Pair{
		{"aa ||", "Unexpected EOF"},
		{"aa &&", "Unexpected EOF"},
		{"(aa", "Unbalanced parentheses"},
		{"aa && dns:", "Bad assertion, no value given (key=dns)"},
		{"bb && foo:a", "Unknown social network: foo"},
		{"&& aa", "Unexpected token: &&"},
		{"|| aa", "Unexpected token: ||"},
		{"aa)", "Found junk at end of input: )"},
		{"()", "Illegal parenthetical expression"},
		{"dns://a", "Invalid hostname: a"},
		{"f@reddit", "Bad username: 'f'"},
		{"a@pgp", "bad hex string: 'a'"},
		{"aBCP@pgp", "bad hex string: 'abcp'"},
		{"jj@pgp", "bad hex string: 'jj'"},
		{"http://", "Bad assertion, no value given (key=http)"},
		{"reddit:", "Bad username: ''"},
		{"gubble.social:", "username must be at least 2 characters, was 0"},
		{"(alice@keybasers.de)@email", "Illegal parenthetical expression"},
		{"twitter://alice&&(alice@keybasers.de)@email", "Found junk at end of input: )"}, // excuse me, now, that's mr. junk for you
	}

	for _, bad := range bads {
		ret, err := AssertionParse(tc.G, bad.k)
		require.Error(t, err, "for %q: ret is: %+v", bad.k, ret)
		require.Equal(t, bad.v, err.Error())
	}
}
