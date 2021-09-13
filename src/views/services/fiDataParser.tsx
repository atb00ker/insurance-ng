import { IFIData } from '../interfaces/IFIData';

var XMLParser = require('react-xml-parser');

export function prepareDataJson(response: IFIData[]) {
  // This should be done in the backend, but I don't want to make
  // golang types all the data right now.
  response.map(element => {
    element.RahasyaData.map(item => {
      if (item.data) item.data = new XMLParser().parseFromString(item.data);
      return item;
    });
  });
  return response;
}
