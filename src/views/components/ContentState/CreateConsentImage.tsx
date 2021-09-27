import React from 'react';
import IContentStateImages from '../../interfaces/IContentStateImages';
import ProcessSummaryImage from '../../assets/images/process-summary.jpg';
import Image from 'react-bootstrap/Image';
import Row from 'react-bootstrap/esm/Row';
import Col from 'react-bootstrap/esm/Col';

const CreateConsentImage: React.FC<IContentStateImages> = ({ height, width, imgHeight }) => {
  return (
      <Row className='justify-content-center' style={{ display: 'flex', height: height }}>
        <Col sm="12" md="10" lg="7" className='mt-4 text-center' style={{ margin: 'auto' }}>
          Please share your phone number to get the lowest possible prices on all
          our insurances.
        </Col>
        <Col sm="12" md="10" lg="7" className='mt-4 text-center'>
          <Image src={ProcessSummaryImage} height={imgHeight} width={width} />
        </Col>
        <Col sm="12" md="10" lg="7" className='mt-4 text-center' style={{ margin: 'auto' }}>
          To provide services, we need your financial information associated with your phone number.
          We will show you all the data we need to access and you'll have the choice to provide consent
          to our data request.
        </Col>
      </Row>
  );
};

export default CreateConsentImage;
