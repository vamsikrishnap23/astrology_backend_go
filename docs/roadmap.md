# Astrology Backend Roadmap & Future Plans

This document outlines the planned architectural additions and feature expansions for the astrology engine. The following three major engines are prioritized for future development.

## 1. Planetary Yogas Engine
A classical rule-based detection engine that scans a calculated Natal Chart to identify auspicious and inauspicious planetary combinations (Yogas).

### Core Features
*   **Rule Evaluation Engine:** An internal logic system that cross-references the `PlanetaryTable` and `HouseTable` to verify complex conditions (e.g., planetary conjunctions, mutual aspects, house lordships, and exaltation/debilitation states).
*   **Classical Yogas:** Detection of major Yogas such as:
    *   *Pancha Mahapurusha Yogas* (Ruchaka, Bhadra, Hansa, Malavya, Sasa)
    *   *Lunar Yogas* (Gajakesari, Sunapha, Anapha, Kemadruma)
    *   *Wealth & Power Yogas* (Dhana Yogas, Raja Yogas, Vipareeta Raja Yoga)
    *   *Kalasarpa Dosha* and its variations.
*   **API Output:** A new endpoint (`/api/yogas`) that returns a list of detected yogas, the specific planets causing them, and a standard description of their effects.

## 2. KP Horary (Prashna) Engine
An extension of the existing Krishnamurti Paddhati (KP) module designed to answer specific questions based on the time of the query and a user-provided seed number.

### Core Features
*   **1-249 Sub-Lord System:** The engine will accept a seed number (1 to 249) provided by the user.
*   **Ascendant Fixing:** Instead of calculating the Ascendant from the current time, the backend will mathematically reverse-engineer and fix the Ascendant to the exact Zodiac degree, Minute, and Second corresponding to the KP Sub-Lord of the provided seed number.
*   **Horary Chart Generation:** Once the Ascendant is fixed, the engine will cast the rest of the chart (Planets and House Cusps) for the exact moment and location of the query.
*   **Ruling Planets Snapshot:** Automatically capture and return the Ruling Planets (Day Lord, Moon Sign Lord, Moon Star Lord, Ascendant Sign Lord, Ascendant Star Lord) at the time of query to assist the astrologer in predicting the outcome.

## 3. Muhurtha (Electional Astrology) Engine
An automated time-search engine that scans forward in time to find astrologically auspicious windows for specific events (marriage, business, travel, etc.).

### Core Features
*   **Time-Range Searching:** Accepts a `start_date`, `end_date`, location, and the user's Natal Moon details (Nakshatra and Rasi).
*   **Panchang Filtering:** Automatically filters out negative daily periods based on the newly implemented Panchang engine:
    *   Eliminates *Durmuhurtam*, *Varjyam* (Tyajyam), *Rahu Kalam*, and *Yamaganda*.
*   **Personalized Bala (Strength):** Evaluates *Tarabala* (Star strength) and *Chandrabala* (Moon strength) against the user's natal placements to ensure the day is personally favorable.
*   **Lagna Shuddhi (Ascendant Purity):** Evaluates the rising sign at various times during the day to ensure benefics are in Kendras (1,4,7,10) and malefics are in Upachayas (3,6,11) for the specific activity.
*   **API Output:** Returns a ranked list of precise `start_time` and `end_time` windows that pass the astrological criteria.
