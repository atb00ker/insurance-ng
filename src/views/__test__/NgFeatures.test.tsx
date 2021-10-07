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
import { NgFeatures } from '../components/Home/NgFeatures';
import { userJohnDoe } from './mock/DataMock';
import { IUser } from '../types/IUser';
import { AuthProviderMock } from '../components/Auth/MockedAuth';
import { appStringText, appTestIds } from './helpers/appIdentifiers';

const ngFeaturesComponent = (isReady: boolean, isAuthenticated: boolean, user: IUser = userJohnDoe) => (
  <AuthProviderMock isReady={isReady} isAuthenticated={isAuthenticated} user={user}>
    <Router>
      <NgFeatures />
    </Router>
  </AuthProviderMock>
);

describe('NgFeatures page tests', () => {
  // beforeEach(() => {});
  afterEach(() => cleanup);

  it('renders without crashing', () => {
    const root = document.createElement('div');
    ReactDOM.render(ngFeaturesComponent(true, false), root);
    ReactDOM.render(ngFeaturesComponent(true, true), root);
  });

  it('test sign up button not exists after auth', () => {
    let page: RenderResult | any;
    act(() => {
      page = render(ngFeaturesComponent(true, true));
    });
    expect(page.getByTestId(appTestIds.featuresPageLogo)).toBeInTheDocument();
    expect(page.queryByText(appStringText.signUp)).not.toBeInTheDocument();
  });

  it('test sign up button exists before auth', () => {
    let page: RenderResult | any;
    act(() => {
      page = render(ngFeaturesComponent(true, false));
    });
    expect(page.getByTestId(appTestIds.featuresPageLogo)).toBeInTheDocument();
    expect(page.getByTestId(appTestIds.featuresPageButton)).toBeInTheDocument();
  });

  it('ngfeature snapshot structure', () => {
    const snapTree = renderer.create(ngFeaturesComponent(true, false));
    expect(snapTree.toJSON()).toMatchSnapshot();
  });
});
