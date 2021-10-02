import React from 'react';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import { heartFill } from '../../helpers/svgIcons';

const Footer = () => {
  const authorGitHub = 'https://github.com/atb00ker';

  return (
    <Row>
      <Col sm='12' className='mb-3 text-center'>
        Made with <span className='text-danger'>{heartFill()}</span> by{' '}
        <a className='href-no-underline' href={authorGitHub}>
          @atb00ker
        </a>
      </Col>
    </Row>
  );
};

export default Footer;
