import React, { useState } from 'react';
import Col from 'react-bootstrap/Col';
import Row from 'react-bootstrap/Row';
import Card from 'react-bootstrap/Card';
import { InsuranceInfo } from '../../enums/Insurance';
import Button from 'react-bootstrap/Button';


const InsuranceFIData: React.FC<any> = ({ data }) => {
  const ngPolicyInfo = InsuranceInfo.find(item => item.policyType == data.children[0].attributes.policyType);
  const [interested, setInterested] = useState(false);

  return (
    <>
    <Col className='mt-4' sm='6'>
      <Card className="border border-warning">
        <Card.Body>
          <Card.Title>{ngPolicyInfo?.title}</Card.Title>
          <Card.Subtitle className="mb-2 text-muted">{data.attributes.maskedAccNumber}</Card.Subtitle>
          <Card.Text>
            You current have insurance provides a cover of
            <span className="text-warning"> ₹{data.children[0].attributes.coverAmount} </span> for
            a premium of ₹{data.children[0].attributes.premiumAmount}, we can offer a cover
            of <span className="text-warning">₹{ngPolicyInfo?.cover}</span> for the same amount.
          </Card.Text>
        </Card.Body>
        <Card.Footer className="text-warning d-flex" style={{ lineHeight: '30px' }}>
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

export default InsuranceFIData;
