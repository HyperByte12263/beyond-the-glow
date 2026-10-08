// Turns the raw downloads in ../raw-data into the small CSV files the web page loads from data/.
// It only filters, totals and renames columns; no values are created or changed.
//
//	cd scripts/prepare_data && go run .   (or: go build -o prepare.exe . && ./prepare.exe)
package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

var (
	raw  = filepath.Join("..", "..", "..", "raw-data")
	data = filepath.Join("..", "..", "data")
)

func main() {
	products := sephoraProducts()
	skinRatings(products)
	trade()
	prices()
	rewindGeoJSON(filepath.Join(raw, "geo", "mys_adm1.geojson"), filepath.Join(data, "malaysia_states.geojson"))
	copyFile(filepath.Join(raw, "geo", "countries-110m.json"), filepath.Join(data, "world-110m.json"))
}

// ---------- helpers ----------

func readCSV(path string) ([]string, [][]string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	r := csv.NewReader(f)
	r.FieldsPerRecord = -1
	r.LazyQuotes = true
	all, err := r.ReadAll()
	if err != nil {
		log.Fatal(path, ": ", err)
	}
	return all[0], all[1:]
}

func cols(header []string) map[string]int {
	m := map[string]int{}
	bom := string([]byte{0xEF, 0xBB, 0xBF})
	for i, h := range header {
		m[strings.TrimPrefix(h, bom)] = i
	}
	return m
}

func writeCSV(name string, header []string, rows [][]string) {
	f, err := os.Create(filepath.Join(data, name))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	w.Write(header)
	w.WriteAll(rows)
	fmt.Printf("%-28s %6d rows\n", name, len(rows))
}

func copyFile(src, dst string) {
	in, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	io.Copy(out, in)
	fmt.Printf("%-28s copied\n", filepath.Base(dst))
}

