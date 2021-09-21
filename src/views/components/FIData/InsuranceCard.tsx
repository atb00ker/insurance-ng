import React, { useState } from 'react';
import Col from 'react-bootstrap/Col';
import Card from 'react-bootstrap/Card';
import { IFIInsurance } from '../../interfaces/IFIData';
import Button from 'react-bootstrap/Button';

const InsuranceCard: React.FC<{ insurance: IFIInsurance }> = ({ insurance }) => {
  const [interested, setInterested] = useState(false);

  const recommend = 'Recommended plan',
    highlyRecommend = 'Highly recommended plan',
    doNotRecommend = 'Not recommended plan',
    notApplicable = 'You cannot purchase this plan';
  let recommendedText = recommend,
    cardColor = 'danger';

  if (insurance.account_id != '') {
    recommendedText = highlyRecommend;
    cardColor = 'primary';
  } else if (insurance.score >= 0.78) {
    recommendedText = highlyRecommend;
    cardColor = 'danger';
  } else if (insurance.score <= 0) {
    recommendedText = notApplicable;
    cardColor = 'secondary';
  } else if (insurance.score <= 0.3) {
    recommendedText = doNotRecommend;
    cardColor = 'secondary';
  } else {
    recommendedText = recommend;
    cardColor = 'danger';
  }

  const computeCover = (cost: number, score: number): string => {
    return Math.ceil(cost + (cost * score) / 15).toLocaleString('en-IN');
  };

  const computePremium = (cost: number, score: number): string => {
    return Math.ceil(cost - (cost * score) / 25).toLocaleString('en-IN');
  };

  return (
    <>
      <Col className='mt-4' sm='5'>
        <Card className={`border border-${cardColor}`}>
          <Card.Body>
            <Card.Title>{insurance.title}</Card.Title>
            <Card.Subtitle className='mb-2 text-muted'>
              {insurance.account_id.length ? insurance.account_id : '-'}
            </Card.Subtitle>
            <Card.Text>
              {insurance.description} <br />
              {insurance.account_id != '' && (
                <>
                  You current have insurance provides a cover of
                  <span className={`text-${cardColor}`}> ₹{insurance.current_cover},000</span> for a premium
                  of <span className={`text-${cardColor}`}>₹{insurance.current_premium},000</span>, we can
                  offer a cover of{' '}
                  <span className={`text-${cardColor}`}>
                    ₹{computeCover(insurance.offer_cover, insurance.score)}
                  </span>{' '}
                  for a premium of{' '}
                  <span className={`text-${cardColor}`}>
                    ₹{computePremium(insurance.offer_premium, insurance.score)}
                  </span>
                  .
                </>
              )}
              {insurance.account_id == '' && (
                <>
                  You are currently not insured, based on your information, we suggest getting a cover of{' '}
                  <span className={`text-${cardColor}`}>
                    ₹{computeCover(insurance.offer_cover, insurance.score)}
                  </span>{' '}
                  only for a premium of{' '}
                  <span className={`text-${cardColor}`}>
                    ₹{computePremium(insurance.offer_premium, insurance.score)}
                  </span>
                  .
                </>
              )}
            </Card.Text>
          </Card.Body>
          <Card.Footer className={`text-${cardColor} d-flex`} style={{ lineHeight: '30px' }}>
            {recommendedText}
            {!interested && insurance.score > 0 && (
              <Button
                className='btn-sm'
                style={{ marginLeft: 'auto' }}
                variant='primary'
                onClick={() => setInterested(true)}
              >
                I am Interested!
              </Button>
            )}
            {!!interested && insurance.score > 0 && (
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

export default InsuranceCard;
