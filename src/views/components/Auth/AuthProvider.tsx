import React from 'react';
import { IAuth } from '../../types/IUser';
import { AuthConfigurations, Provider, ProviderOptions } from './Auth0';

const AuthContext = React.createContext<IAuth>({} as IAuth);
const AuthProvider: React.FC = ({ children }) => {
  return (
    <Provider {...ProviderOptions}>
      <AuthConfigurations>{children}</AuthConfigurations>
    </Provider>
  );
};

export { AuthContext, AuthProvider };
