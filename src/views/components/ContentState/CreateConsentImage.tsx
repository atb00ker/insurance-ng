import React from 'react';
import EmptyInformationDashboard from '../../assets/illustrations/empty-information-dashboard-2.svg';
import IContentStateImages from './IContentStateImages';

const CreateConsentImage: React.FC<IContentStateImages> = ({ height, width, imgHeight }) => {
  return (
    <>
      <div className='text-center' style={{ display: 'flex', height: height }}>
        <div className='text-center' style={{ margin: 'auto' }}>
          <img style={{ height: imgHeight, width: width }} src={EmptyInformationDashboard} alt='Server Down' />
          <div className='text-secondary max-width-720 text-center mt-5'>To provide services, we need your financial information accociated with your phone number.
            <br />When you provide your phone number and submit, we will show you all the data
            we need to access and you'll have the choice to provide consent to our request.</div>
        </div>
      </div>
    </>
  );
};

export default CreateConsentImage;
