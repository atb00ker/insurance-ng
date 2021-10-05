import React from 'react';
import Col from 'react-bootstrap/Col';
import Card from 'react-bootstrap/Card';
import Image from 'react-bootstrap/Image';
import Table from 'react-bootstrap/Table';
import OverlayTrigger from 'react-bootstrap/OverlayTrigger';
import Tooltip from 'react-bootstrap/Tooltip';
import { IFIUserData } from '../../types/IFIData';
import Row from 'react-bootstrap/esm/Row';
import ProfileImageFemale from './../../assets/images/default-profile-picture-female.jpeg';
import ProfileImageMale from './../../assets/images/default-profile-picture.jpeg';
import {
  errorIcon,
  notApplicableIcon,
  questionIcon,
  rightArrowInCircle,
  tickIcon,
  warnIcon,
} from '../../helpers/svgIcons';
import { InsuranceTypes } from '../../enums/Insurance';
import Button from 'react-bootstrap/esm/Button';
import { createConsentRequest, getPathToDashboard } from '../../helpers/axios';
import { PageState } from '../../enums/PageStates';
import { IAuth, IUserProfileScores } from '../../types/IUser';

export type IUserProfile = {
  changePageState: (state: string) => void;
  fiData: IFIUserData;
  auth: IAuth;
};

enum iconCutOff {
  tick = 0.78,
  warn = 0.7,
  error = 0.3,
}

const UserProfile: React.FC<IUserProfile> = ({ changePageState, fiData, auth }) => {
  const dataTypeCount = 5;

  const getIconForScore = (score: number) => {
    if (score > iconCutOff.tick) {
      return tickIcon();
    } else if (score >= iconCutOff.warn) {
      return warnIcon();
    } else if (score > iconCutOff.error) {
      return errorIcon();
    } else {
      return notApplicableIcon();
    }
  };

  const handleCreateConsentSubmit = (event: React.MouseEvent<HTMLElement, MouseEvent>) => {
    event.preventDefault();
    changePageState(PageState.Loading);
    auth.user.jwt().then((jwt: string) => {
      createConsentRequest(fiData.phone, jwt)
        .then(response => {
          const consent_handle = response?.data?.consent_handle;
          if (consent_handle) {
            globalThis.window.location.href = `https://anumati.setu.co/${consent_handle}?redirect_url=${getPathToDashboard()}`;
          }
        })
        .catch(error => {
          console.log(error);
          changePageState(PageState.Error);
        });
    });
  };

  return (
    <>
      <Col className='mt-4' sm='12' md='10'>
        <Card className='border'>
          <Card.Body>
            <Row className='justify-content-center'>
              <Col sm='12' md='4' className='text-center'>
                {fiData.name.includes('Ramkrishna') && (
                  <Image style={{ marginTop: '20px' }} src={ProfileImageMale} height='150px' roundedCircle />
                )}
                {!fiData.name.includes('Ramkrishna') && (
                  <Image
                    style={{ marginTop: '20px' }}
                    src={ProfileImageFemale}
                    height='150px'
                    roundedCircle
                  />
                )}
              </Col>
              <Col sm='12' md='8' className='mt-3'>
                <Card.Title className='roboto-bold'>{fiData.name}</Card.Title>
                <table>
                  <tbody>
                    <tr>
                      <td style={{ width: 160 }}>Date of Birth</td>
                      <td>
                        :{' '}
                        <span style={{ paddingLeft: 30 }}>
                          {new Date(fiData.date_of_birth).toDateString()}
                        </span>
                      </td>
                    </tr>
                    <tr>
                      <td>Pancard</td>
                      <td>
                        : <span style={{ paddingLeft: 30 }}>{fiData.pancard}</span>
                      </td>
                    </tr>
                    <tr>
                      <td>Active Insurances</td>
                      <td>
                        :{' '}
                        <span style={{ paddingLeft: 30 }}>
                          {fiData.insurance.filter(insurance => insurance.account_id != '').length}
                        </span>
                      </td>
                    </tr>
                    <tr>
                      <td>Account Status</td>
                      <td>
                        : <span style={{ paddingLeft: 30 }}>Active {tickIcon()}</span>
                      </td>
                    </tr>
                    <tr>
                      <td>KYC Status</td>
                      <td>
                        :{' '}
                        <span style={{ paddingLeft: 30 }}>
                          {fiData.ckyc_compliance ? (
                            <span>Completed {tickIcon()}</span>
                          ) : (
                            <span>Incomplete {warnIcon()}</span>
                          )}
                        </span>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </Col>
              <Col sm='12' className='mt-4'>
                <Table className='mt-2' hover size='sm'>
                  <thead>
                    <tr>
                      <th>Affecting Factor</th>
                      <th>Score</th>
                      <th></th>
                    </tr>
                  </thead>
                  <tbody>
                    {getUserScoreList(fiData).map(scoreInfo => {
                      return (
                        <tr key={scoreInfo.title}>
                          <td>{scoreInfo.title}</td>
                          <td>
                            {scoreInfo.score.toFixed(2)}{' '}
                            <OverlayTrigger
                              placement='bottom'
                              overlay={<Tooltip id='button-tooltip'>{scoreInfo.explaination}</Tooltip>}>
                              <span>{questionIcon()}</span>
                            </OverlayTrigger>
                          </td>
                          <td>{getIconForScore(scoreInfo.score)}</td>
                        </tr>
                      );
                    })}
                  </tbody>
                </Table>
              </Col>
              {fiData.shared_data_sources < dataTypeCount && (
                <Col sm='12'>
                  <Button
                    onClick={event => handleCreateConsentSubmit(event)}
                    className='float-end btn-sm'
                    variant='outline-primary'>
                    Share more financial information for lower <br />
                    premiums {rightArrowInCircle('0 0 16 16')}
                  </Button>
                </Col>
              )}
            </Row>
          </Card.Body>
        </Card>
      </Col>
    </>
  );
};

const noRecordFound = 0;
const getUserScoreList = (fiData: IFIUserData): IUserProfileScores[] => [
  {
    title: 'Age Score',
    explaination: 'Younger clients score better for term insurance and pension plans.',
    score: fiData.age_score,
  },
  {
    title: 'Wealth Score',
    explaination: 'We suggest you plans based on your income and wealth in your account.',
    score: fiData.wealth_score,
  },
  {
    title: 'Medical Emergency Score',
    explaination: 'Higher score can qualify you for lower premiums on health insurance.',
    score:
      fiData.insurance.find(insurance => insurance.type == InsuranceTypes.MEDICAL_PLAN)?.score ||
      noRecordFound,
  },
  {
    title: 'Debt Score',
    explaination: 'Lower dept can qualify you for lower premiums on many plans.',
    score: fiData.debt_score,
  },
  {
    title: 'Travel Probablity',
    explaination: 'We suggest you plans based on your travel requirements and habits.',
    score:
      fiData.insurance.find(insurance => insurance.type == InsuranceTypes.TRAVEL_PLAN)?.score ||
      noRecordFound,
  },
  {
    title: 'Investment Health',
    explaination:
      'We check your investments (PPF, NPS etc), higher score means lower prenium on pension plan and life insurance.',
    score: fiData.investment_score,
  },
  {
    title: 'Motor Safety Score',
    explaination: 'Based on the safety and maintainence record of your vehicle, we adjust the premium.',
    score:
      fiData.insurance.find(insurance => insurance.type == InsuranceTypes.MOTOR_PLAN)?.score || noRecordFound,
  },
];

export { UserProfile };
