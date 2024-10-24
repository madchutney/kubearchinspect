/*
Copyright 2024 Arm Limited

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package images

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/containers/image/v5/image"
	"github.com/containers/image/v5/transports/alltransports"
	"github.com/containers/image/v5/types"
)

func getDockerConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".docker", "config.json")
}

// CheckLinuxArm64Support checks for the existance of an arm64 linux image in the manifest
func CheckLinuxArm64Support(imgName string) (bool, error) {
	sys := &types.SystemContext{
		ArchitectureChoice:       "arm64",
		OSChoice:                 "linux",
		DockerCompatAuthFilePath: getDockerConfigPath(),
	}

	ref, err := alltransports.ParseImageName(fmt.Sprintf("docker://%s", imgName))
	if err != nil {
		return false, fmt.Errorf("error parsing image name: %w", err)
	}

	src, err := ref.NewImageSource(context.Background(), sys)
	if err != nil {
		return false, fmt.Errorf("error getting image source: %w", err)
	}
	defer src.Close()

	img, err := image.FromUnparsedImage(context.TODO(), sys, image.UnparsedInstance(src, nil))
	if err != nil {
		return false, fmt.Errorf("error parsing manifest: %w", err)
	}

	imgInspect, err := img.Inspect(context.TODO())
	if err != nil {
		return false, fmt.Errorf("error inspecting image: %w", err)
	}

	return imgInspect.Architecture == "arm64", nil
}

func CheckLatestLinuxArm64Support(imgName string) (bool, error) {
	split := strings.Split(imgName, ":")
	if len(split) < 2 {
		return false, fmt.Errorf("invalid image name")
	}
	latestImageName := split[0] + ":latest"
	return CheckLinuxArm64Support(latestImageName)
}
