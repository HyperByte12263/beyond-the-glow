// Shared Vega-Lite config: fonts, text colours and quiet axes, applied to every chart.
// Data colours are set inside each spec so the JSON files stay readable on their own.
// Colour key (keep in sync with css/style.css):
//   trade #c08a1e · prices #1d8072 · products #7656b0 · ingredients #a3442c
//   oily #4f5f0c · dry #5f8ccc · context #c4bcc1
//
// Visual hierarchy (FIT3179 typography and figure-ground notes):
//   figure = chart title (large, dark serif) and annotations (dark);
//   ground = subtitles, axes, legends (smaller, grey). No chart text is smaller than 12 px.
const VEGA_THEME = {
  font: "Figtree, system-ui, sans-serif",
  background: null,
  view: { stroke: null },
  title: {
    font: "Gloock, Georgia, serif",
    fontSize: 21,
    fontWeight: "normal",
    color: "#261b24",
    anchor: "start",
    frame: "bounds",        // title starts at the chart's left edge, on the same sight-line as the page text
    subtitleFont: "Figtree, system-ui, sans-serif",
    subtitleFontSize: 14,
    subtitleColor: "#6e6069",
    subtitlePadding: 6,
    subtitleLineHeight: 19,
    offset: 16
  },
  axis: {
    labelColor: "#6e6069",
    labelFontSize: 13,
    titleColor: "#6e6069",
    titleFontSize: 13,
    titleFontWeight: 500,
    domainColor: "#e4dce1",
    tickColor: "#e4dce1",
    gridColor: "#efe9ec"
  },
  legend: { labelColor: "#6e6069", titleColor: "#6e6069", labelFontSize: 13, titleFontSize: 13, titleFontWeight: 500 },
  header: { labelFontSize: 13, titleFontSize: 13 },
  text: { color: "#261b24", fontSize: 13 }
};
