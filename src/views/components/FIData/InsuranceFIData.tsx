import React, { useState } from 'react';
import Col from 'react-bootstrap/Col';
import Card from 'react-bootstrap/Card';
import { InsuranceInfo } from '../../enums/Insurance';
import Button from 'react-bootstrap/Button';
import { ConsentType } from '../../enums/ConsentInfo';
import { InsuranceTypes } from '../../enums/Insurance';
import { IInsuranceFIData } from '../../interfaces/IFIData';

const InsuranceFIData: React.FC<IInsuranceFIData> = ({ data, userInfo }) => {
  const dataSummary = data.children.find((item: any) => item.name == ConsentType.Summary);
  const ngPolicyInfo = InsuranceInfo.find(item => item.policyType == dataSummary.attributes.policyType);
  const [interested, setInterested] = useState(false);

  const getPremium = () => {
    const premium = ngPolicyInfo?.premium || 1;
    if (dataSummary.attributes.policyType == InsuranceTypes.CHILDREN_PLAN) {
      return (premium + premium * (userInfo.datapoint.wealth.score - userInfo.datapoint.dept.score)).toFixed(
        2,
      );
    }
    return premium;
  };

  const getCover = () => {
    const cover = ngPolicyInfo?.cover || 1;
    if (dataSummary.attributes.policyType == InsuranceTypes.CHILDREN_PLAN) {
      return (cover - cover * (userInfo.datapoint.wealth.score - userInfo.datapoint.dept.score)).toFixed(2);
    }
    return cover;
  };

  return (
    <>
      <Col className='mt-4' sm='5'>
        <Card className='border border-info'>
          <Card.Body>
            <Card.Title>{ngPolicyInfo?.title}</Card.Title>
            <Card.Subtitle className='mb-2 text-muted'>{data.attributes.maskedAccNumber}</Card.Subtitle>
            <Card.Text>
              {ngPolicyInfo?.description} <br />
              You current have insurance provides a cover of
              <span className='text-info'> ₹{dataSummary.attributes.coverAmount}k</span> for a premium of{' '}
              <span className='text-info'>₹{dataSummary.attributes.premiumAmount}k</span>, we can offer a
              cover of <span className='text-info'>₹{getCover()}k</span> for a premium of{' '}
              <span className='text-info'>₹{getPremium()}k</span>.
            </Card.Text>
          </Card.Body>
          <Card.Footer className='text-info d-flex' style={{ lineHeight: '30px' }}>
            We can offer better!
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

export default InsuranceFIData;
