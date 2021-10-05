import React, { useContext, useEffect } from 'react';
import { HTTPResponse, registerUserRequest } from '../../helpers/axios';
import { RouterPath } from '../../enums/UrlPath';
import { AuthContext } from '../../components/Auth/AuthProvider';
import { useHistory } from 'react-router-dom';
import { SectionLoader } from '../../components/ContentState/SectionLoader';
import { IAuth } from '../../types/IUser';
import { IServerActionStatus } from '../../types/IServerResponses';

const RegisterUser: React.FC = () => {
  const history = useHistory();
  const auth: IAuth = useContext(AuthContext);

  useEffect(() => {
    auth.user.jwt().then((jwt: string) => {
      registerUserRequest(jwt).then((_: HTTPResponse<IServerActionStatus>) => {
        history.push(RouterPath.Home);
      });
    });
  });

  return <SectionLoader height='500px' width='100%' />;
};

export { RegisterUser };
