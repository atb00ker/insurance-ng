import React from 'react';
import { AuthContext } from './AuthProvider';
import { IAuth, IUser } from '../../types/IUser';

const AuthProviderMock: React.FC<{
  children: React.ReactNode;
  isReady: boolean;
  isAuthenticated: boolean;
  user: IUser;
}> = ({ children, isReady, isAuthenticated, user }) => {
  const getJwt = async (): Promise<string> => {
    if (isAuthenticated) {
      return '12345';
    }
    return new Promise(() => '');
  };

  const auth: IAuth = {
    loginWithRedirect: () => {
      return new Promise(() => '');
    },
    logout: () => {},
    isAuthenticated: isAuthenticated,
    isReady: isReady,
    user: {
      id: user.id,
      name: user.name,
      email: user.email,
      jwt: getJwt,
    },
  };

  return <AuthContext.Provider value={auth}>{children}</AuthContext.Provider>;
};

export { AuthProviderMock };
