import React, { useContext } from 'react';
import { default as BootstrapNav } from 'react-bootstrap/Navbar';
import Container from 'react-bootstrap/Container';
import Button from 'react-bootstrap/Button';
import Image from 'react-bootstrap/Image';

import { AuthContext } from '../Auth/AuthProvider';
import { Link } from 'react-router-dom';
import LogoThinImage from '../../assets/icons/logo-thin.png';
import { RouterPath } from '../../enums/UrlPath';
import { IAuth } from '../../types/IUser';

const Navbar = () => {
  const auth: IAuth = useContext(AuthContext);

  return (
    <React.Fragment>
      <BootstrapNav bg='light' expand='sm'>
        <Container fluid>
          <BootstrapNav.Brand target='_self' href={RouterPath.Home}>
            <Image src={LogoThinImage} height='26px' />
          </BootstrapNav.Brand>
          <BootstrapNav.Toggle aria-controls='navbar-toggle' />
          <BootstrapNav.Collapse id='navbar-toggle'>
            <Link to={RouterPath.Home}>
              <Button className='m-1 me-1 btn-sm' variant='primary'>
                Home
              </Button>
            </Link>
            {auth.isReady && auth.isAuthenticated && (
              <>
                <Link to={RouterPath.Dashboard}>
                  <Button className='m-1 me-1 btn-sm' variant='primary'>
                    Dashboard
                  </Button>
                </Link>
              </>
            )}
            <Link to={RouterPath.About}>
              <Button className='m-1 me-1 btn-sm' variant='primary'>
                About
              </Button>
            </Link>
            <div style={{ flex: '1 1 auto' }}></div>
            {!auth.isReady && (
              <Button className='m-1 me-4 btn-sm' variant='primary' onClick={() => auth.loginWithRedirect()}>
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

export { Navbar };
