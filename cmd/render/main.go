package main

import (
	"log"
	"math"
	"os"
	"path/filepath"
    "strconv"
    "strings"
	"time"

	"github.com/ziutek/rrd"

	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"

	"github.com/arnehilmann/go-chart-contrib"

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
		for t := 0; t < result.RowCnt - 1; t++ {
			ts = result.Start.Add(time.Duration(t) * result.Step)
			value := result.ValueAt(i, t)
			if !math.IsNaN(value) {
				timelines[i].epochs = append(timelines[i].epochs, ts)
				timelines[i].values = append(timelines[i].values, value)
			}
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
    width := 800
    height := 600
    yaxis_linespacing := 5.0
    axis_color := "0F0"
    result := "aha.png"

	c := chart.Chart{
		Width:  width,
		Height: height,
		XAxis: chart.XAxis{
			Style: chart.Style{
				Show: true,
				StrokeColor: drawing.ColorFromHex(axis_color),
			},
			ValueFormatter: chart.TimeValueFormatterWithFormat("2006-01-02T15:04"),
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show:        true,
				StrokeColor: drawing.ColorFromHex(axis_color),
			},
			ValueFormatter: func(v interface{}) string {
				return chart.FloatValueFormatterWithFormat(v, "%0.1f")
			},
			Range: chartcontrib.ContinuousRangeWithTicksLinespacing(yaxis_linespacing),
		},
        Series: []chart.Series{},
    }

    specs := "res/temperature.rrd:0:MAX:2017-02-12:2017-02-23:00F:3" +
        //""
        " res/temperature.rrd:0:AVERAGE:2017-02-22:2017-02-23:F00:3"
        //" res/temperature.rrd:0:MIN:2017-02-12:2017-02-24:00F:3" +
        //" res/temperature.rrd:0:AVERAGE:2017-02-22:2017-02-24:F00:3"

    for _, spec := range(strings.Split(specs, " ")) {
        token := strings.Split(spec, ":")
        filename := token[0]
        series_index, err := strconv.Atoi(token[1])
        PanicIf(err)
        function := token[2]
        start := token[3]
        end := token[4]
        color := token[5]
        stroke, err := strconv.ParseFloat(token[6], 64)
        PanicIf(err)

        timelines, err := DecodeRrd(filename, function, start, end)
        PanicIf(err)

        timelines[series_index].Dump(log.Println)

        ts := chart.TimeSeries{
                    Name:    timelines[series_index].name,
                    XValues: timelines[series_index].epochs,
                    YValues: timelines[series_index].values,
                    Style: chart.Style{
                        Show:        true,
                        StrokeWidth: float64(stroke),
                        StrokeColor: drawing.ColorFromHex(color),
                    },
                }

        c.Series = append(c.Series, ts)
    }

	f, err := os.Create(result)
	PanicIf(err)
	defer f.Close()
	c.Render(chart.PNG, f)
	log.Println("chart can be found in", f.Name())
}

func main_old() {
	timelines, err := DecodeRrd("res/temperature.rrd",
		"AVERAGE",
		"2017-02-22",
		"2017-02-24")
	PanicIf(err)
	log.Println("--------------------")
	timelines[0].Dump(log.Println)

	linespacing := 5.0

	c := chart.Chart{
		Width:  800,
		Height: 600,
		XAxis: chart.XAxis{
			Style: chart.Style{
				Show: true,
			},
			ValueFormatter: chart.TimeValueFormatterWithFormat("2006-01-02T15:04"),
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show:        true,
				StrokeColor: drawing.Color{0, 0, 255, 255},
			},
			ValueFormatter: func(v interface{}) string {
				return chart.FloatValueFormatterWithFormat(v, "%0.1f")
			},
			Range: chartcontrib.ContinuousRangeWithTicksLinespacing(linespacing),
		},
		Series: []chart.Series{
			chart.TimeSeries{
				Name:    timelines[0].name,
				XValues: timelines[0].epochs,
				YValues: timelines[0].values,
				Style: chart.Style{
					Show:        true,
					StrokeWidth: 3.0,
					StrokeColor: drawing.Color{0, 0, 255, 255},
				},
			},
		},
	}
	filename := filepath.Join(os.TempDir(), "envmonitor-test.png")
	f, err := os.Create(filename)
	PanicIf(err)
	defer f.Close()
	c.Render(chart.PNG, f)
	log.Println("chart can be found in", f.Name())
}
