/**
 * @jest-environment jsdom
 */

import '@testing-library/jest-dom';
import 'regenerator-runtime/runtime';
import React from 'react';
import ReactDOM from 'react-dom';
import { render, cleanup, act, RenderResult } from '@testing-library/react';
import renderer from 'react-test-renderer';
import { BrowserRouter as Router } from 'react-router-dom';
import { fiDataJohnDoe, insuranceMedicalJohnDoe, userJohnDoe } from './mock/DataMock';
import { IUser } from '../types/IUser';
import { AuthProviderMock } from '../components/Auth/MockedAuth';
import { TitleSection } from '../components/InsuranceInfo/TitleSection';
import { IFIData, IFIInsurance } from '../types/IFIData';

const titleSectionComponent = (isReady: boolean, isAuthenticated: boolean,
  fiData: IFIData = fiDataJohnDoe, insuranceInfo: IFIInsurance = insuranceMedicalJohnDoe,
  user: IUser = userJohnDoe) => (
  <AuthProviderMock isReady={isReady} isAuthenticated={isAuthenticated} user={user}>
    <Router>
      <TitleSection fiData={fiData} insuranceInfo={insuranceInfo} />
    </Router>
  </AuthProviderMock>
);

describe('NgFeatures page tests', () => {
  // beforeEach(() => {});
  afterEach(() => cleanup);

  it('renders without crashing', () => {
    const root = document.createElement('div');
    ReactDOM.render(titleSectionComponent(true, true), root);
  });

  it('test sign up button exists before auth', () => {
    let page: RenderResult | any;
    act(() => {
      page = render(titleSectionComponent(true, true));
    });
    expect(page.queryByText(insuranceMedicalJohnDoe.title)).toBeInTheDocument();
  });

  it('title section snapshot structure', () => {
    const snapTree = renderer.create(titleSectionComponent(true, true));
    expect(snapTree.toJSON()).toMatchSnapshot();
  });
});
