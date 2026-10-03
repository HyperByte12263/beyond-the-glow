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
| `data/` | Small, chart-ready data files (about 1.3 MB in total) |
| `scripts/` | How `data/` was made from the raw downloads |
| `sketch/` | PDF of the hand-drawn sketch |

## Colour key

One colour means one thing across every chart: trade amber `#c08a1e`, prices teal `#1d8072`, products violet `#7656b0`, ingredients brick `#a3442c`, oily skin dark olive `#4f5f0c`, dry skin blue `#5f8ccc`, context grey `#c4bcc1`.

## Typography

Headings use **Gloock**, a high-contrast serif that echoes cosmetics packaging. Body text and every chart label use **Figtree**, a friendly sans-serif with clear numbers for prices, percentages and star ratings. Body text is 17 px with 1.6 line height, lines are capped at about 60 characters (7–10 words), and text is left-aligned.

Visual hierarchy follows the FIT3179 figure-ground notes: page title > section titles > chart titles (21 px, dark serif) > side-text headings (18 px) > body text (17 px, near-black). Secondary information (subtitles, axes, legends, captions, sources) is smaller and grey; annotations and tooltip values are dark and bold. No chart text is smaller than 12 px. Key words in the text are coloured with the chart colour they refer to, acting as the colour key.

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
   - averages the 1.09 M reviews into one row per skincare product and skin type,
   - keeps the top 5 brands per category for the treemap,
   - adds everyday country names and map coordinates to the trade data (Indonesia is placed at Jakarta, because its centre falls on Borneo),
   - computes each state's price rise from Aug 2016 to Aug 2026,
   - fixes the ring direction of the state boundaries so Vega draws them correctly.

## Previewing

Browsers block the chart data when `index.html` is opened straight from disk, so run a small local server:

```
powershell -File scripts/serve.ps1
```

and open http://localhost:8000.
