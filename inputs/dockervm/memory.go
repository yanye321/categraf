package dockervm

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func getMemInfoDir() string {
	sysinfo := getSysinfo()
	memInfoDir := ""
	if _, err := os.Stat("/sys/fs/cgroup/memory/system.slice"); err == nil {
		memInfoDir = fmt.Sprintf("/sys/fs/cgroup/memory%s", sysinfo.systemSliceDir)
	} else {
		memInfoDir = "/sys/fs/cgroup/memory"
	}
	return memInfoDir
}

func getMemoryStat() map[string]int64 {
	memInfoDir := getMemInfoDir()
	memStat := make(map[string]int64)
	if _, err := os.Stat(filepath.Join(memInfoDir, "memory.limit_in_bytes")); err != nil {
		return memStat
	}
	memTotalStr, _ := getFileContent(filepath.Join(memInfoDir, "memory.limit_in_bytes"))
	memTotal, _ := strconv.ParseInt(strings.TrimSpace(memTotalStr), 10, 64)
	memUsedStr, _ := getFileContent(filepath.Join(memInfoDir, "memory.usage_in_bytes"))
	memUsed, _ := strconv.ParseInt(strings.TrimSpace(memUsedStr), 10, 64)

	memStat["mem_total"] = memTotal
	memStat["mem_used"] = memUsed

	memStatLines := getFileLines(filepath.Join(memInfoDir, "memory.stat"), true, "")
	for _, line := range memStatLines {
		fields := strings.Fields(line)
		if len(fields) != 2 {
			continue
		}
		if fields[0] == "total_rss" {
			rss, _ := strconv.ParseInt(fields[1], 10, 64)
			memStat["mem_used_rss"] = rss
			memStat["mem_free"] = memTotal - memUsed
			continue
		}
		if fields[0] == "total_cache" {
			cache, _ := strconv.ParseInt(fields[1], 10, 64)
			memStat["mem_cache"] = cache
		}
	}
	return memStat
}
