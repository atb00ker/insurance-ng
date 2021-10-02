import React from 'react';
import Col from 'react-bootstrap/Col';
import Container from 'react-bootstrap/Container';
import Row from 'react-bootstrap/Row';
import Image from 'react-bootstrap/Image';
import CompleteLogoImage from '../../assets/icons/complete-logo.jpg';
import { openInNewTabIcon } from '../../helpers/svgIcons';
import Footer from '../../components/Common/Footer';

const About: React.FC = () => {
  const hackathonUrl: string = 'https://free-your-data.devfolio.co/';
  const authorEmailText: string = 'ajay39in+aa[at]gmail[dot]com';
  const authorEmail: string = 'ajay39in+aa@gmail.com';

  return (
    <Container>
      <Row className='mt-4 mb-5 justify-content-center'>
        <Col sm='12' className='max-width-960'>
          <h3>About</h3>
          <p>
            Insurance NG is the next generation of insurance service with features like variable premium and
            cover, quick in-depth analysis of the financial status of a customer using account aggregator
            framework and much more to outshine the competition. <br />
            Created for the{' '}
            <a className='href-no-underline' href={hackathonUrl}>
              Setu x GitHub Hackathon on Devfolio {openInNewTabIcon('0 0 32 32')}
            </a>
            .
          </p>
        </Col>
        <Col sm='12' className='mt-4 mb-4 max-width-960 text-center'>
          <Image src={CompleteLogoImage} height={130} />
        </Col>
        <Col sm='12' className='max-width-960'>
          <h4>Privacy and Data Policy</h4>
          <p>
            To provide the lowest prices we need your financial information to understand you, your financial
            status and your lifestyle. <br />
            To get the best out of our services, please allow access to your deposit accounts, investment
            accounts, debt information and existing insurances. <br />
          </p>
          <h5>Data Collection</h5>
          <p>
            The data in this application is collected with the help of the Setu{' '}
            <a
              className='href-no-underline'
              href='https://www.thebalance.com/what-is-account-aggregation-1293879'
            >
              Account Agregator {openInNewTabIcon('0 -3 32 32')}
            </a>{' '}
            network in which all your favourite financial institutions participate including apna bank, apna
            insurance, apna bank and apna pension.
          </p>
          <h5>Data Consent Revoke</h5>
          <p>
            You can revoke the data consent from the Setu's Account Aggregator application. (Not released yet)
          </p>
          <h5>Data Processed</h5>
          The data processed to know more about you:
          <ul>
            <li>
              Deposit Account transactions: Your transactions from the past 5 years help us identify your
              lifestyle, social and financial health.
            </li>
            <li>Deposit Account current balance: Helps us get a rough idea about your financial status.</li>
            <li>
              Insurance Account transactions: Helps us understand your risk and payout factor, people with
              lower payout factors get additional discounts and benefits.
            </li>
            <li>
              Personal information: Helps us identify you. (Includes: Name, Date of Birth, Address, Pancard).
            </li>
            <li>
              Existing insurance plans: Helps us identify your what's important to you and what better plan we
              can offer you.
            </li>
            <li>
              Investment (SIP / Mutual Funds / Equities / Term deposit / Reoccurring Deposit / Govt Debts /
              PPF / NPS) plan's summary helps us identify your future planning the financial maturity.
            </li>
            <li>
              Debt (credit cards): Your total debt helps us understand your lifestyle and risk appetite.
            </li>
          </ul>
          <h5>Data Storage</h5>
          We will save the following information about you in our database:
          <ul>
            <li>Name</li>
            <li>Date of Birth</li>
            <li>Address</li>
            <li>Pancard</li>
            <li>Existing insurance plans</li>
          </ul>
          <h5>Data Deletion</h5>
          We will delete all your information from our database if you contact us with delete my data request
          at{' '}
          <a className='href-no-underline' href={`mailto:${authorEmail}`}>
            {authorEmailText}
          </a>{' '}
          or when you initiate another data consent request similar to this one. <br />
          <h5 className='mt-3'>Data Security</h5>
          <p>
            Account Aggregator protocol ensures that the data is not visible to anyone except insurance-ng, it
            uses{' '}
            <a className='href-no-underline' href='https://www.ecdhe.com/'>
              ECDHE
            </a>{' '}
            to ensure that even the account aggregator (Setu) cannot view your data. When the data is reached
            to us, we store it in secure servers and we which follow the{' '}
            <a
              className='href-no-underline'
              href='https://en.wikipedia.org/wiki/Principle_of_least_privilege'
            >
              principle of least privilege
            </a>{' '}
            to ensure your privacy.
          </p>
          <h5 className='mt-3'>Data Sharing</h5>
          <p>Your data is not shared with any third-party organization.</p>
        </Col>
        <Col sm='12' className='max-width-960'>
          <h3>Contact Us</h3>
          If you need any further information or support, please contact us at{' '}
          <a className='href-no-underline' href={`mailto:${authorEmail}`}>
            {authorEmailText}
          </a>
          .
        </Col>
      </Row>
      <Footer />
    </Container>
  );
};

export default About;
