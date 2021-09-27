import React from 'react';
import Col from 'react-bootstrap/Col';
import { useHistory } from 'react-router-dom';
import { IFIData, IFIInsurance } from '../../interfaces/IFIData';
import Button from 'react-bootstrap/esm/Button';
import { RouterPath } from '../../enums/UrlPath';

const InsuranceInfoTitle: React.FC<{ fiData: IFIData; insuranceInfo: IFIInsurance }> = ({
  fiData,
  insuranceInfo,
}) => {
  const history = useHistory();

  return (
    <>
      <Col sm='1'>
        <Button
          onClick={() => {
            history.push({
              pathname: RouterPath.Dashboard,
              state: fiData,
            });
            return;
          }}
          style={{ minWidth: '80px' }}
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
    </>
  );
};

export default InsuranceInfoTitle;
