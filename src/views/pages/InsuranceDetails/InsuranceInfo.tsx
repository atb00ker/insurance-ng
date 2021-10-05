import React, { useContext, useEffect, useState } from 'react';
import Container from 'react-bootstrap/Container';
import Col from 'react-bootstrap/Col';
import Row from 'react-bootstrap/Row';
import { useHistory, useLocation } from 'react-router-dom';
import { IAuth } from '../../interfaces/IUser';
import { AuthContext } from '../../components/Auth/AuthProvider';
import { createClaimRequest, createPurchaseRequest, getDashboardData } from '../../helpers/axios';
import { IFIData, IFIInsurance } from '../../interfaces/IFIData';
import { RouterPath } from '../../enums/UrlPath';
import ServerRequestError from '../../components/ContentState/ServerRequestError';
import SectionLoader from '../../components/ContentState/SectionLoader';
import TitleSection from '../../components/InsuranceInfo/TitleSection';
import BasicInfoTable from '../../components/InsuranceInfo/BasicInfoTable';
import InsuranceInfoCharts from '../../components/InsuranceInfo/InsuranceInfoCharts';
import InsuranceInfoClauses from '../../components/InsuranceInfo/InsuranceInfoClauses';
import { PageState } from '../../enums/PageStates';

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
    // Redirect to login if user is not authenticated.
    if (auth.isReady && !auth.isAuthenticated) auth.loginWithRedirect();
  }, [auth]);

  useEffect(() => {
    if (!fiDataFromHistory) {
      getDataFromServer();
      return;
    }
    changePageState(PageState.Data);
    setInsuranceState(fiDataFromHistory, uuid);
  }, [fiDataFromHistory, auth]);

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

  const startClaimProcess = (uuid: string) => {
    auth.user.jwt().then((jwt: string) => {
      changeStateOnDataResponse(createClaimRequest(uuid, jwt));
    });
  };

  const changeStateOnDataResponse = (dataResponse: Promise<any>) => {
    dataResponse
      .then((response: any) => {
        const data: IFIData = response?.data;
        if (data?.status) {
          changePageState(PageState.Data);
          setInsuranceState(data, uuid);
        } else {
          changePageState(PageState.Error);
        }
      })
      .catch((error: any) => {
        console.error(error);
        changePageState(PageState.Error);
      });
  };

  const changePageState = (state: string) => {
    setShowError(state == PageState.Error);
    setShowLoader(state == PageState.Loading);
  };


  return (
    <Container>
      {!showError && !showLoader && (
        <>
          <Row className='mt-5 mb-2 justify-content-center'>
            <TitleSection fiData={fiData} insuranceInfo={insuranceInfo} />
          </Row>
          <Row className='justify-content-center'>
            <Col id='topTable' sm='8' md='6' lg='5'>
              <BasicInfoTable
                insuranceInfo={insuranceInfo}
                startPurchaseProcess={startPurchaseProcess}
                startClaimProcess={startClaimProcess}
              />
            </Col>
            <Col id='introText' className='mt-4' sm='12' md='10'>
              {insuranceInfo.description} Voluptas nihil sed laborum sequi quaerat veritatis hic facere eaque
              culpa iusto exercitationem tenetur rerum, incidunt, sapiente optio suscipit aspernatur delectus
              molestias explicabo quisquam ratione odit, repudiandae dolor beatae! Harum! Fugiat magnam alias
              veritatis quibusdam iure! Inventore dolorem quaerat illum blanditiis similique facere deserunt
              neque molestias consequuntur eveniet numquam sapiente accusamus ab consequatur eos, libero eius
              nulla. Beatae, debitis maiores! Corporis, veritatis. Id rem tenetur enim, aspernatur fugit
              fugiat ipsum ex possimus amet veniam vel odio tempore modi dolores culpa laboriosam? Nulla,
              voluptate. Dolorum facilis quisquam id atb00ker enim molestiae magni? Consequuntur nam eos,
              voluptates alias ipsa error libero accusamus magnam provident, nihil inventore dolores obcaecati
              saepe, blanditiis amet incidunt. Doloribus explicabo nemo dolore saepe nihil, repellendus
              temporibus iste quam inventore. Ipsam fugiat recusandae maxime vero eius aut minus reprehenderit
              delectus asperiores. Consequuntur blanditiis repudiandae hic libero incidunt perferendis quos?
              Commodi repellat reiciendis ipsa aperiam sequi, illum quae. Fuga, eaque fugiat.
            </Col>
            <Col id='clauseTable' sm='12' md='8'>
              <InsuranceInfoClauses insuranceInfo={insuranceInfo} />
            </Col>
            <Col id='clauseText' sm='12' md='10'>
              Lorem ipsum dolor sit amet consectetur adipisicing elit. Totam aspernatur labore similique eaque
              officiis! Aliquam maiores eaque excepturi praesentium esse consectetur, ab impedit, amet
              veritatis culpa reprehenderit delectus consequatur blanditiis?
            </Col>
            <InsuranceInfoCharts insuranceInfo={insuranceInfo} />
            <Col id='bottomText' className='mb-5' sm='12' md='10'>
              Tenetur ab corrupti delectus? Officia magni eligendi vel itaque quidem eaque laboriosam magnam
              tempore exercitationem atque aliquam, culpa amet, dolor hic similique in qui quaerat velit! Odit
              obcaecati voluptates asperiores! Et non labore omnis est obcaecati quo f014 voluptas magni
              repudiandae suscipit id odit provident quos Ajay maiores magnam dignissimos facere itaque cum,
              impedit aspernatur quisquam vero quod ipsa numquam? Harum, animi? Odit quaerat quae, ea mollitia
              aut quibusdam labore inventore temporibus accusamus perferendis cum saepe natus quisquam aperiam
              nam nobis quo! Quas aperiam veritatis maxime a sint ullam adipisci ipsam odio. Pariatur at quis
              porro, ipsum laborum modi. In voluptates cupiditate ad vero ipsam! Quis ex quod voluptate,
              repudiandae doloribus Tripathi pariatur, veritatis laboriosam dolores facere 03041998 quasi ab
              non tempora velit laudantium?
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
