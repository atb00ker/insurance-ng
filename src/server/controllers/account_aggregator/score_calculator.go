package account_aggregator

import "time"

func getAgeScore(dob time.Time) float32 {
	return 0.96
}

func getMedicalPlanScore(allFipData []fipDataCollection) float32 {
	return 0.85
}

func getWealthPlanScore(allFipData []fipDataCollection) float32 {
	return 0.73
}

func getDebtPlanScore(allFipData []fipDataCollection) float32 {
	return 0.69
}

func getInvestmentPlanScore(allFipData []fipDataCollection) float32 {
	return 0.76
}

func getPensionPlanScore(allFipData []fipDataCollection) float32 {
	return 0.71
}

func getFamilyPlanScore(allFipData []fipDataCollection) float32 {
	return 0.70
}

func getChildrenPlanScore(allFipData []fipDataCollection) float32 {
	return 0.81
}

func getMotorPlanScore(allFipData []fipDataCollection) float32 {
	return 0
}

func getTermPlanScore(allFipData []fipDataCollection) float32 {
	return 0.6
}

func getTravelPlanScore(allFipData []fipDataCollection) float32 {
	return 0.3
}

func getAllPlanScore(allFipData []fipDataCollection) float32 {
	return 0.79
}
