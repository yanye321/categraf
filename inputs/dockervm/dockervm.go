package dockervm

import (
	"time"

	"flashcat.cloud/categraf/config"
	"flashcat.cloud/categraf/inputs"
	"flashcat.cloud/categraf/types"
)

const inputName = "dockervm"
const cpuSetUpdateInterval = time.Second * time.Duration(300)

func init() {
	inputs.Add(inputName, func() inputs.Input {
		return &DockerVmStat{savedCpuStat: map[string]int64{}}
	})
}

func (dv *DockerVmStat) Clone() inputs.Input {
	return &DockerVmStat{}
}

func (dv *DockerVmStat) Name() string {
	return inputName
}

type DockerVmStat struct {
	config.PluginConfig
	cpuSet           []int
	savedCpuStat     map[string]int64
	cpuSetUpdateTime time.Time
}

func (dv *DockerVmStat) Gather(slist *types.SampleList) {
	sysinfo := getSysinfo()
	if !sysinfo.isContainer {
		return
	}
	if len(dv.cpuSet) == 0 || time.Since(dv.cpuSetUpdateTime) > cpuSetUpdateInterval {
		dv.cpuSet = getCPUSet()
		dv.cpuSetUpdateTime = time.Now()
	}

	memStat := getMemoryStat()
	cpuStat := getCPUStat(dv.cpuSet)
	fields := map[string]interface{}{}

	for k, v := range memStat {
		fields[k] = v
	}

	if memStat["mem_total"] >= 0 {
		fields["mem_available_percent"] = 100 * float64(memStat["mem_free"]) / float64(memStat["mem_total"])
	} else {
		fields["mem_available_percent"] = 100
	}

	if len(dv.savedCpuStat) > 0 {
		diffIdle := cpuStat["cpu_idle"] - dv.savedCpuStat["cpu_idle"]
		diffSum := cpuStat["cpu_sum"] - dv.savedCpuStat["cpu_sum"]
		fields["cpu_idle_percent"] = 100 * float64(diffIdle) / float64(diffSum)
	}

	dv.savedCpuStat = cpuStat
	slist.PushSamples(inputName, fields)
}
