# Specifications


## Types

* Int - 4 bytes
* Float 4 bytes
* Boolean 1 bit/byte
* String - variable length, encoding is set by using parameter

## Stores

* Node - 13 bytes

| inUse | nextRelId | nextPropId | nxtLabelId |
|:------:|:---------:|:----------:|:----------:|
| 1 byte | 4 bytes   | 4 bytes    |  4 byte    |

* Relationship - 34 bytes

|  inUse | firstNodeNxtRelId | second NodeNxtRelId | firstNodePrvRelId | secondNodePrvRelId | nxtPropertyId | relTypeId |
|:------:|:-----------------:|:-------------------:|:-----------------:|:------------------:|:-------------:|:---------:|
| 1 byte | 4 bytes           | 4 bytes             | 4 bytes           | 4 bytes            | 4 bytes       | 4 bytes   |

Property - 10 bytes

|  inUse |  type  | keyStringId | valueOrStrPtr |
|:------:|:------:|:-----------:|:-------------:|
| 1 byte | 1 byte | 4 bytes     | 4 bytes       |

