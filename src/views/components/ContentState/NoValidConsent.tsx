import React from 'react';
import { Link } from 'react-router-dom';
import NoValidConsentImage from '../../assets/illustrations/empty-information-dashboard-1.svg';
import { RouterPath } from '../../enums/UrlPath';
import IContentStateImages from '../../interfaces/IContentStateImages';

const NoValidConsent: React.FC<IContentStateImages> = ({ height, width, imgHeight }) => {
  return (
    <>
      <div className='text-center d-flex' style={{ height: height }}>
        <div className='text-center m-auto'>
          <img style={{ height: imgHeight, width: width }} src={NoValidConsentImage} alt='Server Down' />
          <div data-testid='page-error-message' className='mt-4'>
            We need your financial information to provide you with data. <br />
            Please head to the
            <Link className='href-no-underline' to={RouterPath.CreateConsent}>
              {' '}
              home page{' '}
            </Link>{' '}
            and provide your insurance and deposit account details.
          </div>
        </div>
      </div>
    </>
  );
};

export default NoValidConsent;
