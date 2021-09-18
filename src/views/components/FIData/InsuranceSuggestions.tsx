import React, { useEffect, useState } from 'react';
import Col from 'react-bootstrap/Col';
import Card from 'react-bootstrap/Card';
import { InsuranceTypes } from '../../enums/Insurance';
import Button from 'react-bootstrap/Button';
import { IInsuranceFIData } from '../../interfaces/IFIData';

const InsuranceSuggestions: React.FC<IInsuranceFIData> = ({ data, userInfo }) => {
  const [interested, setInterested] = useState(false);
  const [recommendText, setRecommendText] = useState('');
  const [color, setColor] = useState('');
  const [premium, setPremium] = useState(0);
  const [cover, setCover] = useState(0);
  const [score, setScore] = useState(0);

  const scoreSum: () => number = () => {
    if (data.policyType == InsuranceTypes.MOTOR_PLAN) return userInfo.datapoint.motor.score;
    if (data.policyType == InsuranceTypes.TRAVEL_PLAN) return userInfo.datapoint.travel.score;
    if (data.policyType == InsuranceTypes.HOME_PLAN)
      return userInfo.datapoint.age.score + userInfo.datapoint.wealth.score;
    if (data.policyType == InsuranceTypes.PENSION_PLAN)
      return userInfo.datapoint.age.score + userInfo.datapoint.wealth.score + userInfo.datapoint.health.score;
    -userInfo.datapoint.dept.score;
    return (
      userInfo.datapoint.age.score -
      userInfo.datapoint.dept.score +
      userInfo.datapoint.wealth.score +
      userInfo.datapoint.health.score -
      userInfo.datapoint.travel.score +
      userInfo.datapoint.motor.score
    );
  };

  useEffect(() => {
    const recommend = 'We recommend this plan!',
      highlyRecommend = 'We recommend this plan!',
      doNotRecommend = 'We do not think you need this plan!',
      notApplicable = 'You cannot purchase this plan!';
    const cardScore: number = scoreSum();
    setPremium((data.premium + (data.premium * scoreSum()) / 20).toFixed(2))
    setCover((data.cover + (data.cover * scoreSum()) / 20).toFixed(2))

    setScore(cardScore);
    if (data.policyType == InsuranceTypes.ALL_PLAN) {
      setRecommendText(highlyRecommend)
      setColor("danger")
      return
    }
    if (data.policyType == InsuranceTypes.TRAVEL_PLAN) {
      setRecommendText(doNotRecommend)
      setColor("secondary")
      return
    }

    if (cardScore <= 0.001) {
      setRecommendText(notApplicable)
      setColor("secondary")
      return
    }

    setRecommendText(recommend)
    setColor("danger")
    return
  }, [userInfo])

  return (
    <>
      <Col sm='5' className='mt-4'>
        <Card className={`border border-${color}`}>
          <Card.Body>
            <Card.Title>{data?.title}</Card.Title>
            <Card.Subtitle className='mb-2 text-muted'>-</Card.Subtitle>
            <Card.Text>
              {data?.description} <br />
              You are currently not insured, based on your information, we suggest getting a cover of{' '}
              <span className={`text-${color}`}>₹{cover}k</span> only for a premium of{' '}
              <span className={`text-${color}`}>₹{premium}k</span>.
            </Card.Text>
          </Card.Body>
          <Card.Footer className={`text-${color} d-flex`} style={{ lineHeight: '30px' }}>
            {recommendText}
            {!interested && score >= 0.001 && (
              <Button
                className='btn-sm'
                style={{ marginLeft: 'auto' }}
                variant='primary'
                onClick={() => setInterested(true)}>
                I am Interested!
              </Button>
            )}
            {!!interested && score >= 0.001 && (
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
