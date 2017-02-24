package main

import (
    "log"
    "os/exec"
    "os"
    "strconv"
    "strings"
    "time"
    "math"

    "github.com/wcharczuk/go-chart"
    "github.com/wcharczuk/go-chart/drawing"

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
            value, err := strconv.ParseFloat(field, 32)
            PanicIf(err)
            if math.IsNaN(value) {
                continue
            }
            timelines[index].epochs = append(timelines[index].epochs, time.Unix(int64(epoch), 0))
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

    log.Println("assembling graph")

    graph := chart.Chart{
        Width: 400,
        Height: 200,
        XAxis: chart.XAxis{
            Style: chart.Style{
                Show: true,
            },
            ValueFormatter: func (v interface{}) string {
                return chart.TimeValueFormatterWithFormat(v, "2006-01-02T15:04")},
        },
        YAxis: chart.YAxis{
            Style: chart.Style{
                Show: true,
            },
            ValueFormatter: func (v interface{}) string {
                return chart.FloatValueFormatterWithFormat(v, "%0.1f")},
            Ticks: []chart.Tick{
                chart.Tick{4.0, "4"},
                chart.Tick{6.0, "6"},
                chart.Tick{7.0, "7"},
                chart.Tick{14.0, "14"},
            },
        },
        Series: []chart.Series{
            chart.TimeSeries{
                Name: timelines[0].name,
                XValues: timelines[0].epochs,
                YValues: timelines[0].values,
                Style: chart.Style{
                    Show: true,
                    StrokeWidth: 5.0,
                    StrokeColor: drawing.Color{255, 0, 0, 255},
                },
            },
        },
    }
    err = graph.Series[0].Validate()
    PanicIf(err)

    f, err := os.Create("aha.png")
    PanicIf(err)
    defer f.Close()
    graph.Render(chart.PNG, f)

    log.Println("done")
}
