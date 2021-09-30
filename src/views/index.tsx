import React from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import AuthProvider from './components/Auth/AuthProvider';
import Navbar from './components/Navbar/Navbar';
import RegisterUser from './pages/Register/Register';
import CreateConsent from './pages/CreateConsent/CreateConsent';
import Dashboard from './pages/Dashboard/Dashboard';
import InsuranceInfo from './pages/InsuranceDetails/InsuranceInfo';
import About from './pages/About/About';
import { RouterPath } from './enums/UrlPath';
import './index.scss';

const App = () => (
  <AuthProvider>
    <Router>
      <Navbar />
      <Switch>
        <Route path={RouterPath.Register} component={RegisterUser} />
        <Route path={RouterPath.Dashboard} component={Dashboard} />
        <Route path={RouterPath.About} component={About} />
        <Route path={RouterPath.InsuranceDetails} component={InsuranceInfo} />
        <Route path={RouterPath.CreateConsent} component={CreateConsent} />
      </Switch>
    </Router>
  </AuthProvider>
);

ReactDOM.render(<App />, document.getElementById('react-init'));
