import React, { FormEventHandler } from 'react';
import Row from 'react-bootstrap/Row';
import Button from 'react-bootstrap/Button';
import Col from 'react-bootstrap/Col';
import Form from 'react-bootstrap/Form';
import { rightArrowInCircle } from '../../helpers/svgIcons';

export interface ICreateConsentForm {
  validated: boolean | undefined;
  handleCreateConsentSubmit: FormEventHandler<HTMLFormElement> | undefined;
}

const CreateConsentForm: React.FC<ICreateConsentForm> = ({ validated, handleCreateConsentSubmit }) => {
  return (
    <Form id='ProvideConsentForm' noValidate validated={validated} onSubmit={handleCreateConsentSubmit}>
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
            <Button className='d-inline ms-2' style={{ marginBottom: '2px' }} type='submit' variant='primary'>
              {rightArrowInCircle()}
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
