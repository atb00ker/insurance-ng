import React, { useContext, useEffect, useState } from 'react';
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import { getConsentStatus, getDashboardData } from '../../services/axios';
import { prepareDataJson } from '../../services/fiDataParser';
import SectionLoader from '../../components/ContentState/SectionLoader';
import { AuthContext } from '../../components/Auth/AuthProvider';
import { IAuth } from '../../interfaces/IAuth';
import { useHistory } from 'react-router-dom';
import { ConsentStatus } from '../../enums/ConsentStatus';
import { RouterPath } from '../../enums/RouterPath';
import ServerRequestError from '../../components/ContentState/ServerRequestError';
import { IFIData } from '../../interfaces/IFIData';
import { AccountType } from '../../enums/AccountType';
import { InsuranceInfo } from '../../enums/Insurance';
import InsuranceFIData from '../../components/FIData/InsuranceFIData';
import InsuranceSuggestions from '../../components/FIData/InsuranceSuggestions';

const Dashboard: React.FC = () => {
  const history = useHistory();
  const auth: IAuth = useContext(AuthContext);
  const [showError, setShowError] = useState(false);
  const [showLoader, setShowLoader] = useState(true);
  const [dashboardData, setDashboardData] = useState([] as unknown as IFIData[]);

  useEffect(() => {
    auth.user.jwt().then(jwt => {
      getConsentStatus(jwt).then(response => {
          setShowError(false);
          const status = response.data["status"];
          if (![ConsentStatus.UserConsentsReady,
                ConsentStatus.UserConsentsFetched,
                ConsentStatus.UserConsentsActive].includes(status))
            history.push(RouterPath.CreateConsent);
          else {
            getDashboardData(jwt).then(response => {
              const dataJson = prepareDataJson(response.data);
              setDashboardData(dataJson);
              setShowLoader(false);
              setShowError(false);
            }).catch(error => {
              console.error(error);
              setShowLoader(false);
              setShowError(true);
            })
          }
        }).catch(error => {
          console.error(error);
          history.push(RouterPath.CreateConsent);
        });
      });
  }, [auth.isReady])

  return (
    <Container>
      <Row className='mt-4'>
      {!showLoader &&
        <Col sm='12'>
          {!showError &&
          <Row className='mt-1 mb-2'>
            {dashboardData.map(fiData => {
              let existingInsurance: any[] = []
              let dataCards: any[] = []

              const existingData = fiData.RahasyaData.map(item => {
                const accountData: any = item.data;
                if(accountData["attributes"]["type"] == AccountType.Insurance) {
                  existingInsurance.push(accountData.children[0].attributes.policyType)
                  return <InsuranceFIData key={accountData["attributes"]["linkedAccRef"]} data={accountData} />
                }
              });
              dataCards.push(...existingData);
              const morePolicies = InsuranceInfo.filter((item: any) => !existingInsurance.includes(item.policyType));
              const newCards = morePolicies.map(item => <InsuranceSuggestions key={item.policyType} data={item} />)
              dataCards.push(...newCards);
              return dataCards;
            })}
          </Row>
          }
          {!!showError &&
            <ServerRequestError height='500px' imgHeight='250px' width='100%' />
          }
        </Col>
      }
      {!!showLoader &&
        <SectionLoader height='500px' width='100%' />
      }
      </Row>
    </Container>
  );
};

export default Dashboard;
