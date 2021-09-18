import React, { useContext, useEffect, useMemo, useState } from 'react';
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import { getConsentStatus, getDashboardData } from '../../services/axios';
import {
  getInsurancePercents,
  getMedicalBills,
  getMotorTheftScore,
  getTotalWealth,
  getTravelledBills,
  prepareDataJson,
} from '../../services/fiDataParser';
import SectionLoader from '../../components/ContentState/SectionLoader';
import { AuthContext } from '../../components/Auth/AuthProvider';
import { IAuth } from '../../interfaces/IUser';
import { useHistory } from 'react-router-dom';
import { ConsentStatus, ConsentType } from '../../enums/ConsentInfo';
import { RouterPath } from '../../enums/RouterPath';
import ServerRequestError from '../../components/ContentState/ServerRequestError';
import { IFIData } from '../../interfaces/IFIData';
import { AccountType } from '../../enums/AccountType';
import { InsuranceInfo } from '../../enums/Insurance';
import InsuranceFIData from '../../components/FIData/InsuranceFIData';
import InsuranceSuggestions from '../../components/FIData/InsuranceSuggestions';
import UserProfile from '../../components/FIData/UserProfile';
import { IUserMetrics } from '../../interfaces/IUser';
import {
  getAgeScore,
  getDeptScore,
  getInvestmentScore,
  getMedicalBillScore,
  getMotorInsuranceScore,
  getTravelBillScore,
  getWealthScore,
} from '../../services/scoringSystem';

const Dashboard: React.FC = () => {
  const history = useHistory();
  const auth: IAuth = useContext(AuthContext);
  const [showError, setShowError] = useState(false);
  const [showLoader, setShowLoader] = useState(true);
  const [dashboardData, setDashboardData] = useState([] as unknown as IFIData[]);
  let existingInsurance: any[] = [];

  const userInfo: IUserMetrics = useMemo(() => {
    if (dashboardData.length == 0) return {} as IUserMetrics;
    const accountInfo = dashboardData[0].RahasyaData[0].data;
    const dataUserInfo = accountInfo.children
      .find((item: any) => item.name == ConsentType.Profile)
      .children.find((item: any) => item.name == 'Holders').children[0].attributes;

    const totalWealth: number = getTotalWealth(dashboardData);
    const medicalBills: number = getMedicalBills(dashboardData);
    const deptScore: number = getMedicalBills(dashboardData);
    const travelledTrips: number = getTravelledBills(dashboardData);
    const motorTheftScore: number = getMotorTheftScore(dashboardData);
    const [sip_percent, mutualfund_percent, nps_percent, ppf_percent, epf_percent] =
      getInsurancePercents(dashboardData);

    const userInfo: IUserMetrics = {
      name: dataUserInfo.name,
      datapoint: {
        age: {
          title: 'Date of Birth',
          explaination: 'Younger clients score better for term insurance and pension plans.',
          value: dataUserInfo.dob,
          score: getAgeScore(dataUserInfo.dob),
        },
        wealth: {
          title: 'Wealth Score',
          explaination: 'We suggest you plans based on your income and wealth in your account.',
          value: totalWealth,
          score: getWealthScore(totalWealth),
        },
        health: {
          title: 'Medical Emergency Score',
          explaination: 'Higher score can qualify you for lower premiums on health insurance.',
          value: medicalBills,
          score: getMedicalBillScore(medicalBills),
        },
        dept: {
          title: 'Dept Score',
          explaination: 'Lower dept can qualify you for lower premiums on many plans.',
          value: deptScore,
          score: getDeptScore(deptScore),
        },
        travel: {
          title: 'Travel Probablity',
          explaination: 'We suggest you plans based on your travel requirements and habits.',
          value: travelledTrips,
          score: getTravelBillScore(travelledTrips),
        },
        investment: {
          title: 'Investment Health',
          explaination:
            'We check your investments (PPF, NPS etc), higher score means lower prenium on pension plan and life insurance.',
          sip_percent: sip_percent,
          mutualfund_percent: mutualfund_percent,
          nps_percent: nps_percent,
          ppf_percent: ppf_percent,
          epf_percent: epf_percent,
          score: getInvestmentScore(sip_percent, mutualfund_percent, nps_percent, ppf_percent, epf_percent),
        },
        motor: {
          title: 'Motor Safety Score',
          explaination: 'Based on the safety and maintainence record of your vehicle, we adjust the premium.',
          value: motorTheftScore,
          score: getMotorInsuranceScore(motorTheftScore),
        },
      },
    };
    return userInfo;
  }, [dashboardData]);

  useEffect(() => {
    auth.user.jwt().then((jwt: string) => {
      getConsentStatus(jwt)
        .then(response => {
          setShowError(false);
          const status = response.data['status'];
          if (
            ![
              ConsentStatus.UserConsentsReady,
              ConsentStatus.UserConsentsFetched,
              ConsentStatus.UserConsentsActive,
            ].includes(status)
          )
            history.push(RouterPath.CreateConsent);
          else {
            getDashboardData(jwt)
              .then(response => {
                const dataJson = prepareDataJson(response.data);
                console.log(dataJson);
                setDashboardData(dataJson);
                setShowLoader(false);
                setShowError(false);
              })
              .catch(error => {
                console.error(error);
                setShowLoader(false);
                setShowError(true);
              });
          }
        })
        .catch(error => {
          console.error(error);
          history.push(RouterPath.CreateConsent);
        });
    });
  }, [auth.isReady]);

  const suggestInsuranceList = () => {
    const morePolicies = InsuranceInfo.filter((item: any) => !existingInsurance.includes(item.policyType));
    return morePolicies.map(item => (
      <InsuranceSuggestions key={item.policyType} data={item} userInfo={userInfo} />
    ));
  };

  return (
    <Container>
      {!showError && !showLoader && (
        <>
          <Row className='mt-1 mb-2 justify-content-center'>{<UserProfile userInfo={userInfo} />}</Row>
          <Row className='mt-1 mb-5 justify-content-center'>
            {dashboardData.map(fiData => {
              return fiData.RahasyaData.map(item => {
                const accountData: any = item.data;
                if (accountData.attributes.type == AccountType.Insurance) {
                  existingInsurance.push(
                    accountData.children.find((item: any) => item.name == ConsentType.Summary).attributes
                      .policyType,
                  );
                  return (
                    <InsuranceFIData
                      key={accountData.attributes.linkedAccRef}
                      data={accountData}
                      userInfo={userInfo}
                    />
                  );
                }
              });
            })}
            {suggestInsuranceList()}
          </Row>
        </>
      )}

      {!!showError && (
        <Row className='mt-4'>
          <Col sm='12'>
            <ServerRequestError height='500px' imgHeight='250px' width='100%' />
          </Col>
        </Row>
      )}

      {!!showLoader && (
        <Row className='mt-4'>
          <Col sm='12'>
            <SectionLoader height='500px' width='100%' />
          </Col>
        </Row>
      )}
    </Container>
  );
};

export default Dashboard;
