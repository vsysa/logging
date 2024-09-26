package helper

func CopyMapContext(original map[string]interface{}) map[string]interface{} {
	newMap := make(map[string]interface{}, len(original))
	for key, value := range original {
		newMap[key] = value
	}
	return newMap
}
