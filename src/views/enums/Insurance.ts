export enum InsuranceTypes {
  ALL_PLAN = 'ALL_PLAN',
  MEDICAL_PLAN = 'MEDICAL_PLAN',
  PENSION_PLAN = 'PENSION_PLAN',
  HOME_PLAN = 'HOME_PLAN',
  FAMILY_PLAN = 'FAMILY_PLAN',
  CHILDREN_PLAN = 'CHILDREN_PLAN',
  LIFE_PLAN = 'LIFE_PLAN',
  TERM_PLAN = 'TERM_PLAN',
  MOTOR_PLAN = 'MOTOR_PLAN',
  TRAVEL_PLAN = 'TRAVEL_PLAN',
}

export const InsuranceInfo = [
  {
    title: 'Complete Package Plan',
    policyType: 'ALL_PLAN',
    cover: 1000,
    premium: 100,
    description:
      'All insurance packages into one, you pay one premium and enjoy the benefits of all the plans we offer.',
  },
  {
    title: 'Medical Plan',
    policyType: 'MEDICAL_PLAN',
    cover: 100,
    premium: 21,
    description:
      'We cover your medical emergencies, quickly and without the need to difficult and long claim steps.',
  },
  {
    title: 'Pension Plan',
    policyType: 'PENSION_PLAN',
    cover: 400,
    premium: 10,
    description: 'Your income for the old age when you retire, planned for you ahead of time.',
  },
  {
    title: 'Family Plan',
    policyType: 'FAMILY_PLAN',
    cover: 200,
    premium: 25,
    description:
      'The amazing medical plan, but for the entire family to enjoy the benefits from a shared pool of amount.',
  },
  {
    title: "Children's Plan",
    policyType: 'CHILDREN_PLAN',
    cover: 100,
    premium: 20,
    description:
      'You can get seperate insurance for your children. We advice having enough for their adulthhood.',
  },
  {
    title: 'Motor Plan',
    policyType: 'MOTOR_PLAN',
    cover: 150,
    premium: 7.0,
    description: 'Be it your car or bike, we have you covered in the event of an accident.',
  },
  {
    title: 'Term Plan',
    policyType: 'TERM_PLAN',
    cover: 500,
    premium: 35,
    description: 'We recommend atleast x10 of your yearly salary for your family after you.',
  },
  {
    title: 'Travel Plan',
    policyType: 'TRAVEL_PLAN',
    cover: 20,
    premium: 2,
    description: 'When you travel, we ensure that you and your baggage is insured.',
  },
];
