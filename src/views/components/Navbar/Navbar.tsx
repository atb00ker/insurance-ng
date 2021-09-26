import React, { useContext } from 'react';
import { default as BootstrapNav } from 'react-bootstrap/Navbar';
import Container from 'react-bootstrap/Container';
import Button from 'react-bootstrap/Button';
import Image from 'react-bootstrap/Image';

import { AuthContext } from '../Auth/AuthProvider';
import { IAuth } from '../../interfaces/IUser';
import { Link } from 'react-router-dom';
import LogoThinImage from '../../assets/icons/logo-thin.png';
import { RouterPath } from '../../enums/UrlPath';

const Navbar = () => {
  const auth: IAuth = useContext(AuthContext);

  return (
    <React.Fragment>
      <BootstrapNav bg='light' expand='sm'>
        <Container fluid>
          <BootstrapNav.Brand href='#home'>
            <Image src={LogoThinImage} height='26px' />
          </BootstrapNav.Brand>
          <BootstrapNav.Toggle aria-controls='navbar-toggle' />
          <BootstrapNav.Collapse id='navbar-toggle'>
            {auth.isReady && auth.isAuthenticated && (
              <>
                <Link to={RouterPath.CreateConsent}>
                  <Button className='m-1 me-1 btn-sm' variant='primary' disabled>
                    Home
                  </Button>
                </Link>
                <Link to={RouterPath.Dashboard}>
                  <Button className='m-1 me-1 btn-sm' variant='primary' disabled>
                    Dashboard
                  </Button>
                </Link>
                {/* <Link to={RouterPath.About}>
                  <Button className='m-1 me-1 btn-sm' variant='primary' disabled>
                    About
                  </Button>
                </Link> */}
              </>
            )}
            <div style={{ flex: '1 1 auto' }}></div>
            {!auth.isReady && (
              <Button
                className='m-1 me-4 btn-sm'
                variant='primary'
                disabled
                onClick={() => auth.loginWithRedirect()}
              >
                Loading
              </Button>
            )}
            {auth.isReady && !auth.isAuthenticated && (
              <Button className='m-1 me-4 btn-sm' variant='primary' onClick={() => auth.loginWithRedirect()}>
                Login
              </Button>
            )}
            {auth.isReady && auth.isAuthenticated && (
              <Button className='m-1 me-4 btn-sm' variant='danger' onClick={() => auth.logout()}>
                Logout
              </Button>
            )}
          </BootstrapNav.Collapse>
        </Container>
      </BootstrapNav>
    </React.Fragment>
  );
};

export default Navbar;
