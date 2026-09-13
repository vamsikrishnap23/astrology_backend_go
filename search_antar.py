def check():
    target = {
        306: "Vayu",
        307: "Tejo",
        308: "Vayu",
        309: "Akash",
        310: "Tejo"
    }
    
    tatwas = ["Prithvi", "Jala", "Tejo", "Vayu", "Akash"]
    
    # Hypothesis 1: 1-to-1 row mapping with an offset
    for offset in range(30):
        matches = 0
        for row, expected in target.items():
            cycle_row = ((row - 1) % 30) + 1
            idx = (cycle_row - 1 + offset) % 5
            if tatwas[idx] == expected:
                matches += 1
        if matches == 5:
            print("Found 1-to-1 match with offset:", offset)
            return

    # Hypothesis 2: 1-to-1 with Avaroha
    for offset in range(30):
        matches = 0
        for row, expected in target.items():
            cycle_row = ((row - 1) % 30) + 1
            idx = (offset - (cycle_row - 1)) % 5
            if tatwas[idx] == expected:
                matches += 1
        if matches == 5:
            print("Found 1-to-1 reverse match with offset:", offset)
            return

    # Hypothesis 3: Proportional but with different start points
    # Tejo spans from 918 to 936. Jala spans from 906 to 918.
    # What if lengths are equal?
    # Tejo = 6 rows. Jala = 4 rows.
    # We already checked equal lengths.
    
    print("No simple mathematical cycle matches exactly.")

check()
