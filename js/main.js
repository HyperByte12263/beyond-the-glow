// Embeds each chart spec from specs/ into its container on the page,
// and adds a source line under each chart.
const CHARTS = {
  c1: "specs/c1_imports_flow_map.vg.json",
  c2: "specs/c2_import_rank_bump.vg.json",
  c3: "specs/c3_exports_symbol_map.vg.json",
  c4: "specs/c4_state_price_rise.vg.json",
  c5: "specs/c5_price_index.vg.json",
  c6: "specs/c6_price_boxplot.vg.json",
  c7: "specs/c7_price_vs_rating.vg.json",
  c8: "specs/c8_brand_treemap.vg.json",
  c9: "specs/c9_label_heatmap.vg.json",
  c10: "specs/c10_label_waffles.vg.json",
  c11: "specs/c11_oily_vs_dry_dumbbell.vg.json",
  c12: "specs/c12_sunscreen_parallel.vg.json"
};

const SOURCES = {
  c1: "Source: UN Comtrade, Malaysia's imports of beauty, make-up and skincare products (trade code 3304), 2024. Values in US dollars.",
  c2: "Source: UN Comtrade, Malaysia's imports of beauty, make-up and skincare products (trade code 3304), 2015–2024.",
  c3: "Source: UN Comtrade, Malaysia's exports of beauty, make-up and skincare products (trade code 3304), 2024. Values in US dollars.",
  c4: "Source: Department of Statistics Malaysia (OpenDOSM), Consumer Price Index by state, Division 13, August 2016 and August 2026. Boundaries: geoBoundaries.",
  c5: "Source: Department of Statistics Malaysia (OpenDOSM), Consumer Price Index, January 2010 – August 2026.",
  c6: "Source: Sephora Products and Skincare Reviews (Kaggle), 2,420 skincare products from Sephora's US store, collected March 2023.",
  c7: "Source: Sephora Products and Skincare Reviews (Kaggle), collected March 2023. Ratings are customers' average stars.",
  c8: "Source: Sephora Products and Skincare Reviews (Kaggle), all 8,494 products, collected March 2023.",
  c9: "Source: Sephora Products and Skincare Reviews (Kaggle), product highlight labels, collected March 2023.",
  c10: "Source: Sephora Products and Skincare Reviews (Kaggle), product highlight labels, collected March 2023.",
  c11: "Source: Sephora Products and Skincare Reviews (Kaggle), 1.09 million skincare reviews, 2008–2023. Skin type as reported by each reviewer.",
  c12: "Source: Sephora Products and Skincare Reviews (Kaggle), products and reviews, collected March 2023."
};

// Linked views (overview + detail): hovering a country in the bump chart (C2)
// highlights its route on the import flow map (C1).
const views = {};
let linked = false;
function linkCharts() {
  if (linked || !views.c1 || !views.c2) return;
  linked = true;
  views.c2.addSignalListener("hover", (name, value) => {
    const partner = value && value.partner ? [].concat(value.partner)[0] : "";
    views.c1.signal("focus", partner || "").runAsync();
  });
}

function renderCharts(theme) {
  for (const [id, spec] of Object.entries(CHARTS)) {
    const el = document.getElementById(id);
    if (!el) continue;
    if (!spec) {
      el.innerHTML = `<div class="placeholder">${id.toUpperCase()}: chart coming soon</div>`;
      continue;
    }
    vegaEmbed(el, spec, { actions: false, config: theme, renderer: "svg" }).then(result => {
      views[id] = result.view;
      linkCharts();
    }).catch(err => {
      console.error(id, err);
      el.innerHTML = `<div class="placeholder">${id.toUpperCase()} failed to load: ${err.message}</div>`;
    });
    if (SOURCES[id] && !el.nextElementSibling?.classList.contains("source")) {
      const p = document.createElement("p");
      p.className = "source";
      p.textContent = SOURCES[id];
      el.after(p);
    }
  }
}

renderCharts(VEGA_THEME);
