# Matchmaking API (Ashtakoota / Guna Milan)

**Endpoint:** `POST /api/vivaha-pontana`

Compares a groom and a bride from their Moon positions and returns the eight scored kootas (36 points) plus a set of other compatibility checks. All calculations use the **Lahiri ayanamsa**; the `ayanamsa` sent in the request is ignored.

## Request
Two standard `BirthInput` objects under `groom` and `bride`.

```json
{
  "groom": {
    "name": "Groom", "date_of_birth": "2005-11-23", "time_of_birth": "15:35:00",
    "place_of_birth": "Sattenapalle", "latitude": 16.3938, "longitude": 80.1522, "timezone": 5.5
  },
  "bride": {
    "name": "Bride", "date_of_birth": "2002-05-14", "time_of_birth": "10:00:00",
    "place_of_birth": "Bangalore", "latitude": 12.9716, "longitude": 77.5946, "timezone": 5.5
  }
}
```

## Response layout
The response has six sections, always in this order:

| Section | What it holds |
|---|---|
| `rule_set` | `name` (`classical-guna-milan-v2`) and `ayanamsa` (always `Lahiri`) |
| `groom`, `bride` | `name`, `moon` (longitude, sign, degree, nakshatra, pada, lords) and `risks` (four individual chart risk checks) |
| `lagna` | ascendant compatibility (informational, not scored) |
| `ashtakoota` | `kootas` (the eight scored kootas) and `summary` (the totals) |
| `checks` | the other compatibility checks, none of which add to the score |

### A koota (`ashtakoota.kootas.*`)
All eight kootas (`varna`, `vashya`, `tara`, `yoni`, `graha_maitri`, `gana`, `bhakoot`, `nadi`) have exactly this shape:

