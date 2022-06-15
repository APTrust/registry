// charts.js
//
// Utility functions for customizing and rendering charts.
//

// fillColors are used in bar charts, pie charts, etc. rendered
// by chart.js on the front end. Used by helpers/templates.go.
const fillColors = [
	"rgba(187, 149, 102, 1)",
	"rgba(19, 19, 92, 1)",
	"rgba(51, 48, 48, 1)",
	"rgba(96, 130, 146, 1)",
	"rgba(147, 90, 21, 1)",
	"rgba(96, 151, 205, 1)",
]

// barBorders are used in bar charts, pie charts, etc. rendered
// by chart.js on the front end. Used by helpers/templates.go.
const barBorders = [
	"rgba(187, 149, 102, 1)",
	"rgba(19, 19, 92, 1)",
	"rgba(51, 48, 48, 1)",
	"rgba(96, 130, 146, 1)",
	"rgba(147, 90, 21, 1)",
	"rgba(96, 151, 205, 1)",
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
