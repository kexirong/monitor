package plugin

import (
	"time"

	"github.com/shirou/gopsutil/cpu"
)

func CpuStat() []float64 {

	ret, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil
	}
	return ret

}


func CpusStat() []cpu.TimesStat {
	times1, err := cpu.Times(true)
	if err !=nil{
		return nil
	} 
	time.Sleep(time.Second)
	times2, err := cpu.Times(true)
	if len(times1) != len(times2) {
		return nil
	}
	/*
	for i, t1 := range times1 {
		t2 := times2[i]

		all := t2.User - t1.User + t2.System - t1.System + t2.Nice - t1.Nice + t2.Iowait - t1.Iowait +
			t2.Irq - t1.Irq + t2.Softirq - t1.Softirq + t2.Steal - t1.Steal + t2.Guest - t1.Guest +
			t2.GuestNice - t1.GuestNice + t2.Stolen - t1.Stolen + t2.Idle - t1.Idle
		
		

	}*/
    return times1
}
