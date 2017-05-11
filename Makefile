# Copyright 2017 The OPA Authors.  All rights reserved.
# Use of this source code is governed by an Apache2
# license that can be found in the LICENSE file.

.PHONY: build
build:
	@./build/make-all.sh build

.PHONY: push
push:
	@./build/make-all.sh push
