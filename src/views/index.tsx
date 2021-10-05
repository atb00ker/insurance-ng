import React from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import { AuthProvider } from './components/Auth/AuthProvider';
import { Navbar } from './components/Common/Navbar';
import { RouterPath } from './enums/UrlPath';
import { RegisterUser } from './pages/Register/Register';
import { Home } from './pages/Home/Home';
import { Dashboard } from './pages/Dashboard/Dashboard';
import { InsuranceInfo } from './pages/InsuranceDetails/InsuranceInfo';
import { About } from './pages/About/About';
import { NgFeatures } from './components/Home/NgFeatures';
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
        <Route path={RouterPath.Features} component={NgFeatures} />
        <Route path={RouterPath.Home} component={Home} />
      </Switch>
    </Router>
  </AuthProvider>
);

ReactDOM.render(<App />, document.getElementById('react-init'));
