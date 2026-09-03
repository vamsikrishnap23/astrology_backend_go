# Feature 23: Birth Time Rectification (Nadi Tatwa)

## Overview
This module attempts to implement the Nadi Tatwa / Tatwa Shodhana method of Birth Time Rectification. The primary goal is to use precise astronomical state boundaries to determine if a reported birth time correlates properly with classical elemental (Tatwa) rules.

## Process
The pipeline performs a robust set of astronomical evaluations using the `go-swisseph` engine:
1. **LMT and Sunrise:** Converts the reported time into UTC and Local Mean Time (LMT), evaluating exact Sunrise bounds utilizing precise geometrical latitudes/longitudes.
2. **Moon & Star Lord:** Uses sidereal lunar longitude to calculate the Nakshatra and corresponding planetary Star Lord.
3. **Ascendant Modality:** Maps the exact sidereal Ascendant into Movable, Fixed, or Dual modalities.
4. **90-Minute Ruler & Tatwa:** Slices the time elapsed since Sunrise into 90-minute blocks, mapping the Aroha (direct) and Avaroha (reverse) progression of Prithvi, Jala, Tejo, Vayu, and Akash, scaling them proportionally down to 12-second Antar-Tatwas.

## Implementation Status
**Status**: `partially_implemented`

### The Missing Component
The module currently stops at generating candidates with fully loaded states (`90-minute Planet`, `Ascendant Modality`, `Star Lord`, `Gender`, `Tatwa/Antar-Tatwa`).

It **cannot** perform the final `match` validation. 
The Nadi Tatwa BTR system depends on a proprietary, classical lookup table (typically 500+ rows published in closed literature by researchers like D.V. Subbarao). This table explicitly maps the four variables:
`[ 90-min Planet ] + [ Ascendant Type ] + [ Star Lord ] + [ Gender ]`
to a specific, predetermined **Expected Tatwa** or **Antar-Tatwa**.

### Strict Adherence to Accuracy
Per rigorous engineering standards, this module **does not** hallucinate, guess, or synthesize this table. It explicitly tags every candidate evaluation with `match_status: "unresolved"` instead of arbitrarily declaring `false` or inventing numeric weights.

To complete this module, the exact authoritative lookup table dataset must be provided and securely injected into the engine.

## API Endpoint
`POST /api/v1/btr`
Returns the complete `astronomical_context` and a granular scan of boundaries surrounding the given birth time, flagging all as `unresolved`.
