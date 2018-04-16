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

Property - 10 bytes

|  inUse |  type  | keyStringId | valueOrStrPtr |
|:------:|:------:|:-----------:|:-------------:|
| 1 byte | 1 byte | 4 bytes     | 4 bytes       |

String - 64 bytes

|  inUse |  extra |   value  |   nxtPartID |
|:------:|:------:|:--------:|:-----------:|
| 1 byte | 1 byte | 58 bytes |   4 bytes   |

Labels - 9 bytes

|  inUse |   labelStringId | nxtLabelID |
|:------:|:---------------:|:----------:|
| 1 byte |      4 bytes    | 4 bytes    |

InUse - 11 bytes

|  inUse |  type  |  head  |   nodeID | nextFreeRowId |
|:------:|:------:|:------:|:--------:|:-------------:|
| 1 byte | 1 byte | 1 byte | 4 bytes  | 4 bytes       |

LabelString - 21 bytes

|  inUse |  labelString |
|:------:|:------------:|
| 1 byte |    20 bytes  |

PropertyKey - 21 bytes

|  inUse |  propertyKey |
|:------:|:------------:|
| 1 byte |    20 bytes  |

RelationshipType - 21 bytes

|  inUse |  relType  |
|:------:|:---------:|
| 1 byte | 20 bytes  |