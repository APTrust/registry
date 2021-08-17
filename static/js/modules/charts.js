// charts.js
//
// Utility functions for customizing and rendering charts.
//

// fillColors are used in bar charts, pie charts, etc. rendered
// by chart.js on the front end. Used by helpers/templates.go.
const fillColors = [
	"rgba(255, 99, 132, 0.2)",
	"rgba(255, 159, 64, 0.2)",
	"rgba(255, 205, 86, 0.2)",
	"rgba(75, 192, 192, 0.2)",
	"rgba(54, 162, 235, 0.2)",
	"rgba(153, 102, 255, 0.2)",
	"rgba(201, 203, 207, 0.2)",
	"rgba(100, 180, 255, 0.2)",
]

// barBorders are used in bar charts, pie charts, etc. rendered
// by chart.js on the front end. Used by helpers/templates.go.
const barBorders = [
	"rgb(255, 99, 132)",
	"rgb(255, 159, 64)",
	"rgb(255, 205, 86)",
	"rgb(75, 192, 192)",
	"rgb(54, 162, 235)",
	"rgb(153, 102, 255)",
	"rgb(201, 203, 207)",
	"rgb(100, 180, 255)",
]

// fillColor returns a color for a bar, pie slice, etc. in a
// chart.js chart.
function fillColor(i) {
	return fillColors[i % fillColors.length]
}

// borderColor returns a border color for a bar, pie slice, etc.
// in a chart.js chart.
function borderColor(i) {
	return barBorders[i % barBorders.length]
}

// chartColors returns however many chart colors of whichever type
// you ask for.
export function chartColors(whatKind, howMany) {
    let colors = []
    let fn = fillColor
    if (whatKind == 'border') {
        fn = borderColor
    }
    for (let i=0; i < howMany; i++) {
        colors.push(fn(i))
    }
    return colors
}
