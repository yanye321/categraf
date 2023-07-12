package dockervm

import (
	"path/filepath"
	"regexp"
	"strings"
)

type sysInfo struct {
	isContainer    bool
	dockerID       string
	systemSliceDir string
}

var _sysinfo *sysInfo

func getSysinfo() *sysInfo {
	if _sysinfo == nil {
		_sysinfo = new(sysInfo)
		_sysinfo.checkInContainer()
	}
	return _sysinfo
}

func (s *sysInfo) checkInContainer() bool {
	data, err := getFileContent("/proc/1/cgroup")
	if err != nil {
		s.isContainer = false
		return s.isContainer
	}
	if matched, _ := regexp.MatchString(`/docker(/|-[0-9a-f]+\.scope)`, data); matched {
		s.isContainer = true
	} else {
		s.isContainer = false
	}
	if s.isContainer {
		s.setDockerID(data)
	}
	return s.isContainer
}

func (s *sysInfo) setDockerID(data string) string {
	systemSliceDir := ""
	dockerID := ""
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		if strings.Contains(line, ":/system.slice/docker-") {
			fields := strings.SplitN(line, ":", 3)
			systemSliceDir = fields[2]
			dockerID = filepath.Base(systemSliceDir)
			dockerID = strings.TrimSuffix(dockerID, ".scope")
			break
		}
	}
	if dockerID != "" {
		s.dockerID = dockerID
		s.systemSliceDir = systemSliceDir
	}
	return dockerID
}
