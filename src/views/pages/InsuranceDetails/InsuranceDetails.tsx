import React, { useContext, useEffect, useMemo, useState } from 'react';
import Container from 'react-bootstrap/Container';
import Col from 'react-bootstrap/Col';
import Row from 'react-bootstrap/Row';
import { useHistory, useLocation } from 'react-router-dom';
import { IAuth } from '../../interfaces/IUser';
import { AuthContext } from '../../components/Auth/AuthProvider';
import { getDashboardData } from '../../services/axios';
import { IFIData, IFIInsurance } from '../../interfaces/IFIData';
import ServerRequestError from '../../components/ContentState/ServerRequestError';
import SectionLoader from '../../components/ContentState/SectionLoader';
import Table from 'react-bootstrap/esm/Table';
import Button from 'react-bootstrap/esm/Button';
import { tickIcon } from '../../services/svgIcons';
import { Chart } from "react-google-charts";
import { RouterPath } from '../../enums/UrlPath';

const InsuranceDetails: React.FC = () => {
  const history = useHistory();
  const auth: IAuth = useContext(AuthContext);
  const location = useLocation();
  const [insuranceInfo, setInsuranceInfo] = useState({} as IFIInsurance);
  const [fiData, setFiData] = useState({} as IFIData);
  const fiDataFromHistory = location.state as IFIData;
  const [showError, setShowError] = useState(false);
  const [showLoader, setShowLoader] = useState(true);

  const getDataFromServer = (path: string) => {
    auth.user.jwt().then((jwt: string) => {
      getDashboardData(jwt)
        .then(response => {
          const data: IFIData = response?.data;
          if (data?.status) {
            setShowError(false);
            setShowLoader(false);
            getInsuranceFromFIData(data, path)
          } else {
            setShowError(true);
            setShowLoader(false);
          }
        })
        .catch(error => {
          console.error(error);
          setShowError(true);
          setShowLoader(false);
        });
    });
  };

  const getInsuranceFromFIData = (fiData: IFIData, path: string) => {
    const insurance = fiData.data.insurance.find(insurance =>
              insurance.uuid === path.replace('/insurance/', '')) || {} as IFIInsurance;
    setInsuranceInfo(insurance);
    setFiData(fiData);
  };

  useEffect(() => {
    if (!fiDataFromHistory) {
      getDataFromServer(location.pathname);
      return
    }
    setShowError(false);
    setShowLoader(false);
    getInsuranceFromFIData(fiDataFromHistory, location.pathname)
  }, [fiDataFromHistory, auth]);

  const premiumGraph: [number, number][] = useMemo(() => {
    let factor = 5, deduction = 0, graphColumns: [number, number][] = [];
    for (let index = 0; index < 10; index++) {
      graphColumns.push([index, insuranceInfo.offer_premium - deduction*factor]);
      deduction += deduction + Math.floor(Math.random() * 2) + 1;
    }
    return graphColumns;
  }, [insuranceInfo]);

  const coverGraph: [number, number][] = useMemo(() => {
    let factor = 5, deduction = 0, graphColumns: [number, number][] = [];
    for (let index = 0; index < 10; index++) {
      graphColumns.push([index, insuranceInfo.offer_cover - deduction*factor]);
      deduction += deduction + Math.floor(Math.random() * 2) + 1;
    }
    return graphColumns;
  }, [insuranceInfo]);

  return (
    <Container>
      {!showError && !showLoader && (
        <>
          <Row className='mt-5 mb-2 justify-content-center'>
            <Col sm='1'>
              <Button
                onClick={() => {
                  history.push({
                    pathname: RouterPath.Dashboard,
                    state: fiData,
                  });
                  return;
                }}
                style={{ minWidth: '70px' }}
                variant='outline-primary'
                className='btn-sm me-2 mb-2 buttons-height'
              >
                <svg
                  xmlns='http://www.w3.org/2000/svg'
                  width='16'
                  height='16'
                  fill='currentColor'
                  className='bi bi-arrow-return-left'
                  viewBox='0 0 16 16'
                >
                  <path
                    fillRule='evenodd'
                    d='M14.5 1.5a.5.5 0 0 1 .5.5v4.8a2.5 2.5 0 0 1-2.5 2.5H2.707l3.347 3.346a.5.5 0 0 1-.708.708l-4.2-4.2a.5.5 0 0 1 0-.708l4-4a.5.5 0 1 1 .708.708L2.707 8.3H12.5A1.5 1.5 0 0 0 14 6.8V2a.5.5 0 0 1 .5-.5z'
                  />
                </svg>{' '}
                {'  '}
                Back
              </Button>
            </Col>
            <Col sm='10'>
              <h2 className='text-center'>{insuranceInfo.title}</h2>
            </Col>
            <Col sm='1' />
          </Row>
          <Row className='justify-content-center'>
            <Col sm='8' md='6' lg='4'>
              <Table className='mt-4' size='md'>
                <tbody>
                  <tr>
                    <td>Cover</td>
                    <td>: {insuranceInfo.offer_cover}</td>
                    <td></td>
                  </tr>
                  <tr>
                    <td>Yealy Premium</td>
                    <td>: {insuranceInfo.offer_premium}</td>
                    <td></td>
                  </tr>
                  <tr>
                    <td>Purchase</td>
                    {insuranceInfo.score >= 0.8 && (
                      <>
                        <td>: Pre-approved {tickIcon()} </td>
                        <td>
                          <Button variant='outline-primary' className='btn-sm'>
                            {' '}
                            Buy{' '}
                          </Button>
                        </td>
                      </>
                    )}
                    {insuranceInfo.score < 0.8 && (
                      <>
                        <td>
                          :{' '}
                          <Button variant='outline-primary' className='btn-sm'>
                            {' '}
                            Contact Us{' '}
                          </Button>
                        </td>
                        <td></td>
                      </>
                    )}
                  </tr>
                </tbody>
              </Table>
            </Col>
            <Col className='mt-4' sm='12' md='10'>
              {insuranceInfo.description} Voluptas nihil sed laborum sequi quaerat veritatis hic facere eaque culpa iusto exercitationem tenetur rerum, incidunt, sapiente optio suscipit aspernatur delectus molestias explicabo quisquam ratione odit, repudiandae dolor beatae! Harum!
              Fugiat magnam alias veritatis quibusdam iure! Inventore dolorem quaerat illum blanditiis similique facere deserunt neque molestias consequuntur eveniet numquam sapiente accusamus ab consequatur eos, libero eius nulla. Beatae, debitis maiores!
              Corporis, veritatis. Id rem tenetur enim, aspernatur fugit fugiat ipsum ex possimus amet veniam vel odio tempore modi dolores culpa laboriosam? Nulla, voluptate. Dolorum facilis quisquam id enim molestiae magni?
              Consequuntur nam eos, voluptates alias ipsa error libero accusamus magnam provident, nihil inventore dolores obcaecati saepe, blanditiis amet incidunt. Doloribus explicabo nemo dolore saepe nihil, repellendus temporibus iste quam inventore.
              Ipsam fugiat recusandae maxime vero eius aut minus reprehenderit delectus asperiores. Consequuntur blanditiis repudiandae hic libero incidunt perferendis quos? Commodi repellat reiciendis ipsa aperiam sequi, illum quae. Fuga, eaque fugiat.
            </Col>
          </Row>
          <Row className='mt-3 mb-2 justify-content-center'>
            <Col sm='12' md='8' className='mt-3'>
              <Table bordered>
                <thead>
                  <tr>
                    <th className="text-center">#</th>
                    <th className="text-center">Clauses</th>
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td>1.</td>
                    <td>In case of emergency, the transaction will be completed with 24 hours.</td>
                  </tr>
                  <tr>
                    <td>2.</td>
                    <td>In case of emergency, the transaction will be completed with 24 hours.</td>
                  </tr>
                  <tr>
                    <td>3.</td>
                    <td>In case of emergency, the transaction will be completed with 24 hours.</td>
                  </tr>
                </tbody>
              </Table>
            </Col>
          </Row>
          <Row className='mb-4 justify-content-center'>
            <Col sm='12' md='10'>
              Lorem ipsum dolor sit amet consectetur adipisicing elit. Totam aspernatur labore similique eaque
              officiis! Aliquam maiores eaque excepturi praesentium esse consectetur, ab impedit, amet
              veritatis culpa reprehenderit delectus consequatur blanditiis?
            </Col>
          </Row>
          <Row className='mb-4 justify-content-center'>
            <Col sm='12' md='6'>
              <Chart
                height={'400px'}
                chartType="LineChart"
                className="mx-auto"
                loader={<div>Loading Chart</div>}
                data={[
                  ['x', 'premium'],
                  ...premiumGraph
                ]}
                options={{
                  hAxis: {
                    title: 'Time (in years)',
                  },
                  vAxis: {
                    title: 'Premium',
                    viewWindowMode:'explicit',
                    viewWindow: {
                      max: insuranceInfo.offer_premium + insuranceInfo.offer_premium * 0.1,
                      min: insuranceInfo.offer_premium - insuranceInfo.offer_premium * 0.1
                    }
                  }
                }}
                rootProps={{ 'data-testid': 'test-premium-chart' }}
              />
            </Col>
            <Col sm='12' md='6'>
              <Chart
                height={'400px'}
                chartType="LineChart"
                className="mx-auto"
                loader={<div>Loading Chart</div>}
                data={[
                  ['x', 'cover'],
                  ...coverGraph
                ]}
                options={{
                  hAxis: {
                    title: 'Time (in years)',
                  },
                  vAxis: {
                    title: 'Cover',
                    viewWindowMode:'explicit',
                    viewWindow: {
                      max: insuranceInfo.offer_cover + insuranceInfo.offer_cover * 0.001,
                      min: insuranceInfo.offer_cover - insuranceInfo.offer_cover * 0.001
                    }
                  },
                  colors: ["red"]
                }}
                rootProps={{ 'data-testid': 'test-premium-chart' }}
              />
            </Col>
          </Row>
          <Row className='mb-5 justify-content-center'>
            <Col className='mb-5' sm='12' md='10'>
              Tenetur ab corrupti delectus? Officia magni
              eligendi vel itaque quidem eaque laboriosam magnam tempore exercitationem atque aliquam, culpa
              amet, dolor hic similique in qui quaerat velit! Odit obcaecati voluptates asperiores! Et non
              labore omnis est obcaecati quo voluptas magni repudiandae suscipit id odit provident quos
              maiores magnam dignissimos facere itaque cum, impedit aspernatur quisquam vero quod ipsa
              numquam? Harum, animi? Odit quaerat quae, ea mollitia aut quibusdam labore inventore temporibus
              accusamus perferendis cum saepe natus quisquam aperiam nam nobis quo! Quas aperiam veritatis
              maxime a sint ullam adipisci ipsam odio. Pariatur at quis porro, ipsum laborum modi. In
              voluptates cupiditate ad vero ipsam! Quis ex quod voluptate, repudiandae doloribus pariatur,
              veritatis laboriosam dolores facere quasi ab non tempora velit laudantium?
            </Col>
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

export default InsuranceDetails;
