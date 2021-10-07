import React, { useContext } from 'react';
import Col from 'react-bootstrap/Col';
import Row from 'react-bootstrap/Row';
import Image from 'react-bootstrap/Image';
import CompleteLogoImage from '../../assets/icons/complete-logo.jpg';
import './NgFeatures.scss';
import Button from 'react-bootstrap/Button';
import Container from 'react-bootstrap/Container';
import { FeatureDisplayOffWhite } from './FeatureDisplayOffWhite';
import { FeatureDisplayWhite } from './FeatureDisplayWhite';
import EmptyInformationDashboard3 from '../../assets/illustrations/empty-information-dashboard-3.svg';
import SearchPage1 from '../../assets/illustrations/search-page-1.svg';
import SaveMoney1 from '../../assets/illustrations/save-money-1.svg';
import ClauseDiscovery1 from '../../assets/illustrations/clause-discovery-1.svg';
import SteppingToFuture1 from '../../assets/illustrations/stepping-to-future-1.svg';
import { rightArrowInCircle } from '../../helpers/svgIcons';
import { Footer } from '../Common/Footer';
import { AuthContext } from '../Auth/AuthProvider';
import { IAuth } from '../../types/IUser';

const NgFeatures: React.FC = () => {
  const auth: IAuth = useContext(AuthContext);
  return (
    <Container className='overflow-hidden' fluid>
      <Row className='cover-screen'>
        <Col sm='12' className='d-flex-center'>
          <Container className='d-flex-center'>
            <Row className='d-flex-center'>
              <Col sm='12' md='5' className='text-center-md mt-4 max-width-960'>
                <h1>
                  The future of <span className='text-primary'>Insurance Services</span>
                </h1>
                <h5>Insurance management experience tailored to your lifestyle.</h5>
                {auth.isReady && !auth.isAuthenticated && (
                  <Button
                    data-testid='features-page-button'
                    onClick={() => auth.loginWithRedirect()}
                    className='d-inline'
                    variant='primary'>
                    Sign Up {rightArrowInCircle('0 0 16 16')}
                  </Button>
                )}
              </Col>
              <Col
                data-testid='features-page-logo'
                className='text-center'
                xs={{ order: 'first', span: '12' }}
                md={{ order: 'last', span: '7' }}>
                <Image className='d-md-none' src={CompleteLogoImage} width={'90%'} />
                <Image className='d-none d-md-block d-lg-none' src={CompleteLogoImage} width={'400px'} />
                <Image className='d-none d-lg-block' src={CompleteLogoImage} width={'500px'} />
              </Col>
            </Row>
          </Container>
        </Col>
      </Row>
      <FeatureDisplayOffWhite
        image={EmptyInformationDashboard3}
        imageWidth={'450px'}
        title='Variable Premiums'
        description='We offer lower premiums which adjusts to your lifestyle. If you are a low risk taker, why should your premiums be high?'
      />
      <FeatureDisplayWhite
        image={SearchPage1}
        imageWidth={'550px'}
        title='Tailored Dashboard'
        description="Don't waste time searching what you want, we will find it for you and provide an experience created uniquely for you."
      />
      <FeatureDisplayOffWhite
        image={SaveMoney1}
        imageWidth={'550px'}
        title='Find Better Plans'
        description='Do you have an existing insurance, if there exists a better plan in the category, we will suggest it to you.'
      />
      <FeatureDisplayWhite
        image={ClauseDiscovery1}
        imageWidth={'550px'}
        title='Clause Discovery'
        description="No need to go through 18 pages of terms and conditions, clearly see what's applicable for you and for what you are paying."
      />
      <FeatureDisplayOffWhite
        image={SteppingToFuture1}
        imageWidth={'650px'}
        title='See the future'
        description='Your premiums change with time, see the future predictions for upto 5 years, predicted based on your financial status of the past 5 years.'
      />
      <Row className='cover-screen'>
        <Col sm='12' className='d-flex-center'>
          <Container className='d-flex-center'>
            <Row>
              <Col sm='12' className='d-flex-center max-width-960'>
                <h1>Much more...</h1>
              </Col>
              <Col sm='12' className='d-flex-center text-center max-width-960'>
                <h5>
                  No long forms, full control on your data, quick and simple user experience. <br />
                  So, what are you waiting for?
                  <br />
                  {auth.isReady && !auth.isAuthenticated && (
                    <a className='href-no-underline' onClick={() => auth.loginWithRedirect()}>
                      Sign up now and get started {rightArrowInCircle('0 2 20 16')}
                    </a>
                  )}
                </h5>
              </Col>
            </Row>
          </Container>
        </Col>
      </Row>
      <Footer />
    </Container>
  );
};

export { NgFeatures };
