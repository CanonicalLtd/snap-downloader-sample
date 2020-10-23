package cache

import (
	"fmt"
	"github.com/CanonicalLtd/serial-vault/service/log"
	"os"
	"os/exec"
	"path"
)

const xdelta = "xdelta3"

// Delta gets a delta between two revisions of a snap
// The delta file name is in the format: snapname_arch_from_to.delta3
func (c *Cache) Delta(name, arch string, fromRevision, toRevision int) (string, error) {
	// check if we have the requested delta file
	nameDelta := fmt.Sprintf("%s_%s_%d_%d.delta3", name, arch, fromRevision, toRevision)
	fileDelta := path.Join(c.baseDir, name, nameDelta)
	if _, err := os.Stat(fileDelta); err == nil {
		return fileDelta, nil
	}

	// check that the from revision exists
	nameFrom := snapFilename(name, arch, fromRevision)
	fileFrom := path.Join(c.baseDir, name, nameFrom)
	if _, err := os.Stat(fileFrom); err != nil {
		return "", fmt.Errorf("the `from` revision does not exist")
	}

	// check the to revision exists
	nameTo := snapFilename(name, arch, toRevision)
	fileTo := path.Join(c.baseDir, name, nameTo)
	if _, err := os.Stat(fileTo); err != nil {
		return "", fmt.Errorf("the `to` revision does not exist")
	}

	// generate the delta
	if err := generateDelta(fileFrom, fileTo, fileDelta); err != nil {
		return "", err
	}

	return fileDelta, nil
}

func snapFilename(name, arch string, revision int) string {
	return fmt.Sprintf("%s_%d_%s.snap", name, revision, arch)
}

func generateDelta(fromFile, toFile, outFile string) error {
	out, err := exec.Command(xdelta, "-s", fromFile, toFile, outFile).CombinedOutput()
	log.Println(string(out))
	return err
}
