import React from 'react';
import { IFIInsurance } from '../../types/IFIData';
import Button from 'react-bootstrap/Button';
import Table from 'react-bootstrap/Table';
import { hourglassWaitIcon, tickIcon } from '../../helpers/svgIcons';

type BasicInfoTableInput = {
  insuranceInfo: IFIInsurance;
  startPurchaseProcess: (uuid: string) => void;
  startClaimProcess: (uuid: string) => void;
};

const BasicInfoTable: React.FC<BasicInfoTableInput> = ({
  insuranceInfo,
  startPurchaseProcess,
  startClaimProcess,
}) => {
  return (
    <>
      <Table className='mt-4' size='md'>
        <tbody>
          <tr>
            <td>Cover</td>
            <td>: {insuranceInfo.offer_cover}</td>
            <td></td>
          </tr>
          <tr>
            <td>Yealy Premium</td>
            <td>: {insuranceInfo.offer_premium}</td>
            <td></td>
          </tr>
          {!insuranceInfo.is_insurance_ng_acct && (
            <tr>
              <td>Purchase</td>
              {insuranceInfo.score >= 0.8 && (
                <>
                  <td>: Pre-approved {tickIcon()} </td>
                  <td>
                    <Button
                      variant='outline-primary'
                      className='btn-sm'
                      onClick={() => startPurchaseProcess(insuranceInfo.uuid)}>
                      Buy
                    </Button>
                  </td>
                </>
              )}
              {insuranceInfo.score < 0.8 && (
                <>
                  <td>
                    :{' '}
                    <Button
                      variant='outline-primary'
                      className='btn-sm'
                      onClick={() => startPurchaseProcess(insuranceInfo.uuid)}>
                      Talk to an Agent
                    </Button>
                  </td>
                  <td></td>
                </>
              )}
            </tr>
          )}
          {!!insuranceInfo.is_insurance_ng_acct && (
            <tr>
              <td>Status</td>
              {insuranceInfo.is_active && (
                <>
                  {!insuranceInfo.is_claimed && (
                    <>
                      <td>: Active {tickIcon()} </td>
                      <td>
                        <Button
                          variant='outline-primary'
                          className='btn-sm'
                          onClick={() => startClaimProcess(insuranceInfo.uuid)}>
                          Initiate Claim
                        </Button>
                      </td>
                    </>
                  )}
                  {insuranceInfo.is_claimed && (
                    <>
                      <td>: Claim in progress {hourglassWaitIcon('0 0 16 16')}</td>
                      <td></td>
                    </>
                  )}
                </>
              )}
              {!insuranceInfo.is_active && (
                <>
                  <td>
                    :{' '}
                    <Button
                      disabled
                      variant='outline-secondary'
                      className='btn-sm'
                      onClick={() => startPurchaseProcess(insuranceInfo.uuid)}>
                      We will contact you soon
                    </Button>
                  </td>
                  <td></td>
                </>
              )}
            </tr>
          )}
        </tbody>
      </Table>
    </>
  );
};

export { BasicInfoTable };
