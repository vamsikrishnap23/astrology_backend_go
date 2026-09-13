# Frontend Implementation Guide: Birth Time Rectification (BTR)

This document is a comprehensive guide and prompt for building the BTR User Interface. It outlines how to integrate the `/api/btr` backend endpoint across two distinct frontend views.

---

## Overview of the Two Views

You will be building two separate interfaces for the BTR feature:
1. **The Inline Widget (`/chart` page):** A compact, quick-verification tool that displays the summary of the exact birth time, ranked candidate suggestions, and a small snippet of the Nadi table (rows near the inputted time).
2. **The Dedicated BTR Page (`/btr` page):** A full-screen page that renders the entire 480-row daily Nadi Tatwa matrix, allowing power-users to manually scroll through the entire astrological day.

---

## 1. The Inline Widget (For the `/chart` Page)

### Layout & UI Components
This component should be designed as a collapsible card or a modal on the main birth chart page.

**Top Section: Input Controls**
*   `Scan - [MIN]` (Number Input, default: `10`)
*   `to + [MIN]` (Number Input, default: `5`)
*   `Gender` (Dropdown: `Male` | `Female`)
*   `Sign Type` (Dropdown: `Auto-detect` | `Movable` | `Fixed` | `Dual`)
*   `Run Scan` (Primary Button)

**Status Bar**
*   Text string showing the result after clicking scan. E.g., *"Done. sign type: Fixed • nakshatra lord: Mars"*

**Summary Grid**
Display the `input_analysis` data in a clean horizontal flexbox or grid containing:
*   **Weekday**, **Sunrise**, **LMT**, **LMT & Sunrise** (You may need to fetch the raw LMT/Sunrise from the main chart API if needed, or just display the BTR payload).
*   **Star Lord** (From `actual_star_lord`)
*   **Nadi Row** (Highlight this boldly)
*   **Tatwa** (From `calculated_tatwa`)
*   **Planet 90/V** (From `calculated_planet`)
*   **Gender / Star** (Display a ✅/❌ based on `gender_match` and `star_match` booleans).
*   **Ascendant**

**Ranked Candidates Table**
Title: *Ranked candidates (-10/+5 min around your time)*
*   Render the array from `suggested_rectifications`.
*   **Columns:** T1, T2 (clock), Score (Rank), Row, Tatwa, Gen (Male/Female), 90-Min, Vinod, and an action button **[Use]**.
*   Clicking **[Use]** should update the master birth time state in the app and re-trigger the main chart API!

**Nadi Table Snippet**
Title: *Nadi table (rows near your time)*
*   Render a small slice of the `full_table` array. 
*   **Filter Logic:** Find the inputted time's `nadi_row` index. Slice the array to show only ±4 rows around it (e.g., if row is 359, slice from 355 to 363).
*   Highlight the exact row corresponding to the current time and label the action button **[Current]**. Other rows get a **[Use]** button.
*   *Columns to map:* NO, T1, T2, Wed, MonFri, SunTues, Sat, Thur, Movable, Fixed, Dual, M V, F V, D V.

**Bottom Footer Link**
*   Link: *"Open full BTR tool →"* (Routes to the dedicated BTR page).

---

## 2. The Dedicated BTR Page (`/btr` Page)

### Layout & UI Components
This is a full-screen, heavy-data view designed for astrologers who want to manually verify the entire day.

**Top Section**
*   Same input controls as the inline widget.

**The Giant Matrix**
*   This page requires rendering the **entire** 480-element array returned in `full_table`.
*   **Performance Note:** Because rendering 480 rows with 14 columns each can cause DOM lag in React/Vue, consider using a **Virtual List / Windowed List** (like `react-window` or `react-virtualized`) to render the table smoothly.
*   The columns are identical to the Nadi Table Snippet above.

---

## API Integration Details

### API Endpoint
`POST /api/btr`

### Payload Construction
When the user clicks "Run Scan", construct this payload from the app's state:

```javascript
const payload = {
  // Inherited from master chart state
  date_of_birth: appState.date,
  time_of_birth: appState.time,
  latitude: appState.lat,
  longitude: appState.lng,
  timezone: appState.tz,
  ayanamsa: appState.ayanamsa,
  
  // From the BTR inputs
  gender: btrState.gender,
  scan_minus_minutes: parseInt(btrState.scanMinus),
  scan_plus_minutes: parseInt(btrState.scanPlus),
  sign_type_override: btrState.signType === "Auto-detect" ? "" : btrState.signType,
  
  // IMPORTANT: Set to true so we get the full matrix for the table snippet!
  return_full_table: true 
};
```

### JSON Response Mapping Guide

*   `input_analysis.gender_match` ➔ Drives the ✅/❌ in the Gender column.
*   `input_analysis.star_match` ➔ Drives the ✅/❌ in the Star column.
*   `suggested_rectifications` ➔ Map over this array for the **Ranked Candidates** table.
*   `full_table` ➔ Array of 480 objects. 
    *   For the `/chart` widget: `full_table.slice(targetIndex - 4, targetIndex + 5)`
    *   For the `/btr` page: map the entire array into the virtualized table.

---

## Developer Instructions (Prompt)
If you are passing this to an AI coding assistant, use this exact prompt block:

> "Build a React component for a Birth Time Rectification (BTR) tool. It should have two views: an inline widget, and a dedicated full-page route. The API payload must include `scan_minus_minutes`, `scan_plus_minutes`, `sign_type_override`, and `return_full_table: true`. 
> 
> The response contains `input_analysis` (for the summary grid), `suggested_rectifications` (for the ranked candidates list), and `full_table` (an array of 480 rows). 
> 
> In the inline widget, slice the `full_table` to only show the 9 rows surrounding the current input time. Include a 'Use' button on the rows that updates the global birth time state. Implement virtual scrolling for the dedicated full-page table to handle the 480 rows smoothly."
