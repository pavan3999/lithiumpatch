# lithiumpatch

Adds additional functionality to the Lithium EPUB Reader Android app.

## Features

- Requires Android Oreo (8) or higher.
- Latest WebView (at least ~73) required for full functionality.
- Custom icon color.
- Monochrome adaptive icon support.
- Dynamic grid view cover width and custom aspect ratio.
- Optional cover-only grid view.
- Debuggable reader webview.
- Dictionary (works offline).
- Custom fonts.
- Additional font script support (e.g., Thai).
- Smaller minimum font size.
- Additional information in the reader footer.
- Series metadata support.
- Series section in library drawer.
- Increased number of visible actions for the reader toolbar.
- Support for inverted portrait/landscape rotation.
- Expand display settings popup by default.
- Support for hyphenation.
- Additional built-in themes.
- Material You colors on Android 12+.
- Full-bleed background in fullscreen mode on devices with a notch.
- Option to disable page turn animations (e.g., for e-ink screens).
- Option to invert the color of images or the entire page.

## Usage

1. Install JRE 1.8 or newer.
2. Install Go 1.25 or newer.
3. Install zipalign (part of the Android build tools).
4. Optionally run `go generate ./dict/edgedict` to download additional dictionaries.
5. Optionally download additional fonts into the `fonts` directory to add additional fonts (to limit them to a single language, put them in a subdirectory named `latin`/`cyrillic`/`greek`/`thai`).
6. Run `go generate ./app` from the root of the repository to download the APK. If this does not work, you can manually download the Lithium 0.24.5 APK from [here](https://www.apkmirror.com/apk/faultexception/lithium-epub-reader/lithium-epub-reader-0-24-5-release/lithium-epub-reader-0-24-5-android-apk-download/) or extract it from your device.
7. Run `go run . app/Lithium_0.24.5.apk` from the root of the repository. Use `--help` to see additional options including using a custom keystore, setting the tool paths, and adding fonts from an external directory.
8. For Google Drive support, specify a custom keystore with `--keystore whatever.jks`, and create a new Google APIs project with access to the Drive API for the signing key's signature to enable sync.

```
usage: lithiumpatch [options] APK_PATH

options:
  -k, --keystore string              Path to keystore for signing (will be created if does not exist) (default "keystore.jks")
      --keystore-alias string        Keystore alias (default "default")
      --keystore-passphrase string   Keystore passphrase (default "default")
  -o, --output string                Output APK path (default: {basename}.patched.resigned.apk)
  -d, --diff string                  Write diff to the specified file (default: disabled)
      --add-fonts strings            Add extra TTF fonts from a directory (Regular/Roman, Bold, Italic, and BoldItalic variants should be provided) (can be specified multiple times)
      --apktool string               Path to apktool.jar (2.8.1) (default "lib/apktool-2.8.1.jar")
      --apksigner string             Path to apksigner.jar (0.9 or later) (default "lib/apksigner-0.9.jar")
      --zipalign string              zipalign executable (will search PATH) (default "zipalign")
      --keytool string               keytool executable (will search PATH) (default "keytool")
  -q, --quiet                        Do not show the diff
      --help                         Show this help text
```

**Note:** If you get an error from apktool about `No resource identifier found for attribute 'preserveLegacyExternalStorage'`, run `java -jar lib/apktool-2.8.1.jar empty-framework-dir`.


## Custom Dictionaries

`lithiumpatch` supports adding your own custom dictionaries on top of the built-in ones. This lets you use any dictionary you own or find — bilingual dictionaries, specialized/technical glossaries, fictional-language dictionaries, whatever you want — as long as you convert it into the required JSON format below.

### 1. Required JSON format

Each custom dictionary is a single JSON file containing an array of entries. Each entry looks like this:

\`\`\`json
[
  {
    "Terms": ["word", "alternate form"],
    "Name": "word",
    "Pronunciation": "/wɜːrd/",
    "MeaningGroups": [
      {
        "Info": ["noun"],
        "Meanings": [
          {
            "Tags": ["informal"],
            "Text": "A unit of language.",
            "Examples": ["He said a word."]
          }
        ],
        "WordVariants": []
      }
    ],
    "Info": "",
    "Source": "My Custom Dictionary"
  }
]
\`\`\`

Field notes:
| Field | Required? | Description |
|---|---|---|
| `Terms` | Yes | All lookup keys/spellings that should match this entry. |
| `Name` | Yes | The headword shown at the top of the entry. |
| `Pronunciation` | No | IPA or phonetic spelling. |
| `MeaningGroups` | Yes | One or more groups (e.g. for different parts of speech). |
| `MeaningGroups[].Info` | No | Labels like part of speech, word forms. |
| `MeaningGroups[].Meanings[].Text` | Yes | The actual definition. **Must be plain text — no HTML tags or HTML entities** (e.g. `&mdash;`, `&#x0101;`, `&bull;` must be decoded/removed first). |
| `MeaningGroups[].Meanings[].Tags` | No | Short labels for that specific meaning. |
| `MeaningGroups[].Meanings[].Examples` | No | Example sentences. |
| `Source` | No but recommended | Name of the original dictionary, shown to the user. |

This mirrors the `dict.Entry` Go struct in `dict/dict.go` — if you want full field-level detail, that file is the source of truth.

### 2. Converting your dictionary into this format

Your source dictionary could be in any format (StarDict, CSV, plain text word lists, another app's export, etc.). You'll need a small script to turn it into the JSON structure above. General approach:

1. Parse your source file however it's structured.
2. For each entry, build a dict with `Terms`, `Name`, `MeaningGroups`, and `Source`.
3. Strip/decode any markup in the definition text (HTML tags, HTML entities, footnote markers like `[1]`).
4. Write the whole list out as one JSON array with `json.dump(entries, f, ensure_ascii=False)`.

**Example: converting a StarDict dictionary (.ifo/.idx/.dict)**

StarDict dictionaries store a metadata file (`.ifo`), a word index (`.idx`), and the actual definitions (`.dict`). A conversion script needs to:
- Read the `.ifo` file to get `idxoffsetbits` (usually 32).
- Walk the `.idx` file to get each word plus its `(offset, size)` into the `.dict` file.
- Read that byte range out of `.dict` to get the raw definition.
- Clean the definition with `html.unescape()` plus regex to strip tags and footnote markers.
- Emit one JSON entry per word.

**Example: converting a simple CSV (`word,definition`)**

\`\`\`python
import csv, json

entries = []
with open("mydict.csv", encoding="utf-8") as f:
    for row in csv.reader(f):
        word, definition = row[0], row[1]
        entries.append({
            "Terms": [word],
            "Name": word,
            "MeaningGroups": [
                {"Meanings": [{"Text": definition}]}
            ],
            "Source": "My CSV Dictionary"
        })

with open("mydict.json", "w", encoding="utf-8") as f:
    json.dump(entries, f, ensure_ascii=False)
\`\`\`

Adapt the parsing step for whatever format your source data is in — the important part is that the output always matches the JSON structure in section 1.

### 3. Adding your dictionary to the build

1. Place your finished `.json` file in `dict/custom/`, e.g. `dict/custom/mydict.json`.
2. Open (or create) `dict/custom/custom.go` and register it:

\`\`\`go
package custom

import (
    "encoding/json"
    _ "embed"
    "github.com/pgaskin/lithiumpatch/dict"
)

//go:embed mydict.json
var myDictJSON []byte

func init() {
    dict.Register("my_dictionary_id", 100, func() ([]dict.Entry, error) {
        var entries []dict.Entry
        if err := json.Unmarshal(myDictJSON, &entries); err != nil {
            return nil, err
        }
        return entries, nil
    })
}
\`\`\`

- Each `//go:embed` + `var` pair adds one dictionary file into the binary.
- Each `dict.Register(id, priority, parseFunc)` call registers it. The `id` must be unique across all registered dictionaries. The `priority` controls sort order relative to other dictionaries (higher shows first).
- You can register as many dictionaries as you want in the same `custom.go` file — just repeat the `//go:embed` + `dict.Register` pattern for each one.

3. Make sure `main.go` imports the `custom` package so its `init()` runs:

\`\`\`go
import (
    // ...
    _ "github.com/pgaskin/lithiumpatch/dict/custom"
)
\`\`\`

(The blank import `_` is required — the package is loaded only for its side effect of registering dictionaries, so Go won't let you import it normally without using it directly.)

If you want to disable a built-in dictionary (e.g. `webster1913`), just comment out or remove its import line the same way.

### 4. Building

From the repo root:

\`\`\`bash
go run . app/Lithium_0.24.5.apk
\`\`\`

This parses all registered dictionaries (built-in and custom), patches the APK, and produces a signed, patched APK you can install. Your custom dictionary will now show up alongside the others when looking up a word in the app.

### Troubleshooting

- **`imported and not used`** — you forgot the `_` before the import path in `main.go`.
- **Garbled text like `&mdash;` or `&#x0101;` in definitions** — your conversion script isn't decoding HTML entities; run the text through an HTML-unescape step before writing the JSON.
- **Dictionary doesn't show up in the app** — check that the `.json` file is actually embedded (correct filename in `//go:embed`) and that `custom.go`'s package is imported in `main.go`.
- **`panic: ... already exists`** — you used the same dictionary `id` string twice in `dict.Register`; make each one unique.