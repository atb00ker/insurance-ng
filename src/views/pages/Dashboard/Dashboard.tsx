import React, { useContext, useEffect, useMemo, useState } from 'react';
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import { getDashboardData, mockConsentNotification } from '../../helpers/axios';
import SectionLoader from '../../components/ContentState/SectionLoader';
import { AuthContext } from '../../components/Auth/AuthProvider';
import { IAuth } from '../../interfaces/IUser';
import useWebSocket, { ReadyState } from 'react-use-websocket';
import ServerRequestError from '../../components/ContentState/ServerRequestError';
import UserProfile from '../../components/Dashboard/UserProfile';
import { IFIData } from '../../interfaces/IFIData';
import InsuranceCard from '../../components/Dashboard/InsuranceCard';
import { InsuranceTypes } from '../../enums/Insurance';
import FiDataWait from '../../components/ContentState/FIDataWait';
import { ServerPath } from '../../enums/UrlPath';
import NoValidConsent from '../../components/ContentState/InvalidConsent';
import { PageState } from '../../enums/PageStates';

const Dashboard: React.FC = () => {
  const auth: IAuth = useContext(AuthContext);
  const [showError, setShowError] = useState(false);
  const [showLoader, setShowLoader] = useState(true);
  const [showProcessing, setShowProcessing] = useState(false);
  const [showShareInfo, setShowShareInfo] = useState(false);
  const [fiData, setFIData] = useState({} as IFIData);

  const port = location.protocol === 'http:' ? process.env.REACT_APP_PORT : '443';
  const protocol = location.protocol === 'http:' ? 'ws:' : 'wss:';
  const socketUrl = new URL(ServerPath.DataWebsocket, `${protocol}//${location.hostname}:${port}`);
  const { sendMessage, lastMessage, readyState } = useWebSocket(socketUrl.toString());

  useEffect(() => {
    // Notification Hack:
    // This is only a hack to mock Setu Notification.
    // This is not required in a real life application, only added
    // reliability for the hackathon while the notification endpoint
    // is unstable on Setu's end.
    auth.user.jwt().then((jwt: string) => {
      setTimeout(() => {
        mockConsentNotification(jwt).then(_ => {
          getDataFromServer();
        });
      }, 1000);
    });
  }, [auth.isReady]);

  useEffect(() => {
    // If the websocket connection fails for some reason,
    // eg. machine doesn't support it, then we fallback to
    // good ol' periodic polling.
    if (readyState === ReadyState.CLOSED) {
      getDataFromServer();
    } else if (auth.user.id && readyState === ReadyState.OPEN) {
      sendMessage(auth.user.id);
    }
  }, [auth.isReady, readyState]);

  useEffect(() => {
    if (['data-not-shared', 'consent-not-started'].includes(lastMessage?.data)) {
      changePageState(PageState.Waiting);
    } else if (lastMessage?.data === 'false') {
      changePageState(PageState.Processing);
    } else if (lastMessage?.data === 'true') {
      getDataFromServer();
    }
  }, [lastMessage]);

  const sortedFiInsuranceList = useMemo(() => {
    const insuranceList = fiData?.data?.insurance || [];
    return [...insuranceList].sort((a, b) => {
      if (a.type == InsuranceTypes.ALL_PLAN) return -1;
      if (b.type == InsuranceTypes.ALL_PLAN) return 1;
      if (a.account_id) return -1;
      if (b.account_id) return 1;
      return b.score - a.score;
    });
  }, [fiData]);

  const getDataFromServer = () => {
    auth.user.jwt().then((jwt: string) => {
      getDashboardData(jwt)
        .then(response => {
          // Response validation
          if (response === undefined) {
            changePageState(PageState.Error);
            return;
          }

          // If request was successful.
          const data: IFIData = response?.data;
          if (data?.status) {
            if (data.data.name === '') {
              // Sometimes the data is corrupt.
              // We don't want to show that.
              changePageState(PageState.Waiting);
            } else {
              setFIData(data);
              changePageState(PageState.Data);
            }
          } else {
            changePageState(PageState.Processing);
            setTimeout(() => getDataFromServer(), 5000);
          }
        })
        .catch(error => {
          if (error?.response?.data?.error == 'record not found') changePageState(PageState.Waiting);
          else changePageState(PageState.Error);
        });
    });
  };

  const changePageState = (state: string) => {
    setShowShareInfo(state == PageState.Waiting);
    setShowError(state == PageState.Error);
    setShowProcessing(state == PageState.Processing);
    setShowLoader(state == PageState.Loading);
  };

  return (
    <Container>
      {!showError && !showLoader && !showProcessing && !showShareInfo && (
        <>
          <Row className='mt-1 mb-2 justify-content-center'>
            <UserProfile auth={auth} changePageState={changePageState} fiData={fiData.data} />
          </Row>
          <Row className='mt-1 mb-5 justify-content-center'>
            {sortedFiInsuranceList.map(insurance => (
              <InsuranceCard key={insurance.type} fiData={fiData} insurance={insurance} />
            ))}
          </Row>
        </>
      )}

      {!!showShareInfo && (
        <Row className='mt-5'>
          <Col sm='12'>
            <NoValidConsent height='450px' imgHeight='400px' width='100%' />
          </Col>
        </Row>
      )}

      {!!showProcessing && (
        <Row className='mt-4'>
          <Col sm='12'>
            <FiDataWait height='450px' imgHeight='400px' width='100%' />
          </Col>
        </Row>
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
