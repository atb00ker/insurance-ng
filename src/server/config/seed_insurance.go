package config

import (
	"insurance-ng/src/server/models"
)

func initInsuranceSeed() {
	allPolicy := models.Insurance{
		Title:            "Complete Package Plan",
		Type:             models.OfferPlansAllPlan,
		Cover:            10000000,
		Premium:          90000,
		Description:      "All insurance packages into one, you pay one premium and enjoy the benefits of all the plans we offer.",
		Clauses:          models.InsuranceNGMockedClauses,
		YoyDeductionRate: 15.2,
	}
	medicalPolicy := models.Insurance{
		Title:            "Medical Plan",
		Type:             models.OfferPlansMedicalPlan,
		Cover:            100000,
		Premium:          21000,
		Description:      "We cover your medical emergencies, quickly and without the need to difficult and long claim steps.",
		Clauses:          models.InsuranceNGMockedClauses,
		YoyDeductionRate: 4.2,
	}
	motorPolicy := models.Insurance{
		Title:            "Motor Plan",
		Description:      "Be it your car or bike, we have you covered in the event of an accident.",
		Type:             models.OfferPlansMotorPlan,
		Cover:            150000,
		Premium:          7000,
		Clauses:          models.InsuranceNGMockedClauses,
		YoyDeductionRate: 15.2,
	}
	familyPolicy := models.Insurance{
		Title:            "Family Plan",
		Description:      "The amazing medical plan, but for the entire family to enjoy the benefits from a shared pool of amount.",
		Type:             models.OfferPlansFamilyPlan,
		Cover:            200000,
		Premium:          15000,
		Clauses:          models.InsuranceNGMockedClauses,
		YoyDeductionRate: 7.2,
	}
	travelPolicy := models.Insurance{
		Title:            "Travel Plan",
		Description:      "When you travel, we ensure that you and your baggage is insured.",
		Type:             models.OfferPlansTravelPlan,
		Cover:            20000,
		Premium:          2000,
		Clauses:          models.InsuranceNGMockedClauses,
		YoyDeductionRate: 0.3,
	}
	termPolicy := models.Insurance{
		Title:            "Term Plan",
		Description:      "We recommend atleast x10 of your yearly salary for your family after you.",
		Type:             models.OfferPlansTermPlan,
		Cover:            3000000,
		Premium:          28500,
		Clauses:          models.InsuranceNGMockedClauses,
		YoyDeductionRate: 9.2,
	}
	childrenPolicy := models.Insurance{
		Title:            "Children's Plan",
		Type:             models.OfferPlansChildrenPlan,
		Description:      "You can get separate insurance for your children. We advice having enough for their adulthhood.",
		Cover:            1600000,
		Premium:          28500,
		Clauses:          models.InsuranceNGMockedClauses,
		YoyDeductionRate: 4.2,
	}
	pensionPolicy := models.Insurance{
		Title:            "Pension Plan",
		Type:             models.OfferPlansPensionPlan,
		Cover:            400000,
		Premium:          10000,
		Description:      "Your income for the old age when you retire, planned for you ahead of time.",
		Clauses:          models.InsuranceNGMockedClauses,
		YoyDeductionRate: 1.2,
	}

	createInsuranceSeed(allPolicy)
	createInsuranceSeed(medicalPolicy)
	createInsuranceSeed(motorPolicy)
	createInsuranceSeed(familyPolicy)
	createInsuranceSeed(childrenPolicy)
	createInsuranceSeed(termPolicy)
	createInsuranceSeed(pensionPolicy)
	createInsuranceSeed(travelPolicy)
}

func createInsuranceSeed(insurance models.Insurance) {
	if Database.Model(&insurance).Where("type = ?",
		insurance.Type).Updates(&insurance).RowsAffected == 0 {
		Database.Create(&insurance)
	}
}
