# Specifications


## Types

* Int - 4 bytes
* Float 4 bytes
* Boolean 1 bit/byte
* String - fixed length or variable? ASCII or UTF16?

## Stores

* Node - 15 bytes

| In Use | nextRelId | nextPropId | labels? | extra? | 
|:------:|:---------:|:----------:|:-------:|:------:|
| 1 byte | 4 bytes   | 4 bytes    |  hz     |    hz  |