import React, { useContext, useEffect } from 'react';
import { registerUserRequest } from '../../services/axios';
import { RouterPath } from '../../enums/UrlPath';
import { IAuth } from '../../interfaces/IUser';
import { AuthContext } from '../../components/Auth/AuthProvider';
import { useHistory } from 'react-router-dom';
import SectionLoader from '../../components/ContentState/SectionLoader';

const RegisterUser: React.FC = () => {
  const history = useHistory();
  const auth: IAuth = useContext(AuthContext);

  useEffect(() => {
    auth.user.jwt().then(jwt => {
      registerUserRequest(jwt).then(_ => {
        history.push(RouterPath.CreateConsent);
      });
    });
  });

  return <SectionLoader height='500px' width='100%' />;
};

export default RegisterUser;
