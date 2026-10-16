# Exhaustive List of Output Strings (Matchmaking API)

This document lists **each and every possible English string** that the `/api/vivaha-pontana` backend can output. 

---

## 1. Koota Explanations
*These are the baseline explanation strings for the standard 8 Kootas.*

* **Varna:** `Groom varna {varna} (Rank {rank}) vs Bride varna {varna} (Rank {rank})`
* **Vashya:** `Bride's rasi {bride} is vashya to groom's rasi {groom}.` / `Groom's rasi {groom} is vashya to bride's rasi {bride}.` / `No vashya relation is listed between {groom} and {bride}.`
* **Tara:** `Calculated cyclic distances between birth nakshatras.`
* **Yoni:** `Groom yoni {animal} vs Bride yoni {animal}`
* **Graha Maitri:** `Natural relationship between moon sign lords.`
* **Gana:** `Groom gana {gana} vs Bride gana {gana}`
* **Nadi:** `Nadi mapping.`
* **Bhakoot (Base):** `Sign distance between Moons.`
  * *Dynamic Warnings appended if applicable:*
  * ` WARNING: Groom's Moon is 2nd from Bride's Moon (Risk of Vaidhavya/Poverty per Telugu rules).`
  * ` WARNING: Groom's Moon is 5th from Bride's Moon (Risk of Vaidhavya/Loss of children per Telugu rules).`
  * ` WARNING: Groom's Moon is 3rd from Bride's Moon (Risk of Sorrows/Hardships).`
  * ` WARNING: Groom's Moon is 6th from Bride's Moon (Risk of Loss of children/Quarrels).`
* **Rajju:** `Groom Rajju is {rajju}. Bride Rajju is {rajju}.`
  * *Appended if clashing:* ` FATAL DOSHA: Both belong to the same Rajju. This is highly inauspicious.`

---

## 2. Cancellation Reasons (Pariharam)
*Sent in the `cancellation_reason` field when a classical exemption saves a dosha.*

* **Rajju:** `Rajju dosha cancelled because Rasi Lords are same or friendly.`
* **Rajju:** `Rajju dosha cancelled because Lagna Lords are same or friendly.`
* **Nadi:** `Nadi dosha cancelled due to same Nakshatra but different Padas.`
* **Nadi:** `Nadi dosha cancelled due to same Rashi but different Nakshatras.`
* **Bhakoot:** `Bhakoot dosha cancelled because Rashi lords are same or friends (Graha Maitri).`
* **Bhakoot:** `Bhakoot dosha cancelled because this specific sign combination is highly auspicious.`

---

## 3. Check Explanations (Kala Sarpa, Mruga Vairam, Special Tara)
*Sent in the `explanation` of the `checks` entries (`kala_sarpa`, `mruga_vairam`, `special_tara`).*

* **Kala Sarpa (Default):** `Checks if both Bride and Groom have Kala Sarpa Dosha. (If both have it, the match is rejected).`
* **Kala Sarpa (Triggered):** `FATAL DOSHA: Both the Groom and the Bride suffer from Kala Sarpa Dosha in their natal charts. The match is explicitly rejected.`
* **Mruga Vairam (Default):** `Checks if the Yoni (Animal) of the couple are natural enemies.`
* **Mruga Vairam (Triggered):** `MRUGA VAIRAM DOSHA: Groom's Yoni animal is {animal} and Bride's Yoni animal is {animal}. These are natural enemies, indicating extreme incompatibility.`
* **Special Tara (Default):** `Checks for specifically rejected Tara combinations (Ashubha Kurpu).`
* **Special Tara (Triggered):** `ASHUBHA KURPU DOSHA: The pairing of Bride's {star} and Groom's {star} is explicitly forbidden in the Telugu matchmaking text.`

---

## 4. Individual Risk Factors (`risk_factors` array)
*These strings are generated for individual natal charts (both bride and groom) if specific planetary alignments occur.*

