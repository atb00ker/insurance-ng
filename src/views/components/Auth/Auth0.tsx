import React from 'react';
import { useAuth0, Auth0Provider, Auth0ProviderOptions } from '@auth0/auth0-react';
import { AuthContext } from './AuthProvider';
import { IAuth } from '../../types/IUser';

export const Provider = Auth0Provider;
export const ProviderOptions: Auth0ProviderOptions = {
  domain: process.env.AUTH0_DOMAIN || '',
  clientId: process.env.AUTH0_CLIENT_ID || '',
  redirectUri: process.env.AUTH0_REDIRECT_URI || window.location.origin,
  cacheLocation: 'localstorage',
};
const AuthConfigurations: React.FC = ({ children }) => {
  const { isLoading, isAuthenticated, loginWithRedirect, logout, user, getIdTokenClaims } = useAuth0();

  const auth0Logout = () => {
    return logout({ returnTo: process.env.AUTH0_LOGOUT_URI || window.location.origin });
  };

  const getJwt = async (): Promise<string> => {
    if (isAuthenticated) {
      const token = await getIdTokenClaims();
      return token.__raw;
    }
    return new Promise((_, __) => '');
  };

  const auth: IAuth = {
    loginWithRedirect: loginWithRedirect,
    logout: auth0Logout,
    isAuthenticated: isAuthenticated,
    isReady: !isLoading,
    user: {
      id: user?.sub,
      name: user?.name,
      email: user?.email,
      jwt: getJwt,
    },
  };

  // Enable if we want the autologin required feature
  // if (auth.isReady && !auth.isAuthenticated) {
  //   loginWithRedirect();
  //   return <div>Loading...</div>;
  // }
  return (
    <React.Fragment>
      <AuthContext.Provider value={auth}>{children}</AuthContext.Provider>
    </React.Fragment>
  );
};

export { AuthConfigurations };
