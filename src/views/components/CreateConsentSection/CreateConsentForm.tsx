import React, { FormEventHandler } from 'react';
import Row from 'react-bootstrap/Row';
import Button from 'react-bootstrap/Button';
import Col from 'react-bootstrap/Col';
import Form from 'react-bootstrap/Form';

export interface ICreateConsentForm {
  validated: boolean | undefined;
  handleCreateConsentSubmit: FormEventHandler<HTMLFormElement> | undefined;
}

const CreateConsentForm: React.FC<ICreateConsentForm> = ({ validated, handleCreateConsentSubmit }) => {

  return (
      <Form
        id='ProvideConsentForm'
        noValidate
        validated={validated}
        onSubmit={handleCreateConsentSubmit}>
        <Form.Group as={Row} controlId='enterNumber'>
          <Col>
            <div className="create-consent-input-container mx-auto">
              <Form.Control required className="d-inline create-consent-input" placeholder='Enter Phone Number' name='phone' pattern='[0-9]{10}' />
              <Button className='d-inline ms-2' style={{ marginBottom: '2px' }} type='submit' variant='primary'>
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" className="bi bi-arrow-right-circle-fill" viewBox="0 0 16 16">
                  <path d="M8 0a8 8 0 1 1 0 16A8 8 0 0 1 8 0zM4.5 7.5a.5.5 0 0 0 0 1h5.793l-2.147 2.146a.5.5 0 0 0 .708.708l3-3a.5.5 0 0 0 0-.708l-3-3a.5.5 0 1 0-.708.708L10.293 7.5H4.5z"/>
                </svg>
              </Button>
              <Form.Control.Feedback type='invalid'>
                A valid phone number is required for getting your financial information.
              </Form.Control.Feedback>
            </div>
          </Col>
        </Form.Group>
      </Form>
  );
};

export default CreateConsentForm;
