import React, { useState } from 'react';
import Col from 'react-bootstrap/Col';
import Row from 'react-bootstrap/Row';
import Card from 'react-bootstrap/Card';
import { InsuranceInfo } from '../../enums/Insurance';
import Button from 'react-bootstrap/Button';


const InsuranceSuggestions: React.FC<any> = ({ data }) => {
  const [interested, setInterested] = useState(false);

  return (
    <>
    <Col sm='6' className="mt-4" >
      <Card className="border border-danger">
        <Card.Body>
          <Card.Title>{data?.title}</Card.Title>
          <Card.Subtitle className="mb-2 text-muted">-</Card.Subtitle>
          <Card.Text>
            You are currently not insured, based on your financial information, we suggest getting
            a cover of ₹{data?.cover} only for a premium of ₹{data?.premium}.
          </Card.Text>
        </Card.Body>
        <Card.Footer className="text-danger d-flex" style={{ lineHeight: '30px' }}>
            We can offer better.
            {!interested &&
            <Button className="btn-sm" style={{ marginLeft: 'auto' }} variant="primary" onClick={() => setInterested(true)} >I am Interested!</Button>}
            {!!interested &&
            <Button className="btn-sm" style={{ marginLeft: 'auto' }} variant="secondary" disabled>We'll get in touch!</Button>}
        </Card.Footer>
      </Card>
    </Col>
    </>
  );
};

export default InsuranceSuggestions;
