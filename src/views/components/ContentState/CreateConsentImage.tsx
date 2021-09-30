import React from 'react';
import IContentStateImages from '../../interfaces/IContentStateImages';
import ProcessSummaryImage from '../../assets/images/process-summary.jpg';
import Image from 'react-bootstrap/Image';
import Row from 'react-bootstrap/esm/Row';
import Col from 'react-bootstrap/esm/Col';
import { Link } from 'react-router-dom';
import { RouterPath } from '../../enums/UrlPath';

const CreateConsentImage: React.FC<IContentStateImages> = ({ width, imgHeight }) => {
  return (
    <Row className='justify-content-center'>
      <Col sm='12' md='10' lg='7' className='mt-4 text-center'>
        Please share your phone number to get the lowest possible prices on all our insurances.
        To provide services, we need your financial information associated with your phone number. We will
        show you all the data we need to access and you'll have the choice to provide consent to our data
        request.
      </Col>
      <Col sm='12' md='10' lg='7' className='mt-4 text-center'>
        <Image src={ProcessSummaryImage} height={imgHeight} width={width} />
      </Col>
      <Col sm='12' md='10' lg='7' className='mt-4 mb-5 text-center'>
        We use the setu account aggregator to collect the data, if you want more information
        about the data collection, processing, storage, security, sharing or deletion or you{' '}
        want to contact us. Please visit{' '}
        <Link className='href-no-underline' to={RouterPath.About}>our about page</Link>.
      </Col>
    </Row>
  );
};

export default CreateConsentImage;
