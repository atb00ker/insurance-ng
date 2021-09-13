import React, { useContext, useState } from 'react';
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import './CreateConsent.scss';
import CreateConsentImage from '../../components/ContentState/CreateConsentImage';
import CreateConsentForm from '../../components/CreateConsentSection/CreateConsentForm';
import { createConsentRequest, getPathToDashboard } from '../../services/axios';
import SectionLoader from '../../components/ContentState/SectionLoader';
import { AuthContext } from '../../components/Auth/AuthProvider';
import { IAuth } from '../../interfaces/IAuth';
import ServerRequestError from '../../components/ContentState/ServerRequestError';

const CreateConsent: React.FC = () => {
  const auth: IAuth = useContext(AuthContext);
  const [validated, setValidated] = useState(false);
  const [showError, setShowError] = useState(false);
  const [showLoader, setShowLoader] = useState(false);

  const handleCreateConsentSubmit = (event: any) => {
    event.preventDefault();
    const form = event.currentTarget;
    setValidated(true);
    setShowLoader(true);
    if (form.checkValidity() === true) {
      auth.user.jwt().then(jwt => {
        createConsentRequest(form.elements.enterNumber.value, jwt)
          .then(response => {
            const consent_handle = response.data['consent_handle'];
            if (consent_handle)
              globalThis.window.location.href = `https://anumati.setu.co/${consent_handle}?redirect_url=${getPathToDashboard()}`;
            setValidated(false);
            setShowError(false);
          })
          .catch(() => {
            setShowError(true);
            setValidated(false);
            setShowLoader(false);
          });
      });
    } else setShowLoader(false);
  };

  return (
    <Container>
      <Row className='mt-4'>
        {!showLoader && (
          <>
            <Col sm='12'>
              <CreateConsentForm
                validated={validated}
                handleCreateConsentSubmit={handleCreateConsentSubmit}
              />
            </Col>
            <Col sm='12'>
              {!showError && <CreateConsentImage height='500px' imgHeight='250px' width='100%' />}
              {!!showError && <ServerRequestError height='500px' imgHeight='250px' width='100%' />}
            </Col>
          </>
        )}
        {!!showLoader && <SectionLoader height='500px' width='100%' />}
      </Row>
    </Container>
  );
};

export default CreateConsent;
