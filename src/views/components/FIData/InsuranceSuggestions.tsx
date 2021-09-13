import React, { useState } from 'react';
import Col from 'react-bootstrap/Col';
import Row from 'react-bootstrap/Row';
import Card from 'react-bootstrap/Card';
import { InsuranceInfo } from '../../enums/Insurance';
import Button from 'react-bootstrap/Button';

const InsuranceSuggestions: React.FC<any> = ({ data }) => {
  const [interested, setInterested] = useState(false);

  const getRecommendationMessage = () => {
    if (data.policyType == 'ALL_PLAN') return 'We highly recommend this plan!';
    return 'We recommend this plan!';
  };

  return (
    <>
      <Col sm='5' className='mt-4'>
        <Card className='border border-danger'>
          <Card.Body>
            <Card.Title>{data?.title}</Card.Title>
            <Card.Subtitle className='mb-2 text-muted'>-</Card.Subtitle>
            <Card.Text>
              {data?.description} <br />
              You are currently not insured, based on your financial information, we suggest getting a cover
              of <span className='text-danger'>₹{data?.cover}k</span> only for a premium of{' '}
              <span className='text-danger'>₹{data?.premium}k</span>.
            </Card.Text>
          </Card.Body>
          <Card.Footer className='text-danger d-flex' style={{ lineHeight: '30px' }}>
            {getRecommendationMessage()}
            {!interested && (
              <Button
                className='btn-sm'
                style={{ marginLeft: 'auto' }}
                variant='primary'
                onClick={() => setInterested(true)}
              >
                I am Interested!
              </Button>
            )}
            {!!interested && (
              <Button className='btn-sm' style={{ marginLeft: 'auto' }} variant='secondary' disabled>
                We'll get in touch!
              </Button>
            )}
          </Card.Footer>
        </Card>
      </Col>
    </>
  );
};

export default InsuranceSuggestions;
