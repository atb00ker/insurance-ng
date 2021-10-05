import React from 'react';
import FiDataWaitSvg from '../../assets/illustrations/factory-illustration-1.svg';
import { IContentStateImages } from '../../types/react-component-input-types';

const FiDataWait: React.FC<IContentStateImages> = ({ height, width, imgHeight }) => {
  return (
    <>
      <div className='text-center d-flex' style={{ height: height }}>
        <div className='text-center m-auto'>
          <img style={{ height: imgHeight, width: width }} src={FiDataWaitSvg} alt='Processing data' />
          <div data-testid='page-error-message' className='mt-1'>
            We have asked the financial institutions for providing your data. <br />
            This page will refresh as soon as we have the data.
          </div>
        </div>
      </div>
    </>
  );
};

export { FiDataWait };
