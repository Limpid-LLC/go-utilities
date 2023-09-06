package utilities

import "go.mongodb.org/mongo-driver/bson"

var Converter *converterUtility

type converterUtility struct{}

func InitConverter() {
	Converter = &converterUtility{}
}

func (util *converterUtility) ConvertToBsonM(inputData interface{}) (bson.M, error) {
	// Convert the map to a byte slice using bson.Marshal
	data, err := bson.Marshal(inputData)
	if err != nil {
		return nil, err
	}

	// Unmarshal the byte slice into a primitive.M
	var result bson.M
	err = bson.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (util *converterUtility) ConvertToBsonMSlice(inputData interface{}) ([]bson.M, error) {
	inputSlice := inputData.(bson.A)

	result := make([]bson.M, len(inputSlice))

	for iter, elem := range inputSlice {
		elemConverted, errConvert := util.ConvertToBsonM(elem)
		if errConvert != nil {
			return nil, errConvert
		}

		result[iter] = elemConverted
	}

	return result, nil
}

func (util *converterUtility) ConvertBSONAToSliceOfStrings(inputData interface{}) []string {
	if inputData == nil {
		return []string{}
	}

	inputSlice := inputData.(bson.A)

	result := make([]string, len(inputSlice))

	for iter, elem := range inputSlice {
		result[iter] = elem.(string)
	}

	return result
}

func (util *converterUtility) ConvertSliceInterfaceToSliceOfStrings(inputData interface{}) []string {
	if inputData == nil {
		return []string{}
	}

	inputSlice := inputData.([]interface{})

	result := make([]string, len(inputSlice))

	for iter, elem := range inputSlice {
		result[iter] = elem.(string)
	}

	return result
}
