import urllib.request

urls = [
    "https://raw.githubusercontent.com/vedastro/vedastro/master/VedAstro.Library/Astrology/Tatwa.cs",
    "https://raw.githubusercontent.com/jyotisham/jyotisha/master/jyotisha/panchaanga/temporal/tatwa.py"
]

for url in urls:
    try:
        req = urllib.request.Request(url)
        with urllib.request.urlopen(req) as response:
            content = response.read().decode('utf-8')
            print(f"Found at {url}:\n{content[:1000]}\n")
    except Exception as e:
        print(f"Error fetching {url}: {e}")
