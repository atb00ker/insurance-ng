import React from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import AuthProvider from './components/Auth/AuthProvider';
import Navbar from './components/Navbar/Navbar';
import RegisterUser from './pages/Register/Register';
import CreateConsent from './pages/CreateConsent/CreateConsent';
import Dashboard from './pages/Dashboard/Dashboard';
import { RouterPath } from './enums/RouterPath';
import './index.scss';

const App = () => (
  <AuthProvider>
    <Router>
      <Navbar />
      <Switch>
        <Route path={RouterPath.Register} component={RegisterUser} />
        <Route path={RouterPath.Dashboard} component={Dashboard} />
        <Route path={RouterPath.CreateConsent} component={CreateConsent} />
      </Switch>
    </Router>
  </AuthProvider>
);

ReactDOM.render(<App />, document.getElementById('react-init'));
