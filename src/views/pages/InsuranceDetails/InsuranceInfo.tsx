import React, { useContext, useEffect, useMemo, useState } from 'react';
import Container from 'react-bootstrap/Container';
import Col from 'react-bootstrap/Col';
import Row from 'react-bootstrap/Row';
import { useHistory, useLocation } from 'react-router-dom';
import { IAuth } from '../../interfaces/IUser';
import { AuthContext } from '../../components/Auth/AuthProvider';
import { createPurchaseRequest, getDashboardData } from '../../helpers/axios';
import { IFIData, IFIInsurance } from '../../interfaces/IFIData';
import ServerRequestError from '../../components/ContentState/ServerRequestError';
import SectionLoader from '../../components/ContentState/SectionLoader';
import Table from 'react-bootstrap/esm/Table';
import Button from 'react-bootstrap/esm/Button';
import { tickIcon } from '../../helpers/svgIcons';
import { Chart } from 'react-google-charts';
import InsuranceInfoTitle from './InsuranceInfoTitle';
import { RouterPath } from '../../enums/UrlPath';

const InsuranceInfo: React.FC = () => {
  const history = useHistory();
  const auth: IAuth = useContext(AuthContext);
  const location = useLocation();
  const [insuranceInfo, setInsuranceInfo] = useState({} as IFIInsurance);
  const [fiData, setFiData] = useState({} as IFIData);
  const fiDataFromHistory = location.state as IFIData;
  const [showError, setShowError] = useState(false);
  const [showLoader, setShowLoader] = useState(true);
  const uuid = location.pathname.replace('/insurance/', '');

  useEffect(() => {
    if (!fiDataFromHistory) {
      getDataFromServer();
      return;
    }
    setShowError(false);
    setShowLoader(false);
    setInsuranceState(fiDataFromHistory, uuid);
  }, [fiDataFromHistory, auth]);

  const premiumGraph: [number, number][] = useMemo(() => {
    let deductionRate = insuranceInfo.yoy_deduction_rate,
      graphColumns: [number, number][] = [];
    for (let index = 1; index < 6; index++) {
      graphColumns.push([index, insuranceInfo.offer_premium - deductionRate]);
      deductionRate = deductionRate + deductionRate;
    }
    return graphColumns;
  }, [insuranceInfo]);

  const coverGraph: [number, number][] = useMemo(() => {
    let incrementRate = insuranceInfo.yoy_deduction_rate,
      graphColumns: [number, number][] = [];
    for (let index = 1; index < 6; index++) {
      graphColumns.push([index, insuranceInfo.offer_cover + incrementRate]);
      incrementRate = incrementRate + incrementRate;
    }
    return graphColumns;
  }, [insuranceInfo]);

  const setInsuranceState = (fiData: IFIData, uuid: string) => {
    const insurance =
      fiData.data.insurance.find(insurance => insurance.uuid === uuid) || ({} as IFIInsurance);
    setInsuranceInfo(insurance);
    setFiData(fiData);
    history.push({
      pathname: RouterPath.InsuranceDetails.replace(':insurance_uuid', insurance.uuid),
      state: fiData,
    });
  };

  const getDataFromServer = () => {
    auth.user.jwt().then((jwt: string) => {
      changeStateOnDataResponse(getDashboardData(jwt));
    });
  };

  const startPurchaseProcess = (uuid: string) => {
    auth.user.jwt().then((jwt: string) => {
      changeStateOnDataResponse(createPurchaseRequest(uuid, jwt));
    });
  };

  const changeStateOnDataResponse = (dataResponse: Promise<any>) => {
    dataResponse
      .then((response: any) => {
        const data: IFIData = response?.data;
        if (data?.status) {
          // TODO: Add common function
          setShowError(false);
          setShowLoader(false);
          setInsuranceState(data, uuid);
        } else {
          setShowError(true);
          setShowLoader(false);
        }
      })
      .catch((error: any) => {
        console.error(error);
        setShowError(true);
        setShowLoader(false);
      });
  };

  const getClauseTableColumns = (insuranceInfo: IFIInsurance): React.ReactNode => {
    const noOfClauses =
      insuranceInfo.clauses?.length > insuranceInfo.current_clauses?.length
        ? insuranceInfo.clauses?.length
        : insuranceInfo.current_clauses?.length || insuranceInfo.clauses?.length;
    let tableRow = [];
    for (let index = 0; index < noOfClauses; index++) {
      tableRow.push(
        <tr key={index}>
          <td>{index + 1}.</td>
          {!insuranceInfo.is_insurance_ng_acct && insuranceInfo.current_clauses?.length && (
            <td
              dangerouslySetInnerHTML={{
                __html: insuranceInfo.current_clauses[index]?.replace(
                  /(\d+%?)|(\S[A-Z]+\S)/g,
                  function (value) {
                    if (value === 'SLA') return value;
                    return `<span class="text-danger">${value}</span>`;
                  },
                ),
              }}
            ></td>
          )}
          <td
            dangerouslySetInnerHTML={{
              __html: insuranceInfo.clauses[index]?.replace(/(\d+%?)|(\S[A-Z]+\S)/g, function (value) {
                if (value === 'SLA') return value;
                return `<span class="text-success">${value}</span>`;
              }),
            }}
          ></td>
        </tr>,
      );
    }
    return tableRow;
  };

  return (
    <Container>
      {!showError && !showLoader && (
        <>
          <Row className='mt-5 mb-2 justify-content-center'>
            <InsuranceInfoTitle fiData={fiData} insuranceInfo={insuranceInfo} />
          </Row>
          <Row className='justify-content-center'>
            <Col id='topTable' sm='8' md='6' lg='5'>
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
                  {!insuranceInfo.is_insurance_ng_acct && (
                    <tr>
                      <td>Purchase</td>
                      {insuranceInfo.score >= 0.8 && (
                        <>
                          <td>: Pre-approved {tickIcon()} </td>
                          <td>
                            <Button
                              variant='outline-primary'
                              className='btn-sm'
                              onClick={() => startPurchaseProcess(insuranceInfo.uuid)}
                            >
                              Buy
                            </Button>
                          </td>
                        </>
                      )}
                      {insuranceInfo.score < 0.8 && (
                        <>
                          <td>
                            :{' '}
                            <Button
                              variant='outline-primary'
                              className='btn-sm'
                              onClick={() => startPurchaseProcess(insuranceInfo.uuid)}
                            >
                              Talk to an Agent
                            </Button>
                          </td>
                          <td></td>
                        </>
                      )}
                    </tr>
                  )}
                  {insuranceInfo.is_insurance_ng_acct && (
                    <tr>
                      <td>Status</td>
                      {insuranceInfo.is_active && (
                        <>
                          <td>: Active {tickIcon()} </td>
                          <td></td>
                        </>
                      )}
                      {!insuranceInfo.is_active && (
                        <>
                          <td>
                            :{' '}
                            <Button
                              disabled
                              variant='outline-secondary'
                              className='btn-sm'
                              onClick={() => startPurchaseProcess(insuranceInfo.uuid)}
                            >
                              We will contact you soon
                            </Button>
                          </td>
                          <td></td>
                        </>
                      )}
                    </tr>
                  )}
                </tbody>
              </Table>
            </Col>
            <Col id='introText' className='mt-4' sm='12' md='10'>
              {insuranceInfo.description} Voluptas nihil sed laborum sequi quaerat veritatis hic facere eaque
              culpa iusto exercitationem tenetur rerum, incidunt, sapiente optio suscipit aspernatur delectus
              molestias explicabo quisquam ratione odit, repudiandae dolor beatae! Harum! Fugiat magnam alias
              veritatis quibusdam iure! Inventore dolorem quaerat illum blanditiis similique facere deserunt
              neque molestias consequuntur eveniet numquam sapiente accusamus ab consequatur eos, libero eius
              nulla. Beatae, debitis maiores! Corporis, veritatis. Id rem tenetur enim, aspernatur fugit
              fugiat ipsum ex possimus amet veniam vel odio tempore modi dolores culpa laboriosam? Nulla,
              voluptate. Dolorum facilis quisquam id enim molestiae magni? Consequuntur nam eos, voluptates
              alias ipsa error libero accusamus magnam provident, nihil inventore dolores obcaecati saepe,
              blanditiis amet incidunt. Doloribus explicabo nemo dolore saepe nihil, repellendus temporibus
              iste quam inventore. Ipsam fugiat recusandae maxime vero eius aut minus reprehenderit delectus
              asperiores. Consequuntur blanditiis repudiandae hic libero incidunt perferendis quos? Commodi
              repellat reiciendis ipsa aperiam sequi, illum quae. Fuga, eaque fugiat.
            </Col>
            <Col id='clauseTable' sm='12' md='8'>
              <Table className='mt-4 mb-4' bordered>
                <thead>
                  <tr>
                    <th className='text-center'>#</th>
                    {!insuranceInfo.is_insurance_ng_acct && insuranceInfo.current_clauses?.length && (
                      <th className='text-center'>Current Clauses</th>
                    )}
                    <th className='text-center'>Clauses</th>
                  </tr>
                </thead>
                <tbody>{getClauseTableColumns(insuranceInfo)}</tbody>
              </Table>
            </Col>
            <Col id='clauseText' sm='12' md='10'>
              Lorem ipsum dolor sit amet consectetur adipisicing elit. Totam aspernatur labore similique eaque
              officiis! Aliquam maiores eaque excepturi praesentium esse consectetur, ab impedit, amet
              veritatis culpa reprehenderit delectus consequatur blanditiis?
            </Col>
            <Col id='premiumChart' sm='12' md='6'>
              <Chart
                height={'400px'}
                chartType='LineChart'
                className='mx-auto'
                loader={<div>Loading Chart</div>}
                data={[['x', 'premium'], ...premiumGraph]}
                options={{
                  hAxis: {
                    title: 'Time (in years)',
                  },
                  vAxis: {
                    title: 'Premium',
                  },
                }}
                rootProps={{ 'data-testid': 'test-premium-chart' }}
              />
            </Col>
            <Col id='coverChart' sm='12' md='6'>
              <Chart
                height={'400px'}
                chartType='LineChart'
                className='mx-auto'
                loader={<div>Loading Chart</div>}
                data={[['x', 'cover'], ...coverGraph]}
                options={{
                  hAxis: {
                    title: 'Time (in years)',
                  },
                  vAxis: {
                    title: 'Cover',
                  },
                  colors: ['red'],
                }}
                rootProps={{ 'data-testid': 'test-premium-chart' }}
              />
            </Col>
            <Col id='bottomText' className='mb-5' sm='12' md='10'>
              Tenetur ab corrupti delectus? Officia magni eligendi vel itaque quidem eaque laboriosam magnam
              tempore exercitationem atque aliquam, culpa amet, dolor hic similique in qui quaerat velit! Odit
              obcaecati voluptates asperiores! Et non labore omnis est obcaecati quo voluptas magni
              repudiandae suscipit id odit provident quos maiores magnam dignissimos facere itaque cum,
              impedit aspernatur quisquam vero quod ipsa numquam? Harum, animi? Odit quaerat quae, ea mollitia
              aut quibusdam labore inventore temporibus accusamus perferendis cum saepe natus quisquam aperiam
              nam nobis quo! Quas aperiam veritatis maxime a sint ullam adipisci ipsam odio. Pariatur at quis
              porro, ipsum laborum modi. In voluptates cupiditate ad vero ipsam! Quis ex quod voluptate,
              repudiandae doloribus pariatur, veritatis laboriosam dolores facere quasi ab non tempora velit
              laudantium?
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

export default InsuranceInfo;
