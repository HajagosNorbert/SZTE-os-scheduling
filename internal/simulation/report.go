package simulation

import (
	"os"

	"fmt"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"

	"github.com/go-echarts/go-echarts/v2/components"
)

func CreateResultReport(input []Proc, res SimResult, algLongName string) {
	procXAxis := genProcXAxis(res)
	res.algName = algLongName

	// Create an instance of a page
	page := components.NewPage()
	page.PageTitle = res.algName + " Report"

	// chart
	idleBarItems := genIdles(res)
	idleBar := charts.NewBar()
	idleBar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "Idle ticks",
		Subtitle: "Idle ticks for each process",
	}))
	idleBar.SetXAxis(procXAxis).AddSeries("Idle ticks", idleBarItems)

	// chart
	relIdleItems := genRelIdles(res.procResults)
	relIdleBar := charts.NewBar()
	relIdleBar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "Relative Idle ticks",
		Subtitle: "Idle ticks per turn around time (the total ticks the process took) for each process",
	}))
	relIdleBar.SetXAxis(procXAxis).AddSeries("Relative idle ticks", relIdleItems)

	// chart
	pieData := []opts.PieData{{
		Name:  "Idle",
		Value: res.idleTicks,
	}, {
		Name:  "Bussy",
		Value: res.totalTicks,
	}}
	totalIdlePie := charts.NewPie()
	totalIdlePie.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title: algLongName + ": Idle ticks to total ticks",
		Subtitle: "CPU utilization: "+ fmt.Sprintf("%.2f", 100 * float64(res.totalTicks) / float64(res.idleTicks + res.totalTicks)) + "%" ,
	}))
	totalIdlePie.
		AddSeries("Ticks", pieData).
		SetSeriesOptions(charts.
			WithLabelOpts(
				opts.Label{
					Show:      true,
					Formatter: "{b}: {c}",
				},
			))


	// chart
	totalCtxSwitches := 0
	for _, proc := range res.procResults {
		totalCtxSwitches += proc.ctxSwitchCount
	}
	ctxSwitchItems := genCtxSwitches(res.procResults)
	ctxSwitchBar := charts.NewBar()
	ctxSwitchBar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    fmt.Sprintf("Context switches (%d in total)", totalCtxSwitches),
		Subtitle: "Context switches for each procsess",
	}))
	ctxSwitchBar.SetXAxis(procXAxis).AddSeries("Context Switch", ctxSwitchItems)

	//averages 
	avgItems := genAvgs(res.procResults)
	avgBar := charts.NewBar()
	avgBar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "Averages",
		Subtitle: "Turn Around Time",
	}))

	avgBar.SetXAxis([]string{"TAT", "Wait - time"}).AddSeries("Averages", avgItems) 

	page.AddCharts(totalIdlePie, idleBar, relIdleBar, ctxSwitchBar, avgBar)
	// Save the result to a file
	f, err := os.Create("output.html")
	if err != nil {
		panic(err)
	}
	page.Render(f)
	fmt.Println("Created output.html file with the report")
}

func genProcXAxis(res SimResult) []string {
	items := make([]string, 0)
	for i := range res.procResults {
		items = append(items, "#"+fmt.Sprint(i))
	}
	return items
}

func genRelIdles(procResults []ProcResult) []opts.BarData {
	items := make([]opts.BarData, 0)
	for i := 0; i < len(procResults); i++ {
		relIdles := float64(procResults[i].idleTicks) / float64(procResults[i].idleTicks+procResults[i].totalTicks)
		items = append(items, opts.BarData{Value: relIdles})
	}
	return items
}

func genIdles(res SimResult) []opts.BarData {
	idles := make([]opts.BarData, 0)
	for i := 0; i < len(res.procResults); i++ {
		idle := res.procResults[i].idleTicks
		idles = append(idles, opts.BarData{Value: idle})
	}
	return idles
}

func genCtxSwitches(procResults []ProcResult) []opts.BarData {
	items := make([]opts.BarData, 0)
	for i := 0; i < len(procResults); i++ {
		ctxSwitch := procResults[i].ctxSwitchCount
		items = append(items, opts.BarData{Value: ctxSwitch})
	}
	return items}

func genAvgs(procResults []ProcResult) []opts.BarData {

	//TAT
	tatSum := 0
	waitTimeSum := 0
	for i := 0; i < len(procResults); i++ {
		tat := procResults[i].idleTicks + procResults[i].totalTicks 
		waitTimeSum += procResults[i].idleTicks
		tatSum += tat
		// ctxSwitch := procResults[i].ctxSwitchCount
	}
	waitTimeAvg := waitTimeSum / len(procResults)
	tatAvg := tatSum / len(procResults)

	items := []opts.BarData{ {Value: tatAvg}, {Value: waitTimeAvg}}
	return items
}
