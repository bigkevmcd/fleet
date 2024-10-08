package helmdeployer

import (
	"bytes"
	"fmt"
	"strconv"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage/driver"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/rancher/wrangler/v3/pkg/kv"
	"github.com/rancher/wrangler/v3/pkg/yaml"
)

func (h *Helm) EnsureInstalled(bundleID, resourcesID string) (bool, error) {
	releaseName, version, namespace, err := getReleaseNameVersionAndNamespace(bundleID, resourcesID)
	if err != nil {
		return false, err
	}

	if _, err := h.getRelease(releaseName, namespace, version); err == ErrNoRelease {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// Resources returns the resources from the helm release history
func (h *Helm) Resources(bundleID, resourcesID string) (*Resources, error) {
	releaseName, version, namespace, err := getReleaseNameVersionAndNamespace(bundleID, resourcesID)
	if err != nil {
		return &Resources{}, err
	}

	release, err := h.getRelease(releaseName, namespace, version)
	if err == ErrNoRelease {
		return &Resources{}, nil
	} else if err != nil {
		return nil, err
	}

	resources := &Resources{DefaultNamespace: release.Namespace}
	resources.Objects, err = ReleaseToObjects(release)
	return resources, err
}

func (h *Helm) ResourcesFromPreviousReleaseVersion(bundleID, resourcesID string) (*Resources, error) {
	releaseName, version, namespace, err := getReleaseNameVersionAndNamespace(bundleID, resourcesID)
	if err != nil {
		return &Resources{}, err
	}

	release, err := h.getRelease(releaseName, namespace, version-1)
	if err == ErrNoRelease {
		return &Resources{}, nil
	} else if err != nil {
		return nil, err
	}

	resources := &Resources{DefaultNamespace: release.Namespace}
	resources.Objects, err = ReleaseToObjects(release)
	return resources, err
}

func getReleaseNameVersionAndNamespace(bundleID, resourcesID string) (string, int, string, error) {
	// When a bundle is installed a resourcesID is generated. If there is no
	// resourcesID then there isn't anything to lookup.
	if resourcesID == "" {
		return "", 0, "", ErrNoResourceID
	}
	namespace, name := kv.Split(resourcesID, "/")
	releaseName, versionStr := kv.Split(name, ":")
	version, _ := strconv.Atoi(versionStr)

	if releaseName == "" {
		releaseName = bundleID
	}

	return releaseName, version, namespace, nil
}

func (h *Helm) getRelease(releaseName, namespace string, version int) (*release.Release, error) {
	hist := action.NewHistory(&h.globalCfg)

	releases, err := hist.Run(releaseName)
	if err == driver.ErrReleaseNotFound {
		return nil, ErrNoRelease
	} else if err != nil {
		return nil, err
	}

	for _, release := range releases {
		if release.Name == releaseName && release.Version == version && release.Namespace == namespace {
			return release, nil
		}
	}

	return nil, ErrNoRelease
}

func ReleaseToResourceID(release *release.Release) string {
	return fmt.Sprintf("%s/%s:%d", release.Namespace, release.Name, release.Version)
}

func ReleaseToObjects(release *release.Release) ([]runtime.Object, error) {
	var err error

	objs, err := yaml.ToObjects(bytes.NewBufferString(release.Manifest))
	return objs, err
}
