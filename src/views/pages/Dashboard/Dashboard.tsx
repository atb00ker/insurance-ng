import React, { useContext, useEffect, useMemo, useState } from 'react';
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import { getDashboardData } from '../../services/axios';
import SectionLoader from '../../components/ContentState/SectionLoader';
import { AuthContext } from '../../components/Auth/AuthProvider';
import { IAuth } from '../../interfaces/IUser';
import ServerRequestError from '../../components/ContentState/ServerRequestError';
import UserProfile from '../../components/FIData/UserProfile';
import { IFIData } from '../../interfaces/IFIData';
import InsuranceCard from '../../components/FIData/InsuranceCard';
import { InsuranceTypes } from '../../enums/FIData';

const Dashboard: React.FC = () => {
  const auth: IAuth = useContext(AuthContext);
  const [showError, setShowError] = useState(false);
  const [showLoader, setShowLoader] = useState(true);
  const [showProcessing, setShowProcessing] = useState(true);
  const [fiData, setFIData] = useState([] as unknown as IFIData);
  4;
  useEffect(() => {
    auth.user.jwt().then((jwt: string) => {
      getDashboardData(jwt)
        .then(response => {
          const data: IFIData = response.data;
          if (data.status) {
            setFIData(data);
            setShowProcessing(false);
          } else setShowProcessing(true);
          setShowLoader(false);
          setShowError(false);
        })
        .catch(error => {
          console.error(error);
          setShowLoader(false);
          setShowError(true);
        });
    });
  }, [auth.isReady]);

  const sortedFiInsuranceList = useMemo(() => {
    return (
      fiData?.data?.insurance?.sort((a, b) => {
        if (a.account_id || a.type == InsuranceTypes.ALL_PLAN) return -1;
        if (b.account_id || a.type == InsuranceTypes.ALL_PLAN) return 1;
        return b.score - a.score;
      }) || []
    );
  }, [fiData]);

  return (
    <Container>
      {!showError && !showLoader && (
        <>
          <Row className='mt-1 mb-2 justify-content-center'>
            <UserProfile fiData={fiData.data} />
          </Row>
          <Row className='mt-1 mb-5 justify-content-center'>
            {sortedFiInsuranceList.map(insurance => (
              <InsuranceCard key={insurance.type} insurance={insurance} />
            ))}
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
