import React, { useContext } from 'react';
import Container from 'react-bootstrap/Container';
import { CreateConsent } from '../../components/Home/CreateConsent';
import { IAuth } from '../../types/IUser';
import { AuthContext } from '../../components/Auth/AuthProvider';
import { NgFeatures } from '../../components/Home/NgFeatures';
import { SectionLoader } from '../../components/ContentState/SectionLoader';
import { ErrorFallback } from '../../components/ContentState/ErrorFallback';
import { ErrorBoundary } from 'react-error-boundary';

const Home: React.FC = () => {
  const auth: IAuth = useContext(AuthContext);

  return (
    <>
      <Container>
        {!auth.isReady && <SectionLoader height='500px' width='100%' />}
        {auth.isReady && auth.isAuthenticated && (
          <ErrorBoundary FallbackComponent={ErrorFallback}>
            <CreateConsent auth={auth} />
          </ErrorBoundary>
        )}
      </Container>
      {/* NgFeatures page can be an entire page in itself as well,
        hence it has it's own container. */}
      {auth.isReady && !auth.isAuthenticated && <NgFeatures />}
    </>
  );
};

export { Home };
