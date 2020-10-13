// Copyright 2020 Wayne wang<net_use@bzhy.com>.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file

package bzhyserver

import "net/http"

// CreateTestContext returns a fresh engine and context for testing purposes
func CreateTestContext(w http.ResponseWriter) (c *Context, r *Engine) {
	r = New()
	c = r.allocateContext()
	c.reset()
	c.writermem.reset(w)
	return
}
