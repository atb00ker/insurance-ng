import React from 'react';
import Col from 'react-bootstrap/Col';
import Row from 'react-bootstrap/Row';
import Card from 'react-bootstrap/Card';


const InsuranceGraphs: React.FC<any> = ({ data }) => {


  console.log(data)
  return (
  <Row className='mt-1'>

    <Col sm='3'>
      <Card style={{ width: '18rem' }}>
        <Card.Body>
          <Card.Title>Test</Card.Title>
          <Card.Subtitle className="mb-2 text-muted">Card Subtitle</Card.Subtitle>
          <Card.Text>
            Some quick example text to build on the card title and make up the bulk of
            the card's content.
          </Card.Text>
        </Card.Body>
      </Card>
    </Col>
  </Row>
  );
};

export default InsuranceGraphs;
