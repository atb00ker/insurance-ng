// This is just random scores to simulate real algorithms

import { IFIData } from '../interfaces/IFIData';

var XMLParser = require('react-xml-parser');

export function prepareDataJson(response: IFIData[]) {
  // This should be done in the backend, but I don't want to make
  // golang types all the data right now.
  response.map(element => {
    element.RahasyaData.map(item => {
      if (item.data) item.data = new XMLParser().parseFromString(item.data);
      return item;
    });
  });
  return response;
}

export function getMedicalBills(data: any) {
  return 2;
}

export function getDeptBills(data: any) {
  return 3;
}

export function getTravelledBills(data: any) {
  return 1;
}

export function getInsurancePercents(data: any) {
  const sip_percent = 20;
  const mutualfund_percent = 10;
  const nps_percent = 30;
  const ppf_percent = 30;
  const epf_percent = 10;
  return [sip_percent, mutualfund_percent, nps_percent, ppf_percent, epf_percent];
}

export function getTotalWealth(data: any) {
  return 10000;
}

export function getMotorTheftScore(data: any) {
  return 100000000;
}
