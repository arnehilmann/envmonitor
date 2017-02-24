package main

import (
    "log"
    "os/exec"
    "strconv"
    "strings"
    "time"

    //"github.com/wcharczuk/go-chart"

    . "github.com/arnehilmann/goutils"
)

type Timeline struct {
    name string
    epochs []time.Time
    values []float64
}

func NewTimeline(name string) Timeline {
    return Timeline{name, []time.Time{}, []float64{}}
}

func main() {
    out, err := exec.Command("rrdtool", "fetch", "res/temperature.rrd", "AVERAGE",
                             "-s", "02/22/2017", "-e", "02/24/2017").Output()
                             //"-s", "now - 1 day",
                             //"-e",  "now - 12 hours").Output()
    PanicIf(err)

    var timelines []Timeline
    for _, line := range(strings.Split(string(out), "\n")) {
        if line == "" {
            continue
        }
        fields := strings.Fields(line)
        if len(timelines) == 0 {
            for _, name := range(fields) {
                timelines = append(timelines, NewTimeline(name))
            }
            continue
        }
        epoch, err := strconv.Atoi(strings.TrimSuffix(fields[0], ":"))
        PanicIf(err)
        for index, field := range(fields[1:]) {
            timelines[index].epochs = append(timelines[index].epochs, time.Unix(int64(epoch), 0))
            value, err := strconv.ParseFloat(field, 32)
            PanicIf(err)
            timelines[index].values = append(timelines[index].values, float64(value))
        }
    }
    for _, line := range(timelines) {
        log.Println(line.name)
        for i := range(line.epochs) {
            log.Println(line.epochs[i], line.values[i])
        }
        break
    }
}
