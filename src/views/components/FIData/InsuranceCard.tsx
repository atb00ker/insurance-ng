import React from 'react';
import Col from 'react-bootstrap/Col';
import Card from 'react-bootstrap/Card';
import { IFIData, IFIInsurance } from '../../interfaces/IFIData';
import Button from 'react-bootstrap/Button';
import { Link, useHistory } from 'react-router-dom';
import { RouterPath } from '../../enums/UrlPath';

const InsuranceCard: React.FC<{ fiData: IFIData ,insurance: IFIInsurance }> = ({ fiData, insurance }) => {
  const history = useHistory();
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
      <Col className='mt-4' md='10' lg='5'>
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
            {insurance.score > 0 && (
              <Button
                className='btn-sm'
                style={{ marginLeft: 'auto' }}
                variant='primary'
                onClick={() => {
                  history.push({
                    pathname: RouterPath.InsuranceDetails.replace(':insurance_uuid', insurance.uuid),
                    state: fiData,
                  });
                  return;
                }}
              >
                Read More
                <svg
                  xmlns='http://www.w3.org/2000/svg'
                  width='16'
                  height='16'
                  fill='currentColor'
                  className='bi bi-chevron-right'
                  viewBox='0 0 16 18'
                >
                  <path
                    fillRule='evenodd'
                    d='M4.646 1.646a.5.5 0 0 1 .708 0l6 6a.5.5 0 0 1 0 .708l-6 6a.5.5 0 0 1-.708-.708L10.293 8 4.646 2.354a.5.5 0 0 1 0-.708z'
                  />
                </svg>
              </Button>
            )}
          </Card.Footer>
        </Card>
      </Col>
    </>
  );
};

export default InsuranceCard;
