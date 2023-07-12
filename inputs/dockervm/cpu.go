package dockervm

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func parseCPUSet(cpuSet string) []int {
	ret := make([]int, 0)
	if cpuSet == "" {
		return ret
	}
	cpuSet = strings.Trim(cpuSet, "\"' \n")
	cpuRanges := strings.Split(cpuSet, ",")
	for _, cpuRange := range cpuRanges {
		if strings.Contains(cpuRange, "-") {
			rangeBounds := strings.Split(cpuRange, "-")
			a, errA := strconv.Atoi(rangeBounds[0])
			b, errB := strconv.Atoi(rangeBounds[1])
			if errA == nil && errB == nil {
				for i := a; i <= b; i++ {
					ret = append(ret, i)
				}
			}
		} else {
			cpu, err := strconv.Atoi(cpuRange)
			if err == nil {
				ret = append(ret, cpu)
			}
		}
	}
	return ret
}

func getcpuInfoDir() string {
	sysinfo := getSysinfo()
	cpuInfoDir := ""
	if _, err := os.Stat("/sys/fs/cgroup/cpuset/system.slice"); err == nil {
		cpuInfoDir = fmt.Sprintf("/sys/fs/cgroup/cpuset%s", sysinfo.systemSliceDir)
	} else {
		cpuInfoDir = "/sys/fs/cgroup/cpuset"
	}
	return cpuInfoDir
}

func getCPUSet() []int {
	cpuInfoDir := getcpuInfoDir()
	if _, err := os.Stat(cpuInfoDir); err != nil {
		return []int{}
	}
	cpusetStr, _ := getFileContent(filepath.Join(cpuInfoDir, "cpuset.cpus"))
	cpuset := parseCPUSet(cpusetStr)
	return cpuset
}

/*
/proc/stat
cpu0 1132 34 1441 11311718 3675 127 438
cpu1 1123 0 849 11313845 2614 0 18
==============
user: normal processes executing in user mode
nice: niced processes executing in user mode
system: processes executing in kernel mode
idle: twiddling thumbs
iowait: waiting for I/O to complete
irq: servicing interrupts
softirq: servicing softirqs
*/

func getCPUStat(cpuset []int) map[string]int64 {
	cpuStat := make(map[string]int64)
	allCpuStat := make(map[int][]string)
	//ts := time.Now().Unix()
	statLines := getFileLines("/proc/stat", false, "")
	for _, line := range statLines {
		if strings.HasPrefix(line, "cpu") {
			fields := strings.Fields(line)
			cpuIndex, _ := strconv.Atoi(fields[0][3:])
			allCpuStat[cpuIndex] = fields[1:8]
		}
	}
	var cpuSum, cpuIdle int64 = 0, 0
	for _, idx := range cpuset {
		s, ok := allCpuStat[idx]
		if !ok {
			continue
		}
		for i := 0; i < len(s); i++ {
			cpuVal, _ := strconv.ParseInt(s[i], 10, 64)
			cpuSum += cpuVal
			if i == 3 {
				cpuIdle += cpuVal
			}
		}
	}
	cpuStat["cpu_idle"] = cpuIdle
	cpuStat["cpu_sum"] = cpuSum
	return cpuStat
}
