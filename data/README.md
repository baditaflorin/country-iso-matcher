# Country name data

`aliases.csv` includes the Unicode CLDR 48.2.3 territory display names for Arabic (`ar`), German (`de`), Spanish (`es`), French (`fr`), Hindi (`hi`), Italian (`it`), Japanese (`ja`), Dutch (`nl`), Portuguese (`pt`), Russian (`ru`), Turkish (`tr`), Ukrainian (`uk`), Vietnamese (`vi`), and Chinese (`zh`). It also includes the Indonesian (`id`) display name for the United States. The localized names were appended as aliases for matching; the source data was modified for this service. The pinned Vietnamese source fixture is `cldr/48.2.3/vi/territories.json`. The one-field source comment row in `aliases.csv` is skipped by the CSV loader.

Source: [Unicode CLDR JSON, 48.2.3](https://github.com/unicode-org/cldr-json/tree/48.2.3/cldr-json/cldr-localenames-full/main). License: [Unicode License V3](CLDR-LICENSE.txt).

The alias catalog also includes UN M49 country or area names in English, Chinese, Russian, French, Spanish, and Arabic wherever the current aliases did not already resolve the name to the same ISO alpha-2 code. The M49 names are short forms used for day-to-day UN operations and are not necessarily formal country names. The import used the six language tables and ISO alpha-2 keys on 2026-10-20; the service catalog's Taiwan entry remains sourced from ISO/CLDR because it is not in the M49 table.

A review of the Arabic table found four malformed labels: an underscore placeholder for the Netherlands, reversed word order for the British Indian Ocean Territory, a duplicated word for the French Southern Territories, and a truncated United States Minor Outlying Islands label. The Netherlands alias was replaced with the UN Arabic member-state name `مملكة هولندا`; the other three malformed strings were omitted because valid Arabic aliases for those areas were already present.

Sources: [UN Statistics Division M49 overview](https://unstats.un.org/unsd/methodology/m49/overview/) (accessed 2026-10-20) and [UN Arabic member-state list](https://www.un.org/ar/about-us/member-states) for the corrected Netherlands name.
