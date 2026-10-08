# Exhaustive List of Output Strings (Matchmaking API)

This document lists **each and every possible English string** that the `/api/vivaha-pontana` backend can output. 

---

## 1. Koota Explanations
*These are the baseline explanation strings for the standard 8 Kootas.*

* **Varna:** `Groom varna {varna} (Rank {rank}) vs Bride varna {varna} (Rank {rank})`
* **Vashya:** `Groom vashya {vashya} vs Bride vashya {vashya}`
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
* **Nadi:** `Nadi dosha cancelled per Jaimini Maharshi exemptions for these specific Nakshatras.`
* **Nadi:** `Nadi dosha cancelled due to same Nakshatra but different Padas.`
* **Nadi:** `Nadi dosha cancelled due to same Rashi but different Nakshatras.`
* **Bhakoot:** `Bhakoot dosha cancelled because Rashi lords are same or friends (Graha Maitri).`
* **Bhakoot:** `Bhakoot dosha cancelled because this specific sign combination is highly auspicious.`

---

## 3. Mutual Fatal Dosha Explanations
*Sent in the `doshas` block.*

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
