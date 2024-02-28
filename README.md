## Overview

Example of turn images of tables to CSV data. Run OCR on detected table cells.  
The execution of table_ocr.ocr_to_csv in [eihli/image-table-ocr](https://github.com/eihli/image-table-ocr) was not what I expected, so I rewrote it using [otiai10/gosseract]( https://github.com/otiai10/gosseract) and rewrote it using.

## Requirements

* [eihli/image-table-ocr](https://github.com/eihli/image-table-ocr)

```sh
pip3 install table_ocr
```

* [tesseract](https://github.com/tesseract-ocr/tesseract)

```sh
# example of homebrew
brew install tesseract
export LIBRARY_PATH="/opt/homebrew/lib"
export CPATH="/opt/homebrew/include"

## list of languages
tesseract --list-langs
brew list tesseract

## add language tessdata
## copy language tessdata to following directory
## https://github.com/tesseract-ocr/tessdata
/opt/homebrew/Cellar/tesseract/${version}/share/tessdata/
```

## Installation

```sh
go install github.com/kis9a/ocr2csv@latest
```

## Example

```sh
#!/bin/bash

IMAGE_PATH="$1"
python3 -m table_ocr.extract_tables "$IMAGE_PATH" | grep 'table' >/tmp/extracted-tables.txt
xargs -I {} python3 -m table_ocr.extract_cells {} </tmp/extracted-tables.txt | grep 'cells' >/tmp/extracted-cells.txt
xargs -I {} python3 -m table_ocr.ocr_image {} </tmp/extracted-cells.txt
ocr2csv -langs=eng,jpn "${IMAGE_PATH%.*}/cells/ocr_data" >result.csv
```

<br/>
<image width="600px" src="https://raw.githubusercontent.com/eihli/image-table-ocr/master/resources/test_data/simple.png"></image>

./2csv simple.png

```csv
Cell,Format,Formula
B4,Percentage,None
C4,General,None
D4,Accounting,None
E4,Currency,=PMT(B4/12\,C4\,D4)
F4,Currency,=E4*C4
```
