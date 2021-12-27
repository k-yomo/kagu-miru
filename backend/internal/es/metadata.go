package es

type Metadata struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

const (
	MetadataNameFullPath  = "metadata.name"
	MetadataValueFullPath = "metadata.value"
)

const (
	MetadataNameWidthRange  = "幅(cm)"
	MetadataNameDepthRange  = "奥行き(cm)"
	MetadataNameHeightRange = "高さ(cm)"
)

const (
	MetadataValueLengthRangeFrom0to19    = "〜 19cm"
	MetadataValueLengthRangeFrom20to29   = "20 〜 29cm"
	MetadataValueLengthRangeFrom30to39   = "30 〜 39cm"
	MetadataValueLengthRangeFrom40to49   = "40 〜 49cm"
	MetadataValueLengthRangeFrom50to59   = "50 〜 59cm"
	MetadataValueLengthRangeFrom60to69   = "60 〜 69cm"
	MetadataValueLengthRangeFrom70to79   = "70 〜 79cm"
	MetadataValueLengthRangeFrom80to89   = "80 〜 89cm"
	MetadataValueLengthRangeFrom90to99   = "90 〜 99cm"
	MetadataValueLengthRangeFrom100to109 = "100 〜 109cm"
	MetadataValueLengthRangeFrom110to119 = "110 〜 119cm"
	MetadataValueLengthRangeFrom120to129 = "120 〜 129cm"
	MetadataValueLengthRangeFrom130to139 = "130 〜 139cm"
	MetadataValueLengthRangeFrom140to149 = "140 〜 149cm"
	MetadataValueLengthRangeFrom150to159 = "150 〜 159cm"
	MetadataValueLengthRangeFrom160to169 = "160 〜 169cm"
	MetadataValueLengthRangeFrom170to179 = "170 〜 179cm"
	MetadataValueLengthRangeFrom180to189 = "180 〜 189cm"
	MetadataValueLengthRangeFrom190to199 = "190 〜 199cm"
	MetadataValueLengthRangeFrom200      = "200cm 〜"
)

func NewMetadataValueLengthRange(gte int, lte int) string {
	switch {
	case gte >= 0 && lte <= 19:
		return MetadataValueLengthRangeFrom0to19
	case gte >= 20 && lte <= 29:
		return MetadataValueLengthRangeFrom20to29
	case gte >= 30 && lte <= 39:
		return MetadataValueLengthRangeFrom30to39
	case gte >= 40 && lte <= 49:
		return MetadataValueLengthRangeFrom40to49
	case gte >= 50 && lte <= 59:
		return MetadataValueLengthRangeFrom50to59
	case gte >= 60 && lte <= 69:
		return MetadataValueLengthRangeFrom60to69
	case gte >= 70 && lte <= 79:
		return MetadataValueLengthRangeFrom70to79
	case gte >= 80 && lte <= 89:
		return MetadataValueLengthRangeFrom80to89
	case gte >= 90 && lte <= 99:
		return MetadataValueLengthRangeFrom90to99
	case gte >= 100 && lte <= 109:
		return MetadataValueLengthRangeFrom100to109
	case gte >= 110 && lte <= 119:
		return MetadataValueLengthRangeFrom110to119
	case gte >= 120 && lte <= 129:
		return MetadataValueLengthRangeFrom120to129
	case gte >= 130 && lte <= 139:
		return MetadataValueLengthRangeFrom130to139
	case gte >= 140 && lte <= 149:
		return MetadataValueLengthRangeFrom140to149
	case gte >= 150 && lte <= 159:
		return MetadataValueLengthRangeFrom150to159
	case gte >= 160 && lte <= 169:
		return MetadataValueLengthRangeFrom160to169
	case gte >= 170 && lte <= 179:
		return MetadataValueLengthRangeFrom170to179
	case gte >= 180 && lte <= 189:
		return MetadataValueLengthRangeFrom180to189
	case gte >= 190 && lte <= 199:
		return MetadataValueLengthRangeFrom190to199
	case gte >= 200:
		return MetadataValueLengthRangeFrom200
	default:
		return ""
	}
}
