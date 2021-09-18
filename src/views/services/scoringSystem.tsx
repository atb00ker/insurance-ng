export const getAgeScore = (date: string) => {
  // Write age score logic here.
  // this is a random formule
  return 96 / 100;
};

export const getWealthScore = (wealth: number) => {
  // Write wealth score logic here.
  // this is a random formule
  return ((Math.log(wealth) * Math.LOG10E + 1) | 0) / 7;
};

export const getDeptScore = (dept: number) => {
  // Write dept score logic here.
  // this is a random formule
  return ((Math.log(dept) * Math.LOG10E + 1) | 0) / 1.5;
};

export const getMedicalBillScore = (bills: number) => {
  // Write medical bill score logic here.
  // this is a random formule
  return (100 - bills * 4) / 100;
};

export const getTravelBillScore = (bills: number) => {
  // Write travel bill score logic here.
  // this is a random formule
  return 100 / (100 - bills * 2 + 1000);
};

export const getMotorInsuranceScore = (theftProbablity: number) => {
  // Write motor score logic here.
  // this is a random formule
  return 1 / theftProbablity;
};

export const getInvestmentScore = (
  sip_percent: number,
  mutualfund_percent: number,
  nps_percent: number,
  ppf_percent: number,
  epf_percent: number,
) => {
  // Write investment score logic here.
  // this is a random formule
  return 0.93;
};
