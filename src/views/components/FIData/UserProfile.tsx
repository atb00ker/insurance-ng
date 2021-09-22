import React from 'react';
import Col from 'react-bootstrap/Col';
import Card from 'react-bootstrap/Card';
import Image from 'react-bootstrap/Image';
import Table from 'react-bootstrap/Table';
import OverlayTrigger from 'react-bootstrap/OverlayTrigger';
import Tooltip from 'react-bootstrap/Tooltip';
import { IFIUserData } from '../../interfaces/IFIData';
import Row from 'react-bootstrap/esm/Row';
import ProfileImage from './../../assets/images/default-profile-picture.jpeg';
import { errorIcon, notApplicableIcon, questionIcon, tickIcon, warnIcon } from '../../services/svgIcons';

const UserProfile: React.FC<{ fiData: IFIUserData }> = ({ fiData }) => {

  const getIconForScore = (score: number) => {
    if (score > 0.78) {
      return tickIcon();
    } else if (score >= 0.7) {
      return warnIcon();
    } else if (score > 0.3) {
      return errorIcon();
    } else {
      return notApplicableIcon();
    }
  };

  return (
    <>
      <Col className='mt-4' sm='10'>
        <Card className='border'>
          <Card.Body>
            <Row className='justify-content-center'>
              <Col sm='12' md='4' className="vertical-center-relative-image text-center">
                <Image style={{ marginTop: '20px' }} src={ProfileImage}
                       height="150px" roundedCircle />
              </Col>
              <Col sm='12' md='8' className="mt-3">
                <Card.Title>{fiData.name}</Card.Title>
                <table>
                  <tbody>
                  <tr>
                    <td style={{ width: 160 }}>Date of Birth</td>
                    <td>: <span style={{ paddingLeft: 30 }}>
                      {new Date(fiData.date_of_birth).toDateString()}
                    </span></td>
                  </tr>
                  <tr>
                    <td>Pancard</td>
                    <td>: <span style={{ paddingLeft: 30 }}>{fiData.pancard}</span></td>
                  </tr>
                  <tr>
                    <td>Active Insurances</td>
                    <td>: <span style={{ paddingLeft: 30 }}>
                      {fiData.insurance.filter(insurance => insurance.account_id != "").length}
                    </span></td>
                  </tr>
                  <tr>
                    <td>Account Status</td>
                    <td>: <span style={{ paddingLeft: 30 }}>
                      Active {tickIcon()}
                    </span></td>
                  </tr>
                  <tr>
                    <td>KYC Status</td>
                    <td>: <span style={{ paddingLeft: 30 }}>
                      {fiData.ckyc_compliance ? <span>Completed {tickIcon()}</span> :
                      <span>Incomplete {warnIcon()}</span> }
                    </span></td>
                  </tr>
                  </tbody>
                </table>
              </Col>
              <Col sm='12' className="mt-4">
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
                          overlay={<Tooltip id='tooltip'>WEEEEEE</Tooltip>}
                        >
                          <>{questionIcon()}</>
                        </OverlayTrigger>
                      </td>
                      <td>{getIconForScore(insurance.score)}</td>
                    </tr>
                  );
                })}
              </tbody>
            </Table>
              </Col>
            </Row>
          </Card.Body>
        </Card>
      </Col>
    </>
  );
};

export default UserProfile;
