package mta

import (
	"path/filepath"

	"cloud-mta-build-tool/internal/fsys"
	"github.com/pkg/errors"
)

// MtaLocationParameters - MTA tool file properties
type MtaLocationParameters struct {
	// SourcePath - Path to MTA project
	SourcePath string
	// TargetPath - Path to MTA tool results
	TargetPath string
	// MtaFilename - MTA yaml filename "mta.yaml" by default
	MtaFilename string
	// IsDeploymentDescriptor - indicator of deployment descriptor usage (mtad.yaml)
	Descriptor string
}

const (
	dep = "dep"
)

// GetSource - Get Processed Project Path
// If not provided use current directory
func (ep *MtaLocationParameters) GetSource() (string, error) {
	if ep.SourcePath == "" {
		wd, err := dir.GetWorkingDirectory()
		if err != nil {
			return "", errors.Wrap(err, "GetSource failed")
		}
		return wd, nil
	}
	return ep.SourcePath, nil
}

// GetTarget - Get Target Path
// If not provided use path of processed project
func (ep *MtaLocationParameters) GetTarget() (string, error) {
	if ep.TargetPath == "" {
		source, err := ep.GetSource()
		if err != nil {
			return "", errors.Wrap(err, "GetTarget failed")
		}
		return source, nil
	}
	return ep.TargetPath, nil

}

// GetTargetTmpDir - Get Target Temporary Directory path
// Subdirectory in target folder named as source project folder
func (ep *MtaLocationParameters) GetTargetTmpDir() (string, error) {
	source, err := ep.GetSource()
	if err != nil {
		return "", errors.Wrap(err, "GetTargetTmpDir failed")
	}
	_, file := filepath.Split(source)
	target, err := ep.GetTarget()
	if err != nil {
		return "", errors.Wrap(err, "GetTargetTmpDir failed")
	}
	// append to the currentPath the file name
	return filepath.Join(target, file), nil
}

// GetTargetModuleDir - Get path to the packed module directory
// Subdirectory in Target Temporary Directory named by module name
func (ep *MtaLocationParameters) GetTargetModuleDir(moduleName string) (string, error) {
	dir, err := ep.GetTargetTmpDir()
	if err != nil {
		return "", errors.Wrap(err, "GetTargetModuleDir failed")
	}

	return filepath.Join(dir, moduleName), nil
}

// GetTargetModuleZipPath - Get path to the packed module data.zip
// Subdirectory in Target Temporary Directory named by module name
func (ep *MtaLocationParameters) GetTargetModuleZipPath(moduleName string) (string, error) {
	dir, err := ep.GetTargetModuleDir(moduleName)
	if err != nil {
		return "", errors.Wrap(err, "GetTargetModuleZipPath failed")
	}
	return filepath.Join(dir, "data.zip"), nil
}

// GetSourceModuleDir - Get path to module to be packed
// Subdirectory in Source
func (ep *MtaLocationParameters) GetSourceModuleDir(modulePath string) (string, error) {
	source, err := ep.GetSource()
	if err != nil {
		return "", errors.Wrap(err, "GetSourceModuleDir failed")
	}
	return filepath.Join(source, filepath.Clean(modulePath)), nil
}

// getMtaYamlFilename - Get MTA yaml File name
func (ep *MtaLocationParameters) getMtaYamlFilename() string {
	if ep.MtaFilename == "" {
		if ep.Descriptor == dep {
			return "mtad.yaml"
		}
		return "mta.yaml"
	}
	return ep.MtaFilename
}

// GetMtaYamlPath - Get MTA yaml File path
func (ep *MtaLocationParameters) GetMtaYamlPath() (string, error) {
	source, err := ep.GetSource()
	if err != nil {
		return "", errors.Wrap(err, "GetMtaYamlPath failed")
	}
	return filepath.Join(source, ep.getMtaYamlFilename()), nil
}

// GetMetaPath - Get path to generated META-INF directory
func (ep *MtaLocationParameters) GetMetaPath() (string, error) {
	dir, err := ep.GetTargetTmpDir()
	if err != nil {
		return "", errors.Wrap(err, "GetMetaPath failed")
	}
	return filepath.Join(dir, "META-INF"), nil
}

// GetMtadPath - Get path to generated MTAD file
func (ep *MtaLocationParameters) GetMtadPath() (string, error) {
	dir, err := ep.GetMetaPath()
	if err != nil {
		return "", errors.Wrap(err, "GetMtadPath failed")
	}
	return filepath.Join(dir, "mtad.yaml"), nil
}

// GetManifestPath - Get path to generated manifest file
func (ep *MtaLocationParameters) GetManifestPath() (string, error) {
	dir, err := ep.GetMetaPath()
	if err != nil {
		return "", errors.Wrap(err, "GetManifestPath failed")
	}
	return filepath.Join(dir, "MANIFEST.MF"), nil
}

// ValidateDeploymentDescriptor - Validates Deployment Descriptor
func ValidateDeploymentDescriptor(descriptor string) error {
	if descriptor != "" && descriptor != "dev" && descriptor != "dep" {
		return errors.New("Wrong descriptor value. Expected one of [dev, dep]. Default is dev")
	}
	return nil
}

// IsDeploymentDescriptor - Check if flag is related to deployment descriptor
func (ep *MtaLocationParameters) IsDeploymentDescriptor() bool {
	return ep.Descriptor == dep
}