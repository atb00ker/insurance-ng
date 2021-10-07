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
import { userJohnDoe } from './mock/DataMock';
import { IUser } from '../types/IUser';
import { AuthProviderMock } from '../components/Auth/MockedAuth';
import { Home } from '../pages/Home/Home';
import { appTestIds } from './helpers/appIdentifiers';

const homeComponent = (isReady: boolean, isAuthenticated: boolean, user: IUser = userJohnDoe) => (
  <AuthProviderMock isReady={isReady} isAuthenticated={isAuthenticated} user={user}>
    <Router>
      <Home />
    </Router>
  </AuthProviderMock>
);

describe('Home page tests', () => {
  // beforeEach(() => {});
  afterEach(() => cleanup);

  it('renders without crashing', () => {
    const root = document.createElement('div');
    ReactDOM.render(homeComponent(false, false), root);
    ReactDOM.render(homeComponent(true, false), root);
    ReactDOM.render(homeComponent(true, true), root);
  });

  it('loading page before ready', () => {
    let page: RenderResult | any;
    act(() => {
      page = render(homeComponent(false, false));
    });
    expect(page.getByTestId(appTestIds.sectionLoading)).toBeInTheDocument();
  });

  it('home page before login', () => {
    let page: RenderResult | any;
    act(() => {
      page = render(homeComponent(true, false));
    });
    expect(page.getByTestId(appTestIds.featuresPageLogo)).toBeInTheDocument();
    expect(page.getByTestId(appTestIds.featuresPageButton)).toBeInTheDocument();
  });

  it('account aggregator form after login', () => {
    let page: RenderResult | any;
    act(() => {
      page = render(homeComponent(true, true));
    });
    expect(page.getByTestId(appTestIds.homePhoneNumberInput)).toBeInTheDocument();
  });

  it('home snapshot structure after login', () => {
    const snapTree = renderer.create(homeComponent(true, true));
    expect(snapTree.toJSON()).toMatchSnapshot();
  });

  it('home snapshot structure before login', () => {
    const snapTree = renderer.create(homeComponent(true, false));
    expect(snapTree.toJSON()).toMatchSnapshot();
  });
});