// geoBoundaries follows the GeoJSON standard (outer rings counter-clockwise), but Vega's map drawing
// (d3-geo) expects outer rings clockwise and holes counter-clockwise, and otherwise fills everything
// outside each state. Each ring is checked with its signed area and reversed only if it points the wrong way.
func rewindGeoJSON(src, dst string) {
	b, err := os.ReadFile(src)
	if err != nil {
		log.Fatal(err)
	}
	var fc struct {
		Type     string `json:"type"`
		Features []struct {
			Type       string          `json:"type"`
			Properties json.RawMessage `json:"properties"`
			Geometry   struct {
				Type        string          `json:"type"`
				Coordinates json.RawMessage `json:"coordinates"`
			} `json:"geometry"`
		} `json:"features"`
	}
	if err := json.Unmarshal(b, &fc); err != nil {
		log.Fatal(err)
	}
	clockwise := func(ring [][]float64) bool {
		a := 0.0
		for i := 0; i+1 < len(ring); i++ {
			a += ring[i][0]*ring[i+1][1] - ring[i+1][0]*ring[i][1]
		}
		return a < 0
	}
	fix := func(poly [][][]float64) {
		for k, ring := range poly {
			if (k == 0) != clockwise(ring) { // outer ring must be clockwise, holes counter-clockwise
				for i, j := 0, len(ring)-1; i < j; i, j = i+1, j-1 {
					ring[i], ring[j] = ring[j], ring[i]
				}
			}
		}
	}
	for i := range fc.Features {
		g := &fc.Features[i].Geometry
		switch g.Type {
		case "Polygon":
			var p [][][]float64
			json.Unmarshal(g.Coordinates, &p)
			fix(p)
			g.Coordinates, _ = json.Marshal(p)
		case "MultiPolygon":
			var mp [][][][]float64
			json.Unmarshal(g.Coordinates, &mp)
			for _, p := range mp {
				fix(p)
			}
			g.Coordinates, _ = json.Marshal(mp)
		}
	}
	out, _ := json.Marshal(fc)
	if err := os.WriteFile(dst, out, 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-28s rewound\n", filepath.Base(dst))
}

func num(s string) float64 { v, _ := strconv.ParseFloat(s, 64); return v }
func f1(v float64) string  { return strconv.FormatFloat(v, 'f', 1, 64) }
func f2(v float64) string  { return strconv.FormatFloat(v, 'f', 2, 64) }

// ---------- Sephora products ----------

type product struct{ id, primary, secondary string }

// Writes skincare products (C6, C7, C9, C10, C12) and brand counts per category (C8 treemap).
func sephoraProducts() map[string]product {
	h, rows := readCSV(filepath.Join(raw, "sephora", "product_info.csv"))
	c := cols(h)
	products := map[string]product{}
	var skin [][]string
	brandCount := map[string]map[string]int{}

	for _, r := range rows {
		id, primary := r[c["product_id"]], r[c["primary_category"]]
		products[id] = product{id, primary, r[c["secondary_category"]]}

		group := primary
		if group != "Skincare" && group != "Makeup" && group != "Hair" && group != "Fragrance" {
			group = "Other"
		}
		if brandCount[group] == nil {
			brandCount[group] = map[string]int{}
		}
		brandCount[group][r[c["brand_name"]]]++

		if primary != "Skincare" {
			continue
		}
		// highlights look like "['Vegan', 'Hydrating']"; store as "Vegan|Hydrating"
		hl := strings.Trim(r[c["highlights"]], "[]")
		hl = strings.ReplaceAll(strings.ReplaceAll(hl, "', '", "|"), "'", "")
		skin = append(skin, []string{id, r[c["product_name"]], r[c["brand_name"]], r[c["secondary_category"]],
			r[c["price_usd"]], r[c["rating"]], r[c["reviews"]], r[c["loves_count"]], hl})
	}
	writeCSV("skincare_products.csv", []string{"product_id", "product_name", "brand", "category", "price_usd", "rating", "reviews", "loves", "highlights"}, skin)

	// top 5 brands per category, the rest summed as "Other brands"
	var tm [][]string
	for _, g := range []string{"Skincare", "Makeup", "Hair", "Fragrance", "Other"} {
		type bc struct {
			b string
			n int
		}
		var list []bc
		for b, n := range brandCount[g] {
			list = append(list, bc{b, n})
		}
		sort.Slice(list, func(i, j int) bool { return list[i].n > list[j].n || (list[i].n == list[j].n && list[i].b < list[j].b) })
		rest, restBrands := 0, 0
		for i, x := range list {
			if i < 5 {
				tm = append(tm, []string{g, x.b, strconv.Itoa(x.n), strconv.Itoa(i + 1)})
			} else {
				rest += x.n
				restBrands++
			}
		}
		tm = append(tm, []string{g, fmt.Sprintf("Other brands (%d)", restBrands), strconv.Itoa(rest), "6"})
	}
	writeCSV("treemap_brands.csv", []string{"category", "brand", "products", "rank"}, tm)
	return products
}

// ---------- Sephora reviews ----------

type acc struct {
	n, rec, recKnown int
	sum              float64
}

// Averages ~1.09 M reviews into one row per skincare product and skin type (oily, dry, all) for C11 and C12.
func skinRatings(products map[string]product) {
	files, _ := filepath.Glob(filepath.Join(raw, "sephora", "reviews_*.csv"))
	agg := map[string]*acc{}
	add := func(k string, rating float64, rec string) {
		a := agg[k]
		if a == nil {
			a = &acc{}
			agg[k] = a
		}
		a.n++
		a.sum += rating
		if rec != "" {
			a.recKnown++
			if num(rec) == 1 {
				a.rec++
			}
		}
	}
	for _, fn := range files {
		f, _ := os.Open(fn)
		r := csv.NewReader(f)
		r.FieldsPerRecord = -1
		r.LazyQuotes = true
		h, _ := r.Read()
		c := cols(h)
		for {
			row, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil || len(row) <= c["product_id"] {
				continue
			}
			id := row[c["product_id"]]
			if products[id].primary != "Skincare" {
				continue
			}
			rating, rec, st := num(row[c["rating"]]), row[c["is_recommended"]], row[c["skin_type"]]
			add(id+"|all", rating, rec)
			if st == "oily" || st == "dry" || st == "combination" || st == "normal" {
				add(id+"|"+st, rating, rec)
			}
		}
		f.Close()
	}
	keys := make([]string, 0, len(agg))
	for k := range agg {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var out [][]string
	for _, k := range keys {
		a := agg[k]
		p := strings.SplitN(k, "|", 2)
		pct := ""
		if a.recKnown > 0 {
			pct = f1(100 * float64(a.rec) / float64(a.recKnown))
		}
		out = append(out, []string{p[0], p[1], strconv.Itoa(a.n), f2(a.sum / float64(a.n)), pct})
	}
	writeCSV("skin_ratings.csv", []string{"product_id", "skin_type", "reviews", "avg_rating", "pct_recommended"}, out)
}

// ---------- UN Comtrade ----------

// Adds partner names and map coordinates to the trade downloads for C1, C2 and C3.
func trade() {
	// partner code -> name, ISO2
	b, err := os.ReadFile(filepath.Join(raw, "comtrade", "partnerAreas.json"))
	if err != nil {
		log.Fatal(err)
	}
	var ref struct {
		Results []struct {
			PartnerCode int
			PartnerDesc string
			ISO2        string `json:"PartnerCodeIsoAlpha2"`
		}
	}
	json.Unmarshal(b, &ref)
	name, iso2 := map[string]string{}, map[string]string{}
	for _, p := range ref.Results {
		k := strconv.Itoa(p.PartnerCode)
		name[k], iso2[k] = p.PartnerDesc, p.ISO2
	}
	name["490"], iso2["490"] = "Taiwan", "TW" // Comtrade lists Taiwan as "Other Asia, nes"
	// Everyday names for readers instead of Comtrade's official ones
	everyday := map[string]string{"Rep. of Korea": "South Korea", "China, Hong Kong SAR": "Hong Kong", "China, Macao SAR": "Macao",
		"Russian Federation": "Russia", "Lao People's Dem. Rep.": "Laos", "Viet Nam": "Vietnam", "Türkiye": "Turkey",
		"United Rep. of Tanzania": "Tanzania", "Iran": "Iran", "Rep. of Moldova": "Moldova", "Bolivia (Plurinational State of)": "Bolivia",
		"Venezuela": "Venezuela", "Dem. Rep. of the Congo": "DR Congo", "Syria": "Syria", "Brunei Darussalam": "Brunei"}
	for k, v := range name {
		if e, ok := everyday[v]; ok {
			name[k] = e
		}
	}

	// ISO2 -> country centre (Google countries.csv via the FIT3179 studio repo)
	h, rows := readCSV(filepath.Join(raw, "geo", "countryInfo.csv"))
	c := cols(h)
	coord := map[string][2]string{}
	for _, r := range rows {
		coord[r[c["country"]]] = [2]string{r[c["latitude"]], r[c["longitude"]]}
	}
	// Country centres that would mislead on the map: use the capital / main port instead.
	coord["ID"] = [2]string{"-6.21", "106.85"} // Jakarta (Indonesia's centre falls on Borneo, next to Malaysia)

	conv := func(in, out string) {
		h, rows := readCSV(filepath.Join(raw, "comtrade", in))
		c := cols(h)
		total := map[string]float64{}
		for _, r := range rows {
			total[r[c["refYear"]]] += num(r[c["primaryValue"]])
		}
		var res [][]string
		for _, r := range rows {
			y, code, v := r[c["refYear"]], r[c["partnerCode"]], num(r[c["primaryValue"]])
			ll := coord[iso2[code]]
			res = append(res, []string{y, code, name[code], iso2[code], strconv.FormatFloat(v, 'f', 0, 64), f2(100 * v / total[y]), ll[0], ll[1]})
		}
		sort.Slice(res, func(i, j int) bool {
			if res[i][0] != res[j][0] {
				return res[i][0] < res[j][0]
			}
			return num(res[i][4]) > num(res[j][4])
		})
		writeCSV(out, []string{"year", "partner_code", "partner", "iso2", "value_usd", "share_pct", "lat", "lon"}, res)
	}
	conv("imports_3304_2015_2024.csv", "trade_imports.csv")
	conv("exports_3304_2024_2024.csv", "trade_exports_2024.csv")
}

// ---------- DOSM prices ----------

// State price rise for C4 and national price lines for C5.
func prices() {
	// DOSM state name -> geoBoundaries state name
	geo := map[string]string{"W.P. Kuala Lumpur": "Kuala Lumpur", "W.P. Labuan": "Labuan", "W.P. Putrajaya": "Putrajaya", "Melaka": "Malacca", "Pulau Pinang": "Penang"}
	h, rows := readCSV(filepath.Join(raw, "dosm", "cpi_2d_state.csv"))
	c := cols(h)
	idx := map[string]map[string]float64{}
	for _, r := range rows {
		if r[c["division"]] != "13" {
			continue
		}
		s := r[c["state"]]
		if idx[s] == nil {
			idx[s] = map[string]float64{}
		}
		idx[s][r[c["date"]]] = num(r[c["index"]])
	}
	var out [][]string
	for s, m := range idx {
		g := s
		if v, ok := geo[s]; ok {
			g = v
		}
		a, b := m["2016-08-01"], m["2026-08-01"]
		out = append(out, []string{s, g, f1(a), f1(b), f1(100 * (b/a - 1))})
	}
	sort.Slice(out, func(i, j int) bool { return num(out[i][4]) > num(out[j][4]) })
	writeCSV("state_price_rise.csv", []string{"state", "geo_name", "index_aug2016", "index_aug2026", "rise_pct"}, out)

	var nat [][]string
	h, rows = readCSV(filepath.Join(raw, "dosm", "cpi_2d.csv"))
	c = cols(h)
	for _, r := range rows {
		if r[c["division"]] == "overall" && r[c["date"]] >= "2010-01-01" {
			nat = append(nat, []string{r[c["date"]], "All items", r[c["index"]]})
		}
	}
	label := map[string]string{"13120": "Beauty products", "13132": "Salon & grooming"}
	h, rows = readCSV(filepath.Join(raw, "dosm", "cpi_5d.csv"))
	c = cols(h)
	for _, r := range rows {
		if l, ok := label[r[c["subclass"]]]; ok && r[c["date"]] >= "2010-01-01" {
			nat = append(nat, []string{r[c["date"]], l, r[c["index"]]})
		}
	}
	writeCSV("price_index_national.csv", []string{"date", "series", "index"}, nat)
}
