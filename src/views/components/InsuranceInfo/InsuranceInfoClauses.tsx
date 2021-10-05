import React, { useMemo } from 'react';
import Col from 'react-bootstrap/Col';
import Table from 'react-bootstrap/esm/Table';
import Chart from 'react-google-charts';
import { IFIInsurance } from '../../interfaces/IFIData';

const InsuranceInfoClauses: React.FC<{ insuranceInfo: IFIInsurance }> = ({ insuranceInfo }) => {
  const getClauseTableColumns = (insuranceInfo: IFIInsurance): React.ReactNode => {
    const noOfClauses =
      insuranceInfo.clauses?.length > insuranceInfo.current_clauses?.length
        ? insuranceInfo.clauses?.length
        : insuranceInfo.current_clauses?.length || insuranceInfo.clauses?.length;
    let tableRow = [];
    for (let index = 0; index < noOfClauses; index++) {
      tableRow.push(
        <tr key={index}>
          <td>{index + 1}.</td>
          {!insuranceInfo.is_insurance_ng_acct && insuranceInfo.current_clauses?.length && (
            <td
              dangerouslySetInnerHTML={{
                __html: insuranceInfo.current_clauses[index]?.replace(
                  /(\d+%?)|(\S[A-Z]+\S)/g,
                  function (value) {
                    if (value === 'SLA') return value;
                    return `<span class="text-danger">${value}</span>`;
                  },
                ),
              }}
            ></td>
          )}
          <td
            dangerouslySetInnerHTML={{
              __html: insuranceInfo.clauses[index]?.replace(/(\d+%?)|(\S[A-Z]+\S)/g, function (value) {
                if (value === 'SLA') return value;
                return `<span class="text-success">${value}</span>`;
              }),
            }}
          ></td>
        </tr>,
      );
    }
    return tableRow;
  };

  return (
    <Table className='mt-4 mb-4' bordered>
      <thead>
        <tr>
          <th className='text-center'>#</th>
          {!insuranceInfo.is_insurance_ng_acct && insuranceInfo.current_clauses?.length && (
            <th className='text-center'>Current Clauses</th>
          )}
          <th className='text-center'>Clauses</th>
        </tr>
      </thead>
      <tbody>{getClauseTableColumns(insuranceInfo)}</tbody>
    </Table>
  );
};

export default InsuranceInfoClauses;