| Field | Meaning |
|---|---|
| `name` | display name |
| `maximum` | most points this koota can give (1, 2, 3, 4, 5, 6, 7, 8) |
| `raw_score` | score before any cancellation |
| `score` | the score **counted in the total** (equals `raw_score` unless a cancellation applied) |
| `groom`, `bride` | what each person was classified as: varna, rasi, nakshatra, yoni animal, rashi lord, gana, nadi |
| `dosha` | true only for koota with a named dosha (`bhakoot`, `nadi`) when it is present |
| `cancellation_applied`, `cancellation_reason` | whether a dosha was cancelled, and why |
| `suggestion` | (only when present) a rule that may apply but is deliberately **not** counted in the score |
| `explanation` | short text |
| `details` | (only for some) extra fields: `tara` (both directional counts), `graha_maitri` (rashis and the two lords' relationships), `bhakoot` (the two distances) |

### A check (`checks.*`)
`rajju`, `vedha`, `mahendra`, `stree_deergha`, `mruga_vairam`, `special_tara`, `kala_sarpa` all have this shape:

| Field | Meaning |
|---|---|
| `name` | display name |
| `satisfied` | **true when there is no outstanding problem**: the dosha is absent, or a cancellation was applied (`cancellation_applied`, with the reason in `cancellation_reason`). Used by `rajju` and `stree_deergha` |
| `groom`, `bride` | what each person was classified as (rajju group, nakshatra, yoni animal; `Present`/`Absent` for kala sarpa) |
| `cancellation_applied`, `cancellation_reason` | only when a dosha was cancelled (rajju) |
| `explanation` | short text |
| `details` | only for `mahendra` (`bride_to_groom_count`) and `stree_deergha` (count, navaka, quality, alternate opinion, exception) |

## Example (abridged, from the real output for the request above)
```json
{
  "rule_set": {
    "name": "classical-guna-milan-v2",
    "ayanamsa": "Lahiri"
  },
  "groom": {
    "name": "Groom",
    "moon": {
      "longitude": 121.81859918421418,
      "sign": "Leo",
      "degree": 1.8185991842141789,
      "nakshatra": "Magha",
      "pada": 1,
      "nakshatra_lord": "Ketu",
      "rashi_lord": "Sun"
    },
    "risks": "{ vaidhavya_dosham, dwikalatra_dosham, napumsakatvam, dampatya_marakatva: each { has_risk, risk_factors[] } }"
  },
  "bride": "{ same shape as groom }",
  "lagna": {
    "groom_lagna": "Aries",
    "bride_lagna": "Gemini",
    "groom_lagna_lord": "Mars",
    "bride_lagna_lord": "Mercury",
    "distance": 11,
    "quality": "Auspicious",
    "is_auspicious": true,
    "lords_relationship": "Inimical"
  },
  "ashtakoota": {
    "kootas": {
      "varna": {
        "name": "Varna",
        "maximum": 1,
        "raw_score": 1,
        "score": 1,
        "groom": "Kshatriya",
        "bride": "Vaishya",
        "dosha": false,
        "cancellation_applied": false,
        "explanation": "Groom varna Kshatriya (Rank 3) vs Bride varna Vaishya (Rank 2)"
      },
      "tara": {
        "name": "Tara",
        "maximum": 3,
        "raw_score": 1.5,
        "score": 1.5,
        "groom": "Magha",
        "bride": "Rohini",
        "dosha": false,
        "cancellation_applied": false,
        "explanation": "Calculated cyclic distances between birth nakshatras.",
        "details": {
          "groom_to_bride": {
            "nakshatra_count": 22,
            "tara": "Kshema",
            "auspicious": true,
            "score": 1.5
          },
          "bride_to_groom": {
            "nakshatra_count": 7,
            "tara": "Vadha",
            "auspicious": false,
            "score": 0
          }
        }
      },
      "bhakoot": {
        "name": "Bhakoot",
        "maximum": 7,
        "raw_score": 7,
        "score": 7,
        "groom": "Leo",
        "bride": "Taurus",
        "dosha": false,
        "cancellation_applied": false,
        "explanation": "Sign distance between Moons.",
        "details": {
          "groom_to_bride_distance": 10,
          "bride_to_groom_distance": 4,
          "relationship": "10/4"
        }
      },
      "nadi": {
        "name": "Nadi",
        "maximum": 8,
        "raw_score": 0,
        "score": 0,
        "groom": "Antya",
        "bride": "Antya",
        "dosha": true,
        "cancellation_applied": false,
        "suggestion": "Per the Jaimini Maharshi exemption list, Nadi dosha does not apply when both nakshatras (Magha and Rohini) are in that list. This is a suggestion only and is not counted in the score.",
        "explanation": "Nadi mapping."
      },
      "vashya, yoni, graha_maitri, gana": "{ same shape }"
    },
    "summary": {
      "raw_total": 10.5,
      "effective_total": 10.5,
      "maximum": 36,
      "percentage": 29.166666666666668,
      "traditional_threshold_met": false
    }
  },
  "checks": {
    "rajju": {
      "name": "Rajju",
      "satisfied": true,
      "groom": "Pada",
      "bride": "Kantha",
      "explanation": "Groom Rajju is Pada. Bride Rajju is Kantha."
    },
    "stree_deergha": {
      "name": "Stree Deergha",
      "satisfied": false,
      "groom": "Magha",
      "bride": "Rohini",
      "explanation": "Groom's nakshatra is the 7th from the bride's, in navaka 1: Kalaham (quarrel).",
      "details": {
        "bride_to_groom_count": 7,
        "navaka": 1,
        "quality": "Kalaham (quarrel)",
        "satisfied_by_alternate_opinion": false,
        "exception_applies": false
      }
    },
    "vedha, mahendra, mruga_vairam, special_tara, kala_sarpa": "{ same shape }"
  }
}
```

## Rules
- **Varna (1):** groom's varna equal to or above the bride's.
- **Vashya (2):** the Telugu book's rasi-pair table (pp. 69-70): 2 if either rasi lists the other as vashya, else 0.
- **Tara (3):** nakshatras counted both ways; each favourable direction gives 1.5 (see `details`).
- **Yoni (4):** animal compatibility of the birth nakshatras.
- **Graha Maitri (5):** natural friendship of the two Moon-sign lords.
- **Gana (6):** published table read as [bride gana][groom gana]. If the gana is mismatched but the count from the bride's nakshatra to the groom's is more than 14, a `suggestion` says the harm is greatly reduced (book p. 65); the score is unchanged.
- **Bhakoot (7):** Moon-sign distances 2/12, 5/9, 6/8 are a dosha (score 0) unless cancelled by an auspicious sign pair or Graha Maitri of 4 or more.
- **Nadi (8):** same nadi scores 0; cancelled only for same nakshatra with a different pada, or same rasi with a different nakshatra. If both nakshatras are in the Jaimini Maharshi list of 21, a `suggestion` states the rule and the score is **not** changed.
- **Summary:** `raw_total` and `effective_total` are the sums of `raw_score` and `score`; `traditional_threshold_met` is `effective_total >= 18`.
- **Lagna:** `distance` is counted from the bride's lagna to the groom's; `quality` is Auspicious (1, 3, 5, 9, 11), Moderate (2, 7) or Inauspicious (4, 6, 8, 10, 12).
- **Checks:** `mahendra` (count 4, 7, 10, 13, 16, 19, 22 or 25 from the bride's nakshatra), `stree_deergha` (count above 9 in three navakas of 9; some opinions above 7; when the count falls short but Bhakoot and Graha Maitri are good, the shortfall is tolerated: `satisfied` is true with `cancellation_applied` and the reason), `vedha` (13 obstructing nakshatra pairs), `rajju` (same rajju group is a dosha, with cancellations), `mruga_vairam` (enemy yoni animals), `special_tara` (rejected nakshatra pairings), `kala_sarpa` (satisfied unless both persons have it).

## Translating the response (Telugu)
`FRONTEND_TRANSLATIONS.json` (repository root) translates **everything** this endpoint can output. It holds:

| Key | Use |
|---|---|
| `fields` | label for every JSON key name in the response |
| `kootas`, `checks`, `risks` | labels for the 8 koota keys, the 7 check keys and the 4 risk keys (their English label equals the `name` the backend sends) |
| `enums` | Telugu for every enumerated value: `signs`, `nakshatras`, `planets`, `yoni_animals`, `ganas`, `varnas`, `nadis`, `rajju_groups`, `taras`, `lagna_quality`, `lords_relationships`, `maitri_relationships`, `stree_deergha_quality`, `present_absent`, `ayanamsas` |
| `enum_paths` | which enum set translates the value found at a JSON path (`*` = groom or bride), e.g. `ashtakoota.kootas.gana.groom` -> `ganas` |
| `strings` | every fixed text (explanations, cancellation reasons, risk factors), keyed by the exact English text |
| `templates` | texts with values in them, keyed by the English text with `{param}` placeholders; `params` names the enum that translates each value (`null` = a number, copied as is) |
| `rule_sets` | label for `rule_set.name` |

How to translate a response:
1. **Key names**: `fields[key]`; under `ashtakoota.kootas` use `kootas[key]`, under `checks` use `checks[key]`, under `groom.risks` / `bride.risks` use `risks[key]`.
2. **Enumerated values**: find the response path in `enum_paths`, then take `enums[set][value].te`.
3. **Free text** (`explanation`, `cancellation_reason`, `suggestion`, each item of `risk_factors`): look the exact text up in `strings`; otherwise match it against `templates` as below.
4. Numbers, scores, booleans, `relationship` (e.g. `10/4`) and person names are shown as they are.

```js
// translateText(englishText) -> Telugu text, or undefined if the text is not covered
const escapeRe = s => s.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
const templates = Object.entries(tr.templates).map(([key, t]) => {
  const names = [];
  const src = escapeRe(key).replace(/\\\{(\w+)\\\}/g, (_, n) => { names.push(n); return '(.+?)'; });
  return { re: new RegExp('^' + src + '$'), names, t };
});
const enumTe = (set, v) => tr.enums[set]?.[v]?.te;
function translateText(s) {
  if (tr.strings[s] !== undefined) return tr.strings[s];
  for (const { re, names, t } of templates) {
    const m = re.exec(s);
    if (!m) continue;
    const vals = {};
    let ok = true;
    names.forEach((n, i) => {
      const set = t.params[n];
      const te = set ? enumTe(set, m[i + 1]) : m[i + 1];
      if (te === undefined) ok = false; else vals[n] = te;
    });
    if (ok) return t.te.replace(/\{(\w+)\}/g, (_, n) => vals[n]);
  }
  return undefined;
}
const pathRe = p => new RegExp('^' + p.split('.').map(seg => seg === '*' ? '[^.]+' : escapeRe(seg)).join('\\.') + '$');
const enumPaths = Object.entries(tr.enum_paths).map(([p, set]) => ({ re: pathRe(p), set }));
const enumFor = path => enumPaths.find(e => e.re.test(path))?.set;
```

**Keeping it complete.** `TestTranslationsCoverEveryOutputString`, `TestTranslationsCoverEveryEnumeratedValue`, `TestTranslationsCoverEveryKey` and `TestTranslationLabelsMatchEmittedNames` (`internal/astrology/matchmaking/translations_test.go`) enumerate every Moon position (27 nakshatras x 4 padas, 11,664 couples), every koota and check including the cancellation variants, and the risk-factor strings read from the source code. They fail when the backend can output something the file does not cover, when a template is ambiguous, or when the file holds an entry the backend no longer produces. After changing any output text or adding a koota, check or risk rule, update `FRONTEND_TRANSLATIONS.json` until these tests pass.

Code: `internal/astrology/matchmaking/` (`kootas.go`, `checks.go`, `lagna.go`, `matchmaking.go`, `tables.go`); types in `internal/domain/matchmaking.go`.
