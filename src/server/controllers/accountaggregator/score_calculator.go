package accountaggregator

// The formulas used by insurance companies are researched
// over long periods of time by a team of dedicated
// statisticians, I do not know the art of that trade, so
// here, I am mocking the output of the formulas.
// However, in the spirit of completeness, I will provide,
// the datapoints that AA can provide as an input to the
// formulas.

import "time"

func getAgeScore(dob time.Time) float32 {
	// Younger people score higher in age score
	// which factors in medical and term insurance.
	// Datapoints:
	// - deposit:profile.holders.holder.dob
	return 0.81
}

func getMedicalPlanScore(allFipData []fipDataCollection, sharedDataSources int16) float32 {
	// Like it or not, your medical state largely depends on your
	// lifestyle, profession and quality of life.
	// Datapoints:
	// - age score (calculated above)
	// - deposit:transactions.transaction.narration+amount (filtered transactions with Hospital, Pharmacy etc)
	// - insurance_policies:transactions.transaction.txnDate+amount+type (filtered Medical Claims)
	// - deposit:transactions.transaction.type+narration+amount+transactionTimestamp (calculate income and spending)
	// - deposit:transactions.transaction.type+narration+amount (check for spending patterns) -- excessive
	//                                     /hedonistic would rank lower and family spending patterns higher
	// - deposit:profile.holders.holder.address (place of living reflects quality of life, proxymity to
	//                                              emergency services, financial security)
	// - Health records from other sources to fetch known conditions (eg. medical records)
	return 0.85 + float32(sharedDataSources)/50
}

func getWealthPlanScore(allFipData []fipDataCollection, sharedDataSources int16) float32 {
	// Knowning about your financial status in the society can indicate a lot about
	// many aspects of your life, including the length of it.
	// Datapoints:
	// - deposit:summary.currentBalance
	// - deposit:transactions.transaction.type+narration+amount (check for spending patterns) -- excessive
	//                                     /hedonistic would rank lower and family spending patterns higher
	// - credit_card:summary.currentDue (Debt can be a good indictor of financial security and maturity)
	// - term-deposit|reoccuring-deposits|ppf|nfs:summary.currentValue (Investments indicate
	//                                      financial / future security & planning). A good planner is
	//                                      likely to be more mature in other aspects of life as well.
	// - sip|mutual_funds:summary.currentValue (Investments indicate financial / future security)
	return 0.73 + float32(sharedDataSources)/50
}

func getDebtPlanScore(allFipData []fipDataCollection, sharedDataSources int16) float32 {
	// Knowning about your financial status in the society can indicate a lot about
	// many aspects of your life, including the length of it.
	// Datapoints:
	// - deposit:summary.currentBalance
	// - deposit:transactions.transaction.type+narration+amount+transactionTimestamp (spending patterns
	//                                    can indicate the future possibility of debt)
	// - credit_card:summary.currentDue (Debt can be a good indictor of financial security and maturity)
	return 0.69 + float32(sharedDataSources)/50
}

func getInvestmentScore(allFipData []fipDataCollection, sharedDataSources int16) float32 {
	// Knowning about your investments gives insights about the person's future planning and maturity.
	// Datapoints:
	// - term-deposit|reoccuring-deposits|ppf|nfs:summary.currentValue (Investments indicate
	//                                      financial / future security & planning). A good planner is
	//                                      likely to be more mature in other aspects of life as well.
	// - sip|mutual_funds:summary.currentValue (Investments indicate financial / future security)
	return 0.76 + float32(sharedDataSources)/50
}

func getPensionPlanScore(allFipData []fipDataCollection, sharedDataSources int16) float32 {
	// For suggesting the pension we need to know the financial-social situation of the individual.
	// Datapoints:
	// - age score (calculated above)
	// - wealth score (calculated above)
	// - debt score (calculated above)
	// - term-deposit|reoccuring-deposits|ppf|nfs:summary.currentValue (Investments indicate
	//                                      financial / future security & planning). A good planner is
	//                                      likely to be more mature in other aspects of life as well.
	// - sip|mutual_funds:summary.currentValue (Investments indicate financial / future security)
	return 0.71 + float32(sharedDataSources)/50
}

func getFamilyPlanScore(allFipData []fipDataCollection, sharedDataSources int16) float32 {
	// Family plan would be very similar to the medical insurance plan calculated above.
	// Datapoints:
	// - age score (calculated above)
	// - wealth score (calculated above)
	// - medical plan score (calculated above)
	// - debt score (calculated above)
	return 0.70 + float32(sharedDataSources)/50
}

func getChildrenPlanScore(allFipData []fipDataCollection, sharedDataSources int16) float32 {
	// Children's plan would be very similar to the medical & family insurance plan calculated above.
	// Datapoints:
	// - age score (calculated above)
	// - wealth score (calculated above)
	// - medical plan score (calculated above)
	// - debt score (calculated above)
	// - family score (calculated above)
	return 0.81 + float32(sharedDataSources)/50
}

func getMotorPlanScore(allFipData []fipDataCollection, sharedDataSources int16) float32 {
	// Children's plan would be very similar to the medical & family insurance plan calculated above.
	// Datapoints:
	// - wealth score (calculated above) -- will indicate the reasonable amount spend on the vehicle.
	// - insurance_policies:transactions.transaction.txnDate+amount+type (filtered for motor claims)
	// - deposit:profile.holders.holder.address (can help indicate if the person lives in a place where
	//                                              there have been many vehicle theft cases)
	return 0
}

func getTermPlanScore(allFipData []fipDataCollection, sharedDataSources int16) float32 {
	// Term insurance rates boil down to life expentency, which is affected by your quality of life
	// and financial status
	// Datapoints:
	// - age score (calculated above)
	// - wealth score (calculated above)
	// - debt score (calculated above)
	// - deposit:transactions.transaction.type+narration+amount+transactionTimestamp (calculate income and spending)
	// - deposit:transactions.transaction.type+narration+amount (check for spending patterns) -- excessive
	//                                     /hedonistic would rank lower and family spending patterns higher
	// - deposit:profile.holders.holder.address (place of living, reflects quality of life, proxymity to
	//                                              emergency services, financial security)
	return 0.6 + float32(sharedDataSources)/50
}

func getTravelPlanScore(allFipData []fipDataCollection, sharedDataSources int16) float32 {
	// Travel plan depends highly on the previous travel experience, frequency of travel,
	// country of travel.
	// Datapoints:
	// - wealth score (calculated above) -- indicates the type of safety / luxury in the trip.
	// - deposit:transactions.transaction.type+narration+amount+transactionTimestamp -- indicate country
	//                                    or places travelling to in the near future,
	// frequency of travel.
	return 0.3
}

func getAllPlanScore(allFipData []fipDataCollection, sharedDataSources int16) float32 {
	// A plan to cover all plans depends on the factors discussed above.
	// Datapoints:
	// - Age Score (calculated above)
	// - Medical Score (calculated above)
	// - Wealth Score (calculated above)
	// - Debt Score (calculated above)
	// - Investment Score (calculated above)
	// - Pension Score (calculated above)
	// - Family Score (calculated above)
	// - Children Score (calculated above)
	// - Motor Score (calculated above)
	// - Term Score (calculated above)
	// - Travel Score (calculated above)
	return 0.79 + float32(sharedDataSources)/50
}
