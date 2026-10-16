package matchmaking

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/vamsikrishnap23/astrology_backend_go/internal/domain"
)

// FRONTEND_TRANSLATIONS.json is the translation mapping for /api/vivaha-pontana. These tests fail if the
// backend can output a key, an enumerated value or a piece of text that the file does not cover, and also
// if the file holds an entry the backend can no longer produce.

type tText struct{ En, Te string }

type tTemplate struct {
	Te     string             `json:"te"`
	Params map[string]*string `json:"params"` // enum set name, or null for a value copied as is
}

type translations struct {
	RuleSets  map[string]tText            `json:"rule_sets"`
	Fields    map[string]tText            `json:"fields"`
	Kootas    map[string]tText            `json:"kootas"`
	Checks    map[string]tText            `json:"checks"`
	Risks     map[string]tText            `json:"risks"`
	Enums     map[string]map[string]tText `json:"enums"`
	EnumPaths map[string]string           `json:"enum_paths"`
	Strings   map[string]string           `json:"strings"`
	Templates map[string]tTemplate        `json:"templates"`
}

func loadTranslations(t *testing.T) translations {
	t.Helper()
	raw, err := os.ReadFile("../../../FRONTEND_TRANSLATIONS.json")
	if err != nil {
		t.Fatal(err)
	}
	var tr translations
	if err := json.Unmarshal(raw, &tr); err != nil {
		t.Fatalf("FRONTEND_TRANSLATIONS.json is not valid: %v", err)
	}
	return tr
}

var placeholderRE = regexp.MustCompile(`\{(\w+)\}`)

// translate mirrors the algorithm the frontend must use (see docs/api_ashtakoota.md): exact text first,
// then the templates. It returns the Telugu text and the keys of every template that matched validly.
func (tr translations) translate(text string) (string, []string) {
	if te, ok := tr.Strings[text]; ok {
		return te, nil
	}
	var te string
	var matched []string
	keys := make([]string, 0, len(tr.Templates))
	for k := range tr.Templates {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		tpl := tr.Templates[key]
		var names []string
		pattern := placeholderRE.ReplaceAllStringFunc(regexp.QuoteMeta(key), func(m string) string { return m })
		// QuoteMeta escaped the braces; swap each escaped placeholder for a capture group
		pattern = regexp.MustCompile(`\\\{(\w+)\\\}`).ReplaceAllStringFunc(pattern, func(m string) string {
			names = append(names, strings.Trim(m, `\{}`))
			return `(.+?)`
		})
		sub := regexp.MustCompile("^" + pattern + "$").FindStringSubmatch(text)
		if sub == nil {
			continue
		}
		vals := map[string]string{}
		ok := true
		for i, n := range names {
			set := tpl.Params[n]
			if set == nil {
				vals[n] = sub[i+1]
				continue
			}
			e, found := tr.Enums[*set][sub[i+1]]
			if !found {
				ok = false
				break
			}
			vals[n] = e.Te
		}
		if !ok {
			continue
		}
		matched = append(matched, key)
		te = placeholderRE.ReplaceAllStringFunc(tpl.Te, func(m string) string { return vals[m[1:len(m)-1]] })
	}
	return te, matched
}

// ---- collecting everything the backend can output ----

type enumUse struct{ path, value string }

type outputs struct {
	texts map[string]string // text -> where it came from
	enums []enumUse
	names map[string]string // json key -> emitted name
}

func (o *outputs) text(source, s string) {
	if s == "" {
		return
	}
	if _, ok := o.texts[s]; !ok {
		o.texts[s] = source
	}
}

func (o *outputs) enum(path, value string) { o.enums = append(o.enums, enumUse{path, value}) }

func (o *outputs) koota(key string, k domain.Koota) {
	o.names[key] = k.Name
	o.text("koota "+key+" explanation", k.Explanation)
	o.text("koota "+key+" cancellation", k.CancellationReason)
	o.text("koota "+key+" suggestion", k.Suggestion)
	o.enum("ashtakoota.kootas."+key+".groom", k.Groom)
	o.enum("ashtakoota.kootas."+key+".bride", k.Bride)
}

func (o *outputs) check(key string, c domain.Check) {
	o.names["check:"+key] = c.Name
	o.text("check "+key+" explanation", c.Explanation)
	o.text("check "+key+" cancellation", c.CancellationReason)
	o.enum("checks."+key+".groom", c.Groom)
	o.enum("checks."+key+".bride", c.Bride)
}

// allMoons returns one Moon position per (nakshatra, pada): 27 x 4 = 108.
func allMoons() []domain.MoonDetails {
	var out []domain.MoonDetails
	for n := 0; n < 27; n++ {
		for p := 0; p < 4; p++ {
			out = append(out, moonFromLongitude(float64(n)*800.0/60.0+(float64(p)+0.5)*200.0/60.0))
		}
	}
	return out
}

