// Copyright 2020 Wayne wang<net_use@bzhy.com>.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file

package render

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReaderRenderNoHeaders(t *testing.T) {
	content := "test"
	r := Reader{
		ContentLength: int64(len(content)),
		Reader:        strings.NewReader(content),
	}
	err := r.Render(httptest.NewRecorder())
	require.NoError(t, err)
}
