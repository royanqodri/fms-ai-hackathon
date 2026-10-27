package util

func GetListStatusData(statusData string) []string {
	var listStatusData []string

	if statusData == "D" {
		listStatusData = append(listStatusData, "D")
	} else if statusData == "A" {
		listStatusData = append(listStatusData, "I")
		listStatusData = append(listStatusData, "U")
		listStatusData = append(listStatusData, "D")
	} else {
		listStatusData = append(listStatusData, "I")
		listStatusData = append(listStatusData, "U")
	}

	return listStatusData
}

func GetStatusData(status string) []string {
	//if empty default ["I", "U", "D"]
	if status == "" {
		return []string{"I", "U", "D"}
	}

	// case status_data
	switch status {
	case "D":
		return []string{"D"}
	case "A":
		return []string{"I", "U"}
	case "I":
		return []string{"I"}
	case "U":
		return []string{"U"}
	default:
		return []string{"I", "U", "D"}
	}
}
