package config

import "insurance-ng/src/server/models"

func InitInsuranceSeed() {
	all_policy := models.Insurance{
		Type:             models.OfferePlansAllPlan,
		Premium:          91000,
		Cover:            20000,
		YoyDeductionRate: 1.2,
	}
	medical_policy := models.Insurance{
		Type:             models.OfferePlansMedicalPlan,
		Premium:          91000,
		Cover:            20000,
		YoyDeductionRate: 1.2,
	}
	motor_policy := models.Insurance{
		Type:             models.OfferePlansMotorPlan,
		Premium:          91000,
		Cover:            20000,
		YoyDeductionRate: 1.2,
	}
	family_policy := models.Insurance{
		Type:             models.OfferePlansFamilyPlan,
		Premium:          91000,
		Cover:            20000,
		YoyDeductionRate: 1.2,
	}
	home_policy := models.Insurance{
		Type:             models.OfferePlansHomePlan,
		Premium:          91000,
		Cover:            20000,
		YoyDeductionRate: 1.2,
	}
	life_policy := models.Insurance{
		Type:             models.OfferePlansLifePlan,
		Premium:          91000,
		Cover:            20000,
		YoyDeductionRate: 1.2,
	}
	travel_policy := models.Insurance{
		Type:             models.OfferePlansTravelPlan,
		Premium:          91000,
		Cover:            20000,
		YoyDeductionRate: 1.2,
	}
	term_policy := models.Insurance{
		Type:             models.OfferePlansTermPlan,
		Premium:          91000,
		Cover:            20000,
		YoyDeductionRate: 1.2,
	}
	children_policy := models.Insurance{
		Type:             models.OfferePlansChildrenPlan,
		Premium:          91000,
		Cover:            20000,
		YoyDeductionRate: 1.2,
	}
	pension_policy := models.Insurance{
		Type:             models.OfferePlansPensionPlan,
		Premium:          91000,
		Cover:            20000,
		YoyDeductionRate: 1.2,
	}

	createInsuranceSeed(all_policy)
	createInsuranceSeed(medical_policy)
	createInsuranceSeed(motor_policy)
	createInsuranceSeed(family_policy)
	createInsuranceSeed(children_policy)
	createInsuranceSeed(term_policy)
	createInsuranceSeed(pension_policy)
	createInsuranceSeed(travel_policy)
	createInsuranceSeed(home_policy)
	createInsuranceSeed(life_policy)
}

func createInsuranceSeed(insurance models.Insurance) {
	if Database.Model(&insurance).Where("type = ?",
		insurance.Type).Updates(&insurance).RowsAffected == 0 {
		Database.Create(&insurance)
	}
}
