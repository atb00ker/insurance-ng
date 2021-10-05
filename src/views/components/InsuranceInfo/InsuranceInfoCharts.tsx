import React, { useMemo } from 'react';
import Col from 'react-bootstrap/Col';
import Chart from 'react-google-charts';
import { IFIInsurance } from '../../types/IFIData';

const InsuranceInfoCharts: React.FC<{ insuranceInfo: IFIInsurance }> = ({ insuranceInfo }) => {
  const numberOfChartRecords = 6;
  const premiumGraph: [number, number][] = useMemo(() => {
    let deductionRate = insuranceInfo.yoy_deduction_rate;
    const graphColumns: [number, number][] = [];
    for (let index = 1; index < numberOfChartRecords; index++) {
      graphColumns.push([index, insuranceInfo.offer_premium - deductionRate]);
      deductionRate = deductionRate + deductionRate;
    }
    return graphColumns;
  }, [insuranceInfo]);

  const coverGraph: [number, number][] = useMemo(() => {
    let incrementRate = insuranceInfo.yoy_deduction_rate;
    const graphColumns: [number, number][] = [];
    for (let index = 1; index < numberOfChartRecords; index++) {
      graphColumns.push([index, insuranceInfo.offer_cover + incrementRate]);
      incrementRate = incrementRate + incrementRate;
    }
    return graphColumns;
  }, [insuranceInfo]);

  return (
    <>
      <Col id='premiumChart' sm='12' md='6'>
        <Chart
          height={'400px'}
          chartType='LineChart'
          className='mx-auto'
          loader={<div>Loading Chart</div>}
          data={[['x', 'premium'], ...premiumGraph]}
          options={{
            hAxis: {
              title: 'Time (in years)',
            },
            vAxis: {
              title: 'Premium',
            },
          }}
          rootProps={{ 'data-testid': 'test-premium-chart' }}
        />
      </Col>
      <Col id='coverChart' sm='12' md='6'>
        <Chart
          height={'400px'}
          chartType='LineChart'
          className='mx-auto'
          loader={<div>Loading Chart</div>}
          data={[['x', 'cover'], ...coverGraph]}
          options={{
            hAxis: {
              title: 'Time (in years)',
            },
            vAxis: {
              title: 'Cover',
            },
            colors: ['red'],
          }}
          rootProps={{ 'data-testid': 'test-premium-chart' }}
        />
      </Col>
    </>
  );
};

export { InsuranceInfoCharts };
