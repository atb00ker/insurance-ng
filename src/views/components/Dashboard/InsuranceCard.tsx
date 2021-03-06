import React from 'react';
import Col from 'react-bootstrap/Col';
import Card from 'react-bootstrap/Card';
import { IFIData, IFIInsurance } from '../../types/IFIData';
import Button from 'react-bootstrap/Button';
import { useHistory } from 'react-router-dom';
import { RouterPath } from '../../enums/UrlPath';
import { hourglassWaitIcon } from '../../helpers/svgIcons';
import './InsuranceCard.scss';

enum recommendedCutOff {
  notApplicable = 0,
  doNotRecommend = 0.3,
  highlyRecommend = 0.78,
}

const InsuranceCard: React.FC<{ fiData: IFIData; insurance: IFIInsurance }> = ({ fiData, insurance }) => {
  const history = useHistory();

  const recommend = <span>Recommended plan</span>,
    activationInProgress = <span>{hourglassWaitIcon('0 1 18 18')} Application under review</span>,
    activePlan = <span>Active Plan</span>,
    highlyRecommend = <span>Highly recommended plan</span>,
    doNotRecommend = <span>Not recommended plan</span>,
    notApplicable = <span>You cannot purchase this plan</span>;

  let recommendedText = recommend,
    cardColor = 'danger';

  if (insurance.is_insurance_ng_acct && insurance.account_id != '') {
    recommendedText = activePlan;
    cardColor = 'success';
  } else if (insurance.is_insurance_ng_acct && insurance.account_id == '') {
    recommendedText = activationInProgress;
    cardColor = 'success';
  } else if (insurance.account_id != '') {
    recommendedText = highlyRecommend;
    cardColor = 'primary';
  } else if (insurance.score >= recommendedCutOff.highlyRecommend) {
    recommendedText = highlyRecommend;
    cardColor = 'danger';
  } else if (insurance.score <= recommendedCutOff.notApplicable) {
    recommendedText = notApplicable;
    cardColor = 'secondary';
  } else if (insurance.score <= recommendedCutOff.doNotRecommend) {
    recommendedText = doNotRecommend;
    cardColor = 'secondary';
  } else {
    recommendedText = recommend;
    cardColor = 'danger';
  }

  return (
    <>
      <Col className='mt-4' md='10' lg='5'>
        <Card className={`border border-${cardColor}`}>
          <Card.Body>
            <Card.Title>{insurance.title}</Card.Title>
            <Card.Subtitle className='mb-2 text-muted'>
              {insurance.account_id.length ? insurance.account_id : '-'}
            </Card.Subtitle>
            <Card.Text className='insurance-card-text roboto-regular'>
              {insurance.description} <br />
              {!!insurance.is_insurance_ng_acct && (
                <>
                  We are offering you a cover of
                  <span className={`text-${cardColor}`}>
                    {' '}
                    ???{insurance.current_cover.toLocaleString('en-IN')}
                  </span>{' '}
                  for a yearly premium of
                  <span className={`text-${cardColor}`}>
                    {' '}
                    ???{insurance.current_premium.toLocaleString('en-IN')}
                  </span>
                  . For initiating claims, please contact your agent or contact help line.
                </>
              )}
              {!insurance.is_insurance_ng_acct && insurance.account_id != '' && (
                <>
                  You current have insurance provides a cover of
                  <span className={`text-${cardColor}`}>
                    {' '}
                    ???{insurance.current_cover.toLocaleString('en-IN')}
                  </span>{' '}
                  for a premium of{' '}
                  <span className={`text-${cardColor}`}>
                    ???{insurance.current_premium.toLocaleString('en-IN')}
                  </span>
                  , we can offer a cover of{' '}
                  <span className={`text-${cardColor}`}>
                    ???{insurance.offer_cover.toLocaleString('en-IN')}
                  </span>{' '}
                  for a premium of{' '}
                  <span className={`text-${cardColor}`}>
                    ???{insurance.offer_premium.toLocaleString('en-IN')}
                  </span>
                  .
                </>
              )}
              {!insurance.is_insurance_ng_acct && insurance.account_id == '' && (
                <>
                  You are currently not insured, based on your information, we suggest getting a cover of{' '}
                  <span className={`text-${cardColor}`}>
                    ???{insurance.offer_cover.toLocaleString('en-IN')}
                  </span>{' '}
                  only for a premium of{' '}
                  <span className={`text-${cardColor}`}>
                    ???{insurance.offer_premium.toLocaleString('en-IN')}
                  </span>
                  .
                </>
              )}
            </Card.Text>
          </Card.Body>
          <Card.Footer className={`text-${cardColor} d-flex`} style={{ lineHeight: '30px' }}>
            {recommendedText}
            {insurance.score > recommendedCutOff.notApplicable && (
              <Button
                className='btn-sm'
                style={{ marginLeft: 'auto' }}
                variant='primary'
                onClick={() => {
                  // Because we are using history API to change
                  // page here, we need to scroll up manually.
                  window.scrollTo(0, 0);
                  history.push({
                    pathname: RouterPath.InsuranceDetails.replace(':insurance_uuid', insurance.uuid),
                    state: fiData,
                  });
                  return;
                }}>
                Read More
                <svg
                  xmlns='http://www.w3.org/2000/svg'
                  width='16'
                  height='16'
                  fill='currentColor'
                  className='bi bi-chevron-right'
                  viewBox='0 0 16 18'>
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

export { InsuranceCard };
