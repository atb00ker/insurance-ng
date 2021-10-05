import axios, { AxiosError, AxiosResponse } from 'axios';
import { RouterPath, ServerPath } from '../enums/UrlPath';
import { IFIData } from '../types/IFIData';
import { IConsentCreatedResponse, IServerActionStatus } from '../types/IServerResponses';
const port = location.protocol === 'http:' ? process.env.REACT_APP_PORT : '443';

axios.defaults.headers.post['Content-Type'] = 'application/json;charset=utf-8';
axios.defaults.headers.post['Access-Control-Allow-Origin'] = process.env.REACT_CORS;
axios.defaults.baseURL = `${location.protocol}//${location.hostname}:${port}`;

export type HTTPResponse<T> = AxiosResponse<T>;
export type HTTPError = Error | AxiosError;

const getJwtHeader = (jwt: string) => {
  return { headers: { Authorization: `Bearer ${jwt}` } };
};

export const registerUserRequest = (jwt: string): Promise<HTTPResponse<IServerActionStatus>> => {
  return axios.get(ServerPath.Register, getJwtHeader(jwt)).catch((error: HTTPError) => {
    const url = new URL(ServerPath.Register, axios.defaults.baseURL).toString();
    console.error(`Cant't get ${url} because ${error}`);
    throw error;
  });
};

export const createConsentRequest = (
  phoneNumber: string,
  jwt: string,
): Promise<HTTPResponse<IConsentCreatedResponse>> => {
  return axios
    .post(ServerPath.CreateConsent, { phone: phoneNumber }, getJwtHeader(jwt))
    .catch((error: HTTPError) => {
      const url = new URL(ServerPath.CreateConsent, axios.defaults.baseURL).toString();
      console.error(`Cant't get ${url} because ${error}`);
      throw error;
    });
};

export const mockConsentNotification = (jwt: string): Promise<HTTPResponse<string>> => {
  return axios.get(ServerPath.ConsentNotificationMock, getJwtHeader(jwt)).catch((error: HTTPError) => {
    const url = new URL(ServerPath.ConsentNotificationMock, axios.defaults.baseURL).toString();
    console.error(`Cant't get ${url} because ${error}`);
    throw error;
  });
};

export const getDashboardData = (jwt: string): Promise<HTTPResponse<IFIData>> => {
  return axios.get(ServerPath.GetUserData, getJwtHeader(jwt)).catch((error: HTTPError) => {
    const url = new URL(ServerPath.GetUserData, axios.defaults.baseURL).toString();
    console.error(`Cant't get ${url} because ${error}`);
    throw error;
  });
};

export const createPurchaseRequest = (uuid: string, jwt: string): Promise<HTTPResponse<IFIData>> => {
  return axios
    .post(ServerPath.InsurancePurchase, { uuid: uuid }, getJwtHeader(jwt))
    .catch((error: HTTPError) => {
      const url = new URL(ServerPath.InsurancePurchase, axios.defaults.baseURL).toString();
      console.error(`Cant't get ${url} because ${error}`);
      throw error;
    });
};

export const createClaimRequest = (uuid: string, jwt: string): Promise<HTTPResponse<IFIData>> => {
  return axios
    .post(ServerPath.InsuranceClaim, { uuid: uuid }, getJwtHeader(jwt))
    .catch((error: HTTPError) => {
      const url = new URL(ServerPath.InsuranceClaim, axios.defaults.baseURL).toString();
      console.error(`Cant't get ${url} because ${error}`);
      throw error;
    });
};

export const getPathToDashboard = () => {
  return new URL(RouterPath.Dashboard, window.location.origin).toString();
};
