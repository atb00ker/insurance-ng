const axios = require('axios');
const port = location.protocol === 'http:' ? process.env.REACT_APP_PORT : '443';
import { RouterPath } from '../enums/UrlPath';
import { ServerPath } from '../enums/UrlPath';
import { IFIData } from '../interfaces/IFIData';

axios.defaults.headers.post['Content-Type'] = 'application/json;charset=utf-8';
axios.defaults.headers.post['Access-Control-Allow-Origin'] = process.env.REACT_CORS;
axios.defaults.baseURL = `${location.protocol}//${location.hostname}:${port}`;

const getJwtHeader = (jwt: string) => {
  return { headers: { Authorization: `Bearer ${jwt}` } };
};

export const registerUserRequest = (jwt: string): Promise<any> => {
  return axios.get(ServerPath.Register, getJwtHeader(jwt)).catch((error: any) => {
    const url = new URL(ServerPath.Register, axios.defaults.baseURL).toString();
    console.error(`Cant't get ${url} because ${error}`);
  });
};

export const createConsentRequest = (phoneNumber: string, jwt: string): Promise<any> => {
  return axios
    .post(ServerPath.CreateConsent, { phone: phoneNumber }, getJwtHeader(jwt))
    .catch((error: any) => {
      const url = new URL(ServerPath.CreateConsent, axios.defaults.baseURL).toString();
      console.error(`Cant't get ${url} because ${error}`);
    });
};

export const getConsentStatus = (jwt: string): Promise<any> => {
  return axios.get(ServerPath.ConsentStatus, getJwtHeader(jwt)).catch((error: any) => {
    const url = new URL(ServerPath.ConsentStatus, axios.defaults.baseURL).toString();
    console.error(`Cant't get ${url} because ${error}`);
  });
};

export const getDashboardData = (jwt: string): Promise<any> => {
  return axios.get(ServerPath.GetUserData, getJwtHeader(jwt)).catch((error: any) => {
    const url = new URL(ServerPath.GetUserData, axios.defaults.baseURL).toString();
    console.error(`Cant't get ${url} because ${error}`);
  });
};

export const createPurchaseRequest = (uuid: string, jwt: string): Promise<any> => {
  return axios.post(ServerPath.InsurancePurchase, { uuid: uuid }, getJwtHeader(jwt)).catch((error: any) => {
    const url = new URL(ServerPath.InsurancePurchase, axios.defaults.baseURL).toString();
    console.error(`Cant't get ${url} because ${error}`);
  });
};

export const getPathToDashboard = () => {
  return new URL(RouterPath.Dashboard, window.location.origin).toString();
};