func collectOutputs(t *testing.T) *outputs {
	t.Helper()
	o := &outputs{texts: map[string]string{}, names: map[string]string{}}

	moons := allMoons()
	for _, m := range moons {
		for _, who := range []string{"groom", "bride"} {
			o.enum("*.moon.sign", m.Sign)
			o.enum("*.moon.nakshatra", m.Nakshatra)
			o.enum("*.moon.nakshatra_lord", m.NakshatraLord)
			o.enum("*.moon.rashi_lord", m.RashiLord)
			_ = who
		}
	}

	maitriScores := []float64{0, 3, 4, 5}
	lagnaRelations := []string{"Same", "Friendly", "Inimical", "Neutral"}
	for _, g := range moons {
		for _, b := range moons {
			o.koota("varna", varnaKoota(g, b))
			o.koota("vashya", vashyaKoota(g.Sign, b.Sign))
			tara := taraKoota(g, b)
			o.koota("tara", tara)
			td := tara.Details.(domain.TaraDetails)
			o.enum("ashtakoota.kootas.tara.details.groom_to_bride.tara", td.GroomToBride.Tara)
			o.enum("ashtakoota.kootas.tara.details.bride_to_groom.tara", td.BrideToGroom.Tara)
			o.koota("yoni", yoniKoota(g, b))
			gm := grahaMaitriKoota(g, b)
			o.koota("graha_maitri", gm)
			gd := gm.Details.(domain.GrahaMaitriDetails)
			o.enum("ashtakoota.kootas.graha_maitri.details.groom_rashi", gd.GroomRashi)
			o.enum("ashtakoota.kootas.graha_maitri.details.bride_rashi", gd.BrideRashi)
			o.enum("ashtakoota.kootas.graha_maitri.details.groom_lord_relationship_to_bride_lord", gd.GroomLordRelationshipToBrideLord)
			o.enum("ashtakoota.kootas.graha_maitri.details.bride_lord_relationship_to_groom_lord", gd.BrideLordRelationshipToGroomLord)
			o.koota("gana", ganaKoota(g, b))
			o.koota("nadi", nadiKoota(g, b))

			count := nakshatraCount(b.Nakshatra, g.Nakshatra)
			for _, ms := range maitriScores {
				o.koota("bhakoot", bhakootKoota(g, b, ms)) // different scores trigger the cancellation texts
				for _, rel := range lagnaRelations {
					o.check("rajju", rajjuCheck(g, b, ms, rel))
				}
				for _, good := range []bool{false, true} {
					c := streeDeerghaCheck(g, b, count, good, ms)
					o.check("stree_deergha", c)
					o.enum("checks.stree_deergha.details.quality", c.Details.(domain.StreeDeerghaDetails).Quality)
				}
			}
			o.check("vedha", vedhaCheck(g, b))
			o.check("mahendra", mahendraCheck(g, b, count))
			o.check("mruga_vairam", mrugaVairamCheck(g, b))
			o.check("special_tara", specialTaraCheck(g, b))
		}
	}
	for _, gHas := range []bool{false, true} {
		for _, bHas := range []bool{false, true} {
			o.check("kala_sarpa", kalaSarpaResult(gHas, bHas))
		}
	}

	// lagna: the quality and the relationship of the lords
	for d := 1; d <= 12; d++ {
		o.enum("lagna.quality", lagnaQuality(d))
	}
	for _, rel := range lagnaRelations {
		o.enum("lagna.lords_relationship", rel)
	}
	for _, s := range []string{"Aries", "Taurus", "Gemini", "Cancer", "Leo", "Virgo", "Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces"} {
		o.enum("lagna.groom_lagna", s)
		o.enum("lagna.bride_lagna", s)
		o.enum("lagna.groom_lagna_lord", SignLords[s])
		o.enum("lagna.bride_lagna_lord", SignLords[s])
	}
	o.enum("rule_set.ayanamsa", "Lahiri")

	for _, s := range riskFactorStrings(t) {
		o.text("risk factor", s)
	}
	return o
}