### Vaidhavya Dosham (Widowhood / Longevity)
* `Saturn is in the 8th house (risk of Alpayush/Vaidhavya).`
* `Sun is in the 8th house.`
* `Mars and Rahu conjunct in the 8th house.`
* `Exchange (Parivartana) of 7th and 8th lords.`
* `Lagnadhipati (Lagna Lord) is in the 6th, 8th, or 12th house.`
* `Ashtamadhipati (8th Lord) is in an inauspicious house.`
* `7th and 8th lords are conjunct in a Kendra (1,4,7,10).`
* `7th and 8th lords are conjunct in the 5th house.`
* `Mars, Venus, and Rahu conjunction.`
* `Ketu in Aquarius (Madhyamayuvu).`
* `Mars in the 9th house (Alpayush).`
* `Saturn in the 9th house.`

### Dwikalatra Dosham (Second Marriage)
* `Saptamadhipati (7th Lord) is in the 1st, 2nd, 6th, or 8th house.`
* `Mercury and Rahu conjunct in the 2nd house (Kutumba Nasanam).`
* `Lagna Lord, 7th Lord, and Venus are conjunct together in a Dual Sign.`
* `Mars and Venus conjunction outside of own signs.`
* `Exchange (Parivartana) between 2nd and 12th lords.`
* `Mercury and Saturn conjunct in 7th house (Groom chart).`
* `Mars and Venus conjunct in 7th house (Groom chart).`
* `Saturn in the 10th house.`
* `Mars in the 10th house.`
* `Lagna lord placed in the 3rd or 6th house.`

### Napumsakatvam (Impotency - Groom Only)
* `Saturn in 2nd/6th/12th in a water sign with malefic in Lagna.`

### Dampatya Marakatva Dosham (Lethal Marital Flaws)
* `Moon and Venus conjunct with Mars and Saturn in opposition.`
* `5th and 7th lords connection/exchange.`
* `7th and 8th lords placement/exchange.`
* `Saturn and Venus conjunct outside Libra.`
* `Mars and Saturn conjunct in Lagna.`
* `Mars and Saturn conjunct in 7th house.`
* `7th house is in Papa Kartari (hemmed between malefics).`
* `7th Lord is in Papa Kartari (hemmed between malefics).`
* `Venus is in Papa Kartari (hemmed between malefics).`
* `Taurus Lagna with Venus in 7th (Scorpio).`
* `Scorpio Lagna with Mercury in 7th (Taurus).`
* `Cancer Lagna with Jupiter in 7th (Capricorn).`
* `Virgo Lagna with Saturn in 7th and Sun in 1st.`

---

## 4. Nadi suggestion and Lagna quality
* **Nadi `suggestion`** (only when both nakshatras are in the Jaimini Maharshi list; the score is not changed): `Per the Jaimini Maharshi exemption list, Nadi dosha does not apply when both nakshatras ({groom} and {bride}) are in that list. This is a suggestion only and is not counted in the score.`
* **Lagna `quality`:** `Auspicious`, `Moderate`, `Inauspicious`.

---

## 5. Other Check Explanations (`checks` block)
*Sent in the `explanation` fields of `checks`; none of these add to the 36 points.*

* **Kala Sarpa `groom`/`bride`:** `Present` / `Absent`.
* **Gana `suggestion`** (only when the gana is mismatched and the count is more than 14; the score is not changed): `The count from the bride's nakshatra to the groom's is {n} (more than 14), which greatly reduces the harm of a gana mismatch (Telugu muhurta book, p. 65). This is a suggestion only and is not counted in the score.`
* **Mahendra:** `Groom's nakshatra is the {n}th from the bride's, one of 4, 7, 10, 13, 16, 19, 22, 25: auspicious, increases married life.` / `Groom's nakshatra is the {n}th from the bride's, not one of 4, 7, 10, 13, 16, 19, 22, 25.`
* **Stree Deergha:** `Groom's nakshatra is the {n}th from the bride's, in navaka {k}: {quality}.` with `quality` = `Kalaham (quarrel)` / `Samanyam (average)` / `Shubham (auspicious)`; appended when the exception applies: ` Rasi koota and graha maitri are good, so the shortfall can be tolerated.`
* **Vedha:** `{groom} and {bride} obstruct each other (vedha); marriage should proceed only with remedial indications.` / `The two nakshatras do not obstruct each other (no vedha).`
