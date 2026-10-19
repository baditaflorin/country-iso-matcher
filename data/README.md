# Country name data

`aliases.csv` includes the Unicode CLDR 48.2.3 territory display names for Arabic (`ar`), German (`de`), Spanish (`es`), French (`fr`), Hindi (`hi`), Italian (`it`), Japanese (`ja`), Dutch (`nl`), Portuguese (`pt`), Russian (`ru`), Turkish (`tr`), Ukrainian (`uk`), and Chinese (`zh`). The localized names were appended as aliases for matching; the source data was modified for this service. The one-field source comment row in `aliases.csv` is skipped by the CSV loader.

Source: [Unicode CLDR JSON, 48.2.3](https://github.com/unicode-org/cldr-json/tree/48.2.3/cldr-json/cldr-localenames-full/main). License: [Unicode License V3](CLDR-LICENSE.txt).

The alias catalog also includes UN M49 country or area names in English, Chinese, Russian, French, Spanish, and Arabic wherever the current aliases did not already resolve the name to the same ISO alpha-2 code. The M49 names are short forms used for day-to-day UN operations and are not necessarily formal country names. The import used the six language tables and ISO alpha-2 keys on 2026-10-20; the service catalog's Taiwan entry remains sourced from ISO/CLDR because it is not in the M49 table.

Source: [UN Statistics Division M49 overview](https://unstats.un.org/unsd/methodology/m49/overview/).