// riskFactorStrings reads every string evaluateIndividualRisks can append to a RiskFactors list straight from
// the source code, so a newly added rule cannot go untranslated.
func riskFactorStrings(t *testing.T) []string {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "matchmaking.go", nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	var fn *ast.FuncDecl
	for _, d := range file.Decls {
		if f, ok := d.(*ast.FuncDecl); ok && f.Name.Name == "evaluateIndividualRisks" {
			fn = f
		}
	}
	if fn == nil {
		t.Fatal("evaluateIndividualRisks not found")
	}

	// values the local variable houseLabel takes: the second argument of every checkMarsSat call
	idents := map[string][]string{}
	ast.Inspect(fn, func(n ast.Node) bool {
		if c, ok := n.(*ast.CallExpr); ok {
			if id, ok := c.Fun.(*ast.Ident); ok && id.Name == "checkMarsSat" && len(c.Args) == 2 {
				if lit, ok := c.Args[1].(*ast.BasicLit); ok {
					v, _ := strconv.Unquote(lit.Value)
					idents["houseLabel"] = append(idents["houseLabel"], v)
				}
			}
		}
		return true
	})

	var eval func(e ast.Expr) []string
	eval = func(e ast.Expr) []string {
		switch x := e.(type) {
		case *ast.BasicLit:
			v, _ := strconv.Unquote(x.Value)
			return []string{v}
		case *ast.Ident:
			if vals, ok := idents[x.Name]; ok {
				return vals
			}
		case *ast.BinaryExpr:
			if x.Op == token.ADD {
				var out []string
				for _, l := range eval(x.X) {
					for _, r := range eval(x.Y) {
						out = append(out, l+r)
					}
				}
				return out
			}
		}
		t.Fatalf("cannot evaluate a RiskFactors string expression at %v; extend riskFactorStrings", fset.Position(e.Pos()))
		return nil
	}

	var out []string
	ast.Inspect(fn, func(n ast.Node) bool {
		c, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		if id, ok := c.Fun.(*ast.Ident); !ok || id.Name != "append" || len(c.Args) < 2 {
			return true
		}
		if sel, ok := c.Args[0].(*ast.SelectorExpr); ok && sel.Sel.Name == "RiskFactors" {
			for _, a := range c.Args[1:] {
				out = append(out, eval(a)...)
			}
		}
		return true
	})
	if len(out) < 30 {
		t.Fatalf("found only %d risk factor strings; the parser walk is broken", len(out))
	}
	return out
}

// jsonTags returns the JSON key names of a struct type.
func jsonTags(typ reflect.Type) []string {
	var out []string
	for i := 0; i < typ.NumField(); i++ {
		tag := strings.Split(typ.Field(i).Tag.Get("json"), ",")[0]
		if tag != "" && tag != "-" {
			out = append(out, tag)
		}
	}
	return out
}

// ---- the tests ----

func TestTranslationsCoverEveryOutputString(t *testing.T) {
	tr := loadTranslations(t)
	out := collectOutputs(t)

	usedStrings := map[string]bool{}
	usedTemplates := map[string]bool{}
	var problems []string
	for text, source := range out.texts {
		te, matched := tr.translate(text)
		switch {
		case tr.Strings[text] != "":
			usedStrings[text] = true
		case len(matched) == 0:
			problems = append(problems, fmt.Sprintf("no translation for %q (from %s)", text, source))
		case len(matched) > 1:
			problems = append(problems, fmt.Sprintf("%d templates match %q (from %s): %v", len(matched), text, source, matched))
		default:
			usedTemplates[matched[0]] = true
			if strings.Contains(te, "{") {
				problems = append(problems, fmt.Sprintf("unfilled placeholder in the translation of %q: %s", text, te))
			}
		}
	}
	for k := range tr.Strings {
		if !usedStrings[k] {
			problems = append(problems, fmt.Sprintf("stale entry in strings (the backend no longer emits it): %q", k))
		}
	}
	for k := range tr.Templates {
		if !usedTemplates[k] {
			problems = append(problems, fmt.Sprintf("stale or unmatched template: %q", k))
		}
	}
	if len(problems) > 0 {
		sort.Strings(problems)
		if len(problems) > 25 {
			problems = append(problems[:25], fmt.Sprintf("... and %d more", len(problems)-25))
		}
		t.Errorf("FRONTEND_TRANSLATIONS.json is out of step with the backend:\n  %s", strings.Join(problems, "\n  "))
	}
}

func TestTranslationsCoverEveryEnumeratedValue(t *testing.T) {
	tr := loadTranslations(t)
	out := collectOutputs(t)

	usedPaths := map[string]bool{}
	missing := map[string]bool{}
	for _, u := range out.enums {
		set, ok := tr.EnumPaths[u.path]
		if !ok {
			missing["no enum_paths entry for "+u.path] = true
			continue
		}
		usedPaths[u.path] = true
		if _, ok := tr.Enums[set][u.value]; !ok {
			missing[fmt.Sprintf("%s: value %q is not in enums.%s", u.path, u.value, set)] = true
		}
	}
	for p := range tr.EnumPaths {
		if !usedPaths[p] {
			missing["stale enum_paths entry "+p] = true
		}
	}
	if len(missing) > 0 {
		var list []string
		for m := range missing {
			list = append(list, m)
		}
		sort.Strings(list)
		t.Errorf("enumerated values are not covered:\n  %s", strings.Join(list, "\n  "))
	}
}

