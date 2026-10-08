# Beyond the Glow

Where Malaysia's skincare comes from, what it costs us, and whether it works for skin that lives in the humidity.

FIT3179 Data Visualisation 2, Monash University Malaysia, 2026. Author: Tan Le Han.

## Structure

| Path | What it is |
|---|---|
| `index.html` | The page (Pure.css grid + vega-embed, following the FIT3179 studio examples) |
| `css/style.css` | Layout, typography and the colour key |
| `js/theme.js` | Shared Vega-Lite config (fonts, axes) |
| `js/main.js` | Embeds each chart spec into the page |
| `specs/*.vg.json` | One Vega-Lite (or Vega) spec per chart, formatted for reading |
| `data/` | Small, chart-ready data files (about 1.4 MB in total) |
| `scripts/` | How `data/` was made from the raw downloads |
| `sketch/` | PDF of the hand-drawn sketch |

## Colour key

One colour means one thing across every chart: trade amber `#c08a1e`, prices teal `#1d8072`, products violet `#7656b0`, ingredients brick `#a3442c`, oily skin dark olive `#4f5f0c`, dry skin blue `#5f8ccc`, context grey `#c4bcc1`. Map oceans are a pale blue-grey `#e2eaee` and neighbouring land a near-white grey, so the data reads as figure on a quiet ground.

## Typography

Headings use **Gloock**, a high-contrast serif that echoes cosmetics packaging. Body text and every chart label use **Figtree**, a friendly sans-serif with clear numbers for prices, percentages and star ratings. Body text is 17 px with 1.6 line height, lines are capped at about 60 characters (7–10 words), and text is left-aligned.

Visual hierarchy follows the FIT3179 figure-ground notes: page title > section titles > chart titles (21 px, dark serif) > side-text headings (18 px) > body text (17 px, near-black). Secondary information (subtitles, axes, legends, captions, sources) is smaller and grey; annotations and tooltip values are dark and bold. No chart text is smaller than 12 px. Key words in the text are coloured with the chart colour they refer to, acting as the colour key.

## Charts

| Id | Idiom | Question it answers |
|---|---|---|
| C8 | Treemap (Vega) | How many brands share the shelf? |
| C1 | Flow map with curved routes + zoomed inset | Where do our imports come from? |
| C2b | Streamgraph | How have imports grown, and who drove it? |
| C2 | Bump chart | How did each supplier's rank change? |
| C3 | Proportional-symbol map | Who buys Malaysia's exports? |
| C4 | Choropleth map | Which states saw the biggest price rises? |
| C5 | Indexed line chart | Did beauty products rise faster than everything else? |
| C6 | Boxplot with jittered dots | How much does each product type cost? |
| C7 | Strip plot by price band, linked to a bar chart | Do pricier products earn better ratings? |
| C9 | Heatmap | Which ingredients does each product type advertise? |
| C10 | Waffle charts | How common are vegan, cruelty-free and fragrance-free labels? |
| C11a | Ridgeline plot | Do skin types rate products differently? |
| C11 | Dumbbell chart | Which products split oily and dry skin? |
| C12 | Ranked lollipop chart | Which sunscreens suit oily skin? |

## Interaction

Interaction follows Shneiderman's mantra, "overview first, zoom and filter, then details-on-demand" (as discussed by Stephen Few, *The Surest Path to Visual Discovery*, 2006), kept light because the page presents a story rather than an exploration tool:

- **Overview first:** headline numbers and the treemap open the page; each section opens with its overview chart.
- **Zoom and filter:** click a price band in the strip plot (C7) to show its ratings in the bar chart beside it (Shift-click adds bands); choose a product type. The import map (C1) has a zoomed South-East Asia view next to the world view.
- **Details-on-demand:** every chart has tooltips; the price lines (C5) show all three prices for the hovered month.
- **Focus plus context:** hovering highlights one item and fades the rest without hiding them (C2 countries, C4 states, C11 products, C7 price bands).
- **Linked views:** hovering a country in the bump chart (C2) highlights its route on the import flow map (C1), linking the detail back to the overview.

## Data sources

| Data | Source | Licence |
|---|---|---|
| Malaysia's trade in HS 3304 (beauty, make-up and skincare preparations), imports 2015–2024 and exports 2024 | [UN Comtrade](https://comtradeplus.un.org/) public API | UN Comtrade terms of use |
| Consumer Price Index by state (Division 13), national CPI by division and by subclass | Department of Statistics Malaysia, [OpenDOSM](https://open.dosm.gov.my/data-catalogue) | CC BY 4.0 |
| Sephora products and skincare reviews (collected March 2023) | Nady Inky, [Kaggle](https://www.kaggle.com/datasets/nadyinky/sephora-products-and-skincare-reviews) | CC BY 4.0 |
| Malaysia state boundaries | [geoBoundaries](https://www.geoboundaries.org/) gbOpen MYS ADM1 | ODbL |
| World map | [world-atlas](https://github.com/topojson/world-atlas) 110m, from Natural Earth | Public domain |
| Country centre coordinates | Google [countries.csv](https://developers.google.com/public-data/docs/canonical/countries_csv), via the FIT3179 studio examples | See source |

## Rebuilding `data/`

The raw downloads (about 510 MB, mostly Sephora reviews) live in `../raw-data/` and are not part of this repository.

1. `powershell -File scripts/fetch_comtrade.ps1 -Flow M -From 2015 -To 2024` and `-Flow X -From 2024 -To 2024` download the trade data.
2. `cd scripts/prepare_data && go run .` writes everything in `data/`. It only filters, totals and renames columns (no values are created or changed):
   - averages the 1.09 million reviews into one row per skincare product and skin type (all, oily, dry, combination, normal),
   - keeps the top 5 brands per category for the treemap,
   - adds everyday country names and map coordinates to the trade data (Indonesia is placed at Jakarta, because its centre falls on Borneo),
   - computes each state's price rise from August 2016 to August 2026,
   - fixes the ring direction of the state boundaries so Vega draws them correctly.

## Previewing

Browsers block the chart data when `index.html` is opened straight from disk, so run a small local server:

```
powershell -File scripts/serve.ps1
```

and open http://localhost:8000.
