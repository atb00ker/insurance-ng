import React from 'react';
import Col from 'react-bootstrap/Col';
import Row from 'react-bootstrap/Row';
import Image from 'react-bootstrap/Image';
import Container from 'react-bootstrap/esm/Container';
import { FeatureDisplayType } from '../../types/react-component-input-types';

const FeatureDisplayWhite: React.FC<FeatureDisplayType> = ({ image, title, description, imageWidth }) => {
  return (
    <Row className='cover-screen bg-off-white'>
      <Col sm='12'>
        <Container className='h-100 max-width-960'>
          <Row className='h-100'>
            <Col className='cover-screen-flex' xs='12' md='8'>
              <Image className='cover-screen-flex' src={image} alt={title} width={imageWidth} />
            </Col>
            <Col xs='12' md='4'>
              <div className='cover-screen-flex'>
                <Row>
                  <Col xs='12'>
                    <h2 className='text-primary'>{title}</h2>
                  </Col>
                  <Col xs='12'>{description}</Col>
                </Row>
              </div>
            </Col>
          </Row>
        </Container>
      </Col>
    </Row>
  );
};

export default FeatureDisplayWhite;