func TestTranslationsCoverEveryKey(t *testing.T) {
	tr := loadTranslations(t)

	// every JSON key of the structs the endpoint returns, except the three maps translated separately
	used := map[string]bool{}
	for _, typ := range []reflect.Type{
		reflect.TypeOf(domain.MatchmakingResult{}), reflect.TypeOf(domain.RuleSetInfo{}), reflect.TypeOf(domain.PersonMatchDetails{}),
		reflect.TypeOf(domain.MoonDetails{}), reflect.TypeOf(domain.IndividualRiskCheck{}), reflect.TypeOf(domain.LagnaCompatibility{}),
		reflect.TypeOf(domain.Ashtakoota{}), reflect.TypeOf(domain.Koota{}), reflect.TypeOf(domain.TaraDetails{}),
		reflect.TypeOf(domain.TaraDirection{}), reflect.TypeOf(domain.GrahaMaitriDetails{}), reflect.TypeOf(domain.BhakootDetails{}),
		reflect.TypeOf(domain.MatchSummary{}), reflect.TypeOf(domain.Check{}), reflect.TypeOf(domain.CountDetails{}),
		reflect.TypeOf(domain.StreeDeerghaDetails{}),
	} {
		for _, k := range jsonTags(typ) {
			used[k] = true
			if tr.Fields[k].En == "" || tr.Fields[k].Te == "" {
				t.Errorf("fields is missing %q (en and te)", k)
			}
		}
	}
	for k := range tr.Fields {
		if !used[k] {
			t.Errorf("stale entry in fields: %q is not a key of the response", k)
		}
	}

	for _, c := range []struct {
		name string
		typ  reflect.Type
		have map[string]tText
	}{
		{"kootas", reflect.TypeOf(domain.Kootas{}), tr.Kootas},
		{"checks", reflect.TypeOf(domain.Checks{}), tr.Checks},
		{"risks", reflect.TypeOf(domain.IndividualRisks{}), tr.Risks},
	} {
		want := jsonTags(c.typ)
		for _, k := range want {
			if c.have[k].En == "" || c.have[k].Te == "" {
				t.Errorf("%s is missing %q", c.name, k)
			}
		}
		if len(c.have) != len(want) {
			t.Errorf("%s has %d entries, the response has %d", c.name, len(c.have), len(want))
		}
	}

	if tr.RuleSets[ruleSetName].Te == "" {
		t.Errorf("rule_sets is missing %q", ruleSetName)
	}
}

// The name a koota or check puts in its "name" field equals the English label in the mapping, so
// a frontend that matches by name also works.
func TestTranslationLabelsMatchEmittedNames(t *testing.T) {
	tr := loadTranslations(t)
	out := collectOutputs(t)
	for key, label := range tr.Kootas {
		if out.names[key] != label.En {
			t.Errorf("kootas.%s.en = %q but the backend emits name %q", key, label.En, out.names[key])
		}
	}
	for key, label := range tr.Checks {
		if out.names["check:"+key] != label.En {
			t.Errorf("checks.%s.en = %q but the backend emits name %q", key, label.En, out.names["check:"+key])
		}
	}
}

func TestTranslationEntriesAreComplete(t *testing.T) {
	tr := loadTranslations(t)
	for set, entries := range tr.Enums {
		for k, v := range entries {
			if v.En == "" || v.Te == "" {
				t.Errorf("enums.%s[%q] needs both en and te", set, k)
			}
		}
	}
	for k, te := range tr.Strings {
		if te == "" || strings.Contains(k, "{") {
			t.Errorf("strings[%q]: needs a Telugu text and must not be a template", k)
		}
	}
	for key, tpl := range tr.Templates {
		for _, m := range placeholderRE.FindAllStringSubmatch(key, -1) {
			name := m[1]
			set, declared := tpl.Params[name]
			if !declared {
				t.Errorf("template %q: parameter %s is not declared in params", key, name)
			} else if set != nil {
				if _, ok := tr.Enums[*set]; !ok {
					t.Errorf("template %q: parameter %s uses unknown enum %q", key, name, *set)
				}
			}
			if !strings.Contains(tpl.Te, "{"+name+"}") {
				t.Errorf("template %q: the Telugu text does not use {%s}", key, name)
			}
		}
	}
	for path, set := range tr.EnumPaths {
		if _, ok := tr.Enums[set]; !ok {
			t.Errorf("enum_paths[%q] uses unknown enum %q", path, set)
		}
	}
}
