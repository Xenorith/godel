// Copyright 2016 Palantir Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package publish

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/termie/go-shutil"

	"github.com/palantir/godel/apps/distgo/config"
)

type LocalPublishInfo struct {
	Path string
}

func (l LocalPublishInfo) Publish(buildSpec config.ProductBuildSpec, paths ProductPaths, stdout io.Writer) (string, error) {
	productPath := path.Join(l.Path, paths.productPath)
	if err := os.MkdirAll(productPath, 0755); err != nil {
		return "", errors.Wrapf(err, "Failed to create path to %v", productPath)
	}

	if err := copyArtifact(paths.pomFilePath, productPath, stdout); err != nil {
		return "", errors.Wrapf(err, "Failed to copy POM file")
	}

	if err := copyArtifact(paths.artifactPath, productPath, stdout); err != nil {
		return "", errors.Wrapf(err, "Failed to copy artifact file")
	}

	return "", nil
}

func copyArtifact(src, dstDir string, stdout io.Writer) error {
	dst := path.Join(dstDir, path.Base(src))
	fmt.Fprintf(stdout, "Copying %v to %v...\n", src, dst)
	return shutil.CopyFile(src, dst, false)
}