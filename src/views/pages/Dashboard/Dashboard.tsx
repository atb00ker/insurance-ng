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
import { IFIData, RahasyaList } from '../../interfaces/IFIData';
import { AccountType } from '../../enums/AccountType';
import InsuranceGraphs from '../../components/Graphs/InsuranceGraphs';

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
          <>
            {dashboardData.map(fiData => {
              return fiData.RahasyaData.map(item => {
                const accountData: any = item.data;
                if(accountData["attributes"]["type"] == AccountType.Insurance)
                  return <InsuranceGraphs data={accountData} />
              })
            })
            }
          </>
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
