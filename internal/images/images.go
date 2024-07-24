package images

import (
	"context"
	"fmt"
	"strings"

	"github.com/containers/image/v5/manifest"
	"github.com/containers/image/v5/transports/alltransports"
	"github.com/containers/image/v5/types"
)

// ToFullURL converts the provided image name to the full URL, if already not in that format.
// e.g. nginx becomes docker.io/library/nginx
// If already in full URL format, image name is returned as is.
func ToFullURL(imgName string) string {
	splits := strings.Split(imgName, "/")
	if len(splits) > 2 {
		return imgName
	}
	if len(splits) == 2 {
		return fmt.Sprintf("docker.io/%s", imgName)
	}
	return fmt.Sprintf("docker.io/library/%s", imgName)
}

// CheckLinuxArm64Support checks for the existance of an arm64 linux image in the manifest
func CheckLinuxArm64Support(imgName string) (bool, error) {
	sys := &types.SystemContext{
		ArchitectureChoice: "arm64",
		OSChoice:           "linux",
	}

	ref, err := alltransports.ParseImageName(fmt.Sprintf("docker://%s", imgName))
	if err != nil {
		return false, err
	}

	imageSource, err := ref.NewImageSource(context.Background(), sys)
	if err != nil {
		return false, err
	}
	defer imageSource.Close()

	rawManifest, mimeType, err := imageSource.GetManifest(context.Background(), nil)
	if err != nil {
		return false, err
	}

	manifestList, err := manifest.ListFromBlob(rawManifest, mimeType)
	if err != nil {
		return false, err
	}

	// This call will error if it cannot find a instance that supports linux arm64
	_, err = manifestList.ChooseInstance(sys)
	return err == nil, nil
}