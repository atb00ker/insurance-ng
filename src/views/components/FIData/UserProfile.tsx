import React from 'react';
import Col from 'react-bootstrap/Col';
import Card from 'react-bootstrap/Card';
import Table from 'react-bootstrap/Table';
import OverlayTrigger from 'react-bootstrap/OverlayTrigger';
import Tooltip from 'react-bootstrap/Tooltip';
import { IFIUserData } from '../../interfaces/IFIData';
import { getIconForScore, getQuestionIcon } from './getUserProfileSvgIcons';

const UserProfile: React.FC<{ fiData: IFIUserData }> = ({ fiData }) => {
  return (
    <>
      <Col className='mt-4' sm='10'>
        <Card className='border'>
          <Card.Body>
            <Card.Title>{fiData.name}</Card.Title>
            <Table className='mt-2' hover size='sm'>
              <thead>
                <tr>
                  <th>Affecting Factor</th>
                  <th>Score</th>
                  <th></th>
                </tr>
              </thead>
              <tbody>
                {fiData.insurance.map(insurance => {
                  return (
                    <tr key={insurance.type}>
                      <td>{insurance.title}</td>
                      <td>
                        {insurance.score.toFixed(2)}{' '}
                        <OverlayTrigger
                          placement='bottom'
                          overlay={<Tooltip id='tooltip'>"WEEEEEE"</Tooltip>}
                        >
                          <>{getQuestionIcon()}</>
                        </OverlayTrigger>
                      </td>
                      <td>{getIconForScore(insurance.score)}</td>
                    </tr>
                  );
                })}
              </tbody>
            </Table>
          </Card.Body>
        </Card>
      </Col>
    </>
  );
};

export default UserProfile;
