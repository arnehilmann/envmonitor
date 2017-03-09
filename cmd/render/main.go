package main

import (
	"log"
	"time"

	"github.com/ziutek/rrd"

	. "github.com/arnehilmann/goutils"
)

type Timeline struct {
	name   string
	epochs []time.Time
	values []float64
}

func NewTimeline(name string) Timeline {
	return Timeline{name, []time.Time{}, []float64{}}
}

func DecodeRrd(filename, fun, start, end string) ([]Timeline, error) {
	info, err := rrd.Info(filename)
	if err != nil {
		return nil, err
	}
	log.Println(info)
	step := info["step"].(uint)
	log.Println(step)

	startTime, err := time.Parse("2006-01-02", start)
	if err != nil {
		return nil, err
	}
	endTime, err := time.Parse("2006-01-02", end)
	if err != nil {
		return nil, err
	}
	result, err := rrd.Fetch(filename, fun, startTime, endTime, time.Duration(step)*time.Second)
	if err != nil {
		return nil, err
	}

	var timelines []Timeline
	var ts time.Time
	for i, name := range result.DsNames {
		timelines = append(timelines, NewTimeline(name))
		for t := 0; t < result.RowCnt; t++ {
			ts = result.Start.Add(time.Duration(t) * result.Step)
			timelines[i].epochs = append(timelines[i].epochs, ts)
			value := result.ValueAt(i, t)
			timelines[i].values = append(timelines[i].values, value)
		}
	}
	return timelines, nil
}

func (timeline Timeline) Dump(fun func(...interface{})) {
	fun("timeline", timeline.name)
	for i := range timeline.epochs {
		fun(timeline.epochs[i], timeline.values[i])
	}
}

func main() {
	timelines, err := DecodeRrd("res/temperature.rrd",
		"AVERAGE",
		"2017-02-22",
		"2017-02-24")
	PanicIf(err)
    for _, timeline := range timelines {
        timeline.Dump(log.Println)
    }
}
