import React, { useState } from 'react';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import CreateConsentImage from '../ContentState/CreateConsentImage';
import { createConsentRequest, getPathToDashboard } from '../../helpers/axios';
import SectionLoader from '../ContentState/SectionLoader';
import { IAuth } from '../../interfaces/IUser';
import ServerRequestError from '../ContentState/ServerRequestError';
import Button from 'react-bootstrap/Button';
import Form from 'react-bootstrap/Form';
import { rightArrowInCircle } from '../../helpers/svgIcons';
import './CreateConsent.scss';

const CreateConsent: React.FC<{ auth: IAuth }> = ({ auth }) => {
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
    <Row className='mt-4'>
      {!showLoader && (
        <>
          <Col sm='12'>
            <Form
              id='ProvideConsentForm'
              noValidate
              validated={validated}
              onSubmit={handleCreateConsentSubmit}
            >
              <Form.Group as={Row} controlId='enterNumber'>
                <Col>
                  <div className='create-consent-input-container mx-auto'>
                    <Form.Control
                      required
                      className='d-inline create-consent-input'
                      placeholder='Enter Phone Number'
                      name='phone'
                      pattern='[0-9]{10}'
                    />
                    <Button
                      className='d-inline ms-2'
                      style={{ marginBottom: '2px' }}
                      type='submit'
                      variant='primary'
                    >
                      {rightArrowInCircle('0 0 16 16')}
                    </Button>
                    <Form.Control.Feedback type='invalid'>
                      A valid phone number is required for getting your financial information.
                    </Form.Control.Feedback>
                  </div>
                </Col>
              </Form.Group>
            </Form>
          </Col>
          <Col sm='12'>
            {!showError && <CreateConsentImage imgHeight='100%' width='100%' />}
            {!!showError && <ServerRequestError height='500px' imgHeight='250px' width='100%' />}
          </Col>
        </>
      )}
      {!!showLoader && <SectionLoader height='500px' width='100%' />}
    </Row>
  );
};

export default CreateConsent;
