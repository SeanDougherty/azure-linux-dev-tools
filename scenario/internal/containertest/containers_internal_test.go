// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

//go:build scenario

//nolint:testpackage // Intentionally verifies internal container request options.
package containertest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainerImageOptions(t *testing.T) {
	t.Parallel()

	collateral := NewContainerTestCollateral(t).
		WithDockerfile("Dockerfile.image-builder").
		WithTmpfs("/workdir/build", "rw,exec,mode=1777")

	assert.Equal(t, "Dockerfile.image-builder", collateral.dockerfile)
	assert.Equal(t, map[string]string{"/workdir/build": "rw,exec,mode=1777"}, collateral.tmpfs)
	assert.Equal(t, "latest-image-builder", containerImageTag(collateral.dockerfile))
	assert.Equal(t, "latest", containerImageTag(defaultDockerfile))
}

func TestContainerImageOptionsRejectInvalidValues(t *testing.T) {
	t.Parallel()

	tests := map[string]func(){
		"empty Dockerfile": func() {
			NewContainerTestCollateral(t).WithDockerfile("")
		},
		"nested Dockerfile": func() {
			NewContainerTestCollateral(t).WithDockerfile("../Dockerfile")
		},
		"relative tmpfs": func() {
			NewContainerTestCollateral(t).WithTmpfs("workdir", "rw")
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Panics(t, test)
		})
	}
}
