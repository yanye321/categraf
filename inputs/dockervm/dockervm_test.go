package dockervm

import "testing"

func TestSysinfo(t *testing.T) {
	si := getSysinfo()
	t.Log(si.isContainer)
	t.Log(si.systemSliceDir)
	t.Log(si.dockerID)
}

func TestMemoryStat(t *testing.T) {
	memStat := getMemoryStat()
	t.Logf("%+v", memStat)
}

func TestParseCpuSet(t *testing.T) {
	cset := parseCPUSet("0-3")
	t.Logf("%+v", cset)
}

func TestCpuSet(t *testing.T) {
	cset := getCPUSet()
	t.Logf("%+v", cset)
}

func TestCpuStat(t *testing.T) {
	cset := getCPUSet()
	ci := getCPUStat(cset)
	t.Logf("%+v", ci)
}
