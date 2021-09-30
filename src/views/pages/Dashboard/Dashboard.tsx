import React, { useContext, useEffect, useMemo, useState } from 'react';
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import { getDashboardData } from '../../services/axios';
import SectionLoader from '../../components/ContentState/SectionLoader';
import { AuthContext } from '../../components/Auth/AuthProvider';
import { IAuth } from '../../interfaces/IUser';
import useWebSocket, { ReadyState } from 'react-use-websocket';
import ServerRequestError from '../../components/ContentState/ServerRequestError';
import UserProfile from '../../components/FIData/UserProfile';
import { IFIData } from '../../interfaces/IFIData';
import InsuranceCard from '../../components/FIData/InsuranceCard';
import { InsuranceTypes } from '../../enums/Insurance';
import FiDataWait from '../../components/ContentState/FIDataWait';
import { ServerPath } from '../../enums/UrlPath';

const Dashboard: React.FC = () => {
  const auth: IAuth = useContext(AuthContext);
  const [showError, setShowError] = useState(false);
  const [showLoader, setShowLoader] = useState(true);
  const [showProcessing, setShowProcessing] = useState(false);
  const [fiData, setFIData] = useState({} as IFIData);

  const port = location.protocol === 'http:' ? process.env.REACT_APP_PORT : '443';
  const protocol = location.protocol === 'http:' ? 'ws:' : 'wss:';
  const socketUrl = new URL(ServerPath.DataWebsocket, `${protocol}//${location.hostname}:${port}`);
  const { sendMessage, lastMessage, readyState } = useWebSocket(socketUrl.toString());

  useEffect(() => {
    // getDataFromServer();
    if (readyState === ReadyState.CLOSED) {
      setShowError(true);
      setShowLoader(false);
      setShowProcessing(false);
    } else if (auth.user.id && readyState === ReadyState.OPEN) sendMessage(auth.user.id);
  }, [auth.isReady, readyState]);

  useEffect(() => {
    if (lastMessage?.data === 'false') {
      setShowLoader(false);
      setShowProcessing(true);
    } else if (lastMessage?.data === 'true') getDataFromServer();
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
          const data: IFIData = response?.data;
          if (data?.status) {
            setFIData(data);
            setShowProcessing(false);
          } else {
            setShowProcessing(true);
            setTimeout(() => getDataFromServer(), 5000);
          }
          setShowLoader(false);
          setShowError(false);
        })
        .catch(error => {
          console.error(error);
          setShowError(true);
          setShowLoader(false);
          setShowProcessing(false);
        });
    });
  };

  return (
    <Container>
      {!showError && !showLoader && !showProcessing && (
        <>
          <Row className='mt-1 mb-2 justify-content-center'>
            <UserProfile fiData={fiData.data} />
          </Row>
          <Row className='mt-1 mb-5 justify-content-center'>
            {sortedFiInsuranceList.map(insurance => (
              <InsuranceCard key={insurance.type} fiData={fiData} insurance={insurance} />
            ))}
          </Row>
        </>
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
