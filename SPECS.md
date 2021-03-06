# Specifications

## Types

* Int - 4 bytes
* Float 4 bytes
* Boolean 1 bit/byte
* String - variable length, encoding is set by using parameter

IDs start from 1

InUse column is first in everyschema

## Stores

We store label, relationtype, keys

* Node - 13 bytes

| inUse | nextRelId | nextPropId | nxtLabelId |
|:------:|:---------:|:----------:|:----------:|
| 1 byte | 4 bytes   | 4 bytes    |  4 byte    |

* Relationship - 34 bytes

|  inUse |  firstInChain | secondNodeId | firstNodeId | firstNodeNxtRelId | second NodeNxtRelId | firstNodePrvRelId | secondNodePrvRelId | nxtPropertyId | relTypeId |
|:------:|:-------------:|:------------:|:-----------:|:-----------------:|:-------------------:|:-----------------:|:------------------:|:-------------:|:---------:|
| 1 byte |     1 byte    |    4 bytes   |    4 bytes  |       4 bytes     | 4 bytes             | 4 bytes           | 4 bytes            | 4 bytes       | 4 bytes   |

Property - 14 bytes

|  inUse |  type  | keyStringId | nextPropertyID | valueOrStrPtr |
|:------:|:------:|:-----------:|:--------------:|:-------------:|
| 1 byte | 1 byte | 4 bytes     |    4 bytes     | 4 bytes       |

String - 64 bytes

|  inUse |  extra |   nxtPartID |   value  |
|:------:|:------:|:-----------:|:--------:|
| 1 byte | 1 byte |   4 bytes   | 58 bytes |

Label - 9 bytes

|  inUse |   labelStringId | nxtLabelID |
|:------:|:---------------:|:----------:|
| 1 byte |      4 bytes    | 4 bytes    |

InUse - 11 bytes

|  inUse |  type  |  head  |   nodeID | nextFreeRowId |
|:------:|:------:|:------:|:--------:|:-------------:|
| 1 byte | 1 byte | 1 byte | 4 bytes  | 4 bytes       |

LabelString - 21 bytes

|  inUse |  String |
|:------:|:------------:|
| 1 byte |    20 bytes  |

PropertyKey - 21 bytes

|  inUse |  KeyString |
|:------:|:------------:|
| 1 byte |    20 bytes  |

RelationshipType - 21 bytes

|  inUse |  TypeString  |
|:------:|:------------:|
| 1 byte |  20 bytes  |